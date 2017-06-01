package suggest

/*
 * inspired by
 * http://www.aaai.org/ocs/index.php/AAAI/AAAI10/paper/viewFile/1939/2234
 * http://nlp.stanford.edu/IR-book/html/htmledition/k-gram-indexes-for-wildcard-queries-1.html
 * http://bazhenov.me/blog/2012/08/04/autocomplete.html
 * http://www.aclweb.org/anthology/C10-1096
 */

import (
	"container/heap"
)

type invertedListsT map[int][]int

// NGramIndex is structure ... describe me please
type NGramIndex struct {
	clean      *cleaner
	indices    []invertedListsT
	dictionary []WordKey
	config     *IndexConfig
}

// NewNGramIndex returns a new NGramIndex with given config
func NewNGramIndex(config *IndexConfig) *NGramIndex {
	clean := newCleaner(config.alphabet.Chars(), config.pad, config.wrap)
	return &NGramIndex{
		clean, make([]invertedListsT, 0), make([]WordKey, 0), config,
	}
}

// AddWord add given word to invertedList
func (n *NGramIndex) AddWord(word string, key WordKey) {
	prepared := n.prepareString(word)
	set := n.getNGramSet(prepared)
	cardinality := len(set)

	if len(n.indices) <= cardinality {
		tmp := make([]invertedListsT, cardinality+1, cardinality*2)
		copy(tmp, n.indices)
		n.indices = tmp
	}

	invertedLists := n.indices[cardinality]
	if invertedLists == nil {
		invertedLists = make(invertedListsT)
		n.indices[cardinality] = invertedLists
	}

	keyToIndex := len(n.dictionary)
	for _, index := range set {
		invertedLists[index] = append(invertedLists[index], keyToIndex)
	}

	n.dictionary = append(n.dictionary, key)
}

// Suggest returns top-k similar strings
func (n *NGramIndex) Suggest(config *SearchConfig) []Candidate {
	result := make([]Candidate, 0, config.topK)
	preparedQuery := n.prepareString(config.query)
	if len(preparedQuery) < n.config.ngramSize {
		return result
	}

	candidates := n.search(preparedQuery, config)
	for candidates.Len() > 0 {
		r := heap.Pop(candidates).(*rank)
		result = append(
			[]Candidate{{n.dictionary[r.id], r.distance}},
			result...,
		)
	}

	return result
}

func (n *NGramIndex) search(query string, config *SearchConfig) *heapImpl {
	set := n.getNGramSet(query)
	sizeA := len(set)

	metric := config.metric
	similarity := config.similarity
	topK := config.topK

	h := &heapImpl{}
	bMin, bMax := metric.MinY(similarity, sizeA), metric.MaxY(similarity, sizeA)
	rid := make([][]int, 0, sizeA)
	lenIndices := len(n.indices)

	if bMax >= lenIndices {
		bMax = lenIndices - 1
	}

	for sizeB := bMax; sizeB >= bMin; sizeB-- {
		threshold := metric.Threshold(similarity, sizeA, sizeB)
		if threshold == 0 {
			continue
		}

		// reset slice
		rid = rid[:0]
		invertedLists := n.indices[sizeB]
		// maximum allowable ngram miss count
		allowedSkips := sizeA - threshold + 1
		for _, index := range set {
			// there is no reason to continue, because of threshold
			if allowedSkips == 0 {
				break
			}

			list := invertedLists[index]
			if len(list) > 0 {
				rid = append(rid, list)
			} else {
				allowedSkips--
			}
		}

		if len(rid) < threshold {
			continue
		}

		counts := n.getCounts(rid, threshold)
		// use heap search for finding top k items in a list efficiently
		// see http://stevehanov.ca/blog/index.php?id=122
		var r *rank
		for inter := len(counts) - 1; inter >= threshold; inter-- {
			for _, id := range counts[inter] {
				distance := metric.Distance(inter, sizeA, sizeB)

				if h.Len() < topK || h.Top().(*rank).distance > distance {
					if h.Len() == topK {
						r = heap.Pop(h).(*rank)
					} else {
						r = &rank{0, 0}
					}

					r.id = id
					r.distance = distance
					heap.Push(h, r)
				}
			}
		}
	}

	return h
}

// TODO describe me!
func (n *NGramIndex) getCounts(rid [][]int, threshold int) [][]int {
	if threshold == 1 {
		return scanCount(rid, threshold)
	}

	return mergeSkip(rid, threshold)
}

// Return unique ngrams
func (n *NGramIndex) getNGramSet(word string) []int {
	ngrams := SplitIntoNGrams(word, n.config.ngramSize)
	set := make(map[int]struct{}, len(ngrams))
	list := make([]int, 0, len(ngrams))
	for _, ngram := range ngrams {
		index := n.ngramToIndex(ngram)
		_, found := set[index]
		set[index] = struct{}{}
		if !found {
			list = append(list, index)
		}
	}

	return list
}

// Map ngram to int (index)
func (n *NGramIndex) ngramToIndex(ngram string) int {
	index := 0
	alphabet := n.config.alphabet
	size := alphabet.Size()
	for _, char := range ngram {
		i := alphabet.MapChar(char)
		if index == InvalidChar {
			panic("Invalid char was detected")
		}

		index = index*size + i
	}

	return index
}

// Prepare string for indexing
func (n *NGramIndex) prepareString(word string) string {
	return n.clean.cleanAndWrap(word)
}

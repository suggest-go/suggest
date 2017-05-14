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

type NGramIndex struct {
	clean      *cleaner
	indices    []invertedListsT
	dictionary []WordKey
	config     *IndexConfig
}

func NewNGramIndex(config *IndexConfig) *NGramIndex {
	clean := newCleaner(config.alphabet.Chars(), config.pad, config.wrap)
	return &NGramIndex{
		clean, make([]invertedListsT, 0), make([]WordKey, 0), config,
	}
}

// Add given word to invertedList
func (self *NGramIndex) AddWord(word string, key WordKey) {
	prepared := self.prepareString(word)
	set := self.getNGramSet(prepared)
	cardinality := len(set)

	if len(self.indices) <= cardinality {
		tmp := make([]invertedListsT, cardinality+1, cardinality*2)
		copy(tmp, self.indices)
		self.indices = tmp
	}

	invertedLists := self.indices[cardinality]
	if invertedLists == nil {
		invertedLists = make(invertedListsT)
		self.indices[cardinality] = invertedLists
	}

	keyToIndex := len(self.dictionary)
	for _, index := range set {
		invertedLists[index] = append(invertedLists[index], keyToIndex)
	}

	self.dictionary = append(self.dictionary, key)
}

// Return top-k similar strings
func (self *NGramIndex) Suggest(config *SearchConfig) []WordKey {
	result := make([]WordKey, 0, config.topK)
	preparedQuery := self.prepareString(config.query)
	if len(preparedQuery) < self.config.ngramSize {
		return result
	}

	candidates := self.search(preparedQuery, config)
	for candidates.Len() > 0 {
		r := heap.Pop(candidates).(*rank)
		result = append([]WordKey{self.dictionary[r.id]}, result...)
	}

	return result
}

func (self *NGramIndex) search(query string, config *SearchConfig) *heapImpl {
	set := self.getNGramSet(query)
	sizeA := len(set)

	mm := getMeasure(config.measureName)
	similarity := config.similarity
	topK := config.topK

	h := &heapImpl{}
	bMin, bMax := mm.minY(similarity, sizeA), mm.maxY(similarity, sizeA)
	rid := make([][]int, 0, sizeA)
	lenIndices := len(self.indices)

	if bMax >= lenIndices {
		bMax = lenIndices - 1
	}

	for sizeB := bMax; sizeB >= bMin; sizeB-- {
		threshold := mm.threshold(similarity, sizeA, sizeB)
		if threshold == 0 {
			continue
		}

		// reset slice
		rid = rid[:0]
		invertedLists := self.indices[sizeB]
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

		counts := mergeSkip(rid, threshold)
		// use heap search for finding top k items in a list efficiently
		// see http://stevehanov.ca/blog/index.php?id=122
		for inter := len(counts) - 1; inter >= threshold; inter-- {
			for _, id := range counts[inter] {
				distance := mm.distance(inter, sizeA, sizeB)

				if h.Len() < topK || h.Top().(*rank).distance > distance {
					var r *rank
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

// Return unique ngrams
func (self *NGramIndex) getNGramSet(word string) []int {
	ngrams := SplitIntoNGrams(word, self.config.ngramSize)
	set := make(map[int]struct{}, len(ngrams))
	list := make([]int, 0, len(ngrams))
	for _, ngram := range ngrams {
		index := self.ngramToIndex(ngram)
		_, found := set[index]
		set[index] = struct{}{}
		if !found {
			list = append(list, index)
		}
	}

	return list
}

// Map ngram to int (index)
func (self *NGramIndex) ngramToIndex(ngram string) int {
	index := 0
	alphabet := self.config.alphabet
	size := alphabet.Size()
	for _, char := range ngram {
		i := alphabet.MapChar(char)
		if index == INVALID_CHAR {
			panic("Invalid char was detected")
		}

		index = index*size + i
	}

	return index
}

// Prepare string for indexing
func (self *NGramIndex) prepareString(word string) string {
	return self.clean.cleanAndWrap(word)
}

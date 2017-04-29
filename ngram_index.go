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
	k          int
	alphabet   Alphabet
	indices    []invertedListsT
	dictionary []string
	index      int
	config     *conf
}

type conf struct {
	threshold   float64 // 0 - 1
	measureName MeasureT
	pad         string
	wrap        string
}

var defaultConf *conf

func init() {
	defaultConf = &conf{0.5, COSINE, "$", "$"}
}

func NewNGramIndex(k int) *NGramIndex {
	if k < 2 || k > 4 {
		panic("k should be in [2, 4]")
	}

	// TODO declare as constructor argument
	alphabet := NewCompositeAlphabet([]Alphabet{
		NewEnglishAlphabet(),
		NewNumberAlphabet(),
		NewRussianAlphabet(),
		// TODO use config.pad here
		NewSimpleAlphabet([]rune{'$'}),
	})

	return &NGramIndex{
		k, alphabet, make([]invertedListsT, 0), make([]string, 0), 0, defaultConf,
	}
}

// Add given word to invertedList
func (self *NGramIndex) AddWord(word string) {
	prepared := self.prepareString(word)
	ngrams := self.getNGramSet(prepared)
	cardinality := len(ngrams)

	if len(self.indices) <= cardinality {
		tmp := make([]invertedListsT, cardinality, cardinality*2)
		copy(tmp, self.indices)
		self.indices = tmp
	}

	invertedLists := self.indices[cardinality-1]
	if invertedLists == nil {
		invertedLists = make(invertedListsT)
		self.indices[cardinality-1] = invertedLists
	}

	indexes := self.ngramSetToIndexSet(ngrams)
	for _, index := range indexes {
		invertedLists[index] = append(invertedLists[index], self.index)
	}

	self.dictionary = append(self.dictionary, word)
	self.index++
}

// Return top-k similar strings
func (self *NGramIndex) Suggest(word string, topK int) []string {
	result := make([]string, 0, topK)
	preparedWord := self.prepareString(word)
	if len(preparedWord) < self.k {
		return result
	}

	candidates := self.search(preparedWord, topK)
	for candidates.Len() > 0 {
		r := heap.Pop(candidates).(*rank)
		result = append([]string{self.dictionary[r.id]}, result...)
	}

	return result
}

func (self *NGramIndex) search(word string, topK int) *heapImpl {
	set := self.getNGramSet(word)
	sizeA := len(set)
	indexes := self.ngramSetToIndexSet(set)

	h := &heapImpl{}
	mm := getMeasure(self.config.measureName)
	alpha := self.config.threshold
	bMin, bMax := mm.minY(alpha, sizeA), mm.maxY(alpha, sizeA)
	for sizeB := bMax; sizeB >= bMin; sizeB-- {
		if len(self.indices) <= sizeB {
			continue
		}

		// find max word id for memory optimize
		rid := make([][]int, 0)
		invertedLists := self.indices[sizeB]
		for _, index := range indexes {
			list := invertedLists[index]
			if len(list) > 0 {
				rid = append(rid, list)
			}
		}

		t := mm.threshold(alpha, sizeA, sizeB)
		if len(rid) < t || t == 0 {
			continue
		}

		counts := divideSkip(rid, t)
		// use heap search for finding top k items in a list efficiently
		// see http://stevehanov.ca/blog/index.php?id=122
		for inter := len(counts) - 1; inter >= 0; inter-- {
			if len(counts[inter]) == 0 {
				continue
			}

			list := counts[inter]
			for _, id := range list {
				distance := mm.distance(inter, sizeA, sizeB)
				if h.Len() < topK || h.Top().(*rank).distance > distance {
					if h.Len() == topK {
						heap.Pop(h)
					}

					heap.Push(h, &rank{id, distance})
				}
			}
		}
	}

	return h
}

// Return unique ngrams with frequency
func (self *NGramIndex) getNGramSet(word string) []string {
	return GetNGramSet(word, self.k)
}

func (self *NGramIndex) ngramSetToIndexSet(set []string) []int {
	result := make([]int, 0, len(set))
	for _, ngram := range set {
		index := 0
		for _, char := range ngram {
			i := self.alphabet.MapChar(char)
			if index == INVALID_CHAR {
				panic("Invalid char was detected")
			}

			index = index*self.alphabet.Size() + i
		}

		result = append(result, index)
	}

	return result
}

// Prepare string for indexing
func (self *NGramIndex) prepareString(word string) string {
	word = normalizeWord(word)
	return wrapWord(word, self.config.pad, self.config.wrap)
}

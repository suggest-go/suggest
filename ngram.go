package suggest

/*
 * inspired by
 * http://www.aaai.org/ocs/index.php/AAAI/AAAI10/paper/viewFile/1939/2234
 * http://nlp.stanford.edu/IR-book/html/htmledition/k-gram-indexes-for-wildcard-queries-1.html
 * http://bazhenov.me/blog/2012/08/04/autocomplete.html
 */

import "sort"

type ngramsT map[string][]int

type NGramIndex struct {
	k          int
	ngrams     ngramsT
	dictionary map[int]string
	index      int
}

func NewNGramIndex(k int) *NGramIndex {
	if k < 2 || k > 4 {
		panic("k should be in [2, 4]")
	}

	return &NGramIndex{
		k, make(ngramsT), make(map[int]string), 0,
	}
}

func (self *NGramIndex) AddWord(word string) {
	prepared := prepareString(word)
	split := SplitIntoNGrams(prepared, self.k)
	for _, ngram := range split {
		self.ngrams[ngram] = append(self.ngrams[ngram], self.index)
	}

	self.dictionary[self.index] = word
	self.index++
}

func (self *NGramIndex) Suggest(word string, topK int) []string {
	t := 1 // to include all candidates
	corresponding := self.find(word)
	frequence := make(map[int]int)
	/* TODO optimize me */
	for _, c := range corresponding {
		for _, id := range c {
			frequence[id]++
		}
	}

	var candidates pairList
	for id, freq := range frequence {
		/* TODO use strategy */
		if freq >= t {
			candidate := self.dictionary[id]
			//similarity := Levenshtein(candidate, word)
			similarity := freq
			candidates = append(
				candidates,
				pair{candidate, similarity},
			)
		}
	}

	candidatesLen := len(candidates)
	sort.Sort(candidates)
	result := make([]string, candidatesLen)
	for i, pair := range candidates {
		result[i] = pair.word
	}

	if topK > candidatesLen {
		topK = candidatesLen
	}

	return result[:topK]
}

func (self *NGramIndex) find(word string) ngramsT {
	word = prepareString(word)
	result := make(ngramsT)
	split := SplitIntoNGrams(word, self.k)
	for _, ngram := range split {
		result[ngram] = self.ngrams[ngram]
	}

	return result
}

type pair struct {
	word  string
	score int
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].score > p[j].score }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

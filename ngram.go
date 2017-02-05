package suggest

/*
 * inspired by
 * http://www.aaai.org/ocs/index.php/AAAI/AAAI10/paper/viewFile/1939/2234
 * http://nlp.stanford.edu/IR-book/html/htmledition/k-gram-indexes-for-wildcard-queries-1.html
 * http://bazhenov.me/blog/2012/08/04/autocomplete.html
 */

import (
	"sort"
)

type Rank struct {
	word     string
	distance int
}
type RankList []Rank

func (p RankList) Len() int {
	return len(p)
}

func (p RankList) Less(i, j int) bool {
	return p[i].distance < p[j].distance
}

func (p RankList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type invertedListsT map[string][]int

type NGramIndex struct {
	k             int
	invertedLists invertedListsT
	dictionary    map[int]string
	index         int
}

func NewNGramIndex(k int) *NGramIndex {
	if k < 2 || k > 4 {
		panic("k should be in [2, 4]")
	}

	return &NGramIndex{
		k, make(invertedListsT), make(map[int]string), 0,
	}
}

func (self *NGramIndex) AddWord(word string) {
	prepared := prepareString(word)
	ngrams := SplitIntoNGrams(prepared, self.k)
	for _, ngram := range ngrams {
		self.invertedLists[ngram] = append(self.invertedLists[ngram], self.index)
	}

	self.dictionary[self.index] = word
	self.index++
}

/**
 * return top-k similar strings
 */
func (self *NGramIndex) Suggest(word string, topK int) []string {
	candidates := self.FuzzySearch(word)
	if topK > len(candidates) {
		topK = len(candidates)
	}

	sort.Sort(candidates)
	result := make([]string, topK)
	for i, rank := range candidates {
		result[i] = rank.word
		if i == topK-1 {
			break
		}
	}

	return result
}

/**
 * 1. try to receive corresponding inverted list for word's ngrams
 * 2. calculate distance between current word and candidates
 * 3. return RankList
 */
func (self *NGramIndex) FuzzySearch(word string) RankList {
	preparedWord := prepareString(word)
	corresponding := self.find(preparedWord)
	distances := make(map[int]int)
	for _, c := range corresponding {
		for _, id := range c {
			if _, ok := distances[id]; !ok {
				distances[id] = self.distance(
					preparedWord,
					prepareString(self.dictionary[id]),
				)
			}
		}
	}

	candidates := make(RankList, len(distances))
	i := 0
	for id, distance := range distances {
		candidates[i] = Rank{self.dictionary[id], distance}
		i++
	}

	return candidates
}

/*
 * ngram distance between a, b defined in
 * "Approximate string-matching with q-grams and maximal matches"
 *
 * Complexity O(aLen + bLen)
 */
func (self *NGramIndex) distance(a, b string) int {
	profileA, profileB := self.getProfile(a), self.getProfile(b)
	union := make(map[string]byte)
	for k := range profileA {
		union[k] = 1
	}

	for k := range profileB {
		union[k] = 1
	}

	distance := 0
	for key, _ := range union {
		freqA, freqB := 0, 0
		if val, ok := profileA[key]; ok {
			freqA = val
		}

		if val, ok := profileB[key]; ok {
			freqB = val
		}

		d := freqA - freqB
		if d < 0 {
			d = -d
		}

		distance += d
	}

	return distance
}

/*
 * Return unique ngrams with frequency
 */
func (self *NGramIndex) getProfile(word string) map[string]int {
	ngrams := SplitIntoNGrams(word, self.k)
	result := make(map[string]int)
	for _, ngram := range ngrams {
		if _, ok := result[ngram]; ok {
			result[ngram]++
		} else {
			result[ngram] = 1
		}
	}

	return result
}

/*
 * Find corresponding inverted lists by common ngrams
 */
func (self *NGramIndex) find(word string) invertedListsT {
	result := make(invertedListsT)
	profile := self.getProfile(word)
	for ngram, _ := range profile {
		result[ngram] = self.invertedLists[ngram]
	}

	return result
}

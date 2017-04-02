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
	distance float64
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
	dictionary    []string
	index         int
}

func NewNGramIndex(k int) *NGramIndex {
	if k < 2 || k > 4 {
		panic("k should be in [2, 4]")
	}

	return &NGramIndex{
		k, make(invertedListsT), make([]string, 0), 0,
	}
}

// Add given word to invertedList
func (self *NGramIndex) AddWord(word string) {
	prepared := prepareString(word)
	profile := self.getProfile(prepared)
	for _, ngram := range profile.GetNGrams() {
		val := ngram.GetValue()
		self.invertedLists[val] = append(self.invertedLists[val], self.index)
	}

	self.dictionary = append(self.dictionary, word)
	self.index++
}

// Return top-k similar strings
func (self *NGramIndex) Suggest(word string, editDistance EditDistance, topK int) []string {
	candidates := self.FuzzySearch(word, editDistance, topK)
	if topK > len(candidates) {
		topK = len(candidates)
	}

	sort.Sort(candidates)
	result := make([]string, 0)
	for i, rank := range candidates {
		result = append(result, rank.word)
		if i == topK-1 {
			break
		}
	}

	return result
}

//1. try to receive corresponding inverted list for word's ngrams
//2. calculate distance between current word and candidates
//3. return RankList
func (self *NGramIndex) FuzzySearch(word string, editDistance EditDistance, topK int) RankList {
	preparedWord := prepareString(word)
	wordProfile := self.getProfile(preparedWord)

	/* count filtering */
	max := 0
	counts := make([]int, self.index+1)
	for _, ngram := range wordProfile.GetNGrams() {
		for _, id := range self.invertedLists[ngram.GetValue()] {
			counts[id]++
			if max < counts[id] {
				max = counts[id]
			}
		}
	}

	a := make([][]int, max+1)
	for id, count := range counts {
		if count > 0 {
			a[count] = append(a[count], id)
		}
	}

	j := 0
	distances := make(map[int]float64, topK)
	for i := max; max-i < topK && i >= 0 && j < topK; i-- {
		for _, id := range a[i] {
			candidate := self.dictionary[id]
			prepared := prepareString(candidate)
			profile := self.getProfile(prepared)
			distances[id] = editDistance.Calc(wordProfile, profile)
			j++

			if j >= topK {
				break
			}
		}
	}

	candidates := make(RankList, 0, len(distances))
	for id, distance := range distances {
		candidates = append(candidates, Rank{self.dictionary[id], distance})
	}

	return candidates
}

// Return unique ngrams with frequency
func (self *NGramIndex) getProfile(word string) *WordProfile {
	return GetWordProfile(word, self.k)
}

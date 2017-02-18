package suggest

import "sync"

type SuggestService struct {
	sync.RWMutex
	dictionaries map[string]*NGramIndex
	ngramSize    int
	editDistance EditDistance
}

// Creates an empty SuggestService which uses given metric as "edit distance metric"
func NewSuggestService(ngramSize, metric int) *SuggestService {
	editDistance, err := GetEditDistance(metric, ngramSize)
	if err != nil {
		panic(err)
	}

	return &SuggestService{
		dictionaries: make(map[string]*NGramIndex),
		ngramSize:    ngramSize,
		editDistance: editDistance,
	}
}

// add/replace new dictionary with given name
func (self *SuggestService) AddDictionary(name string, words []string) error {
	ngramIndex := NewNGramIndex(self.ngramSize, self.editDistance)
	for _, word := range words {
		ngramIndex.AddWord(word)
	}

	self.Lock()
	self.dictionaries[name] = ngramIndex
	self.Unlock()
	return nil
}

// return Top-k approximate strings for given query in dict
func (self *SuggestService) Suggest(dict string, query string, topK int) []string {
	self.RLock()
	index, ok := self.dictionaries[dict]
	self.RUnlock()
	if ok {
		return index.Suggest(query, topK)
	}

	return make([]string, 0)
}

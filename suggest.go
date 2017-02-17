package suggest

import (
	"errors"
	"sync"
)

type SuggestService struct {
	sync.RWMutex
	dictionaries map[string]*NGramIndex
	ngramSize    int
	editDistance EditDistance
}

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

func (self *SuggestService) AddDictionary(name string, words []string) error {
	self.RLock()
	_, ok := self.dictionaries[name]
	self.RUnlock()

	if ok {
		return errors.New("dictionary already exists")
	}

	ngramIndex := NewNGramIndex(self.ngramSize, self.editDistance)
	for _, word := range words {
		ngramIndex.AddWord(word)
	}

	self.Lock()
	self.dictionaries[name] = ngramIndex
	self.Unlock()
	return nil
}

func (self *SuggestService) Suggest(dict string, query string, topK int) []string {
	self.RLock()
	index, ok := self.dictionaries[dict]
	self.RUnlock()
	if ok {
		return index.Suggest(query, topK)
	}

	return make([]string, 0)
}

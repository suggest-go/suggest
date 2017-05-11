package suggest

import "sync"

type SuggestService struct {
	sync.RWMutex
	indexes      map[string]*NGramIndex
	dictionaries map[string]Dictionary
}

// Creates an empty SuggestService which uses given metric as "edit distance metric"
func NewSuggestService() *SuggestService {
	// fixme
	return &SuggestService{
		indexes:      make(map[string]*NGramIndex),
		dictionaries: make(map[string]Dictionary),
	}
}

// add/replace new dictionary with given name
func (self *SuggestService) AddDictionary(name string, dictionary Dictionary, config *IndexConfig) error {
	ngramIndex := NewNGramIndex(config)

	for {
		key, word := dictionary.Next()
		if key == nil {
			break
		}

		// monkey code, fixme
		if len(word) > config.ngramSize {
			ngramIndex.AddWord(word, key)
		}
	}

	self.Lock()
	self.indexes[name] = ngramIndex
	self.dictionaries[name] = dictionary
	self.Unlock()
	return nil
}

// return Top-k approximate strings for given query in dict
func (self *SuggestService) Suggest(dict string, config *SearchConfig) []string {
	self.RLock()
	index, okIndex := self.indexes[dict]
	dictionary, okDict := self.dictionaries[dict]
	self.RUnlock()
	if okDict && okIndex {
		keys := index.Suggest(config)
		result := make([]string, 0, len(keys))
		for _, key := range keys {
			val, _ := dictionary.Get(key)
			result = append(result, val)
		}

		return result
	}

	return make([]string, 0)
}

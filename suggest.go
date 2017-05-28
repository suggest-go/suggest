package suggest

import "sync"

// Service is a service for topK approximate string search in dictionary
type Service struct {
	sync.RWMutex
	indexes      map[string]*NGramIndex
	dictionaries map[string]Dictionary
}

// NewService creates an empty SuggestService
func NewService() *Service {
	// fixme
	return &Service{
		indexes:      make(map[string]*NGramIndex),
		dictionaries: make(map[string]Dictionary),
	}
}

// AddDictionary add/replace new dictionary with given name
func (s *Service) AddDictionary(name string, dictionary Dictionary, config *IndexConfig) error {
	ngramIndex := NewNGramIndex(config)

	iter := dictionary.Iter()
	for iter.IsValid() {
		key, word := iter.GetPair()
		// monkey code, fixme
		if len(word) > config.ngramSize {
			ngramIndex.AddWord(word, key)
		}

		iter.Next()
	}

	s.Lock()
	s.indexes[name] = ngramIndex
	s.dictionaries[name] = dictionary
	s.Unlock()
	return nil
}

// Suggest returns Top-k approximate strings for given query in dict
func (s *Service) Suggest(dict string, config *SearchConfig) map[WordKey]string {
	s.RLock()
	index, okIndex := s.indexes[dict]
	dictionary, okDict := s.dictionaries[dict]
	s.RUnlock()
	if okDict && okIndex {
		keys := index.Suggest(config)
		result := dictionary.MGet(keys)
		return result
	}

	return make(map[WordKey]string, 0)
}

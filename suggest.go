package suggest

import (
	"sort"
	"sync"
)

type ResultItem struct {
	Candidate
	Value string
}

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
func (s *Service) Suggest(dict string, config *SearchConfig) []ResultItem {
	s.RLock()
	index, okIndex := s.indexes[dict]
	dictionary, okDict := s.dictionaries[dict]
	s.RUnlock()

	if !okDict || !okIndex {
		return nil
	}

	candidates := index.Suggest(config)
	l := len(candidates)
	result := make([]ResultItem, 0, l)
	keys := make([]WordKey, 0, l)
	m := make(map[WordKey]Candidate, l)

	for _, candidate := range candidates {
		keys = append(keys, candidate.Key)
		m[candidate.Key] = candidate
	}

	for key, value := range dictionary.MGet(keys) {
		candidate := m[key]
		result = append(result, ResultItem{candidate, value})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Distance < result[j].Distance
	})

	return result
}

package suggest

import (
	"sort"
	"sync"
)

// ResultItem represents element of top-k similar strings in dictionary for given query
type ResultItem struct {
	Distance float64
	// Value is a string value of candidate
	Value string
}

// Service is a service for topK approximate string search in dictionary
type Service struct {
	sync.RWMutex
	indexes      map[string]NGramIndex
	dictionaries map[string]Dictionary
}

// NewService creates an empty SuggestService
func NewService() *Service {
	return &Service{
		indexes:      make(map[string]NGramIndex),
		dictionaries: make(map[string]Dictionary),
	}
}

// AddRunTimeIndex add/replace new dictionary with given name
func (s *Service) AddRunTimeIndex(name string, config *IndexConfig) error {
	nGramIndex := NewRunTimeBuilder().
		SetAlphabet(config.alphabet).
		SetDictionary(config.dictionary).
		SetNGramSize(config.ngramSize).
		SetWrap(config.wrap).
		SetPad(config.pad).
		Build()

	s.Lock()
	s.indexes[name] = nGramIndex
	s.dictionaries[name] = config.dictionary
	s.Unlock()
	return nil
}

// AddOnDiscIndex add/replace new dictionary with given name
func (s *Service) AddOnDiscIndex(name string, dbPattern string, config *IndexConfig) error {
	nGramIndex := NewBuilder(dbPattern).
		SetAlphabet(config.alphabet).
		SetDictionary(config.dictionary).
		SetNGramSize(config.ngramSize).
		SetWrap(config.wrap).
		SetPad(config.pad).
		Build()

	s.Lock()
	s.indexes[name] = nGramIndex
	s.dictionaries[name] = config.dictionary
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

	for _, candidate := range candidates {
		value, _ := dictionary.Get(candidate.Key)
		result = append(result, ResultItem{candidate.Distance, value})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Distance < result[j].Distance
	})

	return result
}

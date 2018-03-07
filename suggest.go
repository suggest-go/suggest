package suggest

import (
	"fmt"
	"github.com/alldroll/suggest/dictionary"
	"golang.org/x/exp/mmap"
	"runtime"
	"sync"
)

// ResultItem represents element of top-k similar strings in dictionary for given query
type ResultItem struct {
	Distance float64
	// Value is a string value of candidate
	Value string
}

// Service is a service for topK approximate string fuzzySearch in dictionary
type Service struct {
	sync.RWMutex
	indexes      map[string]NGramIndex
	dictionaries map[string]dictionary.Dictionary
}

// NewService creates an empty SuggestService
func NewService() *Service {
	return &Service{
		indexes:      make(map[string]NGramIndex),
		dictionaries: make(map[string]dictionary.Dictionary),
	}
}

// AddRunTimeIndex add/replace new dictionary with given name
func (s *Service) AddRunTimeIndex(name string, config *IndexConfig) error {
	nGramIndex := NewRunTimeBuilder(config).Build()

	s.Lock()
	s.indexes[name] = nGramIndex
	s.dictionaries[name] = config.dictionary
	s.Unlock()
	return nil
}

// AddOnDiscIndex add/replace new dictionary with given name
func (s *Service) AddOnDiscIndex(description IndexDescription) error {
	dictionaryFile, err := mmap.Open(description.GetDictionaryFile())
	if err != nil {
		return err
	}

	dictionary := dictionary.NewCDBDictionary(dictionaryFile)
	runtime.SetFinalizer(dictionary, func(d interface{}) {
		dictionaryFile.Close()
	})

	nGramIndex := NewBuilder(description).Build()

	s.Lock()
	s.indexes[description.Name] = nGramIndex
	s.dictionaries[description.Name] = dictionary
	s.Unlock()
	return nil
}

func (s *Service) GetDictionaries() []string {
	names := make([]string, 0, len(s.dictionaries))

	for name, _ := range s.dictionaries {
		names = append(names, name)
	}

	return names
}

// Suggest returns Top-k approximate strings for given query in dict
func (s *Service) Suggest(dict string, config *SearchConfig) ([]ResultItem, error) {
	s.RLock()
	index, okIndex := s.indexes[dict]
	dictionary, okDict := s.dictionaries[dict]
	s.RUnlock()

	if !okDict || !okIndex {
		return nil, fmt.Errorf("Given dictionary %s is not exists", dict)
	}

	candidates := index.Suggest(config)
	l := len(candidates)
	result := make([]ResultItem, 0, l)

	for _, candidate := range candidates {
		value, err := dictionary.Get(candidate.Key)
		if err != nil {
			return nil, err
		} else {
			result = append(result, ResultItem{candidate.Distance, value})
		}
	}

	return result, nil
}

// AutoComplete returns first limit
func (s *Service) AutoComplete(dict string, query string, limit int) ([]ResultItem, error) {
	s.RLock()
	index, okIndex := s.indexes[dict]
	dictionary, okDict := s.dictionaries[dict]
	s.RUnlock()

	if !okDict || !okIndex {
		return nil, fmt.Errorf("Given dictionary %s is not exists", dict)
	}

	candidates := index.AutoComplete(query, limit)
	result := make([]ResultItem, 0, len(candidates))

	for _, candidate := range candidates {
		value, err := dictionary.Get(candidate.Key)
		if err != nil {
			return nil, err
		} else {
			result = append(result, ResultItem{0, value})
		}
	}

	return result, nil
}

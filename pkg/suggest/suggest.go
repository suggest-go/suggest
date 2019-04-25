package suggest

import (
	"fmt"
	"sync"

	"github.com/alldroll/suggest/pkg/dictionary"
)

// ResultItem represents element of top-k similar strings in dictionary for given query
type ResultItem struct {
	// Score is a float64 value of a candidate
	Score float64
	// Value is a string value of candidate
	Value string
}

// Service provides methods for autocomplete and topK approximate string search
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

// AddIndexByDescription adds a new search index with given description
func (s *Service) AddIndexByDescription(description IndexDescription) error {
	if description.Driver == RAMDriver {
		return s.AddRunTimeIndex(description)
	}

	return s.AddOnDiscIndex(description)
}

// AddRunTimeIndex adds a new RAM search index with the given description
func (s *Service) AddRunTimeIndex(description IndexDescription) error {
	dictionary, err := dictionary.OpenRAMDictionary(description.SourcePath)

	if err != nil {
		return fmt.Errorf("Failed to create RAMDriver builder: %v", err)
	}

	builder, err := NewRAMBuilder(dictionary, description)

	if err != nil {
		return fmt.Errorf("Failed to create RAMDriver builder: %v", err)
	}

	return s.AddIndex(description.Name, dictionary, builder)
}

// AddOnDiscIndex adds a new DISC search index with the given description
func (s *Service) AddOnDiscIndex(description IndexDescription) error {
	dictionary, err := dictionary.OpenCDBDictionary(description.GetDictionaryFile())

	if err != nil {
		return fmt.Errorf("Failed to create CDB dictionary: %v", err)
	}

	builder, err := NewFSBuilder(description)

	if err != nil {
		return fmt.Errorf("Failed to open FS inverted index: %v", err)
	}

	return s.AddIndex(description.Name, dictionary, builder)
}

// AddIndex adds an index with the given name, dictionary and builder
func (s *Service) AddIndex(name string, dict dictionary.Dictionary, builder Builder) error {
	nGramIndex, err := builder.Build()

	if err != nil {
		return fmt.Errorf("Failed to build NGramIndex: %v", err)
	}

	s.Lock()
	s.indexes[name] = nGramIndex
	s.dictionaries[name] = dict
	s.Unlock()

	return nil
}

// GetDictionaries returns the managed list of dictionaries
func (s *Service) GetDictionaries() []string {
	names := make([]string, 0, len(s.dictionaries))

	for name := range s.dictionaries {
		names = append(names, name)
	}

	return names
}

// Suggest returns Top-k approximate strings for the given query in the dict
func (s *Service) Suggest(dict string, config *SearchConfig) ([]ResultItem, error) {
	s.RLock()
	index, okIndex := s.indexes[dict]
	dictionary, okDict := s.dictionaries[dict]
	s.RUnlock()

	if !okDict || !okIndex {
		return nil, fmt.Errorf("Given dictionary %s is not exists", dict)
	}

	candidates, err := index.Suggest(config)

	if err != nil {
		return nil, err
	}

	l := len(candidates)
	result := make([]ResultItem, 0, l)

	for _, candidate := range candidates {
		value, err := dictionary.Get(candidate.Key)
		if err != nil {
			return nil, err
		}

		result = append(result, ResultItem{candidate.Score, value})
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

	candidates, err := index.AutoComplete(query, limit)

	if err != nil {
		return nil, err
	}

	result := make([]ResultItem, 0, len(candidates))

	for _, candidate := range candidates {
		value, err := dictionary.Get(candidate.Key)
		if err != nil {
			return nil, err
		}

		result = append(result, ResultItem{0, value})
	}

	return result, nil
}

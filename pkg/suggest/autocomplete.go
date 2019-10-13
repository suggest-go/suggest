package suggest

import (
	"fmt"

	"github.com/suggest-go/suggest/pkg/analysis"
	"github.com/suggest-go/suggest/pkg/index"
	"golang.org/x/sync/errgroup"
)

// Autocomplete provides autocomplete functionality
// for candidates search
type Autocomplete interface {
	// Autocomplete returns candidates where the query string is a substring of each candidate
	Autocomplete(query string, collectorManager CollectorManager) ([]Candidate, error)
}

// NewAutocomplete creates a new instance of Autocomplete
func NewAutocomplete(
	indices index.InvertedIndexIndices,
	searcher index.Searcher,
	tokenizer analysis.Tokenizer,
) Autocomplete {
	return &nGramAutocomplete{
		indices:   indices,
		searcher:  searcher,
		tokenizer: tokenizer,
	}
}

// nGramAutocomplete implements Autocomplete interface
type nGramAutocomplete struct {
	indices   index.InvertedIndexIndices
	searcher  index.Searcher
	tokenizer analysis.Tokenizer
}

// Autocomplete returns candidates where the query string is a substring of each candidate
func (n *nGramAutocomplete) Autocomplete(query string, collectorManager CollectorManager) ([]Candidate, error) {
	set := n.tokenizer.Tokenize(query)
	lenSet := len(set)
	collectors := []Collector{}
	workerPool := errgroup.Group{}

	for size := lenSet; size < n.indices.Size(); size++ {
		index := n.indices.Get(size)

		if index == nil {
			continue
		}

		collector, err := collectorManager.Create()

		if err != nil {
			return nil, fmt.Errorf("failed to create a collector: %v", err)
		}

		workerPool.Go(func() error {
			if err = n.searcher.Search(index, set, lenSet, collector); err != nil {
				return fmt.Errorf("failed to search posting lists: %v", err)
			}

			return nil
		})

		collectors = append(collectors, collector)
	}

	if err := workerPool.Wait(); err != nil {
		return nil, err
	}

	return collectorManager.Reduce(collectors), nil
}

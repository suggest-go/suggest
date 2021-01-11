package suggest

import (
	"fmt"
	"sync"

	"github.com/suggest-go/suggest/pkg/analysis"
	"github.com/suggest-go/suggest/pkg/index"
	"golang.org/x/sync/errgroup"
)

// Autocomplete provides autocomplete functionality
// for candidates search
type Autocomplete interface {
	// Autocomplete returns candidates where the query string is a substring of each candidate
	Autocomplete(query string, factory CollectorManagerFactory) ([]Candidate, error)
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

// Autocomplete returns candidates where the query string is a prefix of each candidate
func (n *nGramAutocomplete) Autocomplete(query string, factory CollectorManagerFactory) ([]Candidate, error) {
	terms := n.tokenizer.Tokenize(query)
	termsLen := len(terms)
	workerPool := errgroup.Group{}
	collectorManager := factory()
	locker := sync.Mutex{}

	for size := termsLen; size < n.indices.Size(); size++ {
		invertedIndex := n.indices.Get(size)

		if invertedIndex == nil {
			continue
		}

		collector := collectorManager.Create()

		workerPool.Go(func() error {
			if err := n.searcher.Search(invertedIndex, terms, termsLen, collector); err != nil {
				return fmt.Errorf("failed to search posting lists: %w", err)
			}

			locker.Lock()
			defer locker.Unlock()

			if err := collectorManager.Collect(collector); err != nil {
				return err
			}

			return nil
		})
	}

	if err := workerPool.Wait(); err != nil {
		return nil, err
	}

	return collectorManager.GetCandidates(), nil
}

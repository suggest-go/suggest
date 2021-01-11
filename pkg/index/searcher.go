package index

import (
	"fmt"

	"github.com/suggest-go/suggest/pkg/merger"
)

// Searcher is responsible for searching
type Searcher interface {
	// Search performs search for the given index with the terms and threshold
	Search(invertedIndex InvertedIndex, terms []Term, threshold int, collector merger.Collector) error
}

// searcher implements the Searcher interface
type searcher struct {
	merger merger.ListMerger
}

// NewSearcher creates a new Searcher instance
func NewSearcher(merger merger.ListMerger) Searcher {
	return &searcher{
		merger: merger,
	}
}

// Search performs search for the given index with the terms and threshold
func (s *searcher) Search(invertedIndex InvertedIndex, terms []Term, threshold int, collector merger.Collector) error {
	terms = filterTermsByExistence(invertedIndex, terms, threshold)
	n := len(terms)

	if n < threshold {
		return nil
	}

	rid := make([]merger.ListIterator, 0, n)

	for _, term := range terms {
		postingListContext, err := invertedIndex.Get(term)

		if err != nil {
			return fmt.Errorf("failed to retrieve a posting list context: %w", err)
		}

		list := resolvePostingList(postingListContext)

		defer func(list PostingList) {
			if closeErr := releasePostingList(list); err != nil {
				err = closeErr
			}
		}(list)

		if err := list.Init(postingListContext); err != nil {
			return fmt.Errorf("failed to initialize a posting list iterator: %w", err)
		}

		rid = append(rid, list)
	}

	if err := s.merger.Merge(rid, threshold, collector); err != nil {
		return fmt.Errorf("failed to merge posting lists: %w", err)
	}

	return nil
}

func filterTermsByExistence(index InvertedIndex, terms []Term, threshold int) []Term {
	n := len(terms)
	filtered := make([]Term, 0, n)

	for i := 0; i < n && (len(filtered)+n-i) >= threshold; i++ {
		if index.Has(terms[i]) {
			filtered = append(filtered, terms[i])
		}
	}

	return filtered
}

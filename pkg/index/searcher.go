package index

import (
	"fmt"
	"github.com/alldroll/suggest/pkg/analysis"

	"github.com/alldroll/suggest/pkg/merger"
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
	n := len(terms)
	set := make([]analysis.Token, len(terms))
	copy(set, terms)

	for i := 0; i < n && n >= threshold; {
		term := set[i]

		if !invertedIndex.Has(term) {
			set[i], set[n-1] = set[n-1], set[i]
			n--
		} else {
			i++
		}
	}

	if n < threshold {
		return nil
	}

	rid := make([]merger.ListIterator, 0, n)

	for _, term := range set[:n] {
		postingListContext, err := invertedIndex.Get(term)

		if err != nil {
			return fmt.Errorf("failed to retrieve a posting list context: %v", err)
		}

		list := resolvePostingList(postingListContext)

		defer func(list postingList) {
			if closeErr := releasePostingList(list); err != nil {
				err = closeErr
			}
		}(list)

		if err := list.init(postingListContext); err != nil {
			return fmt.Errorf("failed to initialize a posting list iterator: %v", err)
		}

		rid = append(rid, list)
	}

	if err := s.merger.Merge(rid, threshold, collector); err != nil {
		return fmt.Errorf("failed to merge posting lists: %v", err)
	}

	return nil
}

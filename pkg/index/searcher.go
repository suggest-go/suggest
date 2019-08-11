package index

import (
	"fmt"

	"github.com/alldroll/suggest/pkg/merger"
)

// Searcher is responsible for searching
type Searcher interface {
	// Search performs search for the given index with the terms and threshold
	Search(invertedIndex InvertedIndex, terms []Term, threshold int) ([]merger.MergeCandidate, error)
}

// seacher implements the Searcher interface
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
func (s *searcher) Search(invertedIndex InvertedIndex, terms []Term, threshold int) (candidates []merger.MergeCandidate, err error) {
	if threshold > len(terms) {
		return
	}

	allowedSkips := len(terms) - threshold + 1

	for _, term := range terms {
		if allowedSkips == 0 {
			break
		}

		if !invertedIndex.Has(term) {
			allowedSkips--
		}
	}

	if allowedSkips == 0 {
		return
	}

	rid := make([]merger.ListIterator, 0, allowedSkips+threshold-1)

	for _, term := range terms {
		postingListContext, err := invertedIndex.Get(term)

		if err != nil {
			return nil, fmt.Errorf("failed to retrieve a posting list context: %v", err)
		}

		if postingListContext != nil && postingListContext.GetListSize() > 0 {
			list := resolvePostingList(postingListContext)

			defer func(list postingList) {
				if closeErr := releasePostingList(list); err != nil {
					err = closeErr
				}
			}(list)

			if err := list.init(postingListContext); err != nil {
				return nil, fmt.Errorf("failed to initialize a posting list iterator: %v", err)
			}

			rid = append(rid, list)
		}
	}

	candidates, err = s.merger.Merge(rid, threshold)

	if err != nil {
		return nil, fmt.Errorf("failed to merge posting lists: %v", err)
	}

	return
}

package index

import (
	"fmt"

	"sync"

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

// iteratorPool reduces allocation of iterator object
var iteratorPool = sync.Pool{
	New: func() interface{} {
		return &postingListIterator{
			sliceIterator: merger.NewSliceIterator([]Position{}),
			index:         0,
			size:          0,
			current:       uint32(0),
		}
	},
}

// Search performs search for the given index with the terms and threshold
func (s *searcher) Search(invertedIndex InvertedIndex, terms []Term, threshold int) ([]merger.MergeCandidate, error) {
	if threshold > len(terms) {
		return []merger.MergeCandidate{}, nil
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
		return []merger.MergeCandidate{}, nil
	}

	rid := make([]merger.ListIterator, 0, len(terms))

	for _, term := range terms {
		postingListContext, err := invertedIndex.Get(term)

		if err != nil {
			return nil, fmt.Errorf("failed to retrieve a posting list context: %v", err)
		}

		if postingListContext != nil && postingListContext.GetListSize() > 0 {
			iterator := iteratorPool.Get().(*postingListIterator)
			defer iteratorPool.Put(iterator)

			if err := iterator.init(postingListContext); err != nil {
				return nil, fmt.Errorf("failed to initialize a posting list iterator: %v", err)
			}

			rid = append(rid, iterator)
		}
	}

	candidates, err := s.merger.Merge(rid, threshold)

	if err != nil {
		return nil, fmt.Errorf("failed to merge posting lists: %v", err)
	}

	return candidates, nil
}

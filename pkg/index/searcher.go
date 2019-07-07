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

// readerPool reduces allocation of bufio.Reader object
var ridPool = sync.Pool{
	New: func() interface{} {
		rid := make([]merger.ListIterator, 30)

		return &rid
	},
}

// Search performs search for the given index with the terms and threshold
func (s *searcher) Search(invertedIndex InvertedIndex, terms []Term, threshold int) ([]merger.MergeCandidate, error) {
	rid := *(ridPool.Get().(*[]merger.ListIterator))
	rid = rid[:0]

	defer ridPool.Put(&rid)

	// maybe run it concurrent?
	// go func() { buildRid, mergeCandidates, ch <- {candidates, sizeA, sizeB}
	// in main goroutine just collect it
	for _, term := range terms {
		postingList, err := invertedIndex.Get(term)

		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve a posting list: %v", err)
		}

		if postingList != nil && postingList.Len() > 0 {
			rid = append(rid, postingList)
		}
	}

	candidates, err := s.merger.Merge(rid, threshold)

	if err != nil {
		return nil, fmt.Errorf("failed to merge posting lists: %v", err)
	}

	return candidates, nil
}

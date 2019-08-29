package merger

import "errors"

// ErrCollectionTerminated tells to terminate collection of the current workflow
// This error is going to be processed by Intersector/Merger and not to be propagated up
var ErrCollectionTerminated = errors.New("collection terminated")

// Collector collects the doc stream satisfied to a search criteria
type Collector interface {
	// Collect collects the given candidate
	Collect(candidate MergeCandidate) error
}

// SimpleCollector a dummy implementation of Collector
type SimpleCollector struct {
	Candidates []MergeCandidate
}

// Collect collects the given candidate
func (c *SimpleCollector) Collect(candidate MergeCandidate) error {
	c.Candidates = append(c.Candidates, candidate)

	return nil
}

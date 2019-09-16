package suggest

import (
	"github.com/alldroll/suggest/pkg/index"
	"github.com/alldroll/suggest/pkg/merger"
)

// Candidate is an item of Collector
type Candidate struct {
	// Key is a position (docId) in posting list
	Key index.Position
	// Score is a float64 number that represents a score of a document
	Score float64
}

// Less tells is the given candidate is less that the provided
func (c Candidate) Less(o Candidate) bool {
	if c.Score == o.Score {
		return c.Key > o.Key
	}

	return c.Score < o.Score
}

// Collector collects the doc stream satisfied to a search criteria
type Collector interface {
	merger.Collector
	// SetScorer sets a scorer before collection starts
	SetScorer(scorer Scorer)
	// GetCandidates returns the list of collected candidates
	GetCandidates() []Candidate
}

// CollectorManager is responsible for creating collectors and reducing them into the result set
type CollectorManager interface {
	// Create creates a new collector that will be used for a search segment
	Create() (Collector, error)
	// Reduce reduces the result from the given list of collectors
	Reduce(collectors []Collector) []Candidate
}

type firstKCollector struct {
	limit int
	items []merger.MergeCandidate
}

// Collect collects the given merge candidate
func (c *firstKCollector) Collect(item merger.MergeCandidate) error {
	if c.limit == len(c.items) {
		return merger.ErrCollectionTerminated
	}

	c.items = append(c.items, item)

	return nil
}

// SetScorer sets a scorer before collection starts
func (c *firstKCollector) SetScorer(scorer Scorer) {
	return
}

// GetCandidates returns the list of collected candidates
func (c *firstKCollector) GetCandidates() []Candidate {
	result := make([]Candidate, 0, len(c.items))

	for _, item := range c.items {
		result = append(result, Candidate{
			Key: item.Position(),
		})
	}

	return result
}

// NewFirstKCollectorManager creates a new instance of CollectorManager with firstK collectors
func NewFirstKCollectorManager(limit int) CollectorManager {
	return &firstKCollectorManager{
		limit: limit,
	}
}

type firstKCollectorManager struct {
	limit int
}

// Create creates a new collector that will be used for a search segment
func (m *firstKCollectorManager) Create() (Collector, error) {
	return &firstKCollector{
		limit: m.limit,
	}, nil
}

// Reduce reduces the result from the given list of collectors
func (m *firstKCollectorManager) Reduce(collectors []Collector) []Candidate {
	topKQueue := NewTopKQueue(m.limit)

	for _, c := range collectors {
		if collector, ok := c.(*firstKCollector); ok {
			for _, item := range collector.items {
				topKQueue.Add(item.Position(), -float64(item.Position()))
			}
		}
	}

	return topKQueue.GetCandidates()
}

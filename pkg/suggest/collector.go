package suggest

import (
	"errors"
	"math"

	"github.com/suggest-go/suggest/pkg/index"
	"github.com/suggest-go/suggest/pkg/merger"
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
}

// CollectorManager is responsible for creating collectors and reducing them into the result set
type CollectorManager interface {
	// Create creates a new collector that will be used for a search segment
	Create() Collector
	// Collect returns back the given collectors.
	Collect(collectors ...Collector) error
	// GetCandidates returns currently collected candidates.
	GetCandidates() []Candidate
}

// CollectorManagerFactory is a factory method for creating a new instance of CollectorManager.
type CollectorManagerFactory func() CollectorManager

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

// NewFirstKCollectorManager creates a new instance of CollectorManager with firstK collectors
func NewFirstKCollectorManager(limit int, queue TopKQueue) *FirstKCollectorManager {
	return &FirstKCollectorManager{
		limit: limit,
		queue: queue,
	}
}

func newFirstKCollectorManager(limit int) CollectorManagerFactory {
	return func() CollectorManager {
		return NewFirstKCollectorManager(limit, NewTopKQueue(limit))
	}
}

// FirstKCollectorManager represents first k collector manager.
type FirstKCollectorManager struct {
	limit int
	queue TopKQueue
}

// Create creates a new collector that will be used for a search segment
func (m *FirstKCollectorManager) Create() Collector {
	return &firstKCollector{
		limit: m.limit,
	}
}

// Collect returns back the given collectors.
func (m *FirstKCollectorManager) Collect(collectors ...Collector) error {
	for _, item := range collectors {
		collector, ok := item.(*firstKCollector)

		if !ok {
			return errors.New("expected Collector created by FirstKCollectorManager")
		}

		for _, candidate := range collector.items {
			m.queue.Add(candidate.Position(), -float64(candidate.Position()))
		}
	}

	return nil
}

// GetCandidates returns currently collected candidates.
func (m *FirstKCollectorManager) GetCandidates() []Candidate {
	return m.queue.GetCandidates()
}

type fuzzyCollector struct {
	topKQueue TopKQueue
	scorer    Scorer
}

// Collect collects the given merge candidate
// calculates the distance, and tries to add this document with it's score to the collector
func (c *fuzzyCollector) Collect(item merger.MergeCandidate) error {
	c.topKQueue.Add(item.Position(), c.scorer.Score(item))

	return nil
}

// Score returns the score of the given position
func (c *fuzzyCollector) SetScorer(scorer Scorer) {
	c.scorer = scorer
}

// NewFuzzyCollectorManager creates a new instance of FuzzyCollectorManager.
func NewFuzzyCollectorManager(queueFactory func() TopKQueue) *FuzzyCollectorManager {
	return &FuzzyCollectorManager{
		queueFactory: queueFactory,
		globalQueue:  queueFactory(),
	}
}

func newFuzzyCollectorManager(topK int) CollectorManagerFactory {
	return func() CollectorManager {
		return NewFuzzyCollectorManager(func() TopKQueue {
			return NewTopKQueue(topK)
		})
	}
}

// FuzzyCollectorManager represents fuzzy collector manager.
type FuzzyCollectorManager struct {
	queueFactory func() TopKQueue
	globalQueue  TopKQueue
}

// Create creates a new collector that will be used for a search segment
func (m *FuzzyCollectorManager) Create() Collector {
	return &fuzzyCollector{
		topKQueue: m.queueFactory(),
	}
}

// Collect returns back the given collectors.
func (m *FuzzyCollectorManager) Collect(collectors ...Collector) error {
	for _, item := range collectors {
		collector, ok := item.(*fuzzyCollector)

		if !ok {
			return errors.New("expected Collector created by FirstKCollectorManager")
		}

		m.globalQueue.Merge(collector.topKQueue)
	}

	return nil
}

// GetCandidates returns currently collected candidates.
func (m *FuzzyCollectorManager) GetCandidates() []Candidate {
	return m.globalQueue.GetCandidates()
}

// GetLowestScore returns the lowest collected score.
func (m *FuzzyCollectorManager) GetLowestScore() float64 {
	if !m.globalQueue.IsFull() {
		return math.Inf(-1)
	}

	return m.globalQueue.GetLowestScore()
}

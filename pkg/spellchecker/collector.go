package spellchecker

import (
	"errors"

	"github.com/suggest-go/suggest/pkg/lm"
	"github.com/suggest-go/suggest/pkg/merger"
	"github.com/suggest-go/suggest/pkg/suggest"
)

// lmCollector implements Collector interface
type lmCollector struct {
	topKQueue suggest.TopKQueue
	scorer    suggest.Scorer
}

// newCollectorManager creates a new instance of lm CollectorManger.
func newCollectorManager(scorer suggest.Scorer, queueFactory func() suggest.TopKQueue) suggest.CollectorManager {
	return &lmCollectorManager{
		scorer:       scorer,
		queueFactory: queueFactory,
		globalQueue:  queueFactory(),
	}
}

// lmCollectorManager implements CollectorManager interface
type lmCollectorManager struct {
	scorer       suggest.Scorer
	queueFactory func() suggest.TopKQueue
	globalQueue  suggest.TopKQueue
}

// Create creates a new collector that will be used for a search segment
func (l *lmCollectorManager) Create() suggest.Collector {
	return &lmCollector{
		topKQueue: l.queueFactory(),
		scorer:    l.scorer,
	}
}

// Reduce reduces the result from the given list of collectors
func (l *lmCollectorManager) Collect(collectors ...suggest.Collector) error {
	for _, c := range collectors {
		collector, ok := c.(*lmCollector)

		if !ok {
			return errors.New("expected collector created by lmCollectorManager")
		}

		l.globalQueue.Merge(collector.topKQueue)
	}

	return nil
}

func (l *lmCollectorManager) GetCandidates() []suggest.Candidate {
	return l.globalQueue.GetCandidates()
}

// Collect collects the given candidate
func (c *lmCollector) Collect(item merger.MergeCandidate) error {
	doc := item.Position()

	if c.scorer == nil {
		if c.topKQueue.IsFull() {
			return merger.ErrCollectionTerminated
		}

		c.topKQueue.Add(doc, lm.UnknownWordScore)

		return nil
	}

	score := c.scorer.Score(item)
	c.topKQueue.Add(doc, score)

	return nil
}

// SetScorer sets the scorer for calculations
func (c *lmCollector) SetScorer(scorer suggest.Scorer) {
	c.scorer = scorer
}

// GetCandidates returns `top k items`
func (c *lmCollector) GetCandidates() []suggest.Candidate {
	return c.topKQueue.GetCandidates()
}

package spellchecker

import (
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
func newCollectorManager(scorer suggest.Scorer, topK int) suggest.CollectorManager {
	return &lmCollectorManager{
		topK:   topK,
		scorer: scorer,
	}
}

// lmCollectorManager implements CollectorManager interface
type lmCollectorManager struct {
	topK   int
	scorer suggest.Scorer
}

// Create creates a new collector that will be used for a search segment
func (l *lmCollectorManager) Create() (suggest.Collector, error) {
	return &lmCollector{
		topKQueue: suggest.NewTopKQueue(l.topK),
		scorer:    l.scorer,
	}, nil
}

// Reduce reduces the result from the given list of collectors
func (l *lmCollectorManager) Reduce(collectors []suggest.Collector) []suggest.Candidate {
	topKQueue := suggest.NewTopKQueue(l.topK)

	for _, c := range collectors {
		if collector, ok := c.(*lmCollector); ok {
			topKQueue.Merge(collector.topKQueue)
		}
	}

	return topKQueue.GetCandidates()
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

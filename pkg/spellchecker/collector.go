package spellchecker

import (
	"github.com/alldroll/suggest/pkg/lm"
	"github.com/alldroll/suggest/pkg/merger"
	"github.com/alldroll/suggest/pkg/suggest"
	"sort"
)

// lmCollector implements Collector interface
type lmCollector struct {
	topKQueue suggest.TopKQueue
	scorer    suggest.Scorer
	next      []lm.WordCount
}

// lmCollectorManager implements CollectorManager interface
type lmCollectorManager struct {
	topK   int
	scorer *lmScorer
	next   []lm.WordCount
}

// Create creates a new collector that will be used for a search segment
func (l *lmCollectorManager) Create() (suggest.Collector, error) {
	return &lmCollector{
		topKQueue: suggest.NewTopKQueue(l.topK),
		scorer:    l.scorer,
		next:      l.next,
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

	if len(c.next) == 0 {
		if c.topKQueue.IsFull() {
			return merger.ErrCollectionTerminated
		}

		c.topKQueue.Add(doc, lm.UnknownWordScore)

		return nil
	}

	if c.next[0] > item.Position() {
		c.topKQueue.Add(doc, lm.UnknownWordScore)

		return nil
	}

	if c.next[len(c.next)-1] < item.Position() {
		c.topKQueue.Add(doc, lm.UnknownWordScore)
		c.next = c.next[:0]

		return nil
	}

	i := sort.Search(len(c.next), func(i int) bool { return c.next[i] >= doc })
	c.next = c.next[i:]

	if c.next[0] != doc {
		c.topKQueue.Add(doc, lm.UnknownWordScore)

		return nil
	}

	c.topKQueue.Add(doc, c.scorer.Score(item))

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

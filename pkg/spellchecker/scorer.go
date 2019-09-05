package spellchecker

import (
	"sort"

	"github.com/alldroll/suggest/pkg/index"
	lm "github.com/alldroll/suggest/pkg/language-model"
	"github.com/alldroll/suggest/pkg/merger"
	"github.com/alldroll/suggest/pkg/suggest"
)

type lmCollectorManager struct {
	topK   int
	scorer suggest.Scorer
	next   []lm.WordCount
}

func (l *lmCollectorManager) Create() (suggest.Collector, error) {
	return &lmCollector{
		topKQueue: suggest.NewTopKQueue(l.topK),
		scorer:    l.scorer,
		next:      l.next,
	}, nil
}

func (l *lmCollectorManager) Reduce(collectors []suggest.Collector) []suggest.Candidate {
	topKQueue := suggest.NewTopKQueue(l.topK)

	for _, collector := range collectors {
		for _, item := range collector.GetCandidates() {
			topKQueue.Add(item.Key, item.Score)
		}
	}

	return topKQueue.GetCandidates()
}

// lmScorer implements the scorer interface
type lmScorer struct {
	model    lm.LanguageModel
	sentence []lm.WordID
}

// Score returns the score of the given position
func (s *lmScorer) Score(position index.Position) float64 {
	return s.model.ScoreWordIDs(append(s.sentence, position))
}

// lmCollector implements Collector interface
type lmCollector struct {
	topKQueue suggest.TopKQueue
	scorer    suggest.Scorer
	next      []lm.WordCount
}

//
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

	c.topKQueue.Add(doc, c.scorer.Score(doc))

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

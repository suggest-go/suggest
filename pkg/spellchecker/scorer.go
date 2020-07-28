package spellchecker

import (
	"github.com/suggest-go/suggest/pkg/lm"
	"github.com/suggest-go/suggest/pkg/merger"
	"github.com/suggest-go/suggest/pkg/suggest"
)

// lmScorer implements the scorer interface
type lmScorer struct {
	scorer lm.ScorerNext
}

// Score returns the score of the given position
func (s *lmScorer) Score(candidate merger.MergeCandidate) float64 {
	return s.score(candidate.Position())
}

// score returns the lm score for the given word ID
func (s *lmScorer) score(id lm.WordID) float64 {
	return s.scorer.ScoreNext(id)
}

type dummyScorer struct {
}

// Score returns the score of the given position
func (s *dummyScorer) Score(candidate merger.MergeCandidate) float64 {
	return lm.UnknownWordScore
}

// newScorer creates a scorer for the provided lm.ScorerNext
func newScorer(next lm.ScorerNext) suggest.Scorer {
	if next == nil {
		return &dummyScorer{}
	}

	return &lmScorer{
		scorer: next,
	}
}

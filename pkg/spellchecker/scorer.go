package spellchecker

import (
	"github.com/suggest-go/suggest/pkg/lm"
	"github.com/suggest-go/suggest/pkg/merger"
	"github.com/suggest-go/suggest/pkg/suggest"
	"sort"
)

// lmScorer implements the scorer interface
type lmScorer struct {
	scorer  lm.ScorerNext
}

// Score returns the score of the given position
func (s *lmScorer) Score(candidate merger.MergeCandidate) float64 {
	return s.score(candidate.Position())
}

// score returns the lm score for the given word ID
func (s *lmScorer) score(id lm.WordID) float64 {
	return s.scorer.ScoreNext(id)
}

// sortCandidates performs sort of the given candidates using lm
func sortCandidates(scorer *lmScorer, candidates []suggest.Candidate) {
	sort.SliceStable(candidates, func(i, j int) bool {
		return scorer.score(candidates[i].Key) > scorer.score(candidates[j].Key)
	})
}

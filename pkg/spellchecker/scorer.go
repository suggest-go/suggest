package spellchecker

import (
	"github.com/alldroll/suggest/pkg/lm"
	"github.com/alldroll/suggest/pkg/merger"
	"github.com/alldroll/suggest/pkg/suggest"
	"sort"
)

// lmScorer implements the scorer interface
type lmScorer struct {
	model    lm.LanguageModel
	sentence []lm.WordID
}

// Score returns the score of the given position
func (s *lmScorer) Score(candidate merger.MergeCandidate) float64 {
	return s.score(candidate.Position())
}

// score returns the lm score for the given word ID
func (s *lmScorer) score(id lm.WordID) float64 {
	return s.model.ScoreWordIDs(append(s.sentence, id))
}

// sortCandidates performs sort of the given candidates using lm
func sortCandidates(scorer *lmScorer, candidates []suggest.Candidate) {
	sort.SliceStable(candidates, func(i, j int) bool {
		return scorer.score(candidates[i].Key) > scorer.score(candidates[j].Key)
	})
}

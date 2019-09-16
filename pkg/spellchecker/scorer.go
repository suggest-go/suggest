package spellchecker

import (
	"github.com/alldroll/suggest/pkg/lm"
	"github.com/alldroll/suggest/pkg/merger"
)

// lmScorer implements the scorer interface
type lmScorer struct {
	model    lm.LanguageModel
	sentence []lm.WordID
}

// Score returns the score of the given position
func (s *lmScorer) Score(candidate merger.MergeCandidate) float64 {
	return s.model.ScoreWordIDs(append(s.sentence, candidate.Position()))
}

package spellchecker

import (
	"github.com/alldroll/suggest/pkg/index"
	lm "github.com/alldroll/suggest/pkg/language-model"
)

// lmScorer implements the scorer interface
type lmScorer struct {
	model    lm.LanguageModel
	sentence []lm.WordID
}

// Score returns the score of the given position
func (s *lmScorer) Score(position index.Position) float64 {
	return s.model.ScoreWordIDs(append(s.sentence, position))
}

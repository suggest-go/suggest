package spellchecker

import (
	"sort"

	"github.com/alldroll/suggest/pkg/index"
	lm "github.com/alldroll/suggest/pkg/language-model"
)

// lmScorer implements the scorer interface
type lmScorer struct {
	model    lm.LanguageModel
	sentence []lm.WordID
	next     []lm.WordCount
}

// Score returns the score of the given position
func (s *lmScorer) Score(position index.Position) float64 {
	i := sort.Search(len(s.next), func(i int) bool { return s.next[i] >= position })

	if i < 0 || i >= len(s.next) || s.next[i] != position {
		return lm.UnknownWordScore
	}

	return s.model.ScoreWordIDs(append(s.sentence, position))
}

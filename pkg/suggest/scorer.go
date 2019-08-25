package suggest

import "github.com/alldroll/suggest/pkg/index"

// Scorer is responsible for scoring an index position
type Scorer interface {
	// Score returns the score of the given position
	Score(position index.Position) float64
}

type dummyScorer struct {
}

// Score returns the score of the given position
func (d dummyScorer) Score(position index.Position) float64 {
	return -float64(position)
}

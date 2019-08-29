package suggest

import "github.com/alldroll/suggest/pkg/index"

// Scorer is responsible for scoring an index position
type Scorer interface {
	// Score returns the score of the given position
	Score(position index.Position) float64
}


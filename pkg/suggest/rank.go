package suggest

import "github.com/alldroll/suggest/pkg/index"

// Rank is the interface that tells who is the best candidate for a query according to some criterion
type Rank interface {
	// Less tells if the document b is more relevant than the document a
	Less(a, b index.Position) bool
}

// idOrderRank is the rank function, that ranks documents with lower ids
type idOrderRank struct {}

// Less tells if the document b is more relevant than the document a
func (r *idOrderRank) Less(a, b index.Position) bool {
	return a < b
}

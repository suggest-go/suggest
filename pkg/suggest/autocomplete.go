package suggest

import (
	"fmt"

	"github.com/alldroll/suggest/pkg/analysis"
	"github.com/alldroll/suggest/pkg/index"
)

// Autocomplete provides autocomplete functionality
// for candidates search
type Autocomplete interface {
	// Autocomplete returns candidates where the query string is a substring of each candidate
	Autocomplete(query string, limit int, scorer Scorer) ([]Candidate, error)
}

// NewAutocomplete creates a new instance of Autocomplete
func NewAutocomplete(
	index index.InvertedIndex,
	searcher index.Searcher,
	tokenizer analysis.Tokenizer,
) Autocomplete {
	return &nGramAutocomplete{
		index:     index,
		searcher:  searcher,
		tokenizer: tokenizer,
		ranker:    &idOrderRank{},
	}
}

// nGramAutocomplete implements Autocomplete interface
type nGramAutocomplete struct {
	index     index.InvertedIndex
	searcher  index.Searcher
	tokenizer analysis.Tokenizer
	ranker    Rank
}

// Autocomplete returns candidates where the query string is a substring of each candidate
func (n *nGramAutocomplete) Autocomplete(query string, limit int, scorer Scorer) ([]Candidate, error) {
	selector := NewTopKCollectorWithRanker(limit, n.ranker)
	set := n.tokenizer.Tokenize(query)
	candidates, err := n.searcher.Search(n.index, set, len(set))

	if err != nil {
		return nil, fmt.Errorf("failed to search posting lists: %v", err)
	}

	for _, c := range candidates {
		selector.Add(c.Position(), scorer.Score(c.Position()))
	}

	return selector.GetCandidates(), nil
}

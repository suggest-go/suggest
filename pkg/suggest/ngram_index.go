package suggest

// NGramIndex is the interface that provides the access to
// approximate string search and autocomplete
type NGramIndex interface {
	Suggester
	Autocomplete
}

// NewNGramIndex creates a new instance of NGramIndex
func NewNGramIndex(suggester Suggester, autocomplete Autocomplete) NGramIndex {
	return &nGramIndex{
		suggester:    suggester,
		autocomplete: autocomplete,
	}
}

type nGramIndex struct {
	suggester    Suggester
	autocomplete Autocomplete
}

// Suggest returns top-k similar candidates
func (n *nGramIndex) Suggest(config *SearchConfig) ([]Candidate, error) {
	return n.suggester.Suggest(config)
}

// Autocomplete returns candidates where the query string is a substring of each candidate
func (n *nGramIndex) Autocomplete(query string, collectorManager CollectorManager) ([]Candidate, error) {
	return n.autocomplete.Autocomplete(query, collectorManager)
}

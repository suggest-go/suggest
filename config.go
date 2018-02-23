package suggest

import (
	"fmt"
	"github.com/alldroll/suggest/alphabet"
	"github.com/alldroll/suggest/dictionary"
	"github.com/alldroll/suggest/metric"
)

const (
	// MinNGramSize is a minimum allowed size of ngram
	MinNGramSize = 2
	// MaxNGramSize is a maximum allowed size of ngram
	MaxNGramSize = 4
)

// IndexConfig is config for NgramIndex structure
// deprecated
type IndexConfig struct {
	nGramSize  int
	alphabet   alphabet.Alphabet
	wrap       string
	pad        string
	dictionary dictionary.Dictionary
}

// NewIndexConfig returns new instance of IndexConfig
func NewIndexConfig(k int, dictionary dictionary.Dictionary, alphabet alphabet.Alphabet, wrap, pad string) (*IndexConfig, error) {
	if k < MinNGramSize || k > MaxNGramSize {
		return nil, fmt.Errorf("k should be in [%d, %d]", MinNGramSize, MaxNGramSize)
	}

	if alphabet.Size() == 0 {
		return nil, fmt.Errorf("Alphabet should not be empty")
	}

	return &IndexConfig{
		nGramSize:  k,
		alphabet:   alphabet,
		wrap:       wrap,
		pad:        pad,
		dictionary: dictionary,
	}, nil
}

// SearchConfig is a config for NGramIndex Suggest method
type SearchConfig struct {
	query      string
	topK       int
	metric     metric.Metric
	similarity float64
}

// NewSearchConfig returns new instance of SearchConfig
func NewSearchConfig(query string, topK int, metric metric.Metric, similarity float64) (*SearchConfig, error) {
	if topK < 0 {
		return nil, fmt.Errorf("topK is invalid") //TODO fixme
	}

	if similarity <= 0 || similarity > 1 {
		return nil, fmt.Errorf("similarity shouble be in (0.0, 1.0]")
	}

	return &SearchConfig{
		query:      query,
		topK:       topK,
		metric:     metric,
		similarity: similarity,
	}, nil
}

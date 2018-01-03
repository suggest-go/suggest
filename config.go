package suggest

import (
	"fmt"
)

const (
	// MinNGramSize is a minimum allowed size of ngram
	MinNGramSize = 2
	// MaxNGramSize is a maximum allowed size of ngram
	MaxNGramSize = 4
)

// IndexConfig is config for NgramIndex structure
type IndexConfig struct {
	ngramSize  int
	alphabet   Alphabet
	wrap       string
	pad        string
	dictionary Dictionary
}

// NewIndexConfig returns new instance of IndexConfig
func NewIndexConfig(k int, dictionary Dictionary, alphabet Alphabet, wrap, pad string) (*IndexConfig, error) {
	if k < MinNGramSize || k > MaxNGramSize {
		return nil, fmt.Errorf("k should be in [%d, %d]", MinNGramSize, MaxNGramSize)
	}

	if alphabet.Size() == 0 {
		return nil, fmt.Errorf("Alphabet should not be empty")
	}

	return &IndexConfig{
		k,
		alphabet,
		wrap,
		pad,
		dictionary,
	}, nil
}

// SearchConfig is a config for NGramIndex Suggest method
type SearchConfig struct {
	query      string
	topK       int
	metric     Metric
	similarity float64
}

// NewSearchConfig returns new instance of SearchConfig
func NewSearchConfig(query string, topK int, metric Metric, similarity float64) (*SearchConfig, error) {
	if topK < 0 {
		return nil, fmt.Errorf("topK is invalid") //TODO fixme
	}

	if similarity <= 0 || similarity > 1 {
		return nil, fmt.Errorf("similarity shouble be in (0.0, 1.0]")
	}

	return &SearchConfig{
		query,
		topK,
		metric,
		similarity,
	}, nil
}

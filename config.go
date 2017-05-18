package suggest

import (
	"fmt"
)

const (
	MIN_NGRAM_SIZE = 2
	MAX_NGRAM_SIZE = 4
)

//
type IndexConfig struct {
	ngramSize int
	alphabet  Alphabet
	wrap      string
	pad       string
}

//
func NewIndexConfig(k int, alphabet Alphabet, wrap, pad string) (*IndexConfig, error) {
	if k < MIN_NGRAM_SIZE || k > MAX_NGRAM_SIZE {
		return nil, fmt.Errorf("k should be in [%d, %d]", MIN_NGRAM_SIZE, MAX_NGRAM_SIZE)
	}

	if len(alphabet.Chars()) == 0 {
		return nil, fmt.Errorf("Alphabet should not be empty")
	}

	return &IndexConfig{
		k,
		alphabet,
		wrap,
		pad,
	}, nil
}

//
type SearchConfig struct {
	query       string
	topK        int
	measureName MeasureT
	similarity  float64
}

//
func NewSearchConfig(query string, topK int, measureName MeasureT, similarity float64) (*SearchConfig, error) {
	if topK < 0 {
		return nil, fmt.Errorf("topK is invalid") //TODO fixme
	}

	if similarity <= 0 || similarity > 1 {
		return nil, fmt.Errorf("similarity shouble be in (0.0, 1.0]")
	}

	return &SearchConfig{
		query,
		topK,
		measureName,
		similarity,
	}, nil
}

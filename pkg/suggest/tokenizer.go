package suggest

import (
	"github.com/suggest-go/suggest/pkg/alphabet"
	"github.com/suggest-go/suggest/pkg/analysis"
)

// NewSuggestTokenizer creates a tokenizer for suggester service
func NewSuggestTokenizer(d IndexDescription) analysis.Tokenizer {
	filter := analysis.NewNormalizerFilter(alphabet.CreateAlphabet(d.Alphabet), d.Pad)

	return analysis.NewWrapTokenizer(
		analysis.NewFilterTokenizer(
			analysis.NewNGramTokenizer(d.NGramSize),
			filter,
		),
		d.Wrap[0],
		d.Wrap[1],
	)
}

// NewAutocompleteTokenizer creates a tokenizer for autocomplete service
func NewAutocompleteTokenizer(d IndexDescription) analysis.Tokenizer {
	filter := analysis.NewNormalizerFilter(alphabet.CreateAlphabet(d.Alphabet), d.Pad)

	return analysis.NewWrapTokenizer(
		analysis.NewFilterTokenizer(
			analysis.NewNGramTokenizer(d.NGramSize),
			filter,
		),
		d.Wrap[0],
		"", // do not add a wrap symbol to the tail of query
	)
}

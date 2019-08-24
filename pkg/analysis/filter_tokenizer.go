package analysis

import "strings"

// filterTokenizer performs tokenize and filter operations
type filterTokenizer struct {
	tokenizer Tokenizer
	filter    TokenFilter
}

// NewFilterTokenizer creates a new instance of filter tokenizer
func NewFilterTokenizer(tokenizer Tokenizer, filter TokenFilter) Tokenizer {
	return &filterTokenizer{
		tokenizer: tokenizer,
		filter:    filter,
	}
}

// Tokenize splits the given text on a sequence of tokens
func (t *filterTokenizer) Tokenize(text string) []Token {
	text = strings.ToLower(text)
	text = strings.Trim(text, " ")

	tokens := t.tokenizer.Tokenize(text)

	return t.filter.Filter(tokens)
}

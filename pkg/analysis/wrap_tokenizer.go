package analysis

type wrapTokenizer struct {
	tokenizer  Tokenizer
	start, end string
}

// NewWrapTokenizer returns a tokenizer that performs wrap the provided text before tokenization
func NewWrapTokenizer(tokenizer Tokenizer, start, end string) Tokenizer {
	return &wrapTokenizer{
		tokenizer: tokenizer,
		start:     start,
		end:       end,
	}
}

// Tokenize splits the given text on a sequence of tokens
func (t *wrapTokenizer) Tokenize(text string) []Token {
	return t.tokenizer.Tokenize(t.start + text + t.end)
}

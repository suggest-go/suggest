package lm

import (
	"strings"

	"github.com/suggest-go/suggest/pkg/alphabet"
	"github.com/suggest-go/suggest/pkg/analysis"
)

// Token is a string with an assigned and thus identified meaning
type Token = analysis.Token

// NewTokenizer creates a new instance of Tokenizer
func NewTokenizer(alphabet alphabet.Alphabet) analysis.Tokenizer {
	return &tokenizer{
		tokenizer: analysis.NewWordTokenizer(alphabet),
	}
}

// tokenizer implements Tokenizer interface
type tokenizer struct {
	tokenizer analysis.Tokenizer
}

// Tokenize splits the given text on a sequence of tokens
func (t *tokenizer) Tokenize(text string) []Token {
	text = strings.ToLower(text)
	text = strings.Trim(text, " ")

	return t.tokenizer.Tokenize(text)
}

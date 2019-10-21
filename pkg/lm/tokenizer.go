package lm

import (
	"github.com/suggest-go/suggest/pkg/alphabet"
	"github.com/suggest-go/suggest/pkg/analysis"
	"strings"
)

// Token is a string with an assigned and thus identified meaning
type Token = analysis.Token

// NewTokenizer creates a new instance of Tokenizer
func NewTokenizer(alphabet alphabet.Alphabet) analysis.Tokenizer {
	filter := &lmFilter{analysis.NewNormalizerFilter(alphabet, " ")}

	return analysis.NewFilterTokenizer(
		analysis.NewWordTokenizer(alphabet),
		filter,
	)
}

type lmFilter struct {
	filter analysis.TokenFilter
}

func (l *lmFilter) Filter(tokens []Token) []Token {
	tokens = l.filter.Filter(tokens)

	for i, token := range tokens {
		tokens[i] = strings.Trim(token, " -'")
	}

	return tokens
}

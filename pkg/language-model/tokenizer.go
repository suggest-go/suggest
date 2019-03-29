package lm

import (
	"unicode/utf8"

	"github.com/alldroll/suggest/pkg/alphabet"
)

// Token is a string with an assigned and thus identified meaning
type Token = string

// Tokenizer performs splitting the given text on tokens
type Tokenizer interface {
	// Splits the given text on a sequence of tokens
	Tokenize(text string) []Token
}

// NewTokenizer creates tokenizer entity
func NewTokenizer(alphabet alphabet.Alphabet) *tokenizer {
	return &tokenizer{
		alphabet: alphabet,
	}
}

// dummy tokenizer
type tokenizer struct {
	alphabet alphabet.Alphabet
}

// Splits the given text on a sequence of tokens
func (t *tokenizer) Tokenize(text string) []Token {
	words := []Token{}
	wordStart, wordLen := -1, 0

	for i, char := range text {
		if t.alphabet.Has(char) {
			if wordStart == -1 {
				wordStart = i
			}

			wordLen += utf8.RuneLen(char)
		} else {
			if wordStart != -1 {
				words = append(words, text[wordStart:wordStart+wordLen])
			}

			wordStart, wordLen = -1, 0
		}
	}

	if wordStart != -1 {
		words = append(words, text[wordStart:wordStart+wordLen])
	}

	return words
}

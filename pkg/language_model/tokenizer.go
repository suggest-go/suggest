package language_model

import (
	"github.com/alldroll/suggest/pkg/alphabet"
	"unicode/utf8"
)

type Word = string

type Tokenizer interface {
	Tokenize(text string) []Word
}

func NewTokenizer(alphabet alphabet.Alphabet) *tokenizer {
	return &tokenizer{
		alphabet: alphabet,
	}
}

// dummy tokenizer
type tokenizer struct {
	alphabet alphabet.Alphabet
}

func (t *tokenizer) Tokenize(text string) []Word {
	words := []Word{}
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

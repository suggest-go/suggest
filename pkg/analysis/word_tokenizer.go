package analysis

import (
	"unicode/utf8"

	"github.com/suggest-go/suggest/pkg/alphabet"
)

// NewWordTokenizer creates a new instance of Tokenizer
func NewWordTokenizer(alphabet alphabet.Alphabet) Tokenizer {
	return &wordTokenizer{
		alphabet: alphabet,
	}
}

// tokenizer implements Tokenizer interface
type wordTokenizer struct {
	alphabet alphabet.Alphabet
}

// Tokenize splits the given text on a sequence of tokens
func (t *wordTokenizer) Tokenize(text string) []Token {
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

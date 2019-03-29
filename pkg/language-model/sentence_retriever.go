package lm

import (
	"bufio"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/alldroll/suggest/pkg/alphabet"
)

// Sentence is a sequence of tokens
type Sentence = []Token

// SentenceRetriever is an entity that responsible for retrieving
// sentences from the given source
type SentenceRetriever interface {
	// Retrieves and returns the next sentence from the source
	Retrieve() Sentence
}

// NewSentenceRetriever
func NewSentenceRetriever(tokenizer Tokenizer, reader io.Reader, alphabet alphabet.Alphabet) *retriever {
	scanner := bufio.NewScanner(reader)

	r := &retriever{
		scanner:   scanner,
		tokenizer: tokenizer,
		alphabet:  alphabet,
	}

	scanner.Split(r.scanSentence)
	return r
}

// retiever implements SentenceRetriever interface
type retriever struct {
	scanner   *bufio.Scanner
	tokenizer Tokenizer
	alphabet  alphabet.Alphabet
}

// Retrieves and returns the next sentence from the source
func (r *retriever) Retrieve() Sentence {
	if r.scanner.Scan() {
		text := r.scanner.Text()
		return r.tokenizer.Tokenize(strings.ToLower(text))
	}

	return nil
}

//
func (r *retriever) scanSentence(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0
	for width := 0; start < len(data); start += width {
		var char rune
		char, width = utf8.DecodeRune(data[start:])
		if !r.alphabet.Has(char) {
			break
		}
	}

	for width, i := 0, start; i < len(data); i += width {
		var char rune
		char, width = utf8.DecodeRune(data[i:])
		if r.alphabet.Has(char) {
			return i + width, data[start:i], nil
		}
	}

	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	return start, nil, nil
}

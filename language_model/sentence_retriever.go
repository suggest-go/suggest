package language_model

import (
	"bufio"
	"github.com/alldroll/suggest/alphabet"
	"io"
	"unicode/utf8"
)

type Sentence = []Word

type SentenceRetriver interface {
	Retrieve() Sentence
}

func NewSentenceRetriver(tokenizer Tokenizer, reader io.Reader, alphabet alphabet.Alphabet) *retriever {
	scanner := bufio.NewScanner(reader)

	r := &retriever{
		scanner:   scanner,
		tokenizer: tokenizer,
		alphabet:  alphabet,
	}

	scanner.Split(r.scanSentence)
	return r
}

type retriever struct {
	scanner   *bufio.Scanner
	tokenizer Tokenizer
	alphabet  alphabet.Alphabet
}

func (r *retriever) Retrieve() Sentence {
	if r.scanner.Scan() {
		text := r.scanner.Text()
		t := r.tokenizer.Tokenize(text)
		return t
	}

	return nil
}

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

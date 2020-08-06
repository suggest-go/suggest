package analysis

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suggest-go/suggest/pkg/alphabet"
)

func TestTokenize(t *testing.T) {
	testCases := []struct {
		text     string
		expected []Token
	}{
		{"i wanna rock", []Token{"i", "wanna", "rock"}},
		{"", []Token{}},
		{"!!! test $-)", []Token{"test"}},
		{"    ", []Token{}},
		{"hello, привет, 33", []Token{"hello", "привет", "33"}},
	}

	tokenizer := NewWordTokenizer(alphabet.NewCompositeAlphabet(
		[]alphabet.Alphabet{
			alphabet.NewEnglishAlphabet(),
			alphabet.NewRussianAlphabet(),
			alphabet.NewNumberAlphabet(),
		},
	))

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test #%d", i+1), func(t *testing.T) {
			assert.Equal(t, testCase.expected, tokenizer.Tokenize(testCase.text))
		})
	}
}

func BenchmarkWordTokenizer(b *testing.B) {
	tokenizer := NewWordTokenizer(alphabet.NewCompositeAlphabet(
		[]alphabet.Alphabet{
			alphabet.NewEnglishAlphabet(),
			alphabet.NewRussianAlphabet(),
			alphabet.NewNumberAlphabet(),
		},
	))

	for i := 0; i < b.N; i++ {
		tokenizer.Tokenize("Hello, how are you. How to say in russian Привет?")
	}
}

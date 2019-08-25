package analysis

import (
	"reflect"
	"testing"

	"github.com/alldroll/suggest/pkg/alphabet"
)

func TestTokenize(t *testing.T) {
	cases := []struct {
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

	for _, c := range cases {
		actual := tokenizer.Tokenize(c.text)

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("Test fail, expected %v, got %v", c.expected, actual)
		}
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

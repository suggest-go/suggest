package language_model

import (
	"github.com/alldroll/suggest/alphabet"
	"reflect"
	"testing"
)

func TestEnglishTokenize(t *testing.T) {
	cases := []struct {
		text     string
		expected []Word
	}{
		{"i wanna rock", []Word{"i", "wanna", "rock"}},
		{"", []Word{}},
		{"!!! test $-)", []Word{"test"}},
		{"    ", []Word{}},
		{"hello, привет, 33", []Word{"hello", "привет", "33"}},
	}

	tokenizer := NewTokenizer(alphabet.NewCompositeAlphabet(
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

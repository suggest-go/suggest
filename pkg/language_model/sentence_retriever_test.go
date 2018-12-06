package language_model

import (
	"github.com/alldroll/suggest/pkg/alphabet"
	"reflect"
	"strings"
	"testing"
)

func TestSentenceRetrieve(t *testing.T) {
	cases := []struct {
		text     string
		expected []Sentence
	}{
		{
			"i wanna rock. hello my friend. what? dab. чтоооо. ты - не я",
			[]Sentence{
				Sentence{"i", "wanna", "rock"},
				Sentence{"hello", "my", "friend"},
				Sentence{"what"},
				Sentence{"dab"},
				Sentence{"чтоооо"},
				Sentence{"ты", "не", "я"},
			},
		},
	}

	tokenizer := NewTokenizer(alphabet.NewCompositeAlphabet(
		[]alphabet.Alphabet{
			alphabet.NewEnglishAlphabet(),
			alphabet.NewRussianAlphabet(),
			alphabet.NewNumberAlphabet(),
		},
	))

	stopAlphabet := alphabet.NewSimpleAlphabet([]rune{'.', '?', '!'})

	for _, c := range cases {
		retriever := NewSentenceRetriver(
			tokenizer,
			strings.NewReader(c.text),
			stopAlphabet,
		)

		actual := []Sentence{}

		for s := retriever.Retrieve(); s != nil; s = retriever.Retrieve() {
			actual = append(actual, s)
		}

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("Test fail, expected %v, got %v", c.expected, actual)
		}
	}
}

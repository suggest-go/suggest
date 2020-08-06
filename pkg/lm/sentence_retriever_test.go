package lm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suggest-go/suggest/pkg/alphabet"
)

func TestSentenceRetrieve(t *testing.T) {
	testCases := []struct {
		text     string
		expected []Sentence
	}{
		{
			"i wanna rock. hello my friend. what? dab. чтоооо. ты - не я",
			[]Sentence{
				{"i", "wanna", "rock"},
				{"hello", "my", "friend"},
				{"what"},
				{"dab"},
				{"чтоооо"},
				{"ты", "не", "я"},
			},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test #%d", i+1), func(t *testing.T) {
			tokenizer := NewTokenizer(alphabet.NewCompositeAlphabet(
				[]alphabet.Alphabet{
					alphabet.NewEnglishAlphabet(),
					alphabet.NewRussianAlphabet(),
					alphabet.NewNumberAlphabet(),
				},
			))

			stopAlphabet := alphabet.NewSimpleAlphabet([]rune{'.', '?', '!'})

			retriever := NewSentenceRetriever(
				tokenizer,
				strings.NewReader(testCase.text),
				stopAlphabet,
			)

			actual := []Sentence{}

			for s := retriever.Retrieve(); s != nil; s = retriever.Retrieve() {
				actual = append(actual, s)
			}

			assert.Equal(t, testCase.expected, actual)
		})
	}
}

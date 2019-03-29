package lm

import (
	"reflect"
	"strings"
	"testing"

	"github.com/alldroll/suggest/pkg/alphabet"
)

const (
	start = "<S>"
	end   = "</S>"
)

func TestGenerate(t *testing.T) {
	cases := []struct {
		text     string
		order    uint8
		expected []NGrams
	}{
		{
			"london! london is the capital of great britain.",
			uint8(1),
			[]NGrams{
				NGrams{
					NGram{start},
					NGram{"london"},
					NGram{end},
				},

				NGrams{
					NGram{start},
					NGram{"london"},
					NGram{"is"},
					NGram{"the"},
					NGram{"capital"},
					NGram{"of"},
					NGram{"great"},
					NGram{"britain"},
					NGram{end},
				},
			},
		},
		{
			"london! london is the capital of great britain.",
			uint8(2),
			[]NGrams{
				NGrams{
					NGram{start, "london"},
					NGram{"london", end},
				},

				NGrams{
					NGram{start, "london"},
					NGram{"london", "is"},
					NGram{"is", "the"},
					NGram{"the", "capital"},
					NGram{"capital", "of"},
					NGram{"of", "great"},
					NGram{"great", "britain"},
					NGram{"britain", end},
				},
			},
		},
		{
			"london! london is the capital of great britain.",
			uint8(3),
			[]NGrams{
				NGrams{
					NGram{start, "london", end},
				},

				NGrams{
					NGram{start, "london", "is"},
					NGram{"london", "is", "the"},
					NGram{"is", "the", "capital"},
					NGram{"the", "capital", "of"},
					NGram{"capital", "of", "great"},
					NGram{"of", "great", "britain"},
					NGram{"great", "britain", end},
				},
			},
		},
	}

	tokenizer := NewTokenizer(alphabet.NewEnglishAlphabet())
	stopAlphabet := alphabet.NewSimpleAlphabet([]rune{'.', '?', '!'})

	for _, c := range cases {
		retriever := NewSentenceRetriever(
			tokenizer,
			strings.NewReader(c.text),
			stopAlphabet,
		)

		generator := NewGenerator(
			c.order,
			start,
			end,
		)

		for _, expected := range c.expected {
			sentence := retriever.Retrieve()
			actual := generator.Generate(sentence)

			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("Test fail, expected %v, got %v", expected, actual)
			}
		}
	}
}

package index

import (
	"github.com/alldroll/suggest/pkg/alphabet"
	"testing"
)

func TestPrepareString(t *testing.T) {
	cases := []struct {
		word, expected string
	}{
		{"", "$$"},
		{"test", "$test$"},
		{"helLo world", "$hello*world$"},
		{"-+=tesla", "$*tesla$"},
	}

	alphabet := alphabet.NewCompositeAlphabet([]alphabet.Alphabet{
		alphabet.NewEnglishAlphabet(),
		alphabet.NewNumberAlphabet(),
		alphabet.NewRussianAlphabet(),
		alphabet.NewSimpleAlphabet([]rune{'$', '*'}),
	})

	clean := NewCleaner(alphabet.Chars(), "*", [2]string{"$", "$"})
	for _, c := range cases {
		actual := clean.CleanAndWrap(c.word)
		if actual != c.expected {
			t.Errorf(
				"Test Fail, expected %v, got %v",
				c.expected,
				actual,
			)
		}
	}
}

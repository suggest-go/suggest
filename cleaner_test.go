package suggest

import "testing"

func TestPrepareString(t *testing.T) {
	cases := []struct {
		word, expected string
	}{
		{"", "$$"},
		{"test", "$test$"},
		{"helLo world", "$hello*world$"},
		{"-+=tesla", "$*tesla$"},
	}

	alphabet := NewCompositeAlphabet([]Alphabet{
		NewEnglishAlphabet(),
		NewNumberAlphabet(),
		NewRussianAlphabet(),
		NewSimpleAlphabet([]rune{'$', '*'}),
	})

	clean := NewCleaner(alphabet.Chars(), "*", "$")
	for _, c := range cases {
		actual := clean.Clean(c.word)
		if actual != c.expected {
			t.Errorf(
				"Test Fail, expected %v, got %v",
				c.expected,
				actual,
			)
		}
	}
}

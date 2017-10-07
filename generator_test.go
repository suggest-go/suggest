package suggest

import (
	"reflect"
	"testing"
)

func TestSplitIntoNGrams(t *testing.T) {
	cases := []struct {
		word   string
		k      int
		ngrams []string
	}{
		{
			"tet", 2,
			[]string{"te", "et"},
		},
		{
			"te", 2,
			[]string{"te"},
		},
		{
			"testing", 3,
			[]string{"tes", "est", "sti", "tin", "ing"},
		},
		{
			"жигули", 2,
			[]string{"жи", "иг", "гу", "ул", "ли"},
		},
	}

	for _, c := range cases {
		actual := splitIntoNGrams(c.word, c.k)
		if !reflect.DeepEqual(actual, c.ngrams) {
			t.Errorf(
				"Test Fail, expected %v, got %v",
				c.ngrams,
				actual,
			)
		}
	}
}

func TestPrepareString(t *testing.T) {
	cases := []struct {
		word, expected string
	}{
		{"", "$$"},
		{"test", "$test$"},
		{"helLo world", "$hello$world$"},
		{"-+=tesla", "$$tesla$"},
	}

	alphabet := NewCompositeAlphabet([]Alphabet{
		NewEnglishAlphabet(),
		NewNumberAlphabet(),
		NewRussianAlphabet(),
		NewSimpleAlphabet([]rune{'$'}),
	})

	clean := NewCleaner(alphabet.Chars(), "$", "$")
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

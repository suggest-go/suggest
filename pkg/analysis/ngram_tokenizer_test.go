package analysis

import (
	"reflect"
	"testing"
)

func TestTokenizeNGrams(t *testing.T) {
	cases := []struct {
		word   string
		k      int
		ngrams []Token
	}{
		{
			"tet",
			2,
			[]Token{"te", "et"},
		},
		{
			"te",
			2,
			[]Token{"te"},
		},
		{
			"testing",
			3,
			[]Token{"tes", "est", "sti", "tin", "ing"},
		},
		{
			"жигули",
			2,
			[]Token{"жи", "иг", "гу", "ул", "ли"},
		},
		{
			"",
			2,
			[]Token{},
		},
		{
			"lalala",
			2,
			[]Token{"la", "al"},
		},
	}

	for _, c := range cases {
		tokenizer := NewNGramTokenizer(c.k)
		actual := tokenizer.Tokenize(c.word)

		if !reflect.DeepEqual(actual, c.ngrams) {
			t.Errorf(
				"Test Fail, expected %v, got %v",
				c.ngrams,
				actual,
			)
		}
	}
}

func BenchmarkNGramTokenizer(b *testing.B) {
	tokenizer := NewNGramTokenizer(3)

	for i := 0; i < b.N; i++ {
		tokenizer.Tokenize("abcdefghkl123456йцукен")
	}
}

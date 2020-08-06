package analysis

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizeNGrams(t *testing.T) {
	testCases := []struct {
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

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test #%d", i+1), func(t *testing.T) {
			tokenizer := NewNGramTokenizer(testCase.k)
			actual := tokenizer.Tokenize(testCase.word)
			assert.Equal(t, testCase.ngrams, actual)
		})
	}
}

func BenchmarkNGramTokenizer(b *testing.B) {
	tokenizer := NewNGramTokenizer(3)

	for i := 0; i < b.N; i++ {
		tokenizer.Tokenize("abcdefghkl123456йцукен")
	}
}

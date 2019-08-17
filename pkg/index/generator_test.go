package index

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
		generator := NewGenerator(c.k)
		actual := generator.Generate(c.word)

		if !reflect.DeepEqual(actual, c.ngrams) {
			t.Errorf(
				"Test Fail, expected %v, got %v",
				c.ngrams,
				actual,
			)
		}
	}
}

func BenchmarkGenerate(b *testing.B) {
	generator := NewGenerator(3)

	for i := 0; i < b.N; i++ {
		generator.Generate("abcdefghkl123456йцукен")
	}
}

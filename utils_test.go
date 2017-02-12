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
	}

	for _, c := range cases {
		actual := SplitIntoNGrams(c.word, c.k)
		if !reflect.DeepEqual(actual, c.ngrams) {
			t.Errorf(
				"Test Fail, expected %v, got %v",
				c.ngrams,
				actual,
			)
		}
	}
}

func BenchmarkSplitIntoNGrams(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SplitIntoNGrams("TestStringAbabacaMacacaTsaksn", 3)
	}
}

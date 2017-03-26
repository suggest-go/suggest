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

func TestPrepareString(t *testing.T) {
	cases := []struct {
		word, expected string
	}{
		{"", ""},
		{"test", "$test$"},
		{"hello world", "$hello$world$"},
	}

	for _, c := range cases {
		actual := prepareString(c.word)
		if actual != c.expected {
			t.Errorf(
				"Test Fail, expected %v, got %v",
				c.expected,
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

func BenchmarkGetProfile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a := "SomeLongWordsadsadsadsadsadsadsadsadsadssadsada"
		getProfile(a, 3)
	}
}

func BenchmarkGetWordProfile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a := "SomeLongWordsadsadsadsadsadsadsadsadsadssadsada"
		GetWordProfile(a, 3)
	}
}

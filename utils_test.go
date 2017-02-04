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

func TestLevenshtein(t *testing.T) {
	cases := []struct {
		a, b     string
		distance int
	}{
		{
			"tet", "tet",
			0,
		},
		{
			"tes", "pep",
			2,
		},
	}

	for _, c := range cases {
		actual := Levenshtein(c.a, c.b)
		if actual != c.distance {
			t.Errorf(
				"Test Fail, expected %v, got %v",
				c.distance,
				actual,
			)
		}
	}
}

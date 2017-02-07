package suggest

import "testing"

func TestLevenshtein(t *testing.T) {
	cases := []struct {
		a, b     string
		distance float64
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

	editDistance, _ := GetEditDistance(LEVENSHTEIN, 2)
	for _, c := range cases {
		actual := editDistance.Calc(c.a, c.b)
		if actual != c.distance {
			t.Errorf(
				"Test Fail, expected %v, got %v",
				c.distance,
				actual,
			)
		}
	}
}

func TestNGramDistance(t *testing.T) {
	cases := []struct {
		a, b     string
		distance float64
	}{
		{
			"01000", "001111",
			5,
		},
		{
			"ababaca", "ababaca",
			0,
		},
	}

	editDistance, _ := GetEditDistance(NGRAM, 2)
	for _, c := range cases {
		distance := editDistance.Calc(c.a, c.b)
		if distance != c.distance {
			t.Errorf("TestFail, expected {%v}, got {%v}", c.distance, distance)
		}
	}
}

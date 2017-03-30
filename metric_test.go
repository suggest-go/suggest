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

	editDistance := &LevenshteinDistance{}
	for _, c := range cases {
		actual := editDistance.Calc(GetWordProfile(c.a, 3), GetWordProfile(c.b, 3))
		if actual != c.distance {
			t.Errorf(
				"Test Fail, expected %v, got %v",
				c.distance,
				actual,
			)
		}
	}
}

/*
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

func BenchmarkNGramDistance(b *testing.B) {
	editDistance, _ := GetEditDistance(NGRAM, 3)
	for i := 0; i < b.N; i++ {
		editDistance.Calc("SomeLongWordssdasddsadsadsadasdasdsadasdsadasdasdsadsad", "SomeLongWordsadsadsadsadsadsadsadsadsadssadsada")
	}
}


*/
func BenchmarkLevenshtein(b *testing.B) {
	editDistance := &LevenshteinDistance{}
	for i := 0; i < b.N; i++ {
		a, b := "SomeLongWordssdasddsadsadsadasdasdsadasdsadasdasdsadsad", "SomeLongWordsadsadsadsadsadsadsadsadsadssadsada"
		editDistance.Calc(GetWordProfile(a, 3), GetWordProfile(b, 3))
	}
}

func BenchmarkJaccard(b *testing.B) {
	editDistance := &JaccardDistance{3}
	for i := 0; i < b.N; i++ {
		a, b := "SomeLongWordssdasddsadsadsadasdasdsadasdsadasdasdsadsad", "SomeLongWordsadsadsadsadsadsadsadsadsadssadsada"
		editDistance.Calc(GetWordProfile(a, 3), GetWordProfile(b, 3))
	}
}

package suggest

import (
	"reflect"
	"testing"
)

func TestNGramDistance(t *testing.T) {
	cases := []struct {
		a, b     string
		distance int
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

	ngramIndex := NewNGramIndex(2)
	for _, c := range cases {
		distance := ngramIndex.distance(c.a, c.b)
		if distance != c.distance {
			t.Errorf("TestFail, expected {%v}, got {%v}", c.distance, distance)
		}
	}
}

func TestFuzzySearch(t *testing.T) {
	collection := []string{
		"blue",
		"blunder",
		"blunt",
		"flank",
		"flu",
		"fluence",
		"fluent",
		"flunker",
		"test",
		"tes hello",
	}

	expected := map[string]int{
		"flu":     4,
		"fluence": 8,
		"fluent":  7,
		"blue":    9,
		"blunder": 10,
		"blunt":   8,
		"flank":   4,
		"flunker": 4,
	}

	ngramIndex := NewNGramIndex(2)
	for _, word := range collection {
		ngramIndex.AddWord(word)
	}

	candidates := ngramIndex.FuzzySearch("flunk")
	for _, candidate := range candidates {
		if rank, ok := expected[candidate.word]; !ok || rank != candidate.distance {
			t.Errorf("TestFail, expected {%v}, got {%v}", candidates, expected)
		}
	}
}

func TestSuggestAuto(t *testing.T) {
	collection := []string{
		"Nissan March",
		"Nissan Juke",
		"Nissan Maxima",
		"Nissan Murano",
		"Nissan Note",
		"Toyota Mark II",
		"Toyota Corolla",
		"Toyota Corona",
	}

	ngramIndex := NewNGramIndex(3)

	for _, word := range collection {
		ngramIndex.AddWord(word)
	}

	candidates := ngramIndex.Suggest("Nissan mar", 2)
	expected := []string{
		"Nissan March",
		"Nissan Maxima",
	}

	if !reflect.DeepEqual(expected, candidates) {
		t.Errorf(
			"Test Fail, expected %v, got %v",
			expected,
			candidates,
		)
	}
}

package suggest

import (
	"reflect"
	"testing"
)

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

	expected := map[string]float64{
		"flu":     4,
		"fluence": 8,
		"fluent":  7,
		"blue":    9,
		"blunder": 10,
		"blunt":   8,
		"flank":   4,
		"flunker": 4,
	}

	ngramIndex := NewNGramIndex(2, NGRAM)
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

	ngramIndex := NewNGramIndex(3, NGRAM)

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

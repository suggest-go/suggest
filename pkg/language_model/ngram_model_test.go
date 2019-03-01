package language_model

import (
	"math"
	"testing"
)

const tolerance = 0.00001

func TestScore(t *testing.T) {
	model := NewNGramModel(3)
	indexer := NewIndexer()

	data := []struct {
		nGrams []Word
		count  WordCount
	}{
		{[]Word{"<S>"}, 3},
		{[]Word{"</S>"}, 3},
		{[]Word{"I"}, 3},
		{[]Word{"am"}, 2},
		{[]Word{"Sam"}, 2},
		{[]Word{"<S>", "I"}, 1},
		{[]Word{"<S>", "Sam"}, 1},
		{[]Word{"I", "am"}, 2},
		{[]Word{"am", "</S>"}, 1},
		{[]Word{"am", "Sam"}, 1},
		{[]Word{"Sam", "</S>"}, 1},
		{[]Word{"Sam", "I"}, 1},
		{[]Word{"<S>", "I", "am"}, 1},
		{[]Word{"I", "am", "Sam"}, 1},
		{[]Word{"am", "Sam", "</S>"}, 1},
		{[]Word{"<S>", "Sam", "I"}, 1},
		{[]Word{"I", "am", "</S>"}, 1},
		{[]Word{"Sam", "I", "am"}, 1},
	}

	ids := make([]WordId, 0, 3)

	for _, c := range data {
		for _, nGram := range c.nGrams {
			ids = append(ids, indexer.GetOrCreate(nGram))
		}

		model.Put(ids, c.count)
		ids = ids[:0]
	}

	cases := []struct {
		nGrams        []Word
		expectedScore float64
	}{
		{[]Word{"I", "am", "Sam"}, 0.5},
		{[]Word{"I", "am"}, 0.66666},
		{[]Word{"Sam", "I", "am"}, 1},
		{[]Word{"Sam", "am", "I"}, 0.02461},
		{[]Word{"I", "dont", "know"}, 0.03692},
		{[]Word{"no", "one", "word"}, 0.0},
	}

	for _, c := range cases {
		for _, nGram := range c.nGrams {
			ids = append(ids, indexer.GetOrCreate(nGram))
		}

		actual := model.Score(ids)
		ids = ids[:0]

		if diff := math.Abs(actual - c.expectedScore); diff >= tolerance {
			t.Errorf(
				"Test fail, for %v expected score %v, got %v",
				c.nGrams,
				c.expectedScore,
				actual,
			)
		}

	}
}

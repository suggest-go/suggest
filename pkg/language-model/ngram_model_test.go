package lm

import (
	"bytes"
	"encoding/gob"
	"math"
	"reflect"
	"testing"
)

const tolerance = 0.0001

func TestScoreFromFile(t *testing.T) {
	indexer := NewIndexer()
	reader := NewGoogleNGramReader(3, indexer, "testdata/fixtures")

	model, err := reader.Read()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	testModel(t, model, indexer)
}

func TestPredict(t *testing.T) {
	indexer := NewIndexer()
	reader := NewGoogleNGramReader(3, indexer, "testdata/fixtures")
	ids := make([]WordID, 0, 3)

	model, err := reader.Read()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	cases := []struct {
		nGrams   Sentence
		expected []Token
	}{
		{Sentence{"i", "am"}, []Token{"sam", "</S>"}},
		{Sentence{"i"}, []Token{"am", "do"}},
		{Sentence{"green"}, []Token{"eggs"}},
	}

	for _, c := range cases {
		for _, nGram := range c.nGrams {
			ids = append(ids, indexer.Get(nGram))
		}

		list, err := model.Next(ids)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		ids = ids[:0]
		actual := []Token{}

		for _, item := range list {
			token, err := indexer.Find(item)
			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}

			actual = append(actual, token)
		}

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf(
				"Test fail, expected %v, got %v",
				c.expected,
				actual,
			)
		}
	}
}

func TestBinaryMarshalling(t *testing.T) {
	indexer := NewIndexer()
	reader := NewGoogleNGramReader(3, indexer, "testdata/fixtures")

	expected, err := reader.Read()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	var network bytes.Buffer

	// Create an encoder and send a value.
	enc := gob.NewEncoder(&network)
	err = enc.Encode(&expected)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Create a decoder and receive a value.
	dec := gob.NewDecoder(&network)
	var actual NGramModel
	err = dec.Decode(&actual)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	testModel(t, expected, indexer)
	testModel(t, actual, indexer)
}

func testModel(t *testing.T, model NGramModel, indexer Indexer) {
	ids := make([]WordID, 0, 3)

	cases := []struct {
		nGrams        Sentence
		expectedScore float64
	}{
		{Sentence{"i", "am", "sam"}, -0.6931},
		{Sentence{"i", "am"}, -0.4054},
		{Sentence{"sam", "i", "am"}, 0},
		{Sentence{"sam", "am", "i"}, -4.1351},
		{Sentence{"i", "dont", "know"}, -3.7297},
		{Sentence{"no", "one", "word"}, -100},
	}

	for _, c := range cases {
		for _, nGram := range c.nGrams {
			ids = append(ids, indexer.Get(nGram))
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
package lm

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/suggest-go/suggest/pkg/store"
	"math"
	"testing"
)

const tolerance = 0.0001

func TestScoreFromFile(t *testing.T) {
	indexer, err := buildIndexerWithInMemoryDictionary("testdata/fixtures/1-gm")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	directory, err := store.NewFSDirectory("testdata/fixtures")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	reader := NewGoogleNGramReader(3, indexer, directory)

	model, err := reader.Read()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	testModel(t, model, indexer)
}

func TestPredict(t *testing.T) {
	cases := []struct {
		nGrams   Sentence
		word     string
		expected float64
	}{
		{
			nGrams: Sentence{"i", "am"},
			word: "sam",
			expected: -0.6931,
		},
		{
			nGrams: Sentence{"i", "am"},
			word: "</S>",
			expected: -0.6931,
		},
		{
			nGrams: Sentence{"i"},
			word: "am",
			expected: -0.4054,
		},
		{
			nGrams: Sentence{"i"},
			word: "do",
			expected: -1.0986,
		},
		{
			nGrams: Sentence{"green"},
			word: "eggs",
			expected: 0.0,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("predict #%d", i), func(t *testing.T) {
			indexer, err := buildIndexerWithInMemoryDictionary("testdata/fixtures/1-gm")

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			directory, err := store.NewFSDirectory("testdata/fixtures")

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			reader := NewGoogleNGramReader(3, indexer, directory)
			ids := make([]WordID, 0, 3)

			model, err := reader.Read()

			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}

			for _, nGram := range c.nGrams {
				id, _ := indexer.Get(nGram)
				ids = append(ids, id)
			}

			scorerNext, err := model.Next(ids)

			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}

			id, _ := indexer.Get(c.word)
			actual := scorerNext.ScoreNext(id)

			if diff := math.Abs(c.expected - actual); diff >= tolerance {
				t.Errorf(
					"Test fail, for %v expected score %v, got %v",
					c.nGrams,
					c.expected,
					actual,
				)
			}
		})
	}
}

func TestBinaryMarshalling(t *testing.T) {
	indexer, err := buildIndexerWithInMemoryDictionary("testdata/fixtures/1-gm")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	directory, err := store.NewFSDirectory("testdata/fixtures")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	reader := NewGoogleNGramReader(3, indexer, directory)

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
			id, _ := indexer.Get(nGram)
			ids = append(ids, id)
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

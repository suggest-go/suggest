package lm

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suggest-go/suggest/pkg/store"
)

const tolerance = 0.0001

func TestScoreFromFile(t *testing.T) {
	indexer, err := buildIndexerWithInMemoryDictionary("testdata/fixtures/1-gm")
	assert.NoError(t, err)

	directory, err := store.NewFSDirectory("testdata/fixtures")
	assert.NoError(t, err)

	reader := NewGoogleNGramReader(3, indexer, directory)
	model, err := reader.Read()
	assert.NoError(t, err)

	testModel(t, model, indexer)
}

func TestPredict(t *testing.T) {
	testCases := []struct {
		nGrams   Sentence
		word     string
		expected float64
	}{
		{
			nGrams:   Sentence{"i", "am"},
			word:     "sam",
			expected: -0.6931,
		},
		{
			nGrams:   Sentence{"i", "am"},
			word:     "</S>",
			expected: -0.6931,
		},
		{
			nGrams:   Sentence{"i"},
			word:     "am",
			expected: -0.4054,
		},
		{
			nGrams:   Sentence{"i"},
			word:     "do",
			expected: -1.0986,
		},
		{
			nGrams:   Sentence{"green"},
			word:     "eggs",
			expected: 0.0,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("predict #%d", i), func(t *testing.T) {
			indexer, err := buildIndexerWithInMemoryDictionary("testdata/fixtures/1-gm")
			assert.NoError(t, err)

			directory, err := store.NewFSDirectory("testdata/fixtures")
			assert.NoError(t, err)

			reader := NewGoogleNGramReader(3, indexer, directory)
			ids := make([]WordID, 0, 3)

			model, err := reader.Read()
			assert.NoError(t, err)

			for _, nGram := range testCase.nGrams {
				id, _ := indexer.Get(nGram)
				ids = append(ids, id)
			}

			scorerNext, err := model.Next(ids)
			assert.NoError(t, err)

			id, _ := indexer.Get(testCase.word)
			actual := scorerNext.ScoreNext(id)
			assert.Less(t, math.Abs(testCase.expected-actual), tolerance)
		})
	}
}

func TestBinaryMarshalling(t *testing.T) {
	indexer, err := buildIndexerWithInMemoryDictionary("testdata/fixtures/1-gm")
	assert.NoError(t, err)

	directory, err := store.NewFSDirectory("testdata/fixtures")
	assert.NoError(t, err)

	reader := NewGoogleNGramReader(3, indexer, directory)
	expected, err := reader.Read()
	assert.NoError(t, err)

	outDir := store.NewRAMDirectory()
	output, _ := outDir.CreateOutput("lm")

	// Create an encoder and send a value.
	_, err = expected.Store(output)
	assert.NoError(t, err)

	err = output.Close()
	assert.NoError(t, err)

	// Create a decoder and receive a value.
	input, _ := outDir.OpenInput("lm")
	actual := NewNGramModel(nil)

	_, err = actual.Load(input)
	assert.NoError(t, err)

	testModel(t, expected, indexer)
	testModel(t, actual, indexer)
}

func testModel(t *testing.T, model NGramModel, indexer Indexer) {
	ids := make([]WordID, 0, 3)

	testCases := []struct {
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

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test #%d", i+1), func(t *testing.T) {
			for _, nGram := range testCase.nGrams {
				id, _ := indexer.Get(nGram)
				ids = append(ids, id)
			}

			actual := model.Score(ids)
			ids = ids[:0]
			assert.Less(t, math.Abs(testCase.expectedScore-actual), tolerance)
		})
	}
}

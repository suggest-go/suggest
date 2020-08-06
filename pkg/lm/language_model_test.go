package lm

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suggest-go/suggest/pkg/store"
)

func TestScoreSentenceFromFile(t *testing.T) {
	config := Config{
		NGramOrder:  3,
		StartSymbol: "<S>",
		EndSymbol:   "</S>",
		OutputPath:  "testdata/fixtures",
		basePath:    ".",
	}

	indexer, err := buildIndexerWithInMemoryDictionary("testdata/fixtures/1-gm")
	assert.NoError(t, err)

	directory, err := store.NewFSDirectory(config.GetOutputPath())
	assert.NoError(t, err)

	reader := NewGoogleNGramReader(config.NGramOrder, indexer, directory)
	model, err := reader.Read()
	assert.NoError(t, err)

	lm, err := NewLanguageModel(model, indexer, &config)
	assert.NoError(t, err)

	testLM(lm, t)
}

func TestScoreSentenceFromBinary(t *testing.T) {
	config, err := ReadConfig("testdata/config-example.json")
	assert.NoError(t, err)

	directory, err := store.NewFSDirectory(config.GetOutputPath())
	assert.NoError(t, err)

	lm, err := RetrieveLMFromBinary(directory, config)
	assert.NoError(t, err)

	testLM(lm, t)
}

func testLM(lm LanguageModel, t *testing.T) {
	testCases := []struct {
		sentence      Sentence
		expectedScore float64
	}{
		{Sentence{"i", "am", "sam"}, -1.3862},
		{Sentence{"i", "am"}, -1.3862},
		{Sentence{"sam", "i", "am"}, -0.6931},
		{Sentence{"sam", "am", "i"}, -10.2852},
		{Sentence{"i", "dont", "know"}, -105.0514},
		{Sentence{"no", "one", "word"}, -203.7297},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test #%d", i+1), func(t *testing.T) {
			actual, _ := lm.ScoreSentence(testCase.sentence)
			diff := math.Abs(actual - testCase.expectedScore)
			assert.Less(t, diff, tolerance)
		})
	}
}

package lm

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// NGramReader is responsible for creating NGramModel from the files
type NGramReader interface {
	// Read builds NGramModel from the given list of readers
	Read() (NGramModel, error)
}

// googleNGramFormatReader implements NGramReader with google nGram format storage
type googleNGramFormatReader struct {
	indexer    Indexer
	sourcePath string
	nGramOrder uint8
}

// NewGoogleNGramReader creates new instance of NGramReader
func NewGoogleNGramReader(nGramOrder uint8, indexer Indexer, sourcePath string) NGramReader {
	return &googleNGramFormatReader{
		nGramOrder: nGramOrder,
		indexer:    indexer,
		sourcePath: sourcePath,
	}
}

// Read builds NGramModel from the given list of readers
func (gr *googleNGramFormatReader) Read() (NGramModel, error) {
	if gr.nGramOrder == 0 {
		return nil, errors.New("nGramOrder should be >= 1")
	}

	vectors := []NGramVector{}
	nGrams := make([]WordID, 0, int(gr.nGramOrder))

	for i := 0; i < int(gr.nGramOrder); i++ {
		f, err := os.Open(fmt.Sprintf(fileFormat, gr.sourcePath, i+1))

		if err != nil {
			return nil, err
		}

		defer f.Close()
		scanner := bufio.NewScanner(f)
		builder := NewNGramVectorBuilder(vectors)

		for scanner.Scan() {
			line := scanner.Text()
			tabIndex := strings.Index(line, "\t")

			for _, word := range strings.Split(line[:tabIndex], " ") {
				nGrams = append(nGrams, gr.indexer.GetOrCreate(word))
			}

			count, err := strconv.ParseUint(line[tabIndex+1:], 10, 32)
			if err != nil {
				return nil, err
			}

			builder.Put(nGrams, WordCount(count))
			nGrams = nGrams[:0]
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		vectors = append(vectors, builder.Build())
	}

	return NewNGramModel(vectors), nil
}

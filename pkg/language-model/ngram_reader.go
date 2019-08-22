package lm

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/alldroll/suggest/pkg/store"
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
	nGramOrder uint8
	directory store.Directory
}

// NewGoogleNGramReader creates new instance of NGramReader
func NewGoogleNGramReader(nGramOrder uint8, indexer Indexer, directory store.Directory) NGramReader {
	return &googleNGramFormatReader{
		nGramOrder: nGramOrder,
		indexer:    indexer,
		directory: directory,
	}
}

// Read builds NGramModel from the given list of readers
func (gr *googleNGramFormatReader) Read() (NGramModel, error) {
	if gr.nGramOrder == 0 {
		return nil, errors.New("nGramOrder should be >= 1")
	}

	vectors := make([]NGramVector, 0, int(gr.nGramOrder))
	nGrams := make([]WordID, 0, int(gr.nGramOrder))

	for i := 0; i < int(gr.nGramOrder); i++ {
		in, err := gr.directory.OpenInput(fmt.Sprintf(fileFormat, i+1))

		if err != nil {
			return nil, fmt.Errorf("failed to open a ngram input: %v", err)
		}

		scanner := bufio.NewScanner(in)
		builder := NewNGramVectorBuilder(vectors)

		for scanner.Scan() {
			line := scanner.Text()
			tabIndex := strings.Index(line, "\t")

			for _, word := range strings.Split(line[:tabIndex], " ") {
				index, err := gr.indexer.Get(word)

				if err != nil {
					return nil, err
				}

				nGrams = append(nGrams, index)
			}

			count, err := strconv.ParseUint(line[tabIndex+1:], 10, 32)

			if err != nil {
				return nil, fmt.Errorf("ngram file is corrupted, expected number: %v", err)
			}

			if err := builder.Put(nGrams, WordCount(count)); err != nil {
				return nil, fmt.Errorf("failed to add nGrams to a builder: %v", err)
			}

			nGrams = nGrams[:0]
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		vectors = append(vectors, builder.Build())

		if err := in.Close(); err != nil {
			return nil, fmt.Errorf("failed to close an input source: %v", err)
		}
	}

	return NewNGramModel(vectors), nil
}


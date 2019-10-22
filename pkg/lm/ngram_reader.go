package lm

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/suggest-go/suggest/pkg/store"
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
	directory  store.Directory
}

// NewGoogleNGramReader creates new instance of NGramReader
func NewGoogleNGramReader(nGramOrder uint8, indexer Indexer, directory store.Directory) NGramReader {
	return &googleNGramFormatReader{
		nGramOrder: nGramOrder,
		indexer:    indexer,
		directory:  directory,
	}
}

// Read builds NGramModel from the given list of readers
func (gr *googleNGramFormatReader) Read() (NGramModel, error) {
	if gr.nGramOrder == 0 {
		return nil, errors.New("nGramOrder should be >= 1")
	}

	vectors := make([]NGramVector, 0, int(gr.nGramOrder))

	for i := 0; i < int(gr.nGramOrder); i++ {
		builder := NewNGramVectorBuilder(vectors)

		if err := gr.readNGramVector(builder, i+1); err != nil {
			return nil, fmt.Errorf("failed to read %d ngram vector: %v", i+1, err)
		}

		vectors = append(vectors, builder.Build())
	}

	return NewNGramModel(vectors), nil
}

// readNGramVector reads nGram vector for the given order
func (gr *googleNGramFormatReader) readNGramVector(builder NGramVectorBuilder, order int) error {
	in, err := gr.directory.OpenInput(fmt.Sprintf(fileFormat, order))

	if err != nil {
		return fmt.Errorf("failed to open a ngram input: %v", err)
	}

	nGrams := make([]WordID, 0, order)
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		line := scanner.Text()
		tabIndex := strings.Index(line, "\t")

		for _, word := range strings.Split(line[:tabIndex], " ") {
			index, err := gr.indexer.Get(word)

			if err != nil {
				return err
			}

			nGrams = append(nGrams, index)
		}

		count, err := strconv.ParseUint(line[tabIndex+1:], 10, 32)

		if err != nil {
			return fmt.Errorf("ngram file is corrupted, expected number: %v", err)
		}

		if err := builder.Put(nGrams, WordCount(count)); err != nil {
			return fmt.Errorf("failed to add nGrams to a builder: %v", err)
		}

		nGrams = nGrams[:0]
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan ngram file format: %v", err)
	}

	return nil
}

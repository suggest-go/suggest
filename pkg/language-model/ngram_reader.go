package lm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//
type NGramReader interface {
	//
	Read() (NGramModel, error)
}

type googleNGramFormatReader struct {
	indexer    Indexer
	sourcePath string
	nGramOrder uint8
}

//
func NewGoogleNGramReader(nGramOrder uint8, indexer Indexer, sourcePath string) *googleNGramFormatReader {
	return &googleNGramFormatReader{
		nGramOrder: nGramOrder,
		indexer:    indexer,
		sourcePath: sourcePath,
	}
}

func (gr *googleNGramFormatReader) Read() (NGramModel, error) {
	model := NewNGramModel(gr.nGramOrder)
	nGrams := make([]WordID, 0, int(gr.nGramOrder))

	for i := 0; i < int(gr.nGramOrder); i++ {
		f, err := os.Open(fmt.Sprintf(fileFormat, gr.sourcePath, i+1))

		if err != nil {
			return nil, err
		}

		defer f.Close()

		scanner := bufio.NewScanner(f)

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

			model.Put(nGrams, WordCount(count))
			nGrams = nGrams[:0]
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return model, nil
}

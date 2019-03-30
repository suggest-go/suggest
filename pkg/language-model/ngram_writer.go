package lm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	fileFormat  = "%s/%d-gm"
	nGramFormat = "%s\t%d\n"
)

// NGramWriter is the interface that persists the NGram Count Trie to a storage
type NGramWriter interface {
	// Write persists the given trie to a storage
	Write(trie CountTrie) error
}

// NewGoogleNGramWriter creates new instance of NGramWriter that persists the given NGram Count Trie with
// Google NGram Format negotiations
func NewGoogleNGramWriter(indexer Indexer, nGramOrder uint8, outputPath string) NGramWriter {
	return &googleNGramFormatWriter{
		indexer:    indexer,
		nGramOrder: nGramOrder,
		outputPath: outputPath,
	}
}

// googleNGramFormatWriter is the entity that imlements NGramWriter interface
type googleNGramFormatWriter struct {
	indexer    Indexer
	outputPath string
	nGramOrder uint8
}

// Write persists the given trie into a storage
func (gw *googleNGramFormatWriter) Write(trie CountTrie) (err error) {
	bufs := []*bufio.Writer{}

	for i := 0; i < int(gw.nGramOrder); i++ {
		f, err := os.OpenFile(fmt.Sprintf(fileFormat, gw.outputPath, i+1), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

		if err != nil {
			return err
		}

		defer f.Close()

		buf := bufio.NewWriter(f)
		bufs = append(bufs, buf)
		defer buf.Flush()
	}

	grams := make([]string, 0, int(gw.nGramOrder))

	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
		}
	}()

	trie.Walk(func(path []WordID, count WordCount) {
		if len(path) == 0 {
			return
		}

		grams = grams[:0]

		for _, g := range path {
			nGram, err := gw.indexer.Find(g)
			if err != nil {
				panic(err)
			}

			grams = append(grams, nGram)
		}

		fmt.Fprintf(bufs[len(path)-1], nGramFormat, strings.Join(grams, " "), count)
	})

	return nil
}

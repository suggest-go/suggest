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
func NewGoogleNGramWriter(nGramOrder uint8, outputPath string) NGramWriter {
	return &googleNGramFormatWriter{
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
		f, err := openFile(fmt.Sprintf(fileFormat, gw.outputPath, i+1))

		if err != nil {
			return err
		}

		defer f.Close()

		buf := bufio.NewWriter(f)
		bufs = append(bufs, buf)
		defer buf.Flush()
	}

	err = trie.Walk(func(nGrams []Token, count WordCount) {
		if len(nGrams) == 0 {
			return
		}

		joined := strings.Join(nGrams, " ")
		fmt.Fprintf(bufs[len(nGrams)-1], nGramFormat, joined, count)
	})

	if err != nil {
		return err
	}

	return nil
}

// openFile opens a file for writing with necessary flags
func openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
}

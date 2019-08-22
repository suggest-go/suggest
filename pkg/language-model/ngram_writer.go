package lm

import (
	"fmt"
	"github.com/alldroll/suggest/pkg/store"
	"strings"
)

const (
	fileFormat  = "%d-gm"
	nGramFormat = "%s\t%d\n"
)

// NGramWriter is the interface that persists the NGram Count Trie to a storage
type NGramWriter interface {
	// Write persists the given trie to a storage
	Write(trie CountTrie) error
}

// NewGoogleNGramWriter creates new instance of NGramWriter that persists the given NGram Count Trie with
// Google NGram Format negotiations
func NewGoogleNGramWriter(nGramOrder uint8, directory store.Directory) NGramWriter {
	return &googleNGramFormatWriter{
		nGramOrder: nGramOrder,
		directory: directory,
	}
}

// googleNGramFormatWriter is the entity that implements NGramWriter interface
type googleNGramFormatWriter struct {
	indexer    Indexer
	nGramOrder uint8
	directory store.Directory
}

// Write persists the given trie into a storage
func (gw *googleNGramFormatWriter) Write(trie CountTrie) error {
	outs := make([]store.Output, int(gw.nGramOrder))

	for i := range outs {
		out, err := gw.directory.CreateOutput(fmt.Sprintf(fileFormat, i+1))

		if err != nil {
			return fmt.Errorf("failed to create an output: %v", err)
		}

		outs[i] = out
	}

	err := trie.Walk(func(nGrams []Token, count WordCount) error {
		if len(nGrams) == 0 {
			return nil
		}

		joined := strings.Join(nGrams, " ")

		if _, err := fmt.Fprintf(outs[len(nGrams)-1], nGramFormat, joined, count); err != nil {
			return fmt.Errorf("failed to print nGrams: %v", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	for _, out := range outs {
		if err := out.Close(); err != nil {
			return fmt.Errorf("failed to close an output: %v", err)
		}
	}

	return nil
}

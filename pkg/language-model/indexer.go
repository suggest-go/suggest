package lm

import (
	"fmt"

	"github.com/alldroll/suggest/pkg/dictionary"
)

// WordID is an index of the corresponding word
type WordID = uint32

const (
	// UnknownWordID is an index of an unregistered word
	UnknownWordID = uint32(0xffffffff)
	// UnknownWordSymbol is a symbol that is returned on an unregistred word
	UnknownWordSymbol = "<UNK>"
)

// Indexer enumerates words in the vocabulary of a language model. Stores a two-way
// mapping between uint32 and strings.
type Indexer interface {
	// Returns the index for the word, otherwise returns UnknownWordID
	Get(token Token) WordID
	// Find a token by the given index
	Find(id WordID) (Token, error)
}

// indexerImpl implements Indexer interface
type indexerImpl struct {
	dictionary dictionary.Dictionary
	table      map[Token]WordID
}

// Returns the index for the word, otherwise returns UnknownWordID
func (i *indexerImpl) Get(token Token) WordID {
	index, ok := i.table[token]

	if !ok {
		index = UnknownWordID
	}

	return index
}

// Find a token by the given index
func (i *indexerImpl) Find(index WordID) (Token, error) {
	val, err := i.dictionary.Get(index)

	if err != nil {
		return UnknownWordSymbol, fmt.Errorf("Failed to get index from the dictionary: %v", err)
	}

	if val == dictionary.NilValue {
		return UnknownWordSymbol, nil
	}

	return val, nil
}

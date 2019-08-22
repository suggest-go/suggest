package lm

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alldroll/suggest/pkg/dictionary"
	"github.com/alldroll/suggest/pkg/mph"
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
	Get(token Token) (WordID, error)
	// Find a token by the given index
	Find(id WordID) (Token, error)
}

// BuildIndexer builds a indexer from the given dictionary
func BuildIndexer(dict dictionary.Dictionary) (Indexer, error) {
	table := mph.New()

	if err := table.Build(dict); err != nil {
		return nil, err
	}

	return NewIndexer(dict, table), nil
}

// NewIndexer creates a new instance of indexer
func NewIndexer(dict dictionary.Dictionary, table mph.MPH) Indexer {
	return &indexerImpl{
		dictionary: dict,
		table:      table,
	}
}

// indexerImpl implements Indexer interface
type indexerImpl struct {
	dictionary dictionary.Dictionary
	table      mph.MPH
}

// Returns the index for the word, otherwise returns UnknownWordID
func (i *indexerImpl) Get(token Token) (WordID, error) {
	index := i.table.Get(token)
	stored, err := i.dictionary.Get(index)

	if err != nil {
		return UnknownWordID, fmt.Errorf("failed to get index from the dictionary: %v", err)
	}

	if stored != token {
		index = UnknownWordID
	}

	return index, nil
}

// Find a token by the given index
func (i *indexerImpl) Find(index WordID) (Token, error) {
	val, err := i.dictionary.Get(index)

	if err != nil {
		return UnknownWordSymbol, fmt.Errorf("failed to get index from the dictionary: %v", err)
	}

	if val == dictionary.NilValue {
		return UnknownWordSymbol, nil
	}

	return val, nil
}

// buildIndexerWithInMemoryDictionary builds an indexer with a ram dictionary for the given path
func buildIndexerWithInMemoryDictionary(dictPath string) (Indexer, error) {
	f, err := os.Open(dictPath)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	collection := []Token{}

	for scanner.Scan() {
		line := scanner.Text()
		tabIndex := strings.Index(line, "\t")
		collection = append(collection, line[:tabIndex])
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if err := f.Close(); err != nil {
		return nil, err
	}

	return BuildIndexer(dictionary.NewInMemoryDictionary(collection))
}

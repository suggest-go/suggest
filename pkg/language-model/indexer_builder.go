package lm

import (
	"fmt"

	"github.com/alldroll/suggest/pkg/dictionary"
)

// BuildIndexer builds a indexer from the given dictionary
func BuildIndexer(dict dictionary.Dictionary) (indexer Indexer, err error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
		}
	}()

	table := map[Token]WordID{}

	err = dict.Iterate(func(index dictionary.Key, word dictionary.Value) {
		if _, ok := table[word]; ok {
			panic(fmt.Errorf("Dictionary contains nonunique value: %s", word))
		}

		table[word] = index
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to build a indexer from the dictionary: %v", err)
	}

	indexer = &indexerImpl{
		dictionary: dict,
		table:      table,
	}

	return indexer, err
}

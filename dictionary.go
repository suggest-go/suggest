package suggest

import (
	"github.com/alldroll/suggest/dictionary"
	"io"
)

// NewInMemoryDictionary creates new instance of inMemoryDictionary
func NewInMemoryDictionary(words []string) dictionary.Dictionary {
	return dictionary.NewInMemoryDictionary(words)
}

// NewCDBDictionary creates new instance of cdbDictionary
func NewCDBDictionary(reader io.ReaderAt) dictionary.Dictionary {
	return dictionary.NewCDBDictionary(reader)
}

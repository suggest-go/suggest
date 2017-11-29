package suggest

import (
	"errors"
	"github.com/alldroll/go-datastructures/cdb"
	"io"
)

// WordKey represents key in key-value pair for Dictionary
type WordKey interface{}

// Dictionary is an abstract data type composed of a collection of (key, value) pairs
type Dictionary interface {
	// Get returns value associated with a particular key
	Get(key WordKey) (string, error)
	// Iterator returns an iterator over the elements in this dictionary
	Iterator() DictionaryIterator
}

// DictionaryIterator represents Iterator of Dictionary
type DictionaryIterator interface {
	// Next moves iterator to the next item. Returns true on success otherwise false
	Next() bool
	// GetPair returns key-value pair of current item
	GetPair() (WordKey, string)
}

// inMemoryDictionary implements Dictionary with in-memory data access
type inMemoryDictionary struct {
	holder []string
	index  int
}

// NewInMemoryDictionary creates new instance of inMemoryDictionary
func NewInMemoryDictionary(words []string) Dictionary {
	holder := make([]string, len(words))
	copy(holder, words)

	return &inMemoryDictionary{
		holder,
		len(words),
	}
}

func (d *inMemoryDictionary) Get(key WordKey) (string, error) {
	k := key.(int)
	if k < 0 || k >= len(d.holder) {
		return "", errors.New("Key is not exists")
	}

	return d.holder[k], nil
}

func (d *inMemoryDictionary) Iterator() DictionaryIterator {
	return &inMemoryDictIter{d, 0}
}

// inMemoryDictIter implements interface DictionaryIterator for inMemoryDictionary
type inMemoryDictIter struct {
	dict  *inMemoryDictionary
	index int
}

func (i *inMemoryDictIter) Next() bool {
	success := false
	if i.index + 1 < i.dict.index {
		i.index++
		success = true
	}

	return success
}

func (i *inMemoryDictIter) IsValid() bool {
	return i.index < i.dict.index
}

func (i *inMemoryDictIter) GetPair() (WordKey, string) {
	if i.index >= i.dict.index {
		return nil, ""
	}

	val, _ := i.dict.Get(i.index)
	return i.index, val
}

type cdbDictionary struct {
	reader cdb.Reader
}

//
func NewCDBDictionary(r io.ReaderAt) Dictionary {
	handle := cdb.New()
	reader, err := handle.GetReader(r)
	if err != nil {
		panic(err)
	}

	return &cdbDictionary{
		reader,
	}
}

//
func (d *cdbDictionary) Get(key WordKey) (string, error) {
	value, err := d.reader.Get(key.([]byte))
	if err != nil {
		return "", err
	}

	return string(value), nil
}

//
func (d *cdbDictionary) Iterator() DictionaryIterator {
	iterator, err := d.reader.Iterator()
	if err != nil {
		panic(err)
	}

	return &cdbDictIter{iterator}
}

//
type cdbDictIter struct {
	cdbIterator cdb.Iterator
}

//
func (i *cdbDictIter) Next() bool {
	ok, err := i.cdbIterator.Next()
	if err != nil {
		panic(err)
	}

	return ok
}

//
func (i *cdbDictIter) GetPair() (WordKey, string) {
	value, key := i.cdbIterator.Value(), i.cdbIterator.Key()
	if key == nil {
		return nil, ""
	}

	return key, string(value)
}


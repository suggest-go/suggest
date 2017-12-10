package suggest

import (
	"errors"
	"github.com/alldroll/cdb"
	"encoding/binary"
	"io"
)

// WordKey represents key in key-value pair for Dictionary

// Dictionary is an abstract data type composed of a collection of (key, value) pairs
type Dictionary interface {
	// Get returns value associated with a particular key
	Get(key int) (string, error)
	// Iterator returns an iterator over the elements in this dictionary
	Iterator() DictionaryIterator
}

// DictionaryIterator represents Iterator of Dictionary
type DictionaryIterator interface {
	// Next moves iterator to the next item. Returns true on success otherwise false
	Next() bool
	// GetPair returns key-value pair of current item
	GetPair() (int, string)
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

func (d *inMemoryDictionary) Get(key int) (string, error) {
	if key < 0 || key >= len(d.holder) {
		return "", errors.New("Key is not exists")
	}

	return d.holder[key], nil
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

func (i *inMemoryDictIter) GetPair() (int, string) {
	if i.index >= i.dict.index {
		return 0, ""
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
func (d *cdbDictionary) Get(key int) (string, error) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(key))
	value, err := d.reader.Get(bs)
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
func (i *cdbDictIter) GetPair() (int, string) {
	value, key := i.cdbIterator.Value(), i.cdbIterator.Key()
	if key == nil {
		return 0, ""
	}

	return int(binary.LittleEndian.Uint32(key)), string(value)
}

package suggest

import (
	"encoding/binary"
	"errors"
	"github.com/alldroll/cdb"
	"io"
)

// Dictionary is an abstract data type composed of a collection of (key, value) pairs
type Dictionary interface {
	// Get returns value associated with a particular key
	Get(key Position) (string, error)
	// Iterator returns an iterator over the elements in this dictionary
	Iterator() DictionaryIterator
}

// DictionaryIterator represents Iterator of Dictionary
type DictionaryIterator interface {
	// Next moves iterator to the next item. Returns true on success otherwise false
	Next() bool
	// GetPair returns key-value pair of current item
	GetPair() (Position, string)
}

// inMemoryDictionary implements Dictionary with in-memory data access
type inMemoryDictionary struct {
	holder []string
}

// NewInMemoryDictionary creates new instance of inMemoryDictionary
func NewInMemoryDictionary(words []string) Dictionary {
	holder := make([]string, len(words))
	copy(holder, words)

	return &inMemoryDictionary{
		holder: holder,
	}
}

// Get returns value associated with a particular key
func (d *inMemoryDictionary) Get(key Position) (string, error) {
	if key < 0 || int(key) >= len(d.holder) {
		return "", errors.New("Key is not exists")
	}

	return d.holder[key], nil
}

// Iterator returns an iterator over the elements in this dictionary
func (d *inMemoryDictionary) Iterator() DictionaryIterator {
	return &inMemoryDictionaryIterator{
		dict:  d,
		index: 0,
	}
}

// inMemoryDictionaryIterator implements interface DictionaryIterator for inMemoryDictionary
type inMemoryDictionaryIterator struct {
	dict  *inMemoryDictionary
	index Position
}

// Next moves iterator to the next item. Returns true on success otherwise false
func (i *inMemoryDictionaryIterator) Next() bool {
	success := false
	if int(i.index+1) < len(i.dict.holder) {
		i.index++
		success = true
	}

	return success
}

// GetPair returns key-value pair of current item
func (i *inMemoryDictionaryIterator) GetPair() (Position, string) {
	if int(i.index) >= len(i.dict.holder) {
		return 0, ""
	}

	val, _ := i.dict.Get(i.index)
	return i.index, val
}

// cdbDictionary implements Dictionary with cdb as database
type cdbDictionary struct {
	reader cdb.Reader
}

// NewCDBDictionary creates new instance of cdbDictionary
func NewCDBDictionary(r io.ReaderAt) Dictionary {
	handle := cdb.New()
	reader, err := handle.GetReader(r)
	if err != nil {
		panic(err)
	}

	return &cdbDictionary{
		reader: reader,
	}
}

// Get returns value associated with a particular key
func (d *cdbDictionary) Get(key Position) (string, error) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(key))
	value, err := d.reader.Get(bs)
	if err != nil {
		return "", err
	}

	return string(value), nil
}

// Iterator returns an iterator over the elements in this dictionary
func (d *cdbDictionary) Iterator() DictionaryIterator {
	iterator, err := d.reader.Iterator()
	if err != nil {
		panic(err)
	}

	return &cdbDictionaryIterator{cdbIterator: iterator}
}

// cdbDictionaryIterator implements interface DictionaryIterator for cdbDictionary
type cdbDictionaryIterator struct {
	cdbIterator cdb.Iterator
}

// Next moves iterator to the next item. Returns true on success otherwise false
func (i *cdbDictionaryIterator) Next() bool {
	ok, err := i.cdbIterator.Next()
	if err != nil {
		panic(err)
	}

	return ok
}

// GetPair returns key-value pair of current item
func (i *cdbDictionaryIterator) GetPair() (Position, string) {
	value, key := i.cdbIterator.Value(), i.cdbIterator.Key()
	if key == nil {
		return 0, ""
	}

	return Position(binary.LittleEndian.Uint32(key)), string(value)
}

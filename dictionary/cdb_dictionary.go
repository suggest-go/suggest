package dictionary

import (
	"encoding/binary"
	"github.com/alldroll/cdb"
	"io"
)

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
func (d *cdbDictionary) Get(key Key) (Value, error) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, key)
	value, err := d.reader.Get(bs)
	if err != nil {
		return "", err
	}

	return Value(value), nil
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
func (i *cdbDictionaryIterator) GetPair() (Key, Value) {
	record := i.cdbIterator.Record()

	keyReader, keySize := record.Key()
	key := make([]byte, keySize)
	if _, err := keyReader.Read(key); err != nil {
		panic(err)
	}

	if key == nil {
		return 0, ""
	}

	valueReader, valSize := record.Value()
	value := make([]byte, valSize)
	if _, err := valueReader.Read(value); err != nil {
		panic(err)
	}

	return binary.LittleEndian.Uint32(key), Value(value)
}

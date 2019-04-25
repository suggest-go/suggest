package dictionary

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/alldroll/cdb"
)

// cdbDictionary implements Dictionary with cdb as database
type cdbDictionary struct {
	reader cdb.Reader
}

// NewCDBDictionary creates new instance of cdbDictionary
func NewCDBDictionary(r io.ReaderAt) (Dictionary, error) {
	handle := cdb.New()
	reader, err := handle.GetReader(r)

	if err != nil {
		return nil, fmt.Errorf("Fail to create cdb dictionary: %v", err)
	}

	return &cdbDictionary{
		reader: reader,
	}, nil
}

// Get returns value associated with a particular key
func (d *cdbDictionary) Get(key Key) (Value, error) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, key)
	value, err := d.reader.Get(bs)

	if err != nil {
		return NilValue, err
	}

	if value == nil {
		return NilValue, nil
	}

	return Value(value), nil
}

// Size returns the size of the dictionary
func (d *cdbDictionary) Size() int {
	return d.reader.Size()
}

// Iterator returns an iterator over the elements in this dictionary
func (d *cdbDictionary) Iterate(iterator Iterator) error {
	cdbIterator, err := d.reader.Iterator()

	if err != nil {
		return err
	}

	if cdbIterator == nil {
		return nil
	}

	for {
		record := cdbIterator.Record()
		keyReader, keySize := record.Key()
		key := make([]byte, keySize)
		if _, err := keyReader.Read(key); err != nil {
			return err
		}

		valueReader, valSize := record.Value()
		value := make([]byte, valSize)
		if _, err := valueReader.Read(value); err != nil {
			return err
		}

		iterator(binary.LittleEndian.Uint32(key), Value(value))
		ok, err := cdbIterator.Next()

		if err != nil {
			return err
		}

		if !ok {
			break
		}
	}

	return nil
}

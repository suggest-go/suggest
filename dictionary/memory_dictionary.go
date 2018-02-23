package dictionary

import "errors"

// inMemoryDictionary implements Dictionary with in-memory data access
type inMemoryDictionary struct {
	holder []Value
}

// NewInMemoryDictionary creates new instance of inMemoryDictionary
func NewInMemoryDictionary(words []string) Dictionary {
	holder := make([]Value, len(words))
	copy(holder, words)

	return &inMemoryDictionary{
		holder: holder,
	}
}

// Get returns value associated with a particular key
func (d *inMemoryDictionary) Get(key Key) (Value, error) {
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
	index Key
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
func (i *inMemoryDictionaryIterator) GetPair() (Key, Value) {
	if int(i.index) >= len(i.dict.holder) {
		return 0, ""
	}

	val, _ := i.dict.Get(i.index)
	return i.index, val
}

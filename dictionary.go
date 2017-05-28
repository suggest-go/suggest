package suggest

import "errors"

// WordKey represents key in key-value pair for Dictionary
type WordKey interface{}

// DictIter represents Iterator of Dictionary
type DictIter interface {
	// Next moves iterator to the next item. Returns true on success otherwise false
	Next() bool
	// IsValid indicates if the iterator is deferencable
	IsValid() bool
	// GetPair returns key-value pair of current item
	GetPair() (WordKey, string)
}

// Dictionary is an abstract data type composed of a collection of (key, value) pairs
type Dictionary interface {
	// Get returns value associated with a particular key
	Get(key WordKey) (string, error)
	// MGet returns map of founded pairs
	MGet(keys []WordKey) map[WordKey]string
	// Iter returns an iterator over the elements in this dictionary
	Iter() DictIter
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

func (d *inMemoryDictionary) MGet(keys []WordKey) map[WordKey]string {
	m := make(map[WordKey]string)
	for _, key := range keys {
		val, err := d.Get(key)
		if err != nil {
			m[key] = val
		}
	}

	return m
}

func (d *inMemoryDictionary) Iter() DictIter {
	return &inMemoryDictIter{d, 0}
}

// inMemoryDictIter implements interface DictIter for inMemoryDictionary
type inMemoryDictIter struct {
	dict  *inMemoryDictionary
	index int
}

func (i *inMemoryDictIter) Next() bool {
	success := false
	if i.index < i.dict.index {
		i.index++
		success = true
	}

	return success
}

func (i *inMemoryDictIter) IsValid() bool {
	return i.index < i.dict.index
}

func (i *inMemoryDictIter) GetPair() (WordKey, string) {
	if !i.IsValid() {
		panic("Iterator is not deferencable")
	}

	val, _ := i.dict.Get(i.index)
	return i.index, val
}

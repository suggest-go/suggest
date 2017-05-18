package suggest

import "errors"

//
type WordKey interface{}

//
type Dictionary interface {
	Get(key WordKey) (string, error)
	Next() (WordKey, string)
	Reset()
}

//
type InMemoryDictionary struct {
	holder []string
	index  int
}

//
func NewInMemoryDictionary(words []string) *InMemoryDictionary {
	holder := make([]string, len(words))
	copy(holder, words)

	return &InMemoryDictionary{
		holder,
		0,
	}
}

//
func (d *InMemoryDictionary) Get(key WordKey) (string, error) {
	k := key.(int)
	if k < 0 || k >= len(d.holder) {
		return "", errors.New("Key is not exists")
	}

	return d.holder[k], nil
}

//
func (d *InMemoryDictionary) Next() (WordKey, string) {
	if d.index >= len(d.holder) {
		return nil, ""
	}

	index := d.index
	d.index++
	return index, d.holder[index]
}

//
func (d *InMemoryDictionary) Reset() {
	d.index = 0
}

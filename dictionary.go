package suggest

import "errors"

type WordKey interface{}

type Dictionary interface {
	Get(key WordKey) (string, error)
	Next() (WordKey, string)
	Reset()
}

type InMemoryDictionary struct {
	holder []string
	index  int
}

func NewInMemoryDictionary(words []string) *InMemoryDictionary {
	holder := make([]string, len(words))
	copy(holder, words)

	return &InMemoryDictionary{
		holder,
		0,
	}
}

func (self *InMemoryDictionary) Get(key WordKey) (string, error) {
	k := key.(int)
	if k < 0 || k >= len(self.holder) {
		return "", errors.New("Key is not exists")
	}

	return self.holder[k], nil
}

func (self *InMemoryDictionary) Next() (WordKey, string) {
	if self.index >= len(self.holder) {
		return nil, ""
	}

	index := self.index
	self.index++
	return index, self.holder[index]
}

func (self *InMemoryDictionary) Reset() {
	self.index = 0
}

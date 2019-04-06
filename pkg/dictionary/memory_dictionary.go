package dictionary

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
		return "<nil/>", nil
	}

	return d.holder[key], nil
}

// Iterator returns an iterator over the elements in this dictionary
func (d *inMemoryDictionary) Iterate(iterator Iterator) error {
	for key, value := range d.holder {
		iterator(Key(key), Value(value))
	}

	return nil
}

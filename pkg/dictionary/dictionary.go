package dictionary

type (
	Key   = uint32
	Value = string
)

// Dictionary is an abstract data type composed of a collection of (key, value) pairs
type Dictionary interface {
	// Get returns value associated with a particular key
	Get(key Key) (Value, error)
	// Iterator returns an iterator over the elements in this dictionary
	Iterator() DictionaryIterator
}

// DictionaryIterator represents Iterator of Dictionary
type DictionaryIterator interface {
	// Next moves iterator to the next item. Returns true on success otherwise false
	Next() bool
	// GetPair returns key-value pair of current item
	GetPair() (Key, Value)
}

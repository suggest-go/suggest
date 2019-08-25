package dictionary

const (
	// NilValue is a value, that returns when an entry with the given
	// key doesn't exist
	NilValue = "<nil/>"
)

type (
	// Key represents a key of an item
	Key = uint32
	// Value represents a value of an item
	Value = string
)

// Dictionary is an abstract data type composed of a collection of (key, value) pairs
type Dictionary interface {
	Iterable
	// Get returns value associated with a particular key
	Get(key Key) (Value, error)
	// Size returns the size of the dictionary
	Size() int
}

// Iterable tells that an implementation might be an object of for loop
type Iterable interface {
	// Iterate walks through each kv pair and calls iterator on it
	Iterate(iterator Iterator) error
}

// Iterator is a callback that is called on each pair of the dictionary
// Returned error means that Iterable object should stop handling and start propagating the error upper
type Iterator func(docID Key, word Value) error

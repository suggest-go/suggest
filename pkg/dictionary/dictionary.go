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
	// Get returns value associated with a particular key
	Get(key Key) (Value, error)
	// Iterate walks through each kv pair and calls iterator on it
	Iterate(iterator Iterator) error
}

// Iterator is a callback that is called on each pair of the dictionary
type Iterator func(key Key, value Value)

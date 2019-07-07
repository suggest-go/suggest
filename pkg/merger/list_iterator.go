package merger

import (
	"errors"
	"sort"
)

var (
	// ErrIteratorIsNotDereferencable tells than iterator has an invalid state
	ErrIteratorIsNotDereferencable = errors.New("iterator is not dereferencable")
)

// ListIterator is the interface of a posting list iterator
type ListIterator interface {
	// Get returns the current pointed element of the list
	Get() (uint32, error)
	// IsValid returns true if the given iterator is dereferencable, otherwise returns false
	IsValid() bool
	// HasNext tells if the given iterator can be moved to the next record
	HasNext() bool
	// Next moves the given iterator to the next record
	Next() (uint32, error)
	// LowerBound moves the given iterator to the smallest record x
	// in corresponding list such that x >= to
	LowerBound(to uint32) (uint32, error)
	// Len returns the actual size of the list
	Len() int
}

// sliceIterator represents the ListIterator interface for a slice of uint32
type sliceIterator struct {
	slice []uint32
	index int
}

// NewSliceIterator returns a new instance of a slice iterator
func NewSliceIterator(slice []uint32) ListIterator {
	return &sliceIterator{
		slice: slice,
		index: 0,
	}
}

// Get returns the current pointed element of the list
func (i *sliceIterator) Get() (uint32, error) {
	if !i.IsValid() {
		return 0, ErrIteratorIsNotDereferencable
	}

	return i.slice[i.index], nil
}

// IsValid returns true if the given iterator is dereferencable, otherwise returns false
func (i* sliceIterator) IsValid() bool {
	return i.index < len(i.slice)
}

// HasNext tells if the given iterator can be moved to the next record
func (i *sliceIterator) HasNext() bool {
	return i.index+1 < len(i.slice)
}

// Next moves the given iterator to the next record
func (i *sliceIterator) Next() (uint32, error) {
	if !i.HasNext() {
		return 0, ErrIteratorIsNotDereferencable
	}

	i.index++

	return i.slice[i.index], nil
}

// LowerBound moves the given iterator to the smallest record x
// in corresponding list such that x >= to
func (i *sliceIterator) LowerBound(to uint32) (uint32, error) {
	if !i.IsValid() {
		return 0, ErrIteratorIsNotDereferencable
	}

	slice := i.slice[i.index:]
	j := sort.Search(len(slice), func(i int) bool { return slice[i] >= to })

	if j < 0 || j >= len(slice) {
		i.index = len(i.slice)
		return 0, ErrIteratorIsNotDereferencable
	}

	i.index += j

	return slice[j], nil
}

// Len returns the actual size of the list
func (i *sliceIterator) Len() int {
	return len(i.slice)
}

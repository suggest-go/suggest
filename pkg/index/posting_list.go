package index

import (
	"github.com/alldroll/suggest/pkg/merger"
	"github.com/alldroll/suggest/pkg/store"
)

// PostingList represents a list of documents ids, that belongs to the certain index term
type PostingList interface {
	merger.ListIterator
	// Init initialize the given posting list with the provided context
	Init(context PostingListContext) error
}

// postingList is a sequential implementation of PostingList interface
type postingList struct {
	input   store.Input
	index   int
	size    int
	current uint32
}

// Get returns the current pointed element of the list
func (i *postingList) Get() (uint32, error) {
	if !i.isValid() {
		return 0, merger.ErrIteratorIsNotDereferencable
	}

	return i.current, nil
}

// HasNext tells if the given iterator can be moved to the next record
func (i *postingList) HasNext() bool {
	return i.index+1 < i.size
}

// Next moves the given iterator to the next record
func (i *postingList) Next() (uint32, error) {
	if !i.HasNext() {
		return 0, merger.ErrIteratorIsNotDereferencable
	}

	cur, err := i.input.ReadVUInt32()

	if err != nil {
		return 0, err
	}

	i.index++
	i.current += cur

	return i.current, nil
}

// LowerBound moves the given iterator to the smallest record x
// in corresponding list such that x >= to
func (i *postingList) LowerBound(to uint32) (uint32, error) {
	if !i.isValid() {
		return 0, merger.ErrIteratorIsNotDereferencable
	}

	if i.current >= to {
		return i.current, nil
	}

	for i.HasNext() {
		cur, err := i.Next()

		if err != nil {
			return 0, err
		}

		if cur >= to {
			return cur, nil
		}
	}

	i.index = i.size

	return 0, merger.ErrIteratorIsNotDereferencable
}

// Len returns the actual size of the list
func (i *postingList) Len() int {
	return i.size
}

// isValid returns true if the given iterator is dereferencable, otherwise returns false
func (i *postingList) isValid() bool {
	return i.index < i.size
}

// Init initialize the iterator by the given PostingList context
func (i *postingList) Init(context PostingListContext) error {
	i.input = context.Reader
	i.size = context.ListSize
	i.index = 0

	current, err := i.input.ReadVUInt32()

	if err != nil {
		return err
	}

	i.current = current

	return nil
}

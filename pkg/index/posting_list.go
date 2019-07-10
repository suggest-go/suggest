package index

import (
	"github.com/alldroll/suggest/pkg/merger"
	"github.com/alldroll/suggest/pkg/store"
)

// postingListIterator is a dummy implementation of merger.ListIterator
type postingListIterator struct {
	input   store.Input
	index   int
	size    int
	current uint32
}

// Get returns the current pointed element of the list
func (i *postingListIterator) Get() (uint32, error) {
	if !i.isValid() {
		return 0, merger.ErrIteratorIsNotDereferencable
	}

	return i.current, nil
}

// HasNext tells if the given iterator can be moved to the next record
func (i *postingListIterator) HasNext() bool {
	return i.index+1 < i.size
}

// Next moves the given iterator to the next record
func (i *postingListIterator) Next() (uint32, error) {
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
func (i *postingListIterator) LowerBound(to uint32) (uint32, error) {
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
func (i *postingListIterator) Len() int {
	return i.size
}

// isValid returns true if the given iterator is dereferencable, otherwise returns false
func (i *postingListIterator) isValid() bool {
	return i.index < i.size
}

// init initialize the iterator by the given PostingList context
func (i *postingListIterator) init(context PostingListContext) error {
	i.input = context.GetReader()
	i.size = context.GetListSize()
	i.index = 0

	current, err := i.input.ReadVUInt32()

	if err != nil {
		return err
	}

	i.current = current

	return nil
}

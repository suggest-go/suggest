package index

import (
	"github.com/alldroll/suggest/pkg/compression"
	"github.com/alldroll/suggest/pkg/merger"
	"github.com/alldroll/suggest/pkg/store"
)

// postingListIterator is a dummy implementation of merger.ListIterator
type postingListIterator struct {
	sliceIterator *merger.SliceIterator
	input         store.Input
	decoder       compression.Decoder
	index         int
	size          int
	current       uint32
}

// Get returns the current pointed element of the list
func (i *postingListIterator) Get() (uint32, error) {
	if !i.isValid() {
		return 0, merger.ErrIteratorIsNotDereferencable
	}

	if i.sliceIterator.Len() > 0 {
		cur, err := i.sliceIterator.Get()

		if err != nil {
			return 0, err
		}

		return cur + i.current, nil
	}

	return i.current, nil
}

// HasNext tells if the given iterator can be moved to the next record
func (i *postingListIterator) HasNext() bool {
	if i.sliceIterator.Len() > 0 {
		return i.sliceIterator.HasNext()
	}

	return i.index+1 < i.size
}

// Next moves the given iterator to the next record
func (i *postingListIterator) Next() (uint32, error) {
	if !i.HasNext() {
		return 0, merger.ErrIteratorIsNotDereferencable
	}

	if i.sliceIterator.Len() > 0 {
		cur, err := i.sliceIterator.Next()

		if err != nil {
			return 0, err
		}

		return cur + i.current, nil
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

	if !i.HasNext() {
		i.index = i.size
		return 0, merger.ErrIteratorIsNotDereferencable
	}

	if i.sliceIterator.Len() == 0 {
		buf := make([]Position, i.size-i.index-1)
		i.decoder.Decode(i.input, buf)
		i.sliceIterator.Reset(buf)
	}

	cur, err := i.sliceIterator.LowerBound(to - i.current)

	if err != nil {
		return 0, err
	}

	return cur + i.current, nil
}

// Len returns the actual size of the list
func (i *postingListIterator) Len() int {
	return i.size
}

// isValid returns true if the given iterator is dereferencable, otherwise returns false
func (i *postingListIterator) isValid() bool {
	if i.sliceIterator.Len() > 0 {
		return i.sliceIterator.IsValid()
	}

	return i.index < i.size
}

// init initialize the iterator by the given PostingList context
func (i *postingListIterator) init(context PostingListContext) error {
	i.input = context.GetReader()
	i.decoder = context.GetDecoder()
	i.size = context.GetListSize()
	i.index = 0

	current, err := i.input.ReadVUInt32()

	if err != nil {
		return err
	}

	i.current = current
	i.sliceIterator.Reset([]Position{})

	return nil
}

package index

import (
	"errors"
	"io"

	"github.com/alldroll/suggest/pkg/compression"
	"github.com/alldroll/suggest/pkg/merger"
	"github.com/alldroll/suggest/pkg/store"
)

// skippingPostingList TODO describe me
type skippingPostingList struct {
	input               store.Input
	index               int
	size                int
	prev                uint32
	current             uint32
	currentSkipValue    uint32
	currentSkipPosition int
	nextSkipPosition    int
	skippingGap         int
	isLastBlock         bool
}

// Get returns the current pointed element of the list
func (i *skippingPostingList) Get() (uint32, error) {
	if !i.isValid() {
		return 0, merger.ErrIteratorIsNotDereferencable
	}

	return i.current, nil
}

// HasNext tells if the given iterator can be moved to the next record
func (i *skippingPostingList) HasNext() bool {
	return i.index+1 < i.size
}

// Next moves the given iterator to the next record
func (i *skippingPostingList) Next() (uint32, error) {
	if !i.HasNext() {
		return 0, merger.ErrIteratorIsNotDereferencable
	}

	if (i.index+1)%i.skippingGap == 0 {
		if err := i.readSkipping(); err != nil {
			return 0, err
		}
	} else {
		cur, err := i.input.ReadVUInt32()

		if err != nil {
			return 0, err
		}

		i.current = cur + i.prev
		i.prev = i.current
	}

	i.index++

	return i.current, nil
}

// LowerBound moves the given iterator to the smallest record x
// in corresponding list such that x >= to
func (i *skippingPostingList) LowerBound(to uint32) (uint32, error) {
	if !i.isValid() {
		return 0, merger.ErrIteratorIsNotDereferencable
	}

	if i.current >= to {
		return i.current, nil
	}

	for i.HasNext() {
		prev := *i

		if err := i.moveToPosition(i.nextSkipPosition); err != nil {
			return 0, err
		}

		i.index += i.skippingGap - 1

		if i.index >= i.size {
			i.index = i.size - 2
		}

		cur, err := i.Next()

		if err != nil {
			return 0, err
		}

		if cur < to && !i.isLastBlock {
			continue
		}

		// rollback to previus block
		if cur >= to {
			if err := i.moveToPosition(prev.currentSkipPosition); err != nil {
				return 0, err
			}

			if err != nil {
				return 0, err
			}

			*i = prev
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
	}

	i.index = i.size

	return 0, merger.ErrIteratorIsNotDereferencable
}

// Len returns the actual size of the list
func (i *skippingPostingList) Len() int {
	return i.size
}

// isValid returns true if the given iterator is dereferencable, otherwise returns false
func (i *skippingPostingList) isValid() bool {
	return i.index < i.size
}

// moveToPosition moves the given input cursor to the given offset
func (i *skippingPostingList) moveToPosition(position int) error {
	offset := int64(position)
	n, err := i.input.Seek(offset, io.SeekStart)

	if err != nil {
		return err
	}

	if n != offset {
		return errors.New("failed to move to the given position")
	}

	return nil
}

// init initialize the iterator by the given PostingList context
func (i *skippingPostingList) init(context PostingListContext) error {
	i.input = context.GetReader()
	i.size = context.GetListSize()
	i.index = 0
	i.currentSkipValue = 0
	i.nextSkipPosition = 0

	return i.readSkipping()
}

// readSkipping reads a skip pointer and a value
func (i *skippingPostingList) readSkipping() error {
	decodedPosition, err := i.input.ReadUInt16()

	if err != nil {
		return err
	}

	position, isLastBlock := compression.UnpackPos(decodedPosition)
	current, err := i.input.ReadVUInt32()

	if err != nil {
		return err
	}

	currentSkipPosition, err := i.input.Seek(0, io.SeekCurrent)

	if err != nil {
		return err
	}

	i.current = i.currentSkipValue + current
	i.currentSkipPosition = int(currentSkipPosition)
	i.currentSkipValue = current
	i.prev = current
	i.nextSkipPosition += int(position)
	i.isLastBlock = isLastBlock

	return nil
}

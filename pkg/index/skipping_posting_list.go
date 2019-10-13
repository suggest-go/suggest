package index

import (
	"io"

	"github.com/suggest-go/suggest/pkg/compression"
	"github.com/suggest-go/suggest/pkg/merger"
	"github.com/suggest-go/suggest/pkg/store"
)

// skippingPostingList is a posting list which has the ability to use skip pointers for faster intersection
// https://nlp.stanford.edu/IR-book/html/htmledition/faster-postings-list-intersection-via-skip-pointers-1.html
type skippingPostingList struct {
	input               store.Input
	index               int
	size                int
	current             uint32
	currentSkipValue    uint32
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

	currentSkipPosition, err := i.input.Seek(0, io.SeekCurrent)

	if err != nil {
		return 0, err
	}

	if int(currentSkipPosition) == i.nextSkipPosition {
		if err := i.readSkipping(); err != nil {
			return 0, err
		}
	} else {
		cur, err := i.input.ReadVUInt32()

		if err != nil {
			return 0, err
		}

		i.current += cur
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

	// calculate how many skips have already been done
	skips := 0

	if i.index > 0 {
		skips = i.index / i.skippingGap
	}

	// try to use skip pointers to find the corresponding block
	for !i.isLastBlock && i.HasNext() {
		// remember the current state, maybe we will have to restore the state
		prev := *i
		prevPosition, err := i.input.Seek(0, io.SeekCurrent)

		if err != nil {
			return 0, err
		}

		if err := i.moveToPosition(i.nextSkipPosition); err != nil {
			return 0, err
		}

		skips++
		i.index = (skips * i.skippingGap) - 1

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

		// rollback to the previous block
		if cur >= to {
			if err := i.moveToPosition(int(prevPosition)); err != nil {
				return 0, err
			}

			*i = prev
			break
		}
	}

	// here we just should iterate sequentially through the list
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
	_, err := i.input.Seek(offset, io.SeekStart)

	if err != nil {
		return err
	}

	return nil
}

// Init initialize the iterator by the given PostingList context
func (i *skippingPostingList) Init(context PostingListContext) error {
	i.input = context.Reader
	i.size = context.ListSize
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

	i.current = i.currentSkipValue + current
	i.currentSkipValue = i.current
	i.nextSkipPosition += int(position)
	i.isLastBlock = isLastBlock

	return nil
}

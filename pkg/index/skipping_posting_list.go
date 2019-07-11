package index

import (
	"github.com/alldroll/suggest/pkg/merger"
	"github.com/alldroll/suggest/pkg/store"
	"io"
	"log"
)

const skippingGap = 4

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

	if (i.index+1)%skippingGap == 0 || (i.index+1) == i.size {
		log.Printf("Index1: %v\n", i.index)
		nextPosition, err := i.input.ReadUInt16()

		if err != nil {
			log.Printf("Err1: %v %v\n", nextPosition, err)
			return 0, err
		}

		current, err := i.input.ReadVUInt32()

		if err != nil {
			log.Printf("Err2: %v %v\n", current, err)
			return 0, err
		}

		log.Printf("Next1: %v, %v", nextPosition, current)

		currentSkipPosition, err := i.input.Seek(0, io.SeekCurrent)

		if err != nil {
			log.Printf("Err3: %v %v\n", current, err)
			return 0, err
		}

		i.current = i.currentSkipValue + current
		i.currentSkipPosition = int(currentSkipPosition)
		i.currentSkipValue = current
		i.prev = current
		i.nextSkipPosition += int(nextPosition)

		log.Printf("Next2: %v, %v", i.currentSkipPosition, i.current)
	} else {
		log.Printf("Index2: %v\n", i.index)
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

		log.Printf("STEP1: %v\n", i)

		_, err := i.input.Seek(int64(i.nextSkipPosition), io.SeekStart)

		if err != nil {
			return 0, err
		}

		i.index += skippingGap - 1

		if i.index >= i.size {
			i.index = i.size - 2
		}

		cur, err := i.Next()

		if err != nil {
			return 0, err
		}

		log.Printf("STEP1: %v\n", i)

		if cur >= to {
			log.Printf("prev: %v\n", prev)

			_, err := i.input.Seek(int64(prev.currentSkipPosition), io.SeekStart)
			i = &prev

			log.Printf("STEP2: %#v\n", i)

			if err != nil {
				return 0, err
			}

			for i.HasNext() {
				cur, err := i.Next()

				if err != nil {
					return 0, err
				}

				log.Printf("CURR: %v\n", cur)

				if cur >= to {
					return cur, nil
				}
			}

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

// init initialize the iterator by the given PostingList context
func (i *skippingPostingList) init(context PostingListContext) error {
	i.input = context.GetReader()
	i.size = context.GetListSize()

	nextSkipPosition, err := i.input.ReadUInt16()

	if err != nil {
		return nil
	}

	current, err := i.input.ReadVUInt32()

	if err != nil {
		return err
	}

	currentSkipPosition, err := i.input.Seek(0, io.SeekCurrent)

	if err != nil {
		log.Printf("Err3: %v %v\n", current, err)
		return err
	}

	i.index = 0
	i.current = current
	i.currentSkipPosition = int(currentSkipPosition)
	i.currentSkipValue = current
	i.prev = current
	i.nextSkipPosition = int(nextSkipPosition)

	return nil
}

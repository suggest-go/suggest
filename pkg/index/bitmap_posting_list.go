package index

import (
	"fmt"

	"github.com/RoaringBitmap/roaring"
	"github.com/alldroll/suggest/pkg/merger"
	"github.com/alldroll/suggest/pkg/store"
)

// bitmapPostingList is a bitmap PostingList implementation
type bitmapPostingList struct {
	iterator roaring.IntPeekable
	bitmap   *roaring.Bitmap
	current  uint32
	isValid  bool
	length   int
}

// Get returns the current pointed element of the list
func (i *bitmapPostingList) Get() (uint32, error) {
	if !i.isValid {
		return 0, merger.ErrIteratorIsNotDereferencable
	}

	return i.current, nil
}

// HasNext tells if the given iterator can be moved to the next record
func (i *bitmapPostingList) HasNext() bool {
	return i.isValid && i.iterator.HasNext()
}

// Next moves the given iterator to the next record
func (i *bitmapPostingList) Next() (uint32, error) {
	if !i.HasNext() {
		return 0, merger.ErrIteratorIsNotDereferencable
	}

	i.current = i.iterator.Next()

	return i.current, nil
}

// LowerBound moves the given iterator to the smallest record x
// in corresponding list such that x >= to
func (i *bitmapPostingList) LowerBound(to uint32) (uint32, error) {
	if !i.isValid {
		return 0, merger.ErrIteratorIsNotDereferencable
	}

	if i.current >= to {
		return i.current, nil
	}

	i.iterator.AdvanceIfNeeded(to)

	if !i.iterator.HasNext() {
		i.isValid = false
		return 0, merger.ErrIteratorIsNotDereferencable
	}

	i.current = i.iterator.Next()

	return i.current, nil
}

// Len returns the actual size of the list
func (i *bitmapPostingList) Len() int {
	return i.length
}

// Init initialize the iterator by the given PostingList context
func (i *bitmapPostingList) Init(context PostingListContext) error {
	var err error

	if i.bitmap == nil {
		i.bitmap = roaring.NewBitmap()
	}

	reader := context.Reader

	if buf, ok := reader.(store.SliceAccessible); ok {
		_, err = i.bitmap.FromBuffer(buf.Data())
	} else {
		_, err = i.bitmap.ReadFrom(reader)
	}

	if err != nil {
		return fmt.Errorf("failed to create bitmap: %v", err)
	}

	iterator := i.bitmap.Iterator()

	if !iterator.HasNext() {
		return fmt.Errorf("bitmap should not be empty")
	}

	i.length = int(i.bitmap.GetCardinality())
	i.iterator = iterator
	i.isValid = true
	i.current = i.iterator.Next()

	return err
}

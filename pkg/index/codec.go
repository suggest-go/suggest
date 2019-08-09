package index

import (
	"errors"
	"github.com/alldroll/suggest/pkg/compression"
	"github.com/alldroll/suggest/pkg/store"
	"sync"
)

const skippingGapSize = 64
const maxSkippingLen = 1024

var errUnknownPostingListImplementation = errors.New("unknown posting list implementation")

// NewEncoder returns a new instance of Encoder
func NewEncoder() (compression.Encoder, error) {
	skippingEnc, err := compression.SkippingEncoder(skippingGapSize)

	if err != nil {
		return nil, err
	}

	return &encoder{
		vbEnc:       compression.VBEncoder(),
		skippingEnc: skippingEnc,
		bitmapEnc: compression.BitmapEncoder(),
	}, nil
}

type encoder struct {
	vbEnc       compression.Encoder
	skippingEnc compression.Encoder
	bitmapEnc compression.Encoder
}

// Encode encodes the given positing list into the buf array
// Returns number of elements encoded, number of bytes readed
func (e *encoder) Encode(list []uint32, out store.Output) (int, error) {
	n := len(list)

	if n <= (skippingGapSize + 1) {
		return e.vbEnc.Encode(list, out)
	}

	if n <= maxSkippingLen {
		return e.skippingEnc.Encode(list, out)
	}

	return e.bitmapEnc.Encode(list, out)
}

var (
	vbEncPostingListPool = sync.Pool{
		New: func() interface{} {
			return &postingListIterator{}
		},
	}

	skippingPostingListPool = sync.Pool{
		New: func() interface{} {
			return &skippingPostingList{
				skippingGap: skippingGapSize,
			}
		},
	}

	bitmapPostingListPool = sync.Pool{
		New: func() interface{} {
			return &bitmapPostingList{}
		},
	}
)

// resolvePostingList returns the appropriate posting list for the provided context
func resolvePostingList(context PostingListContext) postingList {
	n := context.GetListSize()

	if n <= (skippingGapSize + 1) {
		return vbEncPostingListPool.Get().(postingList)
	}

	if n <= maxSkippingLen {
		return skippingPostingListPool.Get().(postingList)
	}

	return bitmapPostingListPool.Get().(postingList)
}

// releasePostingList puts the given postingList to the corresponding pool
func releasePostingList(list postingList) (err error) {
	switch v := list.(type) {
	case *postingListIterator:
		vbEncPostingListPool.Put(v)
	case *skippingPostingList:
		skippingPostingListPool.Put(v)
	case *bitmapPostingList:
		bitmapPostingListPool.Put(v)
	default:
		err = errUnknownPostingListImplementation
	}

	return
}

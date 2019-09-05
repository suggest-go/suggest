package index

import (
	"errors"
	"sync"

	"github.com/alldroll/suggest/pkg/compression"
	"github.com/alldroll/suggest/pkg/store"
)

const skippingGapSize = 64
const maxSkippingLen = 256

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
		bitmapEnc:   compression.BitmapEncoder(),
	}, nil
}

type encoder struct {
	vbEnc       compression.Encoder
	skippingEnc compression.Encoder
	bitmapEnc   compression.Encoder
}

// Encode encodes the given positing list into the buf array
// Returns number of elements encoded, number of bytes read
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
			return &postingList{}
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
func resolvePostingList(context PostingListContext) PostingList {
	n := context.ListSize

	if n <= (skippingGapSize + 1) {
		return vbEncPostingListPool.Get().(PostingList)
	}

	if n <= maxSkippingLen {
		return skippingPostingListPool.Get().(PostingList)
	}

	return bitmapPostingListPool.Get().(PostingList)
}

// releasePostingList puts the given postingList to the corresponding pool
func releasePostingList(list PostingList) (err error) {
	switch v := list.(type) {
	case *postingList:
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

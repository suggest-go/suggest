package index

import (
	"errors"
	"github.com/alldroll/suggest/pkg/compression"
	"github.com/alldroll/suggest/pkg/store"
	"sync"
)

const skippingGapSize = 64

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
	}, nil
}

type encoder struct {
	vbEnc       compression.Encoder
	skippingEnc compression.Encoder
}

// Encode encodes the given positing list into the buf array
// Returns number of elements encoded, number of bytes readed
func (e *encoder) Encode(list []uint32, out store.Output) (int, error) {
	if len(list) <= (skippingGapSize + 1) {
		return e.vbEnc.Encode(list, out)
	}

	return e.skippingEnc.Encode(list, out)
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
)

// resolvePostingList returns the appropriate posting list for the provided context
func resolvePostingList(context PostingListContext) postingList {
	if context.GetListSize() <= (skippingGapSize + 1) {
		return vbEncPostingListPool.Get().(postingList)
	}

	return skippingPostingListPool.Get().(postingList)
}

// releasePostingList puts the given postingList to the corresponding pool
func releasePostingList(list postingList) (err error) {
	switch v := list.(type) {
	case *postingListIterator:
		vbEncPostingListPool.Put(v)
	case *skippingPostingList:
		skippingPostingListPool.Put(v)
	default:
		err = errUnknownPostingListImplementation
	}

	return
}

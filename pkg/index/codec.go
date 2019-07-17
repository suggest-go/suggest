package index

import (
	"github.com/alldroll/suggest/pkg/compression"
	"github.com/alldroll/suggest/pkg/store"
	"sync"
)

const (
	skippingGapSize = 64
)

// NewEncoder returns a new instance of Encoder
func NewEncoder() compression.Encoder {
	// TODO handle error
	skippingEnc, _ := compression.SkippingEncoder(skippingGapSize)

	return &encoder{
		vbEnc:       compression.VBEncoder(),
		skippingEnc: skippingEnc,
	}
}

type encoder struct {
	vbEnc       compression.Encoder
	skippingEnc compression.Encoder
}

// Encode encodes the given positing list into the buf array
// Returns number of elements encoded, number of bytes readed
func (e *encoder) Encode(list []uint32, out store.Output) (int, error) {
	if len(list) <= skippingGapSize {
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

// resolvePostingList
func resolvePostingList(context PostingListContext) PostingList {

	if context.GetListSize() <= skippingGapSize {
		return vbEncPostingListPool.Get().(PostingList)
	}

	return skippingPostingListPool.Get().(PostingList)
}

// releasePostingList
func releasePostingList(list PostingList) {
	// TODO handle error

	switch v := list.(type) {
	case *postingListIterator:
		vbEncPostingListPool.Put(v)
	case *skippingPostingList:
		skippingPostingListPool.Put(v)
	}
}

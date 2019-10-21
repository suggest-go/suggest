package compression

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/suggest-go/suggest/pkg/store"
)

// BitmapEncoder returns new instance of bitmapEnc which compress the uint32 list with the roaring bitmap library
func BitmapEncoder() Encoder {
	return &bitmapEnc{}
}

// bitmapEnc implements BitmapEncoder
type bitmapEnc struct{}

// Encode encodes the given positing list into the buf array
// Returns a number of written bytes
func (b *bitmapEnc) Encode(list []uint32, out store.Output) (int, error) {
	bitmap := roaring.New()

	for _, i := range list {
		bitmap.Add(i)
	}

	n, err := bitmap.WriteTo(out)

	return int(n), err
}

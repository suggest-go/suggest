package compression

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/suggest-go/suggest/pkg/store"
)

func BitmapEncoder() Encoder {
	return &bitmapEnc{}
}

// vbEnc implements BitmapEncoder
type bitmapEnc struct{}

// Encode encodes the given positing list into the buf array
// Returns number of elements encoded, number of bytes readed
func (b *bitmapEnc) Encode(list []uint32, out store.Output) (int, error) {
	bitmap := roaring.New()

	for _, i := range list {
		bitmap.Add(i)
	}

	n, err := bitmap.WriteTo(out)

	return int(n), err
}

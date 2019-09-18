// Package compression holds the algorithms for compressing the list of sorted lists of integers
package compression

import "github.com/alldroll/suggest/pkg/store"

// Encoder represents entity for encoding given posting list to byte array
type Encoder interface {
	// Encode encodes the given positing list into the buf array
	// Returns a number of written bytes
	Encode(list []uint32, out store.Output) (int, error)
}

// Decoder represents entity for decoding given byte array to posting list
type Decoder interface {
	// Decode decodes the given byte array to the buf list
	// Returns a number of elements encoded
	Decode(in store.Input, buf []uint32) (int, error)
}

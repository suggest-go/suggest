package compression

import "io"

// Encoder represents entity for encoding given posting list to byte array
type Encoder interface {
	// Encode encodes the given positing list into the buf array
	// Returns a number of written bytes
	Encode(list []uint32, buf io.Writer) (int, error)
}

// Decoder represents entity for decoding given byte array to posting list
type Decoder interface {
	// Decode decodes the given byte array to the buf list
	// Returns a number of elements encoded
	Decode(in io.Reader, buf []uint32) (int, error)
}

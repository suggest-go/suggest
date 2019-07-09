package compression

import (
	"encoding/binary"
	"github.com/alldroll/suggest/pkg/store"
	"io"
)

// binaryEnc encode each position in 4 bytes, decode 4 byte to 1 position
type binaryEnc struct{}

// BinaryEncoder returns new instance of binaryEnc which encodes each Position in 4 bytes
func BinaryEncoder() Encoder {
	return &binaryEnc{}
}

// BinaryDecoder returns new instance of binaryEnc
func BinaryDecoder() Decoder {
	return &binaryEnc{}
}

// Encode encodes the given positing list into the buf array
// Returns number of elements encoded, number of bytes readed
func (b *binaryEnc) Encode(list []uint32, out store.Output) (int, error) {
	chunk := make([]byte, 4)
	total := 0

	for _, v := range list {
		binary.LittleEndian.PutUint32(chunk, v)
		n, err := out.Write(chunk)
		total += n

		if err != nil {
			return total, err
		}
	}

	return total, nil
}

// Decode decodes the given byte array to the buf list
// Returns a number of elements encoded
func (b *binaryEnc) Decode(in store.Input, buf []uint32) (n int, err error) {
	for ; n < len(buf); n++ {
		v, err := in.ReadUInt32()

		if err != nil {
			if err == io.EOF {
				err = nil
			}

			return n, err
		}

		buf[n] = v
	}

	return
}

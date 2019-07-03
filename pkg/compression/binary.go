package compression

import (
	"encoding/binary"
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
func (b *binaryEnc) Encode(list []uint32, buf io.Writer) (int, error) {
	chunk := make([]byte, 4)
	total := 0

	for _, v := range list {
		binary.LittleEndian.PutUint32(chunk, v)
		n, err := buf.Write(chunk)
		total += n

		if err != nil {
			return total, err
		}
	}

	return total, nil
}

// Decode decodes the given byte array to the buf list
// Returns a number of elements encoded
func (b *binaryEnc) Decode(in io.Reader, buf []uint32) (n int, err error) {
	total := 0
	chunk := make([]byte, 4)

	for i := 0; i < len(buf); i++ {
		n, err := in.Read(chunk)
		total += n

		if err != nil {
			return total, err
		}

		buf[i] = binary.LittleEndian.Uint32(chunk)
	}

	return total, nil
}

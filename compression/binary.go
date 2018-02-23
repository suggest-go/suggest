package compression

import "encoding/binary"

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

// Encode encodes given positing list in byte array
func (b *binaryEnc) Encode(list []uint32) []byte {
	bytes := make([]byte, len(list)*4)

	for i, x := range list {
		binary.LittleEndian.PutUint32(bytes[4*i:], uint32(x))
	}

	return bytes
}

// Decode decodes given byte array to posting list
func (b *binaryEnc) Decode(bytes []byte) []uint32 {
	if len(bytes) < 4 {
		return nil
	}

	list := make([]uint32, len(bytes)/4)

	for i := range list {
		list[i] = binary.LittleEndian.Uint32(bytes[4*i:])
	}

	return list
}

package suggest

import (
	"encoding/binary"
)

// Encoder represents entity for encoding given posting list to byte array
type Encoder interface {
	// Encode encodes given positing list in byte array
	Encode(list PostingList) []byte
}

// Decoder represents entity for decoding given byte array to posting list
type Decoder interface {
	// Decode decodes given byte array to posting list
	Decode(bytes []byte) PostingList
}

// BinaryEncoder returns new instance of binaryEnc which encodes each Position in 4 bytes
func BinaryEncoder() Encoder {
	return &binaryEnc{}
}

// BinaryDecoder returns new instance of binaryEnc
func BinaryDecoder() Decoder {
	return &binaryEnc{}
}

// binaryEnc encode each position in 4 bytes, decode 4 byte to 1 position
type binaryEnc struct{}

// Encode encodes given positing list in byte array
func (b *binaryEnc) Encode(list PostingList) []byte {
	bytes := make([]byte, len(list)*4)

	for i, x := range list {
		binary.LittleEndian.PutUint32(bytes[4*i:], uint32(x))
	}

	return bytes
}

// Decode decodes given byte array to posting list
func (b *binaryEnc) Decode(bytes []byte) PostingList {
	if len(bytes) < 4 {
		return nil
	}

	list := make(PostingList, len(bytes)/4)

	for i := range list {
		list[i] = Position(binary.LittleEndian.Uint32(bytes[4*i:]))
	}

	return list
}

// VBEncoder returns new instance of vbEnc that encodes posting list using
// delta encoding
// variable length byte string compression
func VBEncoder() Encoder {
	return &vbEnc{}
}

// VBDecoder decodes given bytes array to posting list which was encoded by VBEncoder
func VBDecoder() Decoder {
	return &vbEnc{}
}

// vbEnc implements VBEncoder and VBDecoder
type vbEnc struct{}

func (b *vbEnc) Encode(list PostingList) []byte {
	sum, l, prev, delta := 4, 4, Position(0), Position(0)

	for _, v := range list {
		sum += estimateByteNum(v - prev)
		prev = v
	}

	prev = 0
	encoded := make([]byte, sum)
	binary.LittleEndian.PutUint32(encoded, uint32(len(list)))

	for _, v := range list {
		delta = v - prev
		prev = v

		for ; delta > 0x7F; l++ {
			encoded[l] = 0x80 | uint8(delta&0x7F)
			delta >>= 7
		}

		encoded[l] = uint8(delta)
		l++
	}

	return encoded
}

// inspired by protobuf/master/proto/decode.go
// Decode decodes given byte array to posting list
func (b *vbEnc) Decode(bytes []byte) PostingList {
	if len(bytes) < 4 {
		return nil
	}

	var (
		v    = uint32(0)
		prev = uint32(0)
		s    = uint32(0)
		i    = 0
		j    = 4
	)

	listLen := int(binary.LittleEndian.Uint32(bytes))
	decoded := make(PostingList, listLen)

	if listLen < 10 {
		b.vbDecodeSlow(bytes[4:], decoded)
		return decoded
	}

	for j < len(bytes) {
		if bytes[j] < 0x80 {
			v = uint32(bytes[j])
			j++
			goto done
		}

		// we already checked the first byte
		v = uint32(bytes[j]) - 0x80
		j++

		s = uint32(bytes[j])
		j++
		v += s << 7
		if s&0x80 == 0 {
			goto done
		}
		v -= 0x80 << 7

		s = uint32(bytes[j])
		j++
		v += s << 14
		if s&0x80 == 0 {
			goto done
		}
		v -= 0x80 << 14

		s = uint32(bytes[j])
		j++
		v += s << 21
		if s&0x80 == 0 {
			goto done
		}
		v -= 0x80 << 21

		s = uint32(bytes[j])
		j++
		v += s << 28

	done:
		prev = v + prev
		decoded[i] = Position(prev)
		i++
	}

	return decoded
}

// vbDecodeSlow decodes given byte array to posting list
func (b *vbEnc) vbDecodeSlow(bytes []byte, buf PostingList) {
	var (
		v    = uint32(0)
		prev = uint32(0)
		s    = uint32(0)
		i    = 0
	)

	for _, b := range bytes {
		v |= uint32(b&0x7f) << s

		if b < 0x80 {
			prev = v + prev
			buf[i] = Position(prev)
			s, v = 0, 0
			i++
		} else {
			s += 7
		}
	}
}

// estimateByteNum returns bytes num required for encoding given uint32 (Position)
func estimateByteNum(v Position) int {
	if (1 << 7) > v {
		return 1
	}

	if (1 << 14) > v {
		return 2
	}

	if (1 << 21) > v {
		return 3
	}

	if (1 << 28) > v {
		return 4
	}

	return 5
}

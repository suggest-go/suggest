package suggest

import (
	"encoding/binary"
)

type Encoder interface {
	Encode(list PostingList) []byte
}

type Decoder interface {
	Decode(bytes []byte) PostingList
}

func BinaryEncoder() Encoder {
	return &binaryEnc{}
}

func BinaryDecoder() Decoder {
	return &binaryEnc{}
}

type binaryEnc struct{}

func (b *binaryEnc) Encode(list PostingList) []byte {
	bytes := make([]byte, len(list)*4)

	for i, x := range list {
		binary.LittleEndian.PutUint32(bytes[4*i:], uint32(x))
	}

	return bytes
}

func (b *binaryEnc) Decode(bytes []byte) PostingList {
	list := make(PostingList, len(bytes)/4)

	for i := range list {
		list[i] = Position(binary.LittleEndian.Uint32(bytes[4*i:]))
	}

	return list
}

func VBEncoder() Encoder {
	return &vbEnc{}
}

func VBDecoder() Decoder {
	return &vbEnc{}
}

type vbEnc struct{}

func (b *vbEnc) Encode(list PostingList) []byte {
	sum, l, prev, delta := 0, 0, Position(0), Position(0)

	for _, v := range list {
		sum += estimateByteNum(v - prev)
		prev = v
	}

	prev = 0
	encoded := make([]byte, sum)

	for _, v := range list {
		delta = v - prev
		prev = v

		for delta >= 0x80 {
			encoded[l] = byte(delta) | 0x80
			delta >>= 7
			l++
		}

		encoded[l] = byte(delta)
		l++
	}

	return encoded
}

func (b *vbEnc) Decode(bytes []byte) PostingList {
	v := uint32(0)
	s := uint(0)
	prev := uint32(0)

	decoded := make(PostingList, 0, len(bytes)) //bad (we should to store posting list size

	for _, b := range bytes {
		v |= uint32(b&0x7f) << s

		if b < 0x80 {
			prev = v + prev
			decoded = append(decoded, Position(prev))
			s, v = 0, 0
		} else {
			s += 7
		}
	}

	return decoded
}

func estimateByteNum(v Position) int {
	num := 5

	if (1 << 7) > v {
		num = 1
	} else if (1 << 14) > v {
		num = 2
	} else if (1 << 21) > v {
		num = 3
	} else if (1 << 28) > v {
		num = 4
	}

	return num
}

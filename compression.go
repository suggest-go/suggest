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
	list := make([]int, len(bytes)/4)

	for i := range list {
		list[i] = int(binary.LittleEndian.Uint32(bytes[4*i:]))
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
	sum, l := 0, 0

	for _, v := range list {
		sum += estimateByteNum(v)
	}

	encoded := make([]byte, sum)

	for _, v := range list {
		l += binary.PutUvarint(encoded[l:], uint64(v))
	}

	return encoded
}

func (b *vbEnc) Decode(bytes []byte) PostingList {
	v := uint32(0)
	s := uint(0)
	decoded := make([]int, 0, len(bytes)) //bad

	for _, b := range bytes {
		v |= uint32(b&0x7f) << s

		if b < 0x80 {
			decoded = append(decoded, int(v))
			s, v = 0, 0
		} else {
			s += 7
		}
	}

	return decoded
}

func estimateByteNum(v int) int {
	num := 0

	if (1 << 7) > v {
		num = 1
	} else if (1 << 14) > v {
		num = 2
	} else if (1 << 21) > v {
		num = 3
	} else if (1 << 28) > v {
		num = 4
	} else if (1 << 35) > v {
		num = 5
	}

	return num
}

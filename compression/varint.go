package compression

import "encoding/binary"

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

func (b *vbEnc) Encode(list []uint32) []byte {
	sum, l, prev, delta := 4, 4, uint32(0), uint32(0)

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
func (b *vbEnc) Decode(bytes []byte) []uint32 {
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
	decoded := make([]uint32, listLen)

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
		decoded[i] = prev
		i++
	}

	return decoded
}

// vbDecodeSlow decodes given byte array to posting list
func (b *vbEnc) vbDecodeSlow(bytes []byte, buf []uint32) {
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
			buf[i] = prev
			s, v = 0, 0
			i++
		} else {
			s += 7
		}
	}
}

// estimateByteNum returns bytes num required for encoding given uint32
func estimateByteNum(v uint32) int {
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

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

// Encode encodes the given positing list into the buf array
// Returns number of elements encoded, number of bytes readed
func (b *vbEnc) Encode(list []uint32, buf io.Writer) (int, error) {
	chunk := make([]uint32, 5)

	binary.LittleEndian.PutUint32(chunk, uint32(len(list)))
	n, err := buf.Write(chunk)

	if err != nil {
		return n, err
	}

	total := n

	for i, v := range list {
		j := 0
		delta = v - prev
		prev = v

		for ; delta > 0x7F; j++ {
			chunk[j] = 0x80 | uint8(delta&0x7F)
			delta >>= 7
		}

		chunk[j] = uint8(delta)
		j++

		n, err := buf.Write(chunk[:j])

		if err != nil {
			return n, err
		}

		total += n
	}

	i, j, prev, delta := 0, 4, uint32(0), uint32(0)
	binary.LittleEndian.PutUint32(buf, uint32(len(list)))

	for i < len(list) && j < len(buf); i++
		v := list[i]
		delta = v - prev
		prev = v

		for ; delta > 0x7F; j++ {
			buf[j] = 0x80 | uint8(delta&0x7F)
			delta >>= 7
		}

		buf[j] = uint8(delta)
		j++
	}

	return i, j
}

// inspired by protobuf/master/proto/decode.go
//
// Decode decodes the given byte array to the buf list
// Returns number of bytes readed, number of elements encoded
func (b *vbEnc) Decode(list []byte, buf []uint32) (int, int)
	if len(list) < 4 {
		return 0, 0
	}

	var (
		v    = uint32(0)
		prev = uint32(0)
		s    = uint32(0)
		i    = 4
		j    = 0
	)

	listLen := int(binary.LittleEndian.Uint32(list))

	if listLen < 10 {
		return b.vbDecodeSlow(list[4:], buf)
	}

	for i < len(list) && j < len(buf) {
		if list[i] < 0x80 {
			v = uint32(list[i])
			j++
			goto done
		}

		// we already checked the first byte
		v = uint32(list[i]) - 0x80
		i++

		s = uint32(list[i])
		i++
		v += s << 7
		if s&0x80 == 0 {
			goto done
		}
		v -= 0x80 << 7

		s = uint32(list[i])
		i++
		v += s << 14
		if s&0x80 == 0 {
			goto done
		}
		v -= 0x80 << 14

		s = uint32(list[i])
		i++
		v += s << 21
		if s&0x80 == 0 {
			goto done
		}
		v -= 0x80 << 21

		s = uint32(list[i])
		i++
		v += s << 28

	done:
		prev = v + prev
		list[j] = prev
		j++
	}

	return j, i
}

// vbDecodeSlow decodes given byte array to posting list
func (b *vbEnc) vbDecodeSlow(list []byte, buf []uint32) (int, int) {
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

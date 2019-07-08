package compression

import (
	"bufio"
	"io"
	"sync"
)

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
	var (
		prev  = uint32(0)
		delta = uint32(0)
		total = 0
		chunk = make([]byte, 5)
		j     = 0
	)

	for _, v := range list {
		j = 0
		delta = v - prev
		prev = v

		for ; delta > 0x7F; j++ {
			chunk[j] = 0x80 | uint8(delta&0x7F)
			delta >>= 7
		}

		chunk[j] = uint8(delta)
		j++

		n, err := buf.Write(chunk[:j])
		total += n

		if err != nil {
			if err == io.EOF {
				return total, nil
			}

			return total, err
		}
	}

	return total, nil
}

// readerPool reduces allocation of bufio.Reader object
var readerPool = sync.Pool{
	New: func() interface{} {
		return bufio.NewReaderSize(nil, 128)
	},
}

// inspired by protobuf/master/proto/decode.go
//
// Decode decodes the given byte array to the buf list
// Returns a number of elements encoded
func (b *vbEnc) Decode(in io.Reader, buf []uint32) (int, error) {
	var (
		v      = uint32(0)
		prev   = uint32(0)
		s      = uint32(0)
		total  = 0
		reader io.ByteReader
	)

	if byteReader, ok := in.(io.ByteReader); ok {
		reader = byteReader
	} else {
		r := readerPool.Get().(*bufio.Reader)
		defer readerPool.Put(r)
		r.Reset(in)
		reader = r
	}

	for total < len(buf) {
		b, err := reader.ReadByte()

		if err != nil {
			if err == io.EOF {
				return total, nil
			}

			return total, err
		}

		v |= uint32(b&0x7f) << s

		if b < 0x80 {
			prev = v + prev
			buf[total] = prev
			s, v = 0, 0
			total++
		} else {
			s += 7
		}
	}

	return total, nil
}

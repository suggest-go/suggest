package compression

import (
	"bytes"
	"encoding/binary"
	"github.com/alldroll/suggest/pkg/store"
)

// max 5 byte for var uint32
// gap - max 128, 128 * 5 = 640, uint16 for position
// var uint32 for diff

// 1 13 29 101 506 10003 10004 12000 12901
//
// 1 12 16 72 405 9497 1 1996 901 (just var uint32)
// gap 3

//
// size - 1  1  1  1   2    2  1   2     2   2
// gap  - 2  0  0  0   2    0  0   0     2   2
// star - 0  3  4  5   6   10 12  13    15  19
// vari - 1 12 16 72 505 9497 1 1996 11494 901

// SkippingEncoder creates a new instance of skipping encoder
func SkippingEncoder(gap int) Encoder {
	return &skippingEnc{
		enc: &vbEnc{},
		gap: gap,
	}
}

// SkippingDecoder creates a new instance of skipping decoder
func SkippingDecoder(gap int) Decoder {
	return &skippingEnc{
		gap: gap,
	}
}

// skippingEnc implements skippingEnc
type skippingEnc struct {
	enc *vbEnc
	gap int
}

// Encode encodes the given positing list into the buf array
// Returns number of elements encoded, number of bytes readed
func (b *skippingEnc) Encode(list []uint32, out store.Output) (int, error) {
	var (
		buf     = &bytes.Buffer{} // TODO use estimateByteNum
		prev    = uint32(0)
		pos     = 0
		total   = 0
		listLen = len(list)
		chunk   = make([]uint32, b.gap)
	)

	for i := 0; i < listLen; i += b.gap {
		j := i + b.gap

		if j > listLen {
			j = listLen
		}

		copy(chunk, list[i:j])

		chunk[0] = chunk[0] - prev
		prev = chunk[0]

		n, err := b.enc.Encode(chunk[:j-i], buf)

		if err != nil {
			return 0, err
		}

		total += n + 2
		pos = total - pos

		if err := binary.Write(out, binary.LittleEndian, uint16(pos)); err != nil {
			return 0, err
		}

		if _, err := buf.WriteTo(out); err != nil {
			return 0, err
		}
	}

	return total, nil
}

// Decode decodes the given byte array to the buf list
// Returns a number of elements encoded
func (b *skippingEnc) Decode(in store.Input, buf []uint32) (int, error) {
	var (
		prevV    = uint32(0)
		total    = 0
		prevSkip = uint32(0)
	)

	for total < len(buf) {
		_, err := in.ReadUInt16()

		if err != nil {
			return 0, err
		}

		for i := 0; i < b.gap && total < len(buf); i++ {
			v, err := in.ReadVUInt32()

			if err != nil {
				return 0, err
			}

			if i == 0 {
				buf[total] = prevSkip + v
				prevV = v
				prevSkip = v
			} else {
				buf[total] = prevV + v
				prevV = buf[total]
			}

			total++
		}
	}

	return total, nil
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

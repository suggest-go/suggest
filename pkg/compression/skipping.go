package compression

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/alldroll/suggest/pkg/store"
)

var (
	// ErrGapShouldBeGreaterThanListLen tells that the list length is less or equal to
	// the skipping gap
	ErrGapShouldBeGreaterThanListLen = errors.New("gap should be greater than the list length")
	// ErrGapOverflow tells that it was at attempt to create
	// encoder/decoder skipping gap more than maxSkippingGap
	ErrGapOverflow = errors.New("gap value overflow")
)

const (
	lastBlockFlag  = 1 << 15
	maxSkippingGap = (1 << 14) / 5
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
func SkippingEncoder(gap int) (Encoder, error) {
	if gap >= maxSkippingGap {
		return nil, ErrGapOverflow
	}

	return &skippingEnc{
		enc: &vbEnc{},
		gap: gap,
	}, nil
}

// SkippingDecoder creates a new instance of skipping decoder
func SkippingDecoder(gap int) (Decoder, error) {
	if gap >= maxSkippingGap {
		return nil, ErrGapOverflow
	}

	return &skippingEnc{
		gap: gap,
	}, nil
}

// skippingEnc implements skippingEnc
type skippingEnc struct {
	enc *vbEnc
	gap int
}

// Encode encodes the given positing list into the buf array
// Returns number of elements encoded, number of bytes readed
func (b *skippingEnc) Encode(list []uint32, out store.Output) (int, error) {
	if len(list) < b.gap {
		return 0, ErrGapShouldBeGreaterThanListLen
	}

	var (
		// TODO use estimateByteNum instead of buffer
		buf     = &bytes.Buffer{}
		prev    = uint32(0)
		total   = 0
		chunk   = make([]uint32, b.gap)
		listLen = len(list)
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

		pos := n + 2
		total += pos

		if j == listLen {
			pos = pos | lastBlockFlag
		}

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

// UnpackPos describe me!!
func UnpackPos(packed uint16) (int, bool) {
	return int(packed & uint16(lastBlockFlag-1)), (packed & lastBlockFlag) == lastBlockFlag
}

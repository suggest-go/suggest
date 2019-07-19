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
	// we are keeping the 16th bit as a marker of the last block
	lastBlockFlag  = 1 << 15
	// as the max size of var uint32 is 5 bytes, the maxSkippingGap will be 2^15 / 5
	maxSkippingGap = (1 << 14) / 5
)

// let's imagine we have the next sequence:
// 1 13 29 101 506 10003 10004 12000 12901
//
// with the gap size = 3, we will have the next:
//
//  1  1  1   1   2    2    2    2   2   - the bytes length of var int
//  2  0  0   2   0    0    2    0   0   - the additional 2 bytes for the indicating of block start
//  1 12 16 100 405 9497 9903 1996 901   - delta encoded values of the sequence
// (1 - 0) (101 - 1)    (10004 - 101)..  - deltas for block starts

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
		buf     = bytes.NewBuffer(make([]byte, 0, b.gap * 5)) // max var int * 5
		prev    = uint32(0)
		total   = 0
		listLen = len(list)
	)

	for i := 0; i < listLen; i += b.gap {
		j := i + b.gap

		if j > listLen {
			j = listLen
		}

		// write encoded value into buffer (we should know the encoded size first)
		n, err := varIntEncode(list[i:j], buf, prev)
		prev = list[i]

		if err != nil {
			return 0, err
		}

		pos := n + 2
		total += pos

		// marks the 16 bit with flag, if it the last block
		if j == listLen {
			pos = pos | lastBlockFlag
		}

		// write the start position and the indicator of the last block at the first stage
		if err := binary.Write(out, binary.LittleEndian, uint16(pos)); err != nil {
			return 0, err
		}

		// flush the buffer with the encoded slice at the second stage
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
		prev    = uint32(0)
		i       = 0
		listLen = len(buf)
	)

	for ; i < listLen; i += b.gap {
		if _, err := in.ReadUInt16(); err != nil {
			return 0, err
		}

		j := i + b.gap

		if j > listLen {
			j = listLen
		}

		_, err := varIntDecode(in, buf[i:j], prev)
		prev = buf[i]

		if err != nil {
			return 0, err
		}
	}

	return i, nil
}

// UnpackPos splits the given uint16 packed value on a pair (delta position, is last block flag)
func UnpackPos(packed uint16) (int, bool) {
	return int(packed & uint16(lastBlockFlag-1)), (packed & lastBlockFlag) == lastBlockFlag
}

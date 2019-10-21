package store

// Inspired by https://github.com/golang/protobuf

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	// ErrUInt32Overflow is returned when an integer is too large to be represented
	ErrUInt32Overflow = errors.New("bytesInput: uint32 overflow")
	// ErrNegativeOffset tells that it was an attempt to get access with a negative offset
	ErrNegativeOffset = errors.New("bytesInput: negative offset")
	// ErrOutOfRange tells that it was an attemot to get access out of range
	ErrOutOfRange = errors.New("bytesInput: try to get access out of range")
	// ErrInvalidWhence tells that the whence has invalid value
	ErrInvalidWhence = errors.New("bytesInput: invalid whence")
)

// NewBytesInput creates a new instance of byteInput
func NewBytesInput(buf []byte) Input {
	return &byteInput{
		buf: buf,
		i:   0,
	}
}

// byteInput implements the Input interface for the bytes slice
type byteInput struct {
	buf []byte
	i   int64
}

// Close closes the given byteInput for io operations
func (r *byteInput) Close() error {
	return nil
}

// Read implements the io.Reader interface.
func (r *byteInput) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.buf)) {
		return 0, io.EOF
	}

	n = copy(b, r.buf[r.i:])
	r.i += int64(n)

	return
}

// Data returns the underlying content as byte slice
func (r *byteInput) Data() []byte {
	return r.buf
}

// ReadAt implements the io.ReaderAt interface.
func (r *byteInput) ReadAt(b []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, ErrNegativeOffset
	}

	if off >= int64(len(r.buf)) {
		return 0, io.EOF
	}

	n = copy(b, r.buf[off:])

	if n < len(b) {
		err = io.EOF
	}

	return
}

// ReadByte implements the io.ByteReader interface.
func (r *byteInput) ReadByte() (byte, error) {
	if r.i >= int64(len(r.buf)) {
		return 0, io.EOF
	}

	b := r.buf[r.i]
	r.i++

	return b, nil
}

// Seek sets the offset for the next Read to offset,
// interpreted according to whence: io.SeekStart, io.SeekEnd, io.SeekCurrent
func (r *byteInput) Seek(offset int64, whence int) (int64, error) {
	var abs int64

	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.i + offset
	case io.SeekEnd:
		abs = int64(len(r.buf)) + offset
	default:
		return 0, ErrInvalidWhence
	}

	if abs < 0 {
		return 0, ErrNegativeOffset
	}

	r.i = abs

	return abs, nil
}

// Slice returns a slice of the given Input
func (r *byteInput) Slice(off int64, n int64) (Input, error) {
	if n < 0 || off < 0 || int64(len(r.buf)) < (off+n) {
		return nil, ErrOutOfRange
	}

	data := r.buf[off : off+n]

	return &byteInput{
		buf: data,
		i:   0,
	}, nil
}

// ReadVUInt32 reads a variable-length encoded uint32 number
// Inspired by https://github.com/golang/protobuf/blob/master/proto/decode.go
func (r *byteInput) ReadVUInt32() (uint32, error) {
	var (
		i = r.i
		l = int64(len(r.buf))
		b byte
	)

	if l <= i {
		return 0, ErrUInt32Overflow
	} else if r.buf[i] < 0x80 {
		r.i++
		return uint32(r.buf[i]), nil
	} else if l-i < 5 {
		return r.readVarUInt32Slow()
	}

	v := uint32(r.buf[i]) - 0x80
	i++

	b = r.buf[i]
	i++
	v += uint32(b) << 7

	if b&0x80 == 0 {
		goto done
	}

	v -= 0x80 << 7

	b = r.buf[i]
	i++
	v += uint32(b) << 14

	if b&0x80 == 0 {
		goto done
	}

	v -= 0x80 << 14

	b = r.buf[i]
	i++
	v += uint32(b) << 21

	if b&0x80 == 0 {
		goto done
	}

	v -= 0x80 << 21

	b = r.buf[i]
	i++
	v += uint32(b) << 28

	if b&0x80 == 0 {
		goto done
	}

	return 0, ErrUInt32Overflow

done:
	r.i = i

	return v, nil
}

// ReadUInt32 reads four bytes and returns uint32
func (r *byteInput) ReadUInt32() (uint32, error) {
	if r.i+4 > int64(len(r.buf)) {
		return 0, io.ErrUnexpectedEOF
	}

	v := binary.LittleEndian.Uint32(r.buf[r.i:])
	r.i += 4

	return v, nil
}

// ReadUInt16 reads two bytes and returns uint16
func (r *byteInput) ReadUInt16() (uint16, error) {
	if r.i+2 > int64(len(r.buf)) {
		return 0, io.ErrUnexpectedEOF
	}

	v := binary.LittleEndian.Uint16(r.buf[r.i:])
	r.i += 2

	return v, nil
}

// readVarUInt32Slow decodes VarUInt32 with the loop approach
// Inspired by https://github.com/golang/protobuf/blob/master/proto/decode.go
func (r *byteInput) readVarUInt32Slow() (uint32, error) {
	var (
		i = r.i
		l = int64(len(r.buf))
		v = uint32(0)
		b byte
	)

	for s := uint32(0); s < 35; s += 7 {
		if i >= l {
			return 0, io.ErrUnexpectedEOF
		}

		b = r.buf[i]
		v |= uint32(b&0x7f) << s
		i++

		if b < 0x80 {
			r.i = i
			return v, nil
		}
	}

	return 0, ErrUInt32Overflow
}

package store

import (
	"encoding/binary"
	"errors"
	"io"
)

func NewBytesInput(buf []byte) Input {
	return &byteInput{
		buf: buf,
		i:   0,
	}
}

type byteInput struct {
	buf []byte
	i   int64
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

// ReadAt implements the io.ReaderAt interface.
func (r *byteInput) ReadAt(b []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errors.New("byteInput.ReadAt: negative offset") // TODO declare in var
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

// Slice returns a slice of the given Input
func (r *byteInput) Slice(off int64, n int64) (Input, error) {
	if n < 0 || off < 0 || int64(len(r.buf)) < (off+n) {
		return nil, errors.New("TODO complete me") // TODO declare in var
	}

	data := r.buf[off : off+n]

	return &byteInput{
		buf: data,
		i:   0,
	}, nil
}

// ReadVUInt32 reads a variable-length decoded uint32 number
func (r *byteInput) ReadVUInt32() (uint32, error) {
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

	return 0, errors.New("Overflow")
}

// ReadUInt32 reads a binary decoded uint32 number
func (r *byteInput) ReadUInt32() (uint32, error) {
	if r.i+4 > int64(len(r.buf)) {
		return 0, io.ErrUnexpectedEOF
	}

	v := binary.LittleEndian.Uint32(r.buf[r.i:])
	r.i += 4

	return v, nil
}

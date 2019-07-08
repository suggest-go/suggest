package store

import (
	"bytes"
	"errors"
)

type byteInput struct {
	*bytes.Reader
	buf []byte
}

// Slice returns a slice of the given Input
func (r *byteInput) Slice(off int64, n int64) (Input, error) {
	if n < 0 || off < 0 || int64(r.Len()) < (off+n) {
		return nil, errors.New("TODO complete me")
	}

	data := r.buf[off : off+n]

	return &byteInput{
		buf:    data,
		Reader: bytes.NewReader(data),
	}, nil
}

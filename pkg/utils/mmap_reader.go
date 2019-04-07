package utils

import (
	"errors"
	"io"
	"os"
	"runtime"

	mmap "github.com/edsrzf/mmap-go"
)

var (
	// ErrMMapIsClosed means that it was an attempt to read data
	// from the closed region
	ErrMMapIsClosed = errors.New("MMap file is closed")

	// ErrMMapInvalidOffset means that it was an attempt to read data
	// from the invalid offset
	ErrMMapInvalidOffset = errors.New("Out of range")
)

// NewMMapReader returns new instance of MMapReader
func NewMMapReader(filename string) (*MMapReader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := mmap.Map(file, mmap.RDONLY, 0)
	if err != nil {
		return nil, err
	}

	r := &MMapReader{
		data:    data,
		readPos: 0,
	}

	runtime.SetFinalizer(r, (*MMapReader).Close)
	return r, nil
}

// MMapReader wraps Read methods by using mmap
type MMapReader struct {
	data    mmap.MMap
	readPos int
}

// Close releases unmap of the choosen region
func (r *MMapReader) Close() error {
	if r.data == nil {
		return ErrMMapIsClosed
	}

	data := r.data
	r.data = nil

	runtime.SetFinalizer(r, nil)
	return data.Unmap()
}

// Read reads and copies data to p
func (r *MMapReader) Read(p []byte) (int, error) {
	if r.data == nil {
		return 0, ErrMMapIsClosed
	}

	if r.readPos == len(r.data) {
		return 0, io.EOF
	}

	n, err := r.ReadAt(p, int64(r.readPos))
	r.readPos += n

	return n, err
}

// ReadAt reads and copies data to p from the offset off
func (r *MMapReader) ReadAt(p []byte, off int64) (int, error) {
	if r.data == nil {
		return 0, ErrMMapIsClosed
	}

	if off < 0 || int64(len(r.data)) < off {
		return 0, ErrMMapInvalidOffset
	}

	n := copy(p, r.data[off:])
	if n < len(p) {
		return n, io.EOF
	}

	return n, nil
}

// Data returns the mapped object
func (r *MMapReader) Data() ([]byte, error) {
	return r.data, nil
}

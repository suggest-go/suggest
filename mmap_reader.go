package suggest

import (
	"errors"
	mmap "github.com/edsrzf/mmap-go"
	"io"
	"os"
	"runtime"
)

var (
	mmapIsClosedErr   = errors.New("mmap: closed")
	mmapInvalidOffset = errors.New("mmap: invalid ReadAt offset")
)

type readerAt struct {
	data mmap.MMap
}

func (r *readerAt) Close() error {
	if r.data == nil {
		return mmapIsClosedErr
	}

	data := r.data
	r.data = nil

	runtime.SetFinalizer(r, nil)
	return data.Unmap()
}

func (r *readerAt) ReadAt(p []byte, off int64) (int, error) {
	if r.data == nil {
		return 0, mmapIsClosedErr
	}

	if off < 0 || int64(len(r.data)) < off {
		return 0, mmapInvalidOffset
	}

	n := copy(p, r.data[off:])
	if n < len(p) {
		return n, io.EOF
	}

	return n, nil
}

func NewMmapReader(filename string) (*readerAt, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := mmap.Map(file, mmap.RDONLY, 0)
	if err != nil {
		return nil, err
	}

	r := &readerAt{data}

	runtime.SetFinalizer(r, (*readerAt).Close)
	return r, nil
}

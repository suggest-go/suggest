package index

import (
	"bytes"
	"fmt"
	"io"
)

// ramDirectory is a implementation that stores index
// files in RAM
type ramDirectory struct {
	files map[string]*bytes.Buffer
}

// NewRAMDirectory returns a new instance of RAM Directory
func NewRAMDirectory() Directory {
	return &ramDirectory{
		files: make(map[string]*bytes.Buffer),
	}
}

// CreateOutput creates a new writer in the given directory with the given name
func (rd *ramDirectory) CreateOutput(name string) (Output, error) {
	if _, ok := rd.files[name]; !ok {
		rd.files[name] = &bytes.Buffer{}
	}

	return rd.files[name], nil
}

// OpenInput returns a reader for the given name
func (rd *ramDirectory) OpenInput(name string) (Input, error) {
	if _, ok := rd.files[name]; !ok {
		return nil, fmt.Errorf("Failed to open input reader: there is no such input with the name %v", name)
	}

	buf := rd.files[name].Bytes()

	return &bytesReader{
		buf:    buf,
		reader: bytes.NewReader(buf),
	}, nil
}

type bytesReader struct {
	buf    []byte
	reader io.Reader
}

// Data returns stored bytes
func (b *bytesReader) Data() ([]byte, error) {
	return b.buf, nil
}

// Read reads bytes and copies them to p
func (b *bytesReader) Read(p []byte) (int, error) {
	return b.reader.Read(p)
}

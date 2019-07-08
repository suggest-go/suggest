package store

import (
	"bytes"
	"fmt"
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

	data := rd.files[name].Bytes()

	return &byteInput{
		buf:    data,
		Reader: bytes.NewReader(data),
	}, nil
}

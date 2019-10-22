package store

import (
	"bufio"
	"fmt"
	"os"
	"runtime"

	"github.com/suggest-go/suggest/pkg/utils"
)

// fsDirectory is a implementation that stores index
// files in the file system.
type fsDirectory struct {
	path string
}

// NewFSDirectory creates a new instance of FS Directory
func NewFSDirectory(path string) (Directory, error) {
	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("Given path is not exists")
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to receive stat for the path %v", err)
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("Path should be a directory")
	}

	return &fsDirectory{
		path: path,
	}, nil
}

// CreateOutput creates a new writer in the given directory with the given name
func (fs *fsDirectory) CreateOutput(name string) (Output, error) {
	file, err := os.OpenFile(
		fs.path+"/"+name,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to create output: %v", err)
	}

	return NewBytesOutput(bufio.NewWriter(file)), nil
}

// OpenInput returns a reader for the given name
func (fs *fsDirectory) OpenInput(name string) (Input, error) {
	file, err := utils.NewMMapReader(fs.path + "/" + name)

	if err != nil {
		return nil, fmt.Errorf("Failed to open input: %v", err)
	}

	data, err := file.Bytes()

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch content: %v", err)
	}

	input := NewBytesInput(data)

	runtime.SetFinalizer(input, func(d interface{}) {
		file.Close()
	})

	return input, nil
}

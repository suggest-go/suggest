package index

import "io"

// Directory is a flat list of files.
// Inspired by org.apache.lucene.store.Directory
type Directory interface {
	// CreateOutput creates a new writer in the given directory with the given name
	CreateOutput(name string) (Output, error)
	// OpenInput returns a reader for the given name
	OpenInput(name string) (Input, error)
}

// Input is a wrap for methods Read and retrieving underlying data
type Input interface {
	io.Reader
	// Data returns stored bytes from the reader
	Data() ([]byte, error)
}

// Output is a wrap for method Write
type Output interface {
	io.Writer
}

package store

import "io"

// Input is a wrap for methods Read and retrieving underlying data
type Input interface {
	io.Reader
	io.ReaderAt
	io.ByteReader
	// Slice returns a slice of the given Input
	Slice(off int64, n int64) (Input, error)
}

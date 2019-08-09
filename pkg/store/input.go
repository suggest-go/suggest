package store

import "io"

// Input is a wrap for methods Read and retrieving underlying data
type Input interface {
	io.ReadSeeker
	io.ReaderAt
	io.ByteReader

	// Slice returns a slice of the given Input
	Slice(off int64, n int64) (Input, error)
	// ReadVUInt32 reads a variable-length decoded uint32 number
	ReadVUInt32() (uint32, error)
	// ReadUInt32 reads a binary decoded uint32 number
	ReadUInt32() (uint32, error)
	// ReadUInt16 reads a binary decoded uint16 number
	ReadUInt16() (uint16, error)
}

// SliceAccessible represents the entity with the ability to return underlying byte slice
type SliceAccessible interface {
	// Data returns the underlying content as byte slice
	Data() []byte
}

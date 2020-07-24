package store

import "io"

// Output is a wrap for method Write
type Output interface {
	io.Closer
	io.Writer

	// WriteVUInt32 writes the given uint32 in the variable-length format
	WriteVUInt32(v uint32) (int, error)
	// WriteUInt32 writes the given uint32 in the binary format
	WriteUInt32(v uint32) (int, error)
	// WriteUInt16 writes the given uint16 number in the binary format
	WriteUInt16(v uint16) (int, error)
	// WriteByte writes the given byte
	WriteByte(v byte) error
}

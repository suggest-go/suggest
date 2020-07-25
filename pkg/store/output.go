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

// Marshaler is the interface implemented by an object that can marshal itself into a binary form.
type Marshaler interface {
	// Store encodes the receiver into a binary form and saves the result into the provided Output.
	// Returns the number of written bytes or an error otherwise.
	Store(out Output) (int, error)
}

package store

import (
	"bufio"
	"encoding/binary"
	"io"
)

// NewBytesOutput creates a new instance of byteOutput
func NewBytesOutput(writer io.Writer) Output {
	return &byteOutput{
		writer: writer,
	}
}

type byteOutput struct {
	writer io.Writer
}

// Write wraps io.Writer.Write method
func (b *byteOutput) Write(p []byte) (n int, err error) {
	return b.writer.Write(p)
}

// WriteVUInt32 writes the given uint32 in the variable-length format
func (b *byteOutput) WriteVUInt32(v uint32) (int, error) {
	chunk := make([]byte, 5)
	j := 0

	for ; v > 0x7F; j++ {
		chunk[j] = 0x80 | uint8(v&0x7F)
		v >>= 7
	}

	chunk[j] = uint8(v)

	return b.writer.Write(chunk[:j+1])
}

// WriteUInt32 writes the given uint32 in the binary format
func (b *byteOutput) WriteUInt32(v uint32) (int, error) {
	chunk := make([]byte, 4)
	binary.LittleEndian.PutUint32(chunk, v)

	return b.writer.Write(chunk)
}

// WriteUInt16 writes the given uint32 number in the binary format
func (b *byteOutput) WriteUInt16(v uint16) (int, error) {
	chunk := make([]byte, 2)
	binary.LittleEndian.PutUint16(chunk, v)

	return b.writer.Write(chunk)
}

// WriteByte writes the given byte
func (b *byteOutput) WriteByte(v byte) error {
	_, err := b.writer.Write([]byte{v})

	return err
}

// Close closes the given output for io operations
func (b *byteOutput) Close() error {
	if buf, ok := b.writer.(*bufio.Writer); ok {
		if err := buf.Flush(); err != nil {
			return err
		}
	}

	if closer, ok := b.writer.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}

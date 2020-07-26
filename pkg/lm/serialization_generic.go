// +build !amd64,!386 appengine

package lm

import (
	"encoding/binary"
)

func uint64SliceAsByteSlice(slice []uint64) []byte {
	by := make([]byte, len(slice)*8)

	for i, v := range slice {
		binary.LittleEndian.PutUint64(by[i*8:], v)
	}

	return by
}

func byteSliceAsUint64Slice(slice []byte) []uint64 {
	if len(slice)%8 != 0 {
		panic("Slice size should be divisible by 8")
	}

	b := make([]uint64, len(slice)/8)

	for i := range b {
		b[i] = binary.LittleEndian.Uint64(slice[8*i:])
	}

	return b
}

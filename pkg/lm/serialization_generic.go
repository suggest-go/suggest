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

func rangeContainerSliceAsByteSlice(slice []rangeContainer) []byte {
	by := make([]byte, len(slice)*12)

	for i, v := range slice {
		binary.LittleEndian.PutUint32(by[i*12:], v.context)
		binary.LittleEndian.PutUint32(by[i*12+4:], v.from)
		binary.LittleEndian.PutUint32(by[i*12+8:], v.to)
	}

	return by
}

func byteSliceAsRangeContainerSlice(slice []byte) (result []rangeContainer) {
	if len(slice)%12 != 0 {
		panic("Slice size should be divisible by 12")
	}

	b := make([]rangeContainer, len(slice)/12)

	for i := range b {
		b[i].context = binary.LittleEndian.Uint32(slice[12*i:])
		b[i].from = binary.LittleEndian.Uint32(slice[12*i+4:])
		b[i].to = binary.LittleEndian.Uint32(slice[12*i+8:])
	}

	return b
}

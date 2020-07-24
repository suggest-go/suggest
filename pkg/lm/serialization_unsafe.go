// +build 386 amd64,!appengine

package lm

import (
	"reflect"
	"runtime"
	"unsafe"
)

// inspired by https://github.com/RoaringBitmap/roaring and https://github.com/RoaringBitmap/roaring

func uint64SliceAsByteSlice(slice []uint64) []byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&slice))

	header.Len *= 8
	header.Cap *= 8

	result := *(*[]byte)(unsafe.Pointer(&header))
	runtime.KeepAlive(&slice)

	return result
}

func byteSliceAsUint64Slice(slice []byte) (result []uint64) {
	if len(slice)%8 != 0 {
		panic("Slice size should be divisible by 8")
	}

	bHeader := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	rHeader := (*reflect.SliceHeader)(unsafe.Pointer(&result))

	rHeader.Data = bHeader.Data
	rHeader.Len = bHeader.Len / 8
	rHeader.Cap = bHeader.Cap / 8

	runtime.KeepAlive(&slice)

	return
}

func rangeContainerSliceAsByteSlice(slice []rangeContainer) []byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&slice))

	header.Len *= 12
	header.Cap *= 12

	result := *(*[]byte)(unsafe.Pointer(&header))
	runtime.KeepAlive(&slice)

	return result
}

func byteSliceAsRangeContainerSlice(slice []byte) (result []rangeContainer) {
	if len(slice)%12 != 0 {
		panic("Slice size should be divisible by 8")
	}

	bHeader := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	rHeader := (*reflect.SliceHeader)(unsafe.Pointer(&result))

	rHeader.Data = bHeader.Data
	rHeader.Len = bHeader.Len / 12
	rHeader.Cap = bHeader.Cap / 12

	runtime.KeepAlive(&slice)

	return
}

package utils

import "math"

// Max returns the maximum value
func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// Min returns the minimum value
func Min(a, b int) int {
	if a > b {
		return b
	}

	return a
}

// Pack packes 2 uint32 into uint64
func Pack(a, b uint32) uint64 {
	return (uint64(a) << 32) | uint64(b&math.MaxUint32)
}

// Unpack explodes uint64 into 2 uint32
func Unpack(v uint64) (uint32, uint32) {
	return uint32(v >> 32), uint32(v)
}

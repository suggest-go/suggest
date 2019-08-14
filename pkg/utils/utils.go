package utils

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

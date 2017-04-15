package suggest

type record int

func (self record) Less(other heapItem) bool {
	return self > other.(record)
}

// see Efficient Merging and Filtering Algorithms for
// Approximate String Searches

func mergeSkip(rid [][]int, threshold int) []int {
	return make([]int, 0)
}

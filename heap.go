package suggest

type heapItem interface {
	Less(other heapItem) bool
}

type heapImpl []heapItem

//
func (h heapImpl) Len() int {
	return len(h)
}

//
func (h heapImpl) Less(i, j int) bool {
	return h[i].Less(h[j])
}

//
func (h heapImpl) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

//
func (h *heapImpl) Push(x interface{}) {
	*h = append(*h, x.(heapItem))
}

//
func (h *heapImpl) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

//
func (h *heapImpl) Top() interface{} {
	arr := *h
	return arr[0]
}

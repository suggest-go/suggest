package suggest

// heapItem represents element of heap
type heapItem interface {
	// Less reports whether the given elements less than other
	Less(other heapItem) bool
}

// heapImpl implements heap.Interface
type heapImpl []heapItem

// newHeap returns new instance of heap.Interface with given capacity
func newHeap(capacity int) *heapImpl {
	hp := make(heapImpl, 0, capacity)
	return &hp
}

// Len is the number of elements in the collection.
func (h heapImpl) Len() int {
	return len(h)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (h heapImpl) Less(i, j int) bool {
	return h[i].Less(h[j])
}

// Swap swaps the elements with indexes i and j.
func (h heapImpl) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Push add x as element Len()
func (h *heapImpl) Push(x interface{}) {
	*h = append(*h, x.(heapItem))
}

// Pop remove and return element Len() - 1.
func (h *heapImpl) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Top returns the top element of heap
func (h heapImpl) Top() interface{} {
	if len(h) == 0 {
		panic("Try to get top element on empty heap")
	}

	return h[0]
}

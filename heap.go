package suggest

// heapItem represents element of heap
type heapItem interface {
	// Less reports whether the given elements less than other
	Less(other heapItem) bool
}

// heapImpl implements heap.Interface
type heapImpl struct {
	heap []heapItem
}

// newHeap returns new instance of heap.Interface with given capacity
func newHeap(capacity int) *heapImpl {
	return &heapImpl{
		heap: make([]heapItem, 0, capacity),
	}
}

// Len is the number of elements in the collection.
func (h *heapImpl) Len() int {
	return len(h.heap)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (h *heapImpl) Less(i, j int) bool {
	return h.heap[i].Less(h.heap[j])
}

// Swap swaps the elements with indexes i and j.
func (h *heapImpl) Swap(i, j int) {
	h.heap[i], h.heap[j] = h.heap[j], h.heap[i]
}

// Push add x as element Len()
func (h *heapImpl) Push(x interface{}) {
	h.heap = append(h.heap, x.(heapItem))
}

// Pop remove and return element Len() - 1.
func (h *heapImpl) Pop() interface{} {
	old := h.heap
	n := len(old)
	x := old[n-1]
	h.heap = old[:n-1]
	return x
}

// Top returns the top element of heap
func (h *heapImpl) Top() heapItem {
	if len(h.heap) == 0 {
		panic("Try to get top element on empty heap")
	}

	return h.heap[0]
}

package suggest

type heapItem interface {
	Less(other heapItem) bool
}

type heapImpl []heapItem

func (self heapImpl) Len() int {
	return len(self)
}

func (self heapImpl) Less(i, j int) bool {
	return self[i].Less(self[j])
}

func (self heapImpl) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self *heapImpl) Push(x interface{}) {
	*self = append(*self, x.(heapItem))
}

func (self *heapImpl) Pop() interface{} {
	old := *self
	n := len(old)
	x := old[n-1]
	*self = old[0 : n-1]
	return x
}

func (self *heapImpl) Top() interface{} {
	arr := *self
	return arr[0]
}

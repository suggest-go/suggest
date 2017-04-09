package suggest

// Implement rank heap with distance comparator

type rank struct {
	id       int
	distance float64
}

func (self *rank) Less(other *rank) bool { return self.distance > other.distance }

type rankHeap []*rank

func (self rankHeap) Len() int           { return len(self) }
func (self rankHeap) Less(i, j int) bool { return self[i].Less(self[j]) }
func (self rankHeap) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }

func (self *rankHeap) Push(x interface{}) {
	*self = append(*self, x.(*rank))
}

func (self *rankHeap) Pop() interface{} {
	old := *self
	n := len(old)
	x := old[n-1]
	*self = old[0 : n-1]
	return x
}

func (self *rankHeap) Min() interface{} {
	arr := *self
	return arr[0]
}

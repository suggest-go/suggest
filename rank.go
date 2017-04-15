package suggest

// Implement rank distance comparator

type rank struct {
	id       int
	distance float64
}

func (self *rank) Less(other heapItem) bool {
	o := other.(*rank)
	return self.distance > o.distance
}

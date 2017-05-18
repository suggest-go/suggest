package suggest

// Implement rank distance comparator

type rank struct {
	id       int
	distance float64
}

func (r *rank) Less(other heapItem) bool {
	o := other.(*rank)
	return r.distance > o.distance
}

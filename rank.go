package suggest

// Implement rank distance comparator

// Candidate is representing candidate of the similarity of the query
type Candidate struct {
	Key      int
	Distance float64
}

type rank struct {
	id       int
	overlap  int
	distance float64
}

func (r *rank) Less(other heapItem) bool {
	o := other.(*rank)
	return r.distance > o.distance
}

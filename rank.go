package suggest

// Implement rank distance comparator

// Candidate is representing candidate of the similarity of the query
type Candidate struct {
	Key      Position
	Distance float64
}

type rank struct {
	id       Position
	overlap  int
	distance float64
}

func (r *rank) Less(other heapItem) bool {
	o := other.(*rank)
	return r.distance > o.distance
}

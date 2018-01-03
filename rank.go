package suggest

// Candidate is representing candidate of the similarity of the query
type Candidate struct {
	// Key is position (docId) in posting list
	Key      Position
	// Distance is float64 number from [0, 1]
	Distance float64
}

// rank implements heapItem. simple rank based on distance compare
type rank struct {
	id       Position
	overlap  int
	distance float64
}

// Less reports whether the given elements less than other
func (r *rank) Less(other heapItem) bool {
	o := other.(*rank)
	return r.distance > o.distance
}

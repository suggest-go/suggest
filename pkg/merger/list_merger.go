package merger

// ListMerger solves `threshold`-occurrence problem:
// For given inverted lists find the set of strings ids, that appears at least
// `threshold` times.
type ListMerger interface {
	// Merge returns list of candidates, that appears at least `threshold` times.
	Merge(rid Rid, threshold int) ([]MergeCandidate, error)
}

// Rid represents inverted lists for ListMerger
type Rid []ListIterator

// Len is the number of elements in the collection.
func (p Rid) Len() int { return len(p) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (p Rid) Less(i, j int) bool { return p[i].Len() < p[j].Len() }

// Swap swaps the elements with indexes i and j.
func (p Rid) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// MergeCandidate is result of merging Rid
type MergeCandidate struct {
	Position uint32
	Overlap  int
}

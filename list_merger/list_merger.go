package list_merger

import "sort"

// ListMerger solves `threshold`-occurrence problem:
// For given inverted lists find the set of strings ids, that appears at least
// `threshold` times.
type ListMerger interface {
	// Merge returns list of candidates, that appears at least `threshold` times.
	Merge(rid Rid, threshold int) []*MergeCandidate
}

// ListIntersect
type ListIntersect interface {
	Intersect(rid Rid, max int) []*MergeCandidate
}

type RidItem = []uint32

// Rid represents inverted lists for ListMerger
type Rid []RidItem

// Len is the number of elements in the collection.
func (p Rid) Len() int { return len(p) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (p Rid) Less(i, j int) bool { return len(p[i]) < len(p[j]) }

// Swap swaps the elements with indexes i and j.
func (p Rid) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// MergeCandidate is result of merging Rid
type MergeCandidate struct {
	Position uint32
	Overlap  int
}

// lowerBound returns index for the smallest record t in given arr such that t >= x
// returns -1, if there is not such item
func lowerBound(a RidItem, x uint32) int {
	i := sort.Search(len(a), func(i int) bool { return a[i] >= x })
	if i < 0 || i >= len(a) {
		i = -1
	}

	return i
}

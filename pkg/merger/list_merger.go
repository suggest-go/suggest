// Package merger provides a different set of algorithms for solving T-overlap occurrence problem of sorted lists of integers.
// T-occurrence problem can be described next:
// - Find the set of string ids that appear at least T times on the inverted lists, where T is a constant.
package merger

import "github.com/suggest-go/suggest/pkg/utils"

// MaxOverlap is the largest value of an overlap count for a merge candidate
const MaxOverlap = 0xFFFF

// ListMerger solves `threshold`-occurrence problem:
// For given inverted lists find the set of strings ids, that appears at least
// `threshold` times.
type ListMerger interface {
	// Merge returns list of candidates, that appears at least `threshold` times.
	Merge(rid Rid, threshold int, collector Collector) error
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
type MergeCandidate uint64

// NewMergeCandidate creates a new instance of MergeCandidate
func NewMergeCandidate(position, overlap uint32) MergeCandidate {
	return MergeCandidate(utils.Pack(position, overlap))
}

// Position returns the given position of the candidate
func (m MergeCandidate) Position() uint32 {
	return utils.UnpackLeft(uint64(m))
}

// Overlap returns the current overlap count of the candidate
func (m MergeCandidate) Overlap() int {
	return int(utils.UnpackRight(uint64(m)))
}

// increment increments the overlap value of the candidate
func (m *MergeCandidate) increment() {
	if overlap := utils.UnpackRight(uint64(*m)); overlap == MaxOverlap {
		panic("overlap overflow")
	}

	*m++
}

// mergerOptimizer internal merger that is aimed to optimize merge workflow
type mergerOptimizer struct {
	merger      ListMerger
	intersector ListIntersector
}

func newMerger(merger ListMerger) ListMerger {
	return &mergerOptimizer{
		merger:      merger,
		intersector: Intersector(),
	}
}

// Merge returns list of candidates, that appears at least `threshold` times.
func (m *mergerOptimizer) Merge(rid Rid, threshold int, collector Collector) error {
	n := len(rid)

	if n < threshold || n == 0 || threshold < 0 {
		return nil
	}

	if n == threshold {
		return m.intersector.Intersect(rid, collector)
	}

	return m.merger.Merge(rid, threshold, collector)
}

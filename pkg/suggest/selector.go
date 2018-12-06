package suggest

import (
	"container/heap"
	"github.com/alldroll/suggest/pkg/index"
	"github.com/alldroll/suggest/pkg/list_merger"
)

type Candidate struct {
	// Key is position (docId) in posting list
	Key index.Position
}

// Candidate is representing candidate of the similarity of the query
type FuzzyCandidate struct {
	Candidate
	// Distance is float64 number from [0, 1]
	Distance float64
}

type topKFuzzySearchItem struct {
	candidate *list_merger.MergeCandidate
	distance  float64
}

type topKHeap []*topKFuzzySearchItem

// Len is the number of elements in the collection.
func (h topKHeap) Len() int { return len(h) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (h topKHeap) Less(i, j int) bool { return h[i].distance > h[j].distance }

// Swap swaps the elements with indexes i and j.
func (h topKHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push add x as element Len()
func (h *topKHeap) Push(x interface{}) { *h = append(*h, x.(*topKFuzzySearchItem)) }

// Pop remove and return element Len() - 1.
func (h *topKHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Top returns the top element of heap
func (h topKHeap) top() *topKFuzzySearchItem { return h[0] }

type TopKSelector struct {
	topK int
	h    topKHeap
}

// NewTopKSelector returns instance of TopKSelector
func NewTopKSelector(topK int) *TopKSelector {
	return &TopKSelector{
		topK: topK,
		h:    make(topKHeap, 0, topK),
	}
}

// use heap search for finding top k items in a list efficiently
// see http://stevehanov.ca/blog/index.php?id=122
// Add adds item with given position and distance to collection if item belongs to `top k items`
func (c *TopKSelector) Add(candidate *list_merger.MergeCandidate, distance float64) {
	if c.h.Len() >= c.topK && c.h.top().distance <= distance {
		return
	}

	var r *topKFuzzySearchItem

	if c.h.Len() == c.topK {
		r = heap.Pop(&c.h).(*topKFuzzySearchItem)
	} else {
		r = &topKFuzzySearchItem{
			candidate: nil,
			distance:  0.0,
		}
	}

	r.candidate = candidate
	r.distance = distance
	heap.Push(&c.h, r)
}

//
func (c *TopKSelector) GetLowestRecord() (*list_merger.MergeCandidate, float64) {
	if c.h.Len() > 0 {
		top := c.h.top()
		return top.candidate, top.distance
	}

	return nil, 1
}

//
func (c *TopKSelector) Size() int {
	return c.h.Len()
}

// GetCandidates returns `top k items` (on given moment)
func (c *TopKSelector) GetCandidates() []FuzzyCandidate {
	result := make([]FuzzyCandidate, 0, c.topK)

	for c.h.Len() > 0 {
		r := heap.Pop(&c.h).(*topKFuzzySearchItem)
		result = append(
			[]FuzzyCandidate{{Candidate{r.candidate.Position}, r.distance}},
			result...,
		)
	}

	return result
}

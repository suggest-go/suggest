package suggest

import (
	"container/heap"
	"github.com/alldroll/suggest/index"
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
	position index.Position
	distance float64
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

type TopKCollector struct {
	topK int
	h    topKHeap
}

// NewTopKCollector returns instance of TopKCollector
func NewTopKCollector(topK int) *TopKCollector {
	return &TopKCollector{
		topK: topK,
		h:    make(topKHeap, 0, topK),
	}
}

// use heap search for finding top k items in a list efficiently
// see http://stevehanov.ca/blog/index.php?id=122
// Add adds item with given position and distance to collection if item belongs to `top k items`
func (c *TopKCollector) Add(position index.Position, distance float64) {
	if c.h.Len() >= c.topK && c.h.top().distance <= distance {
		return
	}

	var r *topKFuzzySearchItem

	if c.h.Len() == c.topK {
		r = heap.Pop(&c.h).(*topKFuzzySearchItem)
	} else {
		r = &topKFuzzySearchItem{
			position: 0,
			distance: 0.0,
		}
	}

	r.position = position
	r.distance = distance
	heap.Push(&c.h, r)
}

// GetCandidates returns `top k items` (on given moment)
func (c *TopKCollector) GetCandidates() []FuzzyCandidate {
	result := make([]FuzzyCandidate, 0, c.topK)

	for c.h.Len() > 0 {
		r := heap.Pop(&c.h).(*topKFuzzySearchItem)
		result = append(
			[]FuzzyCandidate{{Candidate{r.position}, r.distance}},
			result...,
		)
	}

	return result
}

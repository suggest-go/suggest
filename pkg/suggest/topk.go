package suggest

import (
	"container/heap"
	"math"

	"github.com/suggest-go/suggest/pkg/index"
)

// TopKQueue is an accumulator that selects the "top k" elements added to it
type TopKQueue interface {
	// Add adds item with given position and distance to collection if item belongs to `top k items`
	Add(candidate index.Position, score float64)
	// GetLowestScore returns the lowest score of the collected candidates. If collection is empty, 0 will be returned
	GetLowestScore() float64
	// CanTakeWithScore returns true if a candidate with the given score can be accepted
	CanTakeWithScore(score float64) bool
	// IsFull tells if selector has collected `top k elements`
	IsFull() bool
	// GetCandidates returns `top k items`
	GetCandidates() []Candidate
	// Merge merges the given queue with the current
	Merge(other TopKQueue)
	// Reset resets the given queue with the provided topK
	Reset(topK int)
}

// topKHeap implements heap.Interface
type topKHeap []Candidate

// Len is the number of elements in the collection.
func (h topKHeap) Len() int { return len(h) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (h topKHeap) Less(i, j int) bool {
	return h[i].Less(h[j])
}

// Swap swaps the elements with indexes i and j.
func (h topKHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push add x as an element Len()
func (h *topKHeap) Push(x interface{}) { *h = append(*h, x.(Candidate)) }

// Pop remove and return element Len() - 1.
func (h *topKHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]

	return x
}

// top returns a pointer to the top element of the heap
func (h topKHeap) top() Candidate { return h[0] }

// top returns a pointer to the top element of the heap
func (h topKHeap) updateTop(candidate Candidate) {
	h[0] = candidate
	heap.Fix(&h, 0)
}

// topKQueue implements TopKQueue interface
type topKQueue struct {
	topK int
	h    topKHeap
}

// NewTopKQueue returns instance of TopKQueue
func NewTopKQueue(topK int) TopKQueue {
	return &topKQueue{
		topK: topK,
		h:    make(topKHeap, 0, topK),
	}
}

// Add adds item with given position and distance to collection if item belongs to `top k items`
// use heap search for finding top k items in a list efficiently
// see http://stevehanov.ca/blog/index.php?id=122
func (c *topKQueue) Add(position index.Position, score float64) {
	if !c.CanTakeWithScore(score) {
		return
	}

	candidate := Candidate{
		Key:   position,
		Score: score,
	}

	if c.h.Len() < c.topK {
		heap.Push(&c.h, candidate)

		return
	}

	if c.h.top().Less(candidate) {
		c.h.updateTop(candidate)
	}
}

// GetLowestScore returns the lowest score of the collected candidates
func (c *topKQueue) GetLowestScore() float64 {
	if c.h.Len() > 0 {
		return c.h.top().Score
	}

	return math.Inf(-1)
}

// CanTakeWithScore returns true if a candidate with the given score can be accepted
func (c *topKQueue) CanTakeWithScore(score float64) bool {
	if !c.IsFull() {
		return true
	}

	return c.h.top().Score <= score
}

// IsFull tells if selector has collected topK elements
func (c *topKQueue) IsFull() bool {
	return c.h.Len() == c.topK
}

// GetCandidates returns `top k items`
func (c *topKQueue) GetCandidates() []Candidate {
	if c.h.Len() == 0 {
		return []Candidate{}
	}

	sorted := make(topKHeap, c.h.Len())

	for c.h.Len() > 0 {
		sorted[c.h.Len()-1] = heap.Pop(&c.h).(Candidate)
	}

	// restore the order of the heap
	c.h = c.h[:len(sorted)]

	for i := len(c.h)/2 - 1; i >= 0; i-- {
		opp := len(c.h) - 1 - i
		c.h.Swap(i, opp)
	}

	return sorted
}

// Merge merges the given queue with the current
func (c *topKQueue) Merge(other TopKQueue) {
	topK, ok := other.(*topKQueue)

	if ok {
		for _, item := range topK.h {
			c.Add(item.Key, item.Score)
		}

		return
	}

	for _, item := range other.GetCandidates() {
		c.Add(item.Key, item.Score)
	}
}

// Reset resets the given queue with the provided topK
func (c *topKQueue) Reset(topK int) {
	c.topK = topK

	if cap(c.h) < topK {
		c.h = make(topKHeap, 0, topK)
	}

	c.h = c.h[:0]
}

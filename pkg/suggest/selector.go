package suggest

import (
	"container/heap"
	"math"
	"sort"

	"github.com/alldroll/suggest/pkg/index"
)

// TopKSelector is an accumulator that selects the "top" k elements added to it
type TopKSelector interface {
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
}

// Candidate is an item of TopKCollector collection
type Candidate struct {
	// Key is a position (docId) in posting list
	Key index.Position
	// Score is a float64 number from [0, 1]
	Score float64
}

// topKHeap implements heap.Interface
type topKHeap []Candidate

// Len is the number of elements in the collection.
func (h topKHeap) Len() int { return len(h) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (h topKHeap) Less(i, j int) bool { return h[i].Score < h[j].Score }

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
func (h topKHeap) top() *Candidate { return &h[0] }

// topKSelector implements TopKSelector interface
type topKSelector struct {
	topK int
	h    topKHeap
	rank Rank
}

// NewTopKSelector returns instance of TopKSelector
func NewTopKSelector(topK int) TopKSelector {
	return &topKSelector{
		topK: topK,
		h:    make(topKHeap, 0, topK),
		rank: &idOrderRank{},
	}
}

// NewTopKSelectorWithRanker returns instance of TopKSelector with ranker
func NewTopKSelectorWithRanker(topK int, rank Rank) TopKSelector {
	return &topKSelector{
		topK: topK,
		h:    make(topKHeap, 0, topK),
		rank: rank,
	}
}

// Add adds item with given position and distance to collection if item belongs to `top k items`
// use heap search for finding top k items in a list efficiently
// see http://stevehanov.ca/blog/index.php?id=122
func (c *topKSelector) Add(candidate index.Position, score float64) {
	if !c.CanTakeWithScore(score) {
		return
	}

	if c.h.Len() < c.topK {
		heap.Push(&c.h, Candidate{
			Key:   candidate,
			Score: score,
		})
		return
	}

	top := c.h.top()

	if top.Score < score || c.rank.Less(top.Key, candidate) {
		top.Key = candidate
		top.Score = score
		heap.Fix(&c.h, 0)
	}
}

// GetLowestScore returns the lowest score of the collected candidates
func (c *topKSelector) GetLowestScore() float64 {
	if c.h.Len() > 0 {
		return c.h.top().Score
	}

	return math.Inf(-1)
}

// CanTakeWithScore returns true if a candidate with the given score can be accepted
func (c *topKSelector) CanTakeWithScore(score float64) bool {
	if !c.IsFull() {
		return true
	}

	return c.h.top().Score < score
}

// IsFull tells if selector has collected topK elements
func (c *topKSelector) IsFull() bool {
	return c.h.Len() == c.topK
}

// GetCandidates returns `top k items`
func (c *topKSelector) GetCandidates() []Candidate {
	if c.h.Len() == 0 {
		return []Candidate{}
	}

	sorted := make(topKHeap, c.h.Len())
	copy(sorted, c.h)
	sort.Sort(sort.Reverse(sorted))

	return sorted
}

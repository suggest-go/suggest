package merger

import "container/heap"

type record struct {
	ridID    int
	position uint32
}

type recordHeap []*record

// Len is the number of elements in the collection.
func (h recordHeap) Len() int { return len(h) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (h recordHeap) Less(i, j int) bool { return h[i].position < h[j].position }

// Swap swaps the elements with indexes i and j.
func (h recordHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push add x as element Len()
func (h *recordHeap) Push(x interface{}) { *h = append(*h, x.(*record)) }

// Pop remove and return element Len() - 1.
func (h *recordHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// top returns the top element of heap
func (h recordHeap) top() *record { return h[0] }

// MergeSkip was described in paper
// "Efficient Merging and Filtering Algorithms for Approximate String Searches"
// Formally, main idea is to skip on the lists those record ids that cannot be in
// the answer to the query, by utilizing the threshold
func MergeSkip() ListMerger {
	return &mergeSkip{}
}

// mergeSkip implements MergeSkip algorithm
type mergeSkip struct{}

// Merge returns list of candidates, that appears at least `threshold` times.
func (ms *mergeSkip) Merge(rid Rid, threshold int) ([]MergeCandidate, error) {
	var (
		lenRid      = len(rid)
		h           = make(recordHeap, 0, lenRid)
		poppedItems = make([]*record, 0, lenRid)
		tops        = make([]record, lenRid)
		result      = make([]MergeCandidate, 0, lenRid)
		item        *record
	)

	for i := 0; i < lenRid; i++ {
		item = &tops[i]
		r, err := rid[i].Get()

		if err != nil && err != ErrIteratorIsNotDereferencable {
			return nil, err
		}

		item.ridID, item.position = i, r
		h.Push(item)
	}

	heap.Init(&h)
	item = nil

	for h.Len() > 0 {
		// reset slice
		poppedItems = poppedItems[:0]
		t := h.top()

		for h.Len() > 0 && t.position >= h.top().position {
			item = heap.Pop(&h).(*record)
			poppedItems = append(poppedItems, item)
		}

		n := len(poppedItems)

		if n >= threshold {
			result = append(result, MergeCandidate{
				Position: t.position,
				Overlap:  n,
			})

			for _, item := range poppedItems {
				cur := rid[item.ridID]

				if cur.HasNext() {
					r, err := cur.Next()

					if err != nil {
						return nil, err
					}

					item.position = r
					heap.Push(&h, item)
				}
			}
		} else {
			for j := threshold - 1 - n; j > 0 && h.Len() > 0; j-- {
				item = heap.Pop(&h).(*record)
				poppedItems = append(poppedItems, item)
			}

			if h.Len() == 0 {
				break
			}

			topPos := h.top().position

			for _, item := range poppedItems {
				cur := rid[item.ridID]

				if cur.Len() == 0 {
					continue
				}

				r, err := cur.LowerBound(topPos)

				if err != nil && err != ErrIteratorIsNotDereferencable {
					return nil, err
				}

				if err == nil {
					item.position = r
					heap.Push(&h, item)
				}
			}
		}
	}

	return result, nil
}

package merger

import "container/heap"

type record struct {
	iterator ListIterator
	position uint32
}

type recordHeap struct {
	slice []record
	size  int
}

// Len is the number of elements in the collection.
func (h recordHeap) Len() int { return h.size }

// Less reports whether the element with
// index i should sort before the element with index j.
func (h recordHeap) Less(i, j int) bool { return h.slice[i].position < h.slice[j].position }

// Swap swaps the elements with indexes i and j.
func (h recordHeap) Swap(i, j int) { h.slice[i], h.slice[j] = h.slice[j], h.slice[i] }

// Push add x as element Len()
func (h *recordHeap) Push(x interface{}) {
	h.size++
}

// Pop remove and return element Len() - 1.
func (h *recordHeap) Pop() interface{} {
	h.size--

	return nil
}

// top returns the top element of heap
func (h recordHeap) top() record { return h.slice[0] }

// MergeSkip was described in paper
// "Efficient Merging and Filtering Algorithms for Approximate String Searches"
// Formally, main idea is to skip on the lists those record ids that cannot be in
// the answer to the query, by utilizing the threshold
func MergeSkip() ListMerger {
	return newMerger(&mergeSkip{})
}

// mergeSkip implements MergeSkip algorithm
type mergeSkip struct{}

// Merge returns list of candidates, that appears at least `threshold` times.
func (ms *mergeSkip) Merge(rid Rid, threshold int, collector Collector) error {
	var (
		lenRid = len(rid)
		h      = recordHeap{
			slice: make([]record, lenRid),
			size:  lenRid,
		}
		poppedItems = 0
	)

	for i := 0; i < lenRid; i++ {
		r, err := rid[i].Get()

		if err != nil && err != ErrIteratorIsNotDereferencable {
			return err
		}

		h.slice[i].iterator, h.slice[i].position = rid[i], r
	}

	heap.Init(&h)

	for h.Len() > 0 {
		poppedItems = 0
		t := h.top()

		for h.Len() > 0 && t.position >= h.top().position {
			_ = heap.Pop(&h)
			poppedItems++
		}

		if poppedItems >= threshold {
			err := collector.Collect(NewMergeCandidate(t.position, uint32(poppedItems)))

			if err == ErrCollectionTerminated {
				return nil
			}

			if err != nil {
				return err
			}

			start := h.size

			for _, item := range h.slice[start : start+poppedItems] {
				if item.iterator.HasNext() {
					r, err := item.iterator.Next()

					if err != nil {
						return err
					}

					h.slice[h.size].iterator = item.iterator
					h.slice[h.size].position = r
					heap.Push(&h, nil)
				}
			}
		} else {
			for j := threshold - 1 - poppedItems; j > 0 && h.Len() > 0; j-- {
				_ = heap.Pop(&h)
				poppedItems++
			}

			if h.Len() == 0 {
				break
			}

			topPos := h.top().position
			start := h.size

			for _, item := range h.slice[start : start+poppedItems] {
				if item.iterator.Len() == 0 {
					continue
				}

				r, err := item.iterator.LowerBound(topPos)

				if err != nil && err != ErrIteratorIsNotDereferencable {
					return err
				}

				if err == nil {
					h.slice[h.size].iterator = item.iterator
					h.slice[h.size].position = r
					heap.Push(&h, nil)
				}
			}
		}
	}

	return nil
}

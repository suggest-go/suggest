package suggest

import (
	"container/heap"
	"math"
	"sort"
)

// ListMerger solves `threshold`-occurrence problem:
// For given inverted lists find the set of strings ids, that appears at least
// `threshold` times.
type ListMerger interface {
	// Merge returns list of candidates, that appears at least `threshold` times.
	Merge(rid Rid, threshold int) []*MergeCandidate
}

// Rid represents inverted lists for ListMerger
type Rid []PostingList

// Len is the number of elements in the collection.
func (p Rid) Len() int { return len(p) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (p Rid) Less(i, j int) bool { return len(p[i]) < len(p[j]) }

// Swap swaps the elements with indexes i and j.
func (p Rid) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// MergeCandidate is result of merging Rid
type MergeCandidate struct {
	Pos     Position
	Overlap int
}

type record struct {
	ridID int
	pos Position
}

func (r *record) Less(other heapItem) bool { return r.pos < other.(*record).pos }

// scanCount scan the N inverted lists one by one.
// For each string id on each list, we increment the count
// corresponding to the string by 1. We report the string ids that
// appear at least `threshold` times on the lists.
type ScanCount struct {}

// Merge returns list of candidates, that appears at least `threshold` times.
func (lm *ScanCount) Merge(rid Rid, threshold int) []*MergeCandidate {
	size := len(rid)
	candidates := make([]*MergeCandidate, 0, size)
	tmp := make([]*MergeCandidate, 0, size)
	j, k, endMergeCandidate, endRid := 0, 0, 0, 0

	for _, list := range rid {
		j, k = 0, 0
		tmp = tmp[:0]
		endMergeCandidate, endRid = len(candidates), len(list)

		for j < endMergeCandidate || k < endRid {
			if j >= endMergeCandidate || (k < endRid && candidates[j].Pos > list[k]) {
				tmp = append(tmp, &MergeCandidate{list[k], 1})
				k++
			} else if k >= endRid || (j < endMergeCandidate && candidates[j].Pos < list[k]) {
				tmp = append(tmp, candidates[j])
				j++
			} else {
				candidates[j].Overlap++
				tmp = append(tmp, candidates[j])
				j++
				k++
			}
		}

		candidates, tmp = tmp, candidates
	}

	tmp = tmp[:0]

	for _, c := range candidates {
		if c.Overlap >= threshold {
			tmp = append(tmp, c)
		}
	}

	candidates = tmp
	return candidates
}

// CPMerge was described in paper
// "Simple and Efficient Algorithm for Approximate Dictionary Matching"
// inspired by https://github.com/chokkan/simstring
type CPMerge struct {}

// Merge returns list of candidates, that appears at least `threshold` times.
func (cp *CPMerge) Merge(rid Rid, threshold int) []*MergeCandidate {
	sort.Sort(rid)

	lenRid := len(rid)
	minQueries := lenRid - threshold + 1
	candidates := make([]*MergeCandidate, 0, lenRid)
	tmp := make([]*MergeCandidate, 0, lenRid)
	j, k, endMergeCandidate, endRid := 0, 0, 0, 0

	for _, list := range rid[:minQueries] {
		j, k = 0, 0
		tmp = tmp[:0]
		endMergeCandidate, endRid = len(candidates), len(list)

		for j < endMergeCandidate || k < endRid {
			if j >= endMergeCandidate || (k < endRid && candidates[j].Pos > list[k]) {
				tmp = append(tmp, &MergeCandidate{list[k], 1})
				k++
			} else if k >= endRid || (j < endMergeCandidate && candidates[j].Pos < list[k]) {
				tmp = append(tmp, candidates[j])
				j++
			} else {
				candidates[j].Overlap++
				tmp = append(tmp, candidates[j])
				j++
				k++
			}
		}

		candidates, tmp = tmp, candidates
	}

	if len(candidates) == 0 {
		return candidates
	}

	for i := minQueries; i < lenRid; i++ {
		tmp = tmp[:0]

		for _, c := range candidates {
			j := binarySearchLowerBound(rid[i], c.Pos)
			if j != -1 {
				if rid[i][j] == c.Pos {
					c.Overlap++
				}

				rid[i] = rid[i][j:]
			}

			if c.Overlap + (lenRid - i - 1) >= threshold {
				tmp = append(tmp, c)
			}
		}

		candidates, tmp = tmp, candidates

		if len(candidates) == 0 {
			break
		}
	}

	tmp = tmp[:0]

	for _, c := range candidates {
		if c.Overlap >= threshold {
			tmp = append(tmp, c)
		}
	}

	candidates = tmp
	return candidates
}


// DivideSkip was described in paper
// "Efficient Merging and Filtering Algorithms for Approximate String Searches"
// We have to choose `good` parameter mu, for improving speed. So, mu depends
// only on given dictionary, so we can find it
type DivideSkip struct { mu float64; merger ListMerger }

// Merge returns list of candidates, that appears at least `threshold` times.
func (ds *DivideSkip) Merge(rid Rid, threshold int) []*MergeCandidate {
	sort.Reverse(rid)

	M := float64(len(rid[0]))
	l := int(float64(threshold) / (ds.mu * math.Log(M) + 1))

	lLong := rid[:l]
	lShort := rid[l:]

	if len(lShort) == 0 {
		return ds.merger.Merge(rid, threshold)
	}

	var (
		r Position
		candidates = ds.merger.Merge(lShort, threshold - l)
		result = make([]*MergeCandidate, 0, len(candidates))
	)

	for _, c := range candidates {
		r = c.Pos

		for _, longList := range lLong {
			idx := binarySearchLowerBound(longList, r)
			if idx != -1 && longList[idx] == r {
				c.Overlap++
			}
		}

		if c.Overlap >= threshold {
			result = append(result, c)
		}
	}

	return result
}

// MergeSkip was described in paper
// "Efficient Merging and Filtering Algorithms for Approximate String Searches"
// Formally, main idea is to skip on the lists those record ids that cannot be in
// the answer to the query, by utilizing the threshold
type MergeSkip struct {}

// Merge returns list of candidates, that appears at least `threshold` times.
func (ms *MergeSkip) Merge(rid Rid, threshold int) []*MergeCandidate {
	lenRid := len(rid)
	h := newHeap(lenRid)
	poppedItems := make([]*record, 0, lenRid)
	tops := make([]record, lenRid)
	result := make([]*MergeCandidate, 0, lenRid)
	var item *record

	for i := 0; i < lenRid; i++ {
		item = &tops[i]
		item.ridID, item.pos = i, rid[i][0]
		h.Push(item)
	}

	heap.Init(h)
	item = nil

	for h.Len() > 0 {
		// reset slice
		poppedItems = poppedItems[:0]
		t := h.Top().(*record)
		for h.Len() > 0 && !t.Less(h.Top()) {
			item = heap.Pop(h).(*record)
			poppedItems = append(poppedItems, item)
		}

		n := len(poppedItems)
		if n >= threshold {
			result = append(result, &MergeCandidate{
				Pos: t.pos,
				Overlap: n,
			})

			for _, item := range poppedItems {
				cur := rid[item.ridID]
				if len(cur) > 1 {
					cur = cur[1:]
					rid[item.ridID] = cur
					item.pos = cur[0]
					heap.Push(h, item)
				}
			}
		} else {
			for j := threshold - 1 - n; j > 0 && h.Len() > 0; j-- {
				item = heap.Pop(h).(*record)
				poppedItems = append(poppedItems, item)
			}

			if h.Len() == 0 {
				break
			}

			topPos := h.Top().(*record).pos
			for _, item := range poppedItems {
				cur := rid[item.ridID]
				if len(cur) == 0 {
					continue
				}

				r := binarySearchLowerBound(cur, topPos)
				if r != -1 {
					cur = cur[r:]
					rid[item.ridID] = cur
					item.pos = cur[0]
					heap.Push(h, item)
				}
			}
		}
	}

	return result
}

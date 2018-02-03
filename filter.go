package suggest

// Algorithms given below solve `threshold`-occurrence problem:
// For given inverted lists find the set of strings ids, that appears at least
// `threshold` times.

import (
	"container/heap"
	"math"
	"sort"
)

type ListMerger interface {
	Merge(rid Rid, threshold int) []*MergeCandidate
}

type Rid []PostingList

func (p Rid) Len() int {
	return len(p)
}

func (p Rid) Less(i, j int) bool {
	return len(p[i]) < len(p[j])
}

func (p Rid) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type MergeCandidate struct {
	Pos     Position
	Overlap int
}

type record struct {
	ridID int
	pos Position
}

func (r *record) Less(other heapItem) bool {
	return r.pos < other.(*record).pos
}

// scanCount scan the N inverted lists one by one.
// For each string id on each list, we increment the count
// corresponding to the string by 1. We report the string ids that
// appear at least `threshold` times on the lists.
type ScanCount struct {}

func (lm *ScanCount) Merge(rid Rid, threshold int) []*MergeCandidate {
	size := len(rid)
	MergeCandidates := make([]*MergeCandidate, 0, size)
	tmp := make([]*MergeCandidate, 0, size)
	j, k, endMergeCandidate, endRid := 0, 0, 0, 0

	for _, list := range rid {
		j, k = 0, 0
		tmp = tmp[:0]
		endMergeCandidate, endRid = len(MergeCandidates), len(list)

		for j < endMergeCandidate || k < endRid {
			if j >= endMergeCandidate || (k < endRid && MergeCandidates[j].Pos > list[k]) {
				tmp = append(tmp, &MergeCandidate{list[k], 1})
				k++
			} else if k >= endRid || (j < endMergeCandidate && MergeCandidates[j].Pos < list[k]) {
				tmp = append(tmp, MergeCandidates[j])
				j++
			} else {
				MergeCandidates[j].Overlap++
				tmp = append(tmp, MergeCandidates[j])
				j++
				k++
			}
		}

		MergeCandidates, tmp = tmp, MergeCandidates
	}

	return MergeCandidates
}

type CPMerge struct {}

// cpMerge was described in paper
// "Simple and Efficient Algorithm for Approximate Dictionary Matching"
// inspired by https://github.com/chokkan/simstring
func (cp *CPMerge) Merge(rid Rid, threshold int) []*MergeCandidate {
	sort.Sort(rid)

	lenRid := len(rid)
	minQueries := lenRid - threshold + 1
	MergeCandidates := make([]*MergeCandidate, 0, lenRid)
	tmp := make([]*MergeCandidate, 0, lenRid)
	j, k, endMergeCandidate, endRid := 0, 0, 0, 0

	for _, list := range rid[:minQueries] {
		j, k = 0, 0
		tmp = tmp[:0]
		endMergeCandidate, endRid = len(MergeCandidates), len(list)

		for j < endMergeCandidate || k < endRid {
			if j >= endMergeCandidate || (k < endRid && MergeCandidates[j].Pos > list[k]) {
				tmp = append(tmp, &MergeCandidate{list[k], 1})
				k++
			} else if k >= endRid || (j < endMergeCandidate && MergeCandidates[j].Pos < list[k]) {
				tmp = append(tmp, MergeCandidates[j])
				j++
			} else {
				MergeCandidates[j].Overlap++
				tmp = append(tmp, MergeCandidates[j])
				j++
				k++
			}
		}

		MergeCandidates, tmp = tmp, MergeCandidates
	}

	if len(MergeCandidates) == 0 {
		return MergeCandidates
	}

	for i := minQueries; i < lenRid; i++ {
		tmp = tmp[:0]

		for _, c := range MergeCandidates {
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

		MergeCandidates, tmp = tmp, MergeCandidates

		if len(MergeCandidates) == 0 {
			break
		}
	}

	return MergeCandidates
}


// divideSkip was described in paper
// "Efficient Merging and Filtering Algorithms for Approximate String Searches"
// We have to choose `good` parameter mu, for improving speed. So, mu depends
// only on given dictionary, so we can find it
type DivideSkip struct { mu float64; merger ListMerger }

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
		MergeCandidates = ds.merger.Merge(lShort, threshold - l)
	)

	for _, c := range MergeCandidates {
		r = c.Pos

		for _, longList := range lLong {
			idx := binarySearchLowerBound(longList, r)
			if idx != -1 && longList[idx] == r {
				c.Overlap++
			}
		}
	}

	return MergeCandidates
}

// mergeSkip was described in paper
// "Efficient Merging and Filtering Algorithms for Approximate String Searches"
// Formally, main idea is to skip on the lists those record ids that cannot be in
// the answer to the query, by utilizing the threshold
type MergeSkip struct {}

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
			result = append(result, &MergeCandidate{t.pos, n})

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

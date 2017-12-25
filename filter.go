package suggest

// Algorithms given below solve `threshold`-occurrence problem:
// For given inverted lists find the set of strings ids, that appears at least
// `threshold` times.
// All filters returns [][]int - [intersection][corresponding string ids]

import (
	"container/heap"
	"math"
	"sort"
)

type record struct {
	ridID int
	pos   Position
}

func (r *record) Less(other heapItem) bool {
	return r.pos < other.(*record).pos
}

// scanCount scan the N inverted lists one by one.
// For each string id on each list, we increment the count
// corresponding to the string by 1. We report the string ids that
// appear at least `threshold` times on the lists.
func scanCount(rid []PostingList, threshold int) []PostingList {
	size := len(rid)
	result := make([]PostingList, size+1)

	mapSize := 0
	for i := range rid {
		mapSize += len(rid[i])
	}

	counts := make(map[Position]int, mapSize+1)
	for _, curRid := range rid {
		for _, pos := range curRid {
			counts[pos]++
		}
	}

	for pos, count := range counts {
		if count < threshold {
			continue
		}

		result[count] = append(result[count], pos)
	}

	return result
}

// cpMerge was described in paper
// "Simple and Efficient Algorithm for Approximate Dictionary Matching"
// inspired by https://github.com/chokkan/simstring
func cpMerge(rid []PostingList, threshold int) []PostingList {
	sort.SliceStable(rid, func(i, j int) bool {
		return len(rid[i]) < len(rid[j])
	})

	type candidate struct {
		pos Position
		overlap int
	}

	lenRid := len(rid)
	minQueries := lenRid - threshold + 1
	candidates := make([]*candidate, 0, lenRid)
	tmp := make([]*candidate, 0, lenRid)
	i, j, k, endCandidate, endRid := 0, 0, 0, 0, 0

	for ; i < minQueries; i++ {
		j, k = 0, 0
		tmp = tmp[:0]
		endCandidate, endRid = len(candidates), len(rid[i])

		for j < endCandidate || k < endRid {
			if j >= endCandidate || (k < endRid && candidates[j].pos > rid[i][k]) {
				tmp = append(tmp, &candidate{rid[i][k], 1})
				k++
			} else if k >= endRid || (j < endCandidate && candidates[j].pos < rid[i][k]) {
				tmp = append(tmp, candidates[j])
				j++
			} else {
				candidates[j].overlap++
				tmp = append(tmp, candidates[j])
				j++
				k++
			}
		}

		candidates, tmp = tmp, candidates
	}

	if len(candidates) == 0 {
		return nil
	}

	result := make([]PostingList, len(rid)+1)

	for ; i < lenRid; i++ {
		tmp = tmp[:0]
		j, k = 0, 0

		for _, c := range candidates {
			if binarySearch(rid[i], c.pos) != -1 {
				c.overlap++
			}

			// Modify algorithm: we should to know exact overlap count, so leave candidate
			/*
			if c.overlap >= threshold {
				result[c.overlap] = append(result[c.overlap], c.pos)
			}*/

			if c.overlap + (lenRid - i - 1) >= threshold {
				tmp = append(tmp, c)
			}
		}

		candidates, tmp = tmp, candidates

		if len(candidates) == 0 {
			break;
		}
	}

	if len(candidates) > 0 {
		for _, c := range candidates {
			if c.overlap >= threshold {
				result[c.overlap] = append(result[c.overlap], c.pos)
			}
		}
	}

	return result
}

// divideSkip was described in paper
// "Efficient Merging and Filtering Algorithms for Approximate String Searches"
// We have to choose `good` parameter mu, for improving speed. So, mu depends
// only on given dictionary, so we can find it
func divideSkip(rid []PostingList, threshold int, mu float64) []PostingList {
	sort.SliceStable(rid, func(i, j int) bool {
		return len(rid[i]) > len(rid[j])
	})

	M := float64(len(rid[0]))
	l := int(float64(threshold) / (mu*math.Log(M) + 1))

	lLong := rid[:l]
	lShort := rid[l:]
	result := make([]PostingList, len(rid)+1)

	if len(lShort) == 0 {
		return mergeSkip(rid, threshold)
	}

	for count, list := range mergeSkip(lShort, threshold-l) {
		for _, r := range list {
			j := count
			for _, longList := range lLong {
				idx := binarySearchUpperBound(longList, r)
				if idx != -1 && longList[idx] == r {
					j++
				}
			}

			if j >= threshold {
				result[j] = append(result[j], r)
			}
		}
	}

	return result
}

// mergeSkip was described in paper
// "Efficient Merging and Filtering Algorithms for Approximate String Searches"
// Formally, main idea is to skip on the lists those record ids that cannot be in
// the answer to the query, by utilizing the threshold
func mergeSkip(rid []PostingList, threshold int) []PostingList {
	lenRid := len(rid)
	h := newHeap(lenRid)
	result := make([]PostingList, lenRid+1)
	poppedItems := make([]*record, 0, lenRid)

	for i := 0; i < lenRid; i++ {
		heap.Push(h, &record{i, rid[i][0]})
	}

	for h.Len() > 0 {
		// reset slice
		poppedItems = poppedItems[:0]
		t := h.Top()
		for h.Len() > 0 && h.Top().(*record).pos == t.(*record).pos {
			item := heap.Pop(h)
			poppedItems = append(poppedItems, item.(*record))
		}

		n := len(poppedItems)
		if n >= threshold {
			result[n] = append(result[n], t.(*record).pos)
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
			j := threshold - 1 - n
			for j > 0 && h.Len() > 0 {
				item := heap.Pop(h)
				j--
				poppedItems = append(poppedItems, item.(*record))
			}

			if h.Len() == 0 {
				break
			}

			t = h.Top()
			for _, item := range poppedItems {
				cur := rid[item.ridID]
				if len(cur) == 0 {
					continue
				}

				r := binarySearchUpperBound(cur, t.(*record).pos)
				if r == -1 {
					continue
				}

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

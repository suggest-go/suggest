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
	ridID, strID int
}

func (r *record) Less(other heapItem) bool {
	return r.strID < other.(*record).strID
}

// scanCount scan the N inverted lists one by one.
// For each string id on each list, we increment the count
// corresponding to the string by 1. We report the string ids that
// appear at least `threshold` times on the lists.
func scanCount(rid [][]int, threshold int) [][]int {
	size := len(rid)
	result := make([][]int, size+1)

	mapSize := 0
	for i := range rid {
		mapSize += len(rid[i])
	}

	counts := make(map[int]int, mapSize+1)
	for _, curRid := range rid {
		for _, strID := range curRid {
			counts[strID]++
		}
	}

	for strID, count := range counts {
		if count < threshold {
			continue
		}

		result[count] = append(result[count], strID)
	}

	return result
}

// cpMerge was described in paper
// "Simple and Efficient Algorithm for Approximate Dictionary Matching"
func cpMerge(rid [][]int, threshold int) [][]int {
	panic("Implement me")
	return nil
}

// divideSkip was described in paper
// "Efficient Merging and Filtering Algorithms for Approximate String Searches"
// We have to choose `good` parameter mu, for improving speed. So, mu depends
// only on given dictionary, so we can find it
func divideSkip(rid [][]int, threshold int, mu float64) [][]int {
	sort.Slice(rid, func(i, j int) bool {
		return len(rid[i]) > len(rid[j])
	})

	M := float64(len(rid[0]))
	l := int(float64(threshold) / (mu*math.Log(M) + 1))

	lLong := rid[:l]
	lShort := rid[l:]
	result := make([][]int, len(rid)+1)

	if len(lShort) == 0 {
		return mergeSkip(rid, threshold)
	}

	for count, list := range mergeSkip(lShort, threshold-l) {
		for _, r := range list {
			j := count
			for _, longList := range lLong {
				idx := binarySearch(longList, r)
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
func mergeSkip(rid [][]int, threshold int) [][]int {
	h := &heapImpl{}
	lenRid := len(rid)
	result := make([][]int, lenRid+1)

	for i := 0; i < lenRid; i++ {
		heap.Push(h, &record{i, rid[i][0]})
	}

	poppedItems := make([]*record, 0, lenRid)
	for h.Len() > 0 {
		// reset slice
		poppedItems = poppedItems[:0]
		t := h.Top()
		for h.Len() > 0 && h.Top().(*record).strID == t.(*record).strID {
			item := heap.Pop(h)
			poppedItems = append(poppedItems, item.(*record))
		}

		n := len(poppedItems)
		if n >= threshold {
			result[n] = append(result[n], t.(*record).strID)
			for _, item := range poppedItems {
				cur := rid[item.ridID]
				if len(cur) > 1 {
					cur = cur[1:]
					rid[item.ridID] = cur
					item.strID = cur[0]
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

				r := binarySearch(cur, t.(*record).strID)
				if r == -1 {
					continue
				}

				if r != -1 {
					cur = cur[r:]
					rid[item.ridID] = cur
					item.strID = cur[0]
					heap.Push(h, item)
				}
			}
		}
	}

	return result
}

// binarySearch find the smallest record t in given arr such that t >= value
func binarySearch(arr []int, value int) int {
	i := 0
	j := len(arr)
	if i == j || arr[j-1] < value {
		return -1
	}

	if arr[i] >= value {
		return i
	}

	if arr[j-1] == value {
		return j - 1
	}

	for i < j {
		mid := i + (j-i)>>1
		if arr[mid] < value {
			i = mid + 1
		} else if arr[mid] > value {
			j = mid - 1
		} else {
			return mid
		}
	}

	if i > len(arr)-1 {
		return -1
	}

	if j < 0 {
		return 0
	}

	if arr[i] >= value {
		return i
	}

	return j + 1
}

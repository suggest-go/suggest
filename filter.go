package suggest

import (
	"container/heap"
	"math"
	"sort"
)

type record struct {
	ridId, strId int
}

func (self *record) Less(other heapItem) bool {
	return self.strId < other.(*record).strId
}

func cpMerge(rid [][]int, threshold int) [][]int {
	sort.Slice(rid, func(i, j int) bool {
		return len(rid[i]) < len(rid[j])
	})

	size := len(rid)
	result := make([][]int, size+1)
	mapSize := 0
	k := size - threshold
	for i, _ := range rid {
		if i >= k {
			break
		}

		mapSize += len(rid[i])
	}

	counts := make(map[int]int, mapSize+1)
	i := 0
	for ; i < k; i++ {
		for _, strId := range rid[i] {
			counts[strId]++
		}
	}

	for strId, count := range counts {
		for j := i; j < size; j++ {
			idx := binarySearch(rid[j], 0, strId)
			if idx != -1 && rid[j][idx] == strId {
				count++
			}
		}

		if count >= threshold {
			result[count] = append(result[count], strId)
		}
	}

	return result
}

// TODO придумать как выбирать параметр mu
func divideSkip(rid [][]int, threshold int, mu float64) [][]int {
	sort.Slice(rid, func(i, j int) bool {
		return len(rid[i]) > len(rid[j])
	})

	M := float64(len(rid[0]))
	l := int(float64(threshold) / (mu*math.Log(M) + 1))

	lLong := rid[:l]
	lShort := rid[l:]
	result := make([][]int, len(rid)+1)

	for count, list := range mergeSkip(lShort, threshold-l) {
		for _, r := range list {
			j := count
			for _, longList := range lLong {
				idx := binarySearch(longList, 0, r)
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

// see Efficient Merging and Filtering Algorithms for
// Approximate String Searches
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
		for h.Len() > 0 && h.Top().(*record).strId == t.(*record).strId {
			item := heap.Pop(h)
			poppedItems = append(poppedItems, item.(*record))
		}

		n := len(poppedItems)
		if n >= threshold {
			result[n] = append(result[n], t.(*record).strId)
			for _, item := range poppedItems {
				cur := rid[item.ridId]
				if len(cur) > 1 {
					cur = cur[1:]
					rid[item.ridId] = cur
					item.strId = cur[0]
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
				cur := rid[item.ridId]
				if len(cur) == 0 {
					continue
				}

				r := binarySearch(cur, 0, t.(*record).strId)
				if r == -1 {
					continue
				}

				if r != -1 {
					cur = cur[r:]
					rid[item.ridId] = cur
					item.strId = cur[0]
					heap.Push(h, item)
				}
			}
		}
	}

	return result
}

func binarySearch(arr []int, i int, value int) int {
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

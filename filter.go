package suggest

import (
	"container/heap"
	//"log"
	"math"
	"sort"
)

type record struct {
	ridId, strId int
}

func (self *record) Less(other heapItem) bool {
	return self.strId < other.(*record).strId
}

func divideSkip(rid [][]int, threshold int) map[int][]int {
	sort.Slice(rid, func(i, j int) bool {
		return len(rid[i]) > len(rid[j])
	})

	m := len(rid[0])
	mu := 0.0085
	l := int(float64(threshold) / (mu*math.Log2(float64(m)) + 1))
	lLong := rid[:l]
	lShort := rid[l:]

	//	log.Printf("%v, %v\n", lShort, lLong)

	result := make(map[int][]int)
	for count, list := range mergeSkip(lShort, threshold-l) {
		//log.Printf("list: %v\n", list)
		for _, r := range list {
			j := count
			for _, longList := range lLong {
				idx := binarySearch(longList, 0, r)
				if idx != -1 && longList[idx] == r {
					j++
				}
			}

			result[j] = append(result[j], r)
			//log.Printf("count: %d, val: %d\n", j, r)
		}
	}

	return result
}

// see Efficient Merging and Filtering Algorithms for
// Approximate String Searches

func mergeSkip(rid [][]int, threshold int) map[int][]int {
	iters := make([]int, len(rid))
	h := &heapImpl{}
	result := make(map[int][]int)

	for i, iter := range iters {
		heap.Push(h, &record{i, rid[i][iter]})
	}

	for h.Len() > 0 {
		t := h.Top()
		poppedItems := make([]*record, 0, len(rid))
		for h.Len() > 0 && h.Top().(*record).strId == t.(*record).strId {
			item := heap.Pop(h)
			poppedItems = append(poppedItems, item.(*record))
		}

		n := len(poppedItems)
		if n >= threshold {
			result[n] = append(result[n], t.(*record).strId)
			for _, item := range poppedItems {
				iters[item.ridId]++
				if iters[item.ridId] < len(rid[item.ridId]) {
					item.strId = rid[item.ridId][iters[item.ridId]]
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
				for _, item := range poppedItems {
					iters[item.ridId]++
					if iters[item.ridId] < len(rid[item.ridId]) {
						item.strId = rid[item.ridId][iters[item.ridId]]
						heap.Push(h, item)
					}
				}
				continue
			}

			t = h.Top()
			for _, item := range poppedItems {
				i := item.ridId
				if len(rid[i]) <= iters[i] {
					continue
				}

				r := binarySearch(rid[i], iters[i], t.(*record).strId)
				if r == -1 {
					continue
				}

				iters[i] = r
				val := rid[i][r]
				if r != -1 {
					item.strId = val
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

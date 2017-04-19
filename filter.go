package suggest

import (
	"container/heap"
	//"log"
	"math"
)

type record struct {
	ridId, strId int
}

func (self *record) Less(other heapItem) bool {
	return self.strId < other.(*record).strId
}

func divideSkip(rid [][]int, threshold int) map[int][]int {
	m := len(rid[0])
	mu := 0.0085
	l := int(float64(threshold) / (mu*math.Log(float64(m)) + 1))

	if l <= 0 {
		return mergeSkip(rid, threshold)
	}

	lLong := rid[:l]
	lShort := rid[l:]

	//	log.Printf("%v, %v\n", lShort, lLong)
	result := make(map[int][]int, len(rid)+1)
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

			if j >= threshold {
				result[j] = append(result[j], r)
			}
			//log.Printf("count: %d, val: %d\n", j, r)
		}
	}

	return result
}

// see Efficient Merging and Filtering Algorithms for
// Approximate String Searches

func mergeSkip(rid [][]int, threshold int) map[int][]int {
	h := &heapImpl{}
	result := make(map[int][]int, len(rid)+1)
	iters := make([]int, len(rid))

	for i, iter := range iters {
		heap.Push(h, &record{i, rid[i][iter]})
	}

	poppedItems := make([]*record, len(rid))
	for h.Len() > 0 {
		poppedIter := 0
		t := h.Top()
		for h.Len() > 0 && h.Top().(*record).strId == t.(*record).strId {
			item := heap.Pop(h)
			poppedItems[poppedIter] = item.(*record)
			poppedIter++
		}

		n := poppedIter
		if n >= threshold {
			result[n] = append(result[n], t.(*record).strId)
			for j := 0; j < poppedIter; j++ {
				item := poppedItems[j]
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
				poppedItems[poppedIter] = item.(*record)
				poppedIter++
			}

			if h.Len() == 0 {
				break
			}

			t = h.Top()
			for j := 0; j < poppedIter; j++ {
				item := poppedItems[j]
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

package suggest

import (
	"container/heap"
	"fmt"
)

type record struct {
	ridId, strId int
}

func (self *record) Less(other heapItem) bool {
	return self.strId < other.(*record).strId
}

func (self *record) String() string {
	return fmt.Sprintf("{strId: %d, ridId: %d}", self.strId, self.ridId)
}

// see Efficient Merging and Filtering Algorithms for
// Approximate String Searches

func mergeSkip(rid [][]int, threshold int) []int {
	iters := make([]int, len(rid))
	h := &heapImpl{}
	result := make([]int, 0)

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
			result = append(result, t.(*record).strId)
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

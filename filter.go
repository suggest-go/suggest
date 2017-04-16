package suggest

import (
	"container/heap"
	"fmt"
	"log"
	//"time"
)

type record struct {
	ridId, recId int
}

func (self *record) Less(other heapItem) bool {
	return self.recId < other.(*record).recId
}

func (self *record) String() string {
	return fmt.Sprintf("{recId: %d, ridId: %d}", self.recId, self.ridId)
}

// see Efficient Merging and Filtering Algorithms for
// Approximate String Searches

func mergeSkip(rid [][]int, threshold int) []int {
	iters := make([]int, len(rid))
	h := &heapImpl{}
	result := make([]int, 0)

	for i, iter := range iters {
		heap.Push(h, &record{i, rid[i][iter]})
		iters[i]++
	}

	for h.Len() > 0 {
		log.Printf("HEAP0 %v\n", h)
		t := h.Top()
		log.Printf("TOP %v\n", t)
		poppedItems := make([]*record, 0, len(rid))
		for h.Len() > 0 && h.Top().(*record).recId == t.(*record).recId {
			item := heap.Pop(h)
			poppedItems = append(poppedItems, item.(*record))
		}

		n := len(poppedItems)
		log.Printf("HEAP1 %v\n", h)
		log.Printf("ITERS %v\n", iters)
		log.Printf("N %v\n", n)
		if n >= threshold {
			result = append(result, t.(*record).recId)

			for _, item := range poppedItems {
				if iters[item.ridId] < len(rid[item.ridId]) {
					item.recId = rid[item.ridId][iters[item.ridId]]
					heap.Push(h, item)
					iters[item.ridId]++
				}
			}
		} else {
			j := threshold - 1 - n
			log.Printf("J %v\n", j)
			for j > 0 && h.Len() > 0 {
				item := heap.Pop(h)
				j--
				poppedItems = append(poppedItems, item.(*record))
			}

			log.Printf("HEAP2 %v\n", h)

			if h.Len() == 0 {
				for _, item := range poppedItems {
					if iters[item.ridId] < len(rid[item.ridId]) {
						item.recId = rid[item.ridId][iters[item.ridId]]
						heap.Push(h, item)
						iters[item.ridId]++
					}
				}
				continue
			}

			t = h.Top()
			log.Printf("T1 %v\n", t)
			log.Printf("POPPED %v\n", poppedItems)
			for _, item := range poppedItems {
				log.Printf("ITEM %v\n", item)
				i := item.ridId
				if len(rid[i]) <= iters[i] {
					log.Printf("SKIP %d\n", i)
					continue
				}

				r := binarySearch(rid[i], iters[i], item.recId)
				if r == -1 {
					continue
				}

				iters[i] = r + 1
				val := rid[i][r]
				log.Printf("R %v\n", val)
				if r != -1 {
					item.recId = val
					heap.Push(h, item)
				}
			}
		}

		//time.Sleep(3000 * time.Millisecond)
	}

	log.Printf("%v\n", result)

	return result
}

/*
func mergeSkip(rid [][]int, threshold int) []int {
	iters := make([]int, len(rid))
	h := &heapImpl{}
	result := make([]int, 0)

	for i, iter := range iters {
		heap.Push(h, (record)(rid[i][iter]))
		iters[i]++
	}

	for h.Len() > 0 {
		log.Printf("HEAP0 %v\n", h)
		t := h.Top()
		log.Printf("TOP %v\n", t)
		n := 0
		for h.Len() > 0 && h.Top() == t {
			n++
			heap.Pop(h)
		}

		log.Printf("HEAP1 %v\n", h)
		log.Printf("ITERS %v\n", iters)
		log.Printf("N %v\n", n)
		if n >= threshold {
			result = append(result, int(t.(record)))
			for i, iter := range iters {
				iter++
				if iter < len(rid[i]) {
					iters[i] = iter
					heap.Push(h, record(rid[i][iter]))
				}
			}
		} else {
			j := threshold - 1 - n
			log.Printf("J %v\n", j)
			for j > 0 && h.Len() > 0 {
				heap.Pop(h)
				j--
			}

			log.Printf("HEAP2 %v\n", h)

			if h.Len() == 0 {
				log.Printf("I WAS HERE\n")
				break
			}

			t = h.Top()
			log.Printf("T1 %v\n", t)
			j = 0
			for i, iter := range iters {
				if len(rid[i]) <= iter {
					log.Printf("SKIP %d\n", i)
					continue
				}

				if j == threshold-1 {
					break
				}

				j++

				r := binarySearch(rid[i], iter, int(t.(record)))
				if r == -1 {
					continue
				}

				iters[i] = r + 1
				val := rid[i][r]
				log.Printf("R %v\n", val)
				if r != -1 {
					heap.Push(h, record(val))
				}
			}
		}

		//time.Sleep(3000 * time.Millisecond)
	}

	log.Printf("%v\n", result)

	return result
}
*/

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

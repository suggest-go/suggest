package suggest

const (
	exactSearch = 0
	upperBoundSearch = 1
)

func binarySearchImpl(arr PostingList, value Position, mode int) int {
	i := 0
	j := len(arr)
	if i == j || arr[j-1] < value {
		return -1
	}

	if (mode == upperBoundSearch && arr[i] >= value) || arr[i] == value {
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

	if mode == upperBoundSearch {
		if arr[i] >= value {
			return i
		}

		return j + 1
	}

	if arr[i] == value {
		return i
	}

	if arr[j] == value {
		return j
	}

	return -1
}

// upperBound find the smallest record t in given arr such that t >= value
func binarySearchUpperBound(arr PostingList, value Position) int {
	return binarySearchImpl(arr, value, upperBoundSearch)
}

func binarySearch(arr PostingList, value Position) int {
	return binarySearchImpl(arr, value, exactSearch)
}

package suggest

// binarySearchLowerBound find index for the smallest record t in given arr such that t >= value
func binarySearchLowerBound(arr PostingList, value Position) int {
	i := 0
	j := len(arr)
	mid := 0
	midVal := Position(0)

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
		mid = i + (j-i)>>1
		midVal = arr[mid]

		if midVal >= value {
			j = mid
		} else {
			i = mid + 1
		}
	}

	return i
}

// binarySearch find index for given value, returns -1 if values is not in arr
func binarySearch(arr PostingList, value Position) int {
	i := 0
	j := len(arr)
	mid := 0
	midVal := Position(0)

	if i == j || arr[j-1] < value {
		return -1
	}

	if arr[i] == value {
		return i
	}

	if arr[j-1] == value {
		return j - 1
	}

	for i <= j {
		mid = i + (j-i)>>1
		midVal = arr[mid]

		if midVal < value {
			i = mid + 1
		} else if midVal > value {
			j = mid - 1
		} else {
			return mid
		}
	}

	return -1
}

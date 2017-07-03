package suggest

const maxN = 8

// SplitIntoNGrams split given word on k-gram list
// inspired by https://github.com/Lazin/go-ngram
func SplitIntoNGrams(word string, k int) []string {
	sliceLen := len(word) - k + 1
	if sliceLen <= 0 || sliceLen > len(word) {
		panic("Invalid word length for spliting")
	}

	var prevIndexes [maxN]int
	result := make([]string, 0, sliceLen)
	i := 0
	for index := range word {
		i++
		if i > k {
			top := prevIndexes[(i-k)%k]
			substr := word[top:index]
			result = append(result, substr)
		}

		prevIndexes[i%k] = index
	}

	top := prevIndexes[(i+1)%k]
	substr := word[top:]
	result = append(result, substr)

	return result
}

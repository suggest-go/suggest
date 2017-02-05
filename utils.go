package suggest

import (
	"strings"
)

func SplitIntoNGrams(word string, k int) []string {
	sliceLen := len(word) - k + 1
	if sliceLen <= 0 || sliceLen > len(word) {
		panic("Invalid word length for spliting")
	}

	result := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		result[i] = word[i : i+k]
	}

	return result
}

func Levenshtein(a, b string) int {
	aLen, bLen := len(a), len(b)
	if aLen == 0 {
		return bLen
	}

	if bLen == 0 {
		return aLen
	}

	r1, r2 := []rune(a), []rune(b)
	column := make([]int, aLen+1)

	for i := 1; i < aLen+1; i++ {
		column[i] = i
	}

	for j := 1; j < bLen+1; j++ {
		column[0] = j
		prev := j - 1
		for i := 1; i < aLen+1; i++ {
			tmp := column[i]
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}

			column[i] = min3(
				column[i]+1,
				column[i-1]+1,
				prev+cost,
			)
			prev = tmp
		}
	}

	return column[aLen]
}

func prepareString(word string) string {
	word = strings.ToLower(word)
	word = strings.Trim(word, " ")
	word = strings.Replace(word, " ", "$", -1)
	return "$" + word + "$"
}

func min3(a, b, c int) int {
	if a < b && a < c {
		return a
	}

	if b < c {
		return b
	}

	return c
}

package suggest

import (
	"regexp"
	"strings"
)

var reg *regexp.Regexp

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

/*
 * Return unique ngrams with frequency
 */
func getProfile(word string, k int) map[string]int {
	ngrams := SplitIntoNGrams(word, k)
	result := make(map[string]int)
	for _, ngram := range ngrams {
		if _, ok := result[ngram]; ok {
			result[ngram]++
		} else {
			result[ngram] = 1
		}
	}

	return result
}

func prepareString(word string) string {
	word = strings.ToLower(word)
	word = strings.Trim(word, " ")
	word = reg.ReplaceAllString(word, "$")
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

func init() {
	var err error
	reg, err = regexp.Compile("[^a-z0-9а-яё ]+")
	if err != nil {
		panic(err)
	}
}

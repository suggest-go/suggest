package suggest

import (
	"regexp"
	"strings"
)

type profile struct {
	frequencies map[string]int
	ngrams      []string
}

var reg *regexp.Regexp

func SplitIntoNGrams(word string, k int) []string {
	sliceLen := len(word) - k + 1
	if sliceLen <= 0 || sliceLen > len(word) {
		panic("Invalid word length for spliting")
	}

	result := make([]string, 0, sliceLen)
	for i := 0; i < sliceLen; i++ {
		result = append(result, word[i: i+k])
	}

	return result
}

/*
 * Return unique ngrams with frequency
 */
func getProfile(word string, k int) *profile {
	ngrams := SplitIntoNGrams(word, k)
	frequencies := make(map[string]int, len(ngrams))
	for _, ngram := range ngrams {
		frequencies[ngram]++
	}

	unique := make([]string, 0, len(frequencies))
	for ngram, _ := range frequencies {
		unique = append(unique, ngram)
	}

	return &profile{frequencies, unique}
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

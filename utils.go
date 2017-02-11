package suggest

import (
	"fmt"
	"regexp"
	"strings"
)

type profile map[string]int

var (
	reg          *regexp.Regexp
	profileCache map[string]profile
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

/*
 * Return unique ngrams with frequency
 */
func getProfile(word string, k int) profile {
	key := fmt.Sprintf("%s_%d", word, k)
	if prof, ok := profileCache[key]; ok {
		return prof
	}

	ngrams := SplitIntoNGrams(word, k)
	result := make(profile)
	for _, ngram := range ngrams {
		if _, ok := result[ngram]; ok {
			result[ngram]++
		} else {
			result[ngram] = 1
		}
	}

	profileCache[key] = result
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

	profileCache = make(map[string]profile)
}

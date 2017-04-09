package suggest

import (
	"bufio"
	"bytes"
	"os"
	"regexp"
	"strings"
)

const maxN = 8

var reg *regexp.Regexp

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

func GetNGramSet(word string, k int) []string {
	ngrams := SplitIntoNGrams(word, k)
	set := make(map[string]struct{}, len(ngrams))
	list := make([]string, 0, len(ngrams))
	for _, ngram := range ngrams {
		_, found := set[ngram]
		set[ngram] = struct{}{}
		if !found {
			list = append(list, ngram)
		}
	}

	return list
}

/*
* inspired by https://github.com/jprichardson/readline-go/blob/master/readline.go
 */
func GetWordsFromFile(fileName string) []string {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	var result []string
	defer f.Close()
	buf := bufio.NewReader(f)
	line, err := buf.ReadBytes('\n')
	for err == nil {
		line = bytes.TrimRight(line, "\n")
		if len(line) > 0 {
			if line[len(line)-1] == 13 { //'\r'
				line = bytes.TrimRight(line, "\r")
			}

			result = append(result, string(line))
		}

		line, err = buf.ReadBytes('\n')
	}

	if len(line) > 0 {
		result = append(result, string(line))
	}

	return result
}

func prepareString(word string) string {
	if len(word) < 2 {
		return word
	}

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
	reg, err = regexp.Compile("[^a-z0-9а-яё]+")
	if err != nil {
		panic(err)
	}
}

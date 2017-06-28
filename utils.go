package suggest

import (
	"bufio"
	"bytes"
	"os"
)

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

// GetWordsFromFile is a helper for getting list of string from given fileName
// inspired by https://github.com/jprichardson/readline-go/blob/master/readline.go
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

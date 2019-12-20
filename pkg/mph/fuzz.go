// +build gofuzz

package mph

import (
	"fmt"
	"github.com/suggest-go/suggest/pkg/dictionary"
	"github.com/suggest-go/suggest/pkg/utils"
	"math/rand"
)

const (
	chunk       = 20
	maxSize     = (1 << 32) - 1
	maxTestSize = 100
)

func Fuzz(data []byte) int {
	n := len(data)

	if n >= maxSize {
		return -1
	}

	var corpus []string

	for i := 0; i < n; i += chunk {
		end := i + chunk

		if end > n {
			end = n
		}

		corpus = append(corpus, string(data[i:end]))
	}

	dict := dictionary.NewInMemoryDictionary(corpus)
	mph := New()

	if err := mph.Build(dict); err != nil {
		panic(err)
	}

	testN := utils.Min(maxTestSize, len(corpus))

	for i := 0; i < testN; i++ {
		v := rand.Intn(len(corpus))

		expected := corpus[v]
		key := mph.Get(expected)

		if int(key) != v {
			panic(fmt.Sprintf("not equal, e: %v, a: %v, v: %v", expected, v, key))
		}
	}

	return 1
}

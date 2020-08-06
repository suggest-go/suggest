package lm

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlow(t *testing.T) {
	trie := NewCountTrie()

	trie.Put(Sentence{"1", "2", "3"}, 3)
	trie.Put(Sentence{"1", "2", "3"}, 0)
	trie.Put(Sentence{"1", "2", "4"}, 2)
	trie.Put(Sentence{"1", "2", "3"}, 2)
	trie.Put(Sentence{"2", "3", "4", "5"}, 7)
	trie.Put(Sentence{"1", "2"}, 7)
	trie.Put(Sentence{"1"}, 12)
	trie.Put(Sentence{"4"}, 8)
	trie.Put(Sentence{"4"}, 0)
	trie.Put(Sentence{"1", "2", "3", "4"}, 7)
	trie.Put(Sentence{"3"}, 2)
	trie.Put(Sentence{"3", "2"}, 3)

	type row struct {
		path  string
		count WordCount
	}

	expected := []row{
		{"1", 12},
		{"1 2", 7},
		{"1 2 3", 5},
		{"1 2 3 4", 7},
		{"1 2 4", 2},
		{"2 3 4 5", 7},
		{"3", 2},
		{"3 2", 3},
		{"4", 8},
	}

	actual := make([]row, 0)

	err := trie.Walk(func(path []Token, count WordCount) error {
		actual = append(
			actual,
			row{
				path:  strings.Trim(strings.Replace(fmt.Sprint(path), " ", " ", -1), "[]"),
				count: count,
			},
		)

		return nil
	})

	assert.NoError(t, err)

	sort.Slice(actual, func(i, j int) bool {
		return actual[i].path < actual[j].path
	})

	assert.Equal(t, expected, actual)
}

func BenchmarkWalk(b *testing.B) {
	trie := NewCountTrie()
	path := make(Sentence, 0, 4)

	for i := 0; i < 10000; i++ {
		for j := 0; j < rand.Intn(3)+1; j++ {
			path = append(path, Token(rand.Int31n(10000)))
		}

		trie.Put(path, 100)
		path = path[:0]
	}

	b.StartTimer()

	j := 0
	for i := 0; i < b.N; i++ {
		err := trie.Walk(func(path Sentence, count WordCount) error {
			j++
			return nil
		})

		if err != nil {
			b.Errorf("Unexpected error occurs: %v", err)
		}
	}
}

package mph

import (
	"math/rand"
	"testing"

	"github.com/suggest-go/suggest/pkg/dictionary"
)

func TestFlow(t *testing.T) {
	collection := []string{
		"Hello",
		"This",
		"is",
		"mph",
		"package",
		"!",
	}

	dict := dictionary.NewInMemoryDictionary(collection)
	mph := New()

	if err := mph.Build(dict); err != nil {
		t.Errorf("Unexpected error occurs: %v", err)
	}

	list := rand.Perm(len(collection))

	for _, v := range list {
		expected := collection[v]
		key := mph.Get(expected)
		actual, err := dict.Get(key)

		if err != nil {
			t.Errorf("Unexpected error occurs: %v", err)
		}

		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	}
}

func BenchmarkMPHGet(b *testing.B) {
	dict, err := dictionary.OpenRAMDictionary("testdata/words.dict")

	if err != nil {
		b.Fatalf("Unexpected error occurs: %v", err)
	}

	mph := New()

	if err := mph.Build(dict); err != nil {
		b.Errorf("Unexpected error occurs: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	size := dict.Size()

	for i := 0; i < b.N; i++ {
		expected := dictionary.Key(i % size)
		word, err := dict.Get(expected)

		if err != nil {
			b.Fatalf("Unexpected error occurs: %v", err)
		}

		actual := mph.Get(word)

		if actual != expected {
			b.Fatalf("Test fail, expected %d, got %d", expected, actual)
		}
	}
}

func BenchmarkGoMapGet(b *testing.B) {
	dict, err := dictionary.OpenRAMDictionary("testdata/words.dict")

	if err != nil {
		b.Fatalf("Unexpected error occurs: %v", err)
	}

	table := make(map[dictionary.Value]dictionary.Key, dict.Size())

	_ = dict.Iterate(func(docID dictionary.Key, word dictionary.Value) error {
		table[word] = docID

		return nil
	})

	b.ReportAllocs()
	b.ResetTimer()

	size := dict.Size()

	for i := 0; i < b.N; i++ {
		expected := dictionary.Key(i % size)
		word, err := dict.Get(expected)

		if err != nil {
			b.Fatalf("Unexpected error occurs: %v", err)
		}

		actual := table[word]

		if actual != expected {
			b.Fatalf("Test fail, expected %d, got %d", expected, actual)
		}
	}
}

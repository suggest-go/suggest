package mph

import (
	"math/rand"
	"testing"

	"github.com/alldroll/suggest/pkg/dictionary"
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
	mph, err := BuildMPH(dict)

	if err != nil {
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

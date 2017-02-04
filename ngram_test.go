package suggest

import (
	"log"
	"testing"
)

func TestSuggest(t *testing.T) {
	collection := []string{
		"blue",
		"blunder",
		"blunt",
		"flank",
		"flu",
		"fluence",
		"fluent",
		"flunker",
		"test",
	}

	ngramIndex := NewNGramIndex(2)

	for _, word := range collection {
		ngramIndex.AddWord(word)
	}

	candidates := ngramIndex.Suggest("flunk", 3)
	log.Printf("%v", candidates)

	candidates = ngramIndex.Suggest("tes", 5)
	log.Printf("%v", candidates)
}

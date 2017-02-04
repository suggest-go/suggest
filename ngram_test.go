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
		"tes hello",
	}

	ngramIndex := NewNGramIndex(2)

	for _, word := range collection {
		ngramIndex.AddWord(word)
	}

	candidates := ngramIndex.Suggest("flunk", 3)
	log.Printf("%v", candidates)
}

func TestSuggestAuto(t *testing.T) {
	collection := []string{
		"Nissan March",
		"Nissan Juke",
		"Nissan Maxima",
		"Nissan Murano",
		"Nissan Moco",
		"Toyota Mark II",
	}

	ngramIndex := NewNGramIndex(3)

	for _, word := range collection {
		ngramIndex.AddWord(word)
	}

	candidates := ngramIndex.Suggest("Nissan", 5)
	log.Printf("%v", candidates)
}

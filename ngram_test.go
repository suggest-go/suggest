package suggest

import (
	"reflect"
	"testing"
)

func TestSuggestAuto(t *testing.T) {
	collection := []string{
		"Nissan March",
		"Nissan Juke",
		"Nissan Maxima",
		"Nissan Murano",
		"Nissan Note",
		"Toyota Mark II",
		"Toyota Corolla",
		"Toyota Corona",
	}

	ngramIndex := NewNGramIndex(3)
	for _, word := range collection {
		ngramIndex.AddWord(word)
	}

	candidates := ngramIndex.Suggest("Nissan mar", 2)
	expected := []string{
		"Nissan March",
		"Nissan Maxima",
	}

	if !reflect.DeepEqual(expected, candidates) {
		t.Errorf(
			"Test Fail, expected %v, got %v",
			expected,
			candidates,
		)
	}
}

func BenchmarkSuggest(b *testing.B) {
	b.StopTimer()
	collection := []string{
		"Nissan March",
		"Nissan Juke",
		"Nissan Maxima",
		"Nissan Murano",
		"Nissan Note",
		"Toyota Mark II",
		"Toyota Corolla",
		"Toyota Corona",
	}

	ngramIndex := NewNGramIndex(3)
	for _, word := range collection {
		ngramIndex.AddWord(word)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ngramIndex.Suggest("Nissan mar", 2)
	}
}

func BenchmarkRealExample(b *testing.B) {
	b.StopTimer()
	collection := GetWordsFromFile("cmd/web/cars.dict")

	ngramIndex := NewNGramIndex(3)

	for _, word := range collection {
		ngramIndex.AddWord(word)
	}

	queries := [...]string{
		"Nissan Mar",
		"Hnda Fi",
		"Mersdes Benz",
		"Tayota carolla",
		"Nssan Skylike",
		"Nissan Juke",
		"Dodje iper",
		"Hummer",
		"tayota",
	}

	qLen := len(queries)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		word := queries[i%qLen]
		ngramIndex.Suggest(word, 5)
	}
}

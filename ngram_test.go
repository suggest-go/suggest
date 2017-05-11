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

	ngramIndex := NewNGramIndex(getIndexConfWithBaseAlphabet(3))
	for i, word := range collection {
		ngramIndex.AddWord(word, i)
	}

	conf, err := NewSearchConfig("Nissan ma", 2, JACCARD, 0.5)
	if err != nil {
		panic(err)
	}

	candidates := ngramIndex.Suggest(conf)
	expected := []WordKey{
		2,
		0,
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

	ngramIndex := NewNGramIndex(getIndexConfWithBaseAlphabet(3))
	for i, word := range collection {
		ngramIndex.AddWord(word, i)
	}

	b.StartTimer()
	conf, err := NewSearchConfig("Nissan mar", 2, JACCARD, 0.5)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		ngramIndex.Suggest(conf)
	}
}

func BenchmarkRealExample(b *testing.B) {
	b.StopTimer()
	collection := GetWordsFromFile("cmd/web/cars.dict")

	ngramIndex := NewNGramIndex(getIndexConfWithBaseAlphabet(3))

	for i, word := range collection {
		ngramIndex.AddWord(word, i)
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

	conf, err := NewSearchConfig("Nissan mar", 5, COSINE, 0.3)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		conf.query = queries[i%qLen]
		ngramIndex.Suggest(conf)
	}
}

func getIndexConfWithBaseAlphabet(ngramSize int) *IndexConfig {
	alphabet := NewCompositeAlphabet([]Alphabet{
		NewEnglishAlphabet(),
		NewNumberAlphabet(),
		NewRussianAlphabet(),
		NewSimpleAlphabet([]rune{'$'}),
	})

	conf, err := NewIndexConfig(ngramSize, alphabet, "$", "$")
	if err != nil {
		panic(err)
	}

	return conf
}

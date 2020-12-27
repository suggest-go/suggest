package suggest

import (
	"bufio"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suggest-go/suggest/pkg/dictionary"
	"github.com/suggest-go/suggest/pkg/index"
	"github.com/suggest-go/suggest/pkg/metric"
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

	nGramIndex := buildNGramIndex(collection)

	candidates, err := nGramIndex.Suggest("Nissan ma", 0.5, metric.JaccardMetric(), newFuzzyCollectorManager(2))
	assert.NoError(t, err)

	actual := make([]index.Position, 0, len(candidates))

	for _, candidate := range candidates {
		actual = append(actual, candidate.Key)
	}

	expected := []index.Position{2, 0}
	assert.Equal(t, expected, actual)
}

func TestAutoComplete(t *testing.T) {
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

	nGramIndex := buildNGramIndex(collection)

	candidates, err := nGramIndex.Autocomplete("Niss", newFirstKCollectorManager(5))
	assert.NoError(t, err)

	actual := make([]index.Position, 0, len(candidates))

	for _, candidate := range candidates {
		actual = append(actual, candidate.Key)
	}

	expected := []index.Position{0, 1, 2, 3, 4}
	assert.Equal(t, expected, actual)
}

func BenchmarkSuggest(b *testing.B) {
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

	nGramIndex := buildNGramIndex(collection)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		nGramIndex.Suggest("Nissan mar", 0.5, metric.CosineMetric(), newFuzzyCollectorManager(5))
	}
}

func BenchmarkAutoComplete(b *testing.B) {
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

	qLen := len(collection)

	nGramIndex := buildNGramIndex(collection)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		nGramIndex.Autocomplete(collection[i%qLen], newFirstKCollectorManager(5))
	}
}

func BenchmarkRealExampleInMemory(b *testing.B) {
	file, err := os.Open("testdata/cars.dict")

	if err != nil {
		b.Errorf("Unexpected error: %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	collection := make([]string, 0)

	for scanner.Scan() {
		collection = append(collection, scanner.Text())
	}

	nGramIndex := buildNGramIndex(collection)

	b.ResetTimer()
	benchmarkRealExample(b, nGramIndex)
}

func BenchmarkRealExampleOnDisc(b *testing.B) {
	nGramIndex := buildOnDiscNGramIndex(0)

	b.ResetTimer()
	benchmarkRealExample(b, nGramIndex)
}

func BenchmarkSuggestWordsOnDisc(b *testing.B) {
	index := buildOnDiscNGramIndex(1)

	b.ResetTimer()
	b.ReportAllocs()

	queries := [...]string{
		"testing",
		"Acuracacy",
		"Indpendence",
		"Villictiy",
		"Velocity",
		"matehmatica",
		"acationally",
		"misleading",
		"litter",
		"arthroendoscopy",
	}

	qLen := len(queries)

	for i := 0; i < b.N; i++ {
		index.Suggest(queries[i%qLen], 0.5, metric.CosineMetric(), newFuzzyCollectorManager(5))
	}
}

func BenchmarkAutocompleteWordsOnDisc(b *testing.B) {
	index := buildOnDiscNGramIndex(1)

	b.ResetTimer()
	b.ReportAllocs()

	queries := [...]string{
		"testing",
		"Acuracacy",
		"Indpendence",
		"Villictiy",
		"Velocity",
		"matehmatica",
		"acationally",
		"misleading",
		"litter",
		"arthroendoscopy",
	}

	qLen := len(queries)

	for i := 0; i < b.N; i++ {
		index.Autocomplete(queries[i%qLen], newFirstKCollectorManager(5))
	}
}

func benchmarkRealExample(b *testing.B, index NGramIndex) {
	b.ReportAllocs()

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

	for i := 0; i < b.N; i++ {
		index.Suggest(queries[i%qLen], 0.5, metric.CosineMetric(), newFuzzyCollectorManager(5))
	}
}

func buildNGramIndex(collection []string) NGramIndex {
	config := IndexDescription{
		Driver:    RAMDriver,
		Name:      "index",
		NGramSize: 3,
		Pad:       "$",
		Wrap:      [2]string{"$", "$"},
		Alphabet:  []string{"english", "russian", "numbers", "$"},
	}

	dict := dictionary.NewInMemoryDictionary(collection)
	builder, err := NewRAMBuilder(dict, config)

	if err != nil {
		log.Fatal(err)
	}

	index, err := builder.Build()

	if err != nil {
		log.Fatal(err)
	}

	return index
}

func buildOnDiscNGramIndex(off int) NGramIndex {
	description, err := ReadConfigs("testdata/config.json")

	if err != nil {
		log.Fatal(err)
	}

	builder, err := NewFSBuilder(description[off])

	if err != nil {
		log.Fatal(err)
	}

	index, err := builder.Build()

	if err != nil {
		log.Fatal(err)
	}

	return index
}

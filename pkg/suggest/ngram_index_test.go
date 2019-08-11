package suggest

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/alldroll/suggest/pkg/dictionary"
	"github.com/alldroll/suggest/pkg/index"
	"github.com/alldroll/suggest/pkg/metric"
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
	conf, err := NewSearchConfig("Nissan ma", 2, metric.JaccardMetric(), 0.5)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	candidates, err := nGramIndex.Suggest(conf)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	actual := make([]index.Position, 0, len(candidates))

	for _, candidate := range candidates {
		actual = append(actual, candidate.Key)
	}

	expected := []index.Position{
		2,
		0,
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"Test Fail, expected %v, got %v",
			expected,
			actual,
		)
	}
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
	conf, err := NewSearchConfig("Nissan mar", 2, metric.JaccardMetric(), 0.5)

	if err != nil {
		b.Errorf("Unexpected error: %v", err)
	}

	for i := 0; i < b.N; i++ {
		nGramIndex.Suggest(conf)
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
	nGramIndex := buildOnDiscNGramIndex()

	b.ResetTimer()
	benchmarkRealExample(b, nGramIndex)
}

//
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
	conf, err := NewSearchConfig("Nissan mar", 5, metric.CosineMetric(), 0.3)

	if err != nil {
		b.Errorf("Unexpected error: %v", err)
	}

	for i := 0; i < b.N; i++ {
		conf.query = queries[i%qLen]
		index.Suggest(conf)
	}
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
	candidates, err := nGramIndex.AutoComplete("Niss", 5)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	actual := make([]index.Position, 0, len(candidates))

	for _, candidate := range candidates {
		actual = append(actual, candidate.Key)
	}

	expected := []index.Position{
		0, 1, 2, 3, 4,
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"Test Fail, expected %v, got %v",
			expected,
			actual,
		)
	}
}

func BenchmarkAutoCompleteOnDisc(b *testing.B) {
	b.StopTimer()

	nGramIndex := buildOnDiscNGramIndex()

	queries := [...]string{
		"Nissan Mar",
		"Fit",
		"Benz",
		"Toyo",
		"Nissan Juke",
		"Hummer",
		"Corol",
	}

	qLen := len(queries)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		nGramIndex.AutoComplete(queries[i%qLen], 5)
	}
}

func buildNGramIndex(collection []string) NGramIndex {
	builder, err := NewRAMBuilder(
		dictionary.NewInMemoryDictionary(collection),
		IndexDescription{
			Name:      "index",
			NGramSize: 3,
			Pad:       "$",
			Wrap:      [2]string{"$", "$"},
			Alphabet:  []string{"english", "russian", "numbers", "$"},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	index, err := builder.Build()

	if err != nil {
		log.Fatal(err)
	}

	return index
}

func buildOnDiscNGramIndex() NGramIndex {
	description, err := ReadConfigs("testdata/config.json")

	if err != nil {
		log.Fatal(err)
	}

	builder, err := NewFSBuilder(description[0])

	if err != nil {
		log.Fatal(err)
	}

	index, err := builder.Build()

	if err != nil {
		log.Fatal(err)
	}

	return index
}

package suggest

import (
	"reflect"
	"os"
	"bufio"
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

	nGramIndex := buildNGramIndex(NewInMemoryDictionary(collection), 3)

	conf, err := NewSearchConfig("Nissan ma", 2, JaccardMetric(), 0.5)
	if err != nil {
		panic(err)
	}

	candidates := nGramIndex.Suggest(conf)
	actual := make([]int, 0, len(candidates))
	for _, candidate := range candidates {
		actual = append(actual, candidate.Key)
	}

	expected := []int{
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

	nGramIndex := buildNGramIndex(NewInMemoryDictionary(collection), 3)

	b.StartTimer()
	conf, err := NewSearchConfig("Nissan mar", 2, JaccardMetric(), 0.5)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		nGramIndex.Suggest(conf)
	}
}

func BenchmarkRealExample(b *testing.B) {
	b.StopTimer()

	file, err := os.Open("testdata/cars.dict")
	defer file.Close()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	collection := make([]string, 0)
	for scanner.Scan() {
		collection = append(collection, scanner.Text())
	}

	nGramIndex := buildNGramIndex(NewInMemoryDictionary(collection), 3)

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

	conf, err := NewSearchConfig("Nissan mar", 5, CosineMetric(), 0.3)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		conf.query = queries[i%qLen]
		nGramIndex.Suggest(conf)
	}
}

func buildNGramIndex(dictionary Dictionary, ngramSize int) NGramIndex {
	alphabet := NewCompositeAlphabet([]Alphabet{
		NewRussianAlphabet(),
		NewEnglishAlphabet(),
		NewNumberAlphabet(),
		NewSimpleAlphabet([]rune{'$'}),
	})

	return NewRunTimeBuilder().
		SetDictionary(dictionary).
		SetNGramSize(ngramSize).
		SetAlphabet(alphabet).
		Build()
}

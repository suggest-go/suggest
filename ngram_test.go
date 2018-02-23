package suggest

import (
	"bufio"
	"github.com/alldroll/suggest/alphabet"
	"github.com/alldroll/suggest/dictionary"
	"github.com/alldroll/suggest/index"
	"github.com/alldroll/suggest/metric"
	"golang.org/x/exp/mmap"
	"io"
	"log"
	"os"
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

	nGramIndex := buildNGramIndex(collection, 3)

	conf, err := NewSearchConfig("Nissan ma", 2, metric.JaccardMetric(), 0.5)
	if err != nil {
		panic(err)
	}

	candidates := nGramIndex.Suggest(conf)
	actual := make(index.PostingList, 0, len(candidates))
	for _, candidate := range candidates {
		actual = append(actual, candidate.Key)
	}

	expected := index.PostingList{
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

	nGramIndex := buildNGramIndex(collection, 3)

	b.StartTimer()
	conf, err := NewSearchConfig("Nissan mar", 2, metric.JaccardMetric(), 0.5)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		nGramIndex.Suggest(conf)
	}
}

func BenchmarkRealExampleInMemory(b *testing.B) {
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

	nGramIndex := buildNGramIndex(collection, 3)

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

	conf, err := NewSearchConfig("Nissan mar", 5, metric.CosineMetric(), 0.3)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		conf.query = queries[i%qLen]
		nGramIndex.Suggest(conf)
	}
}

func BenchmarkRealExampleOnDisc(b *testing.B) {
	b.StopTimer()

	f, err := mmap.Open("testdata/db/cars.cdb")
	if err != nil {
		log.Fatal(err)
	}

	nGramIndex := buildOnDiscNGramIndex(f, 3)

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

	conf, err := NewSearchConfig("Nissan mar", 5, metric.CosineMetric(), 0.3)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		conf.query = queries[i%qLen]
		nGramIndex.Suggest(conf)
	}
}

func buildNGramIndex(collection []string, nGramSize int) NGramIndex {
	alphabet := alphabet.NewCompositeAlphabet([]alphabet.Alphabet{
		alphabet.NewRussianAlphabet(),
		alphabet.NewEnglishAlphabet(),
		alphabet.NewNumberAlphabet(),
		alphabet.NewSimpleAlphabet([]rune{'$'}),
	})

	return NewRunTimeBuilder().
		SetDictionary(dictionary.NewInMemoryDictionary(collection)).
		SetNGramSize(nGramSize).
		SetAlphabet(alphabet).
		Build()
}

func buildOnDiscNGramIndex(reader io.ReaderAt, nGramSize int) NGramIndex {
	alphabet := alphabet.NewCompositeAlphabet([]alphabet.Alphabet{
		alphabet.NewRussianAlphabet(),
		alphabet.NewEnglishAlphabet(),
		alphabet.NewNumberAlphabet(),
		alphabet.NewSimpleAlphabet([]rune{'$'}),
	})

	return NewBuilder("testdata/db/cars.hd", "testdata/db/cars.dl").
		SetAlphabet(alphabet).
		SetDictionary(dictionary.NewCDBDictionary(reader)).
		SetNGramSize(nGramSize).
		SetWrap("$").
		SetPad("$").
		Build()
}

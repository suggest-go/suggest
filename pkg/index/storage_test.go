package index

import (
	"bufio"
	"github.com/alldroll/suggest/pkg/alphabet"
	"github.com/alldroll/suggest/pkg/compression"
	"github.com/alldroll/suggest/pkg/dictionary"
	"log"
	"os"
	"testing"
)

func TestOnDiscWriter_Save(t *testing.T) {
	headerFile := "../suggest/testdata/db/test.hd"
	defer func() {
		os.Remove(headerFile)
	}()

	docListFile := "../suggest/testdata/db/test.dl"
	defer func() {
		os.Remove(docListFile)
	}()

	dict, err := os.Open("../suggest/testdata/cars.dict")
	defer dict.Close()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(dict)
	collection := make([]string, 0)
	for scanner.Scan() {
		collection = append(collection, scanner.Text())
	}

	indices := buildIndex()
	writer := NewOnDiscIndicesWriter(compression.VBEncoder(), headerFile, docListFile)
	err = writer.Save(indices)
	if err != nil {
		t.Error(err)
	}

	reader := NewOnDiscIndicesReader(compression.VBDecoder(), headerFile, docListFile)
	_, err = reader.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func buildIndex() Indices {
	dict, err := os.Open("../suggest/testdata/cars.dict")
	defer dict.Close()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(dict)
	collection := make([]string, 0)
	for scanner.Scan() {
		collection = append(collection, scanner.Text())
	}

	dictionary := dictionary.NewInMemoryDictionary(collection)

	alphabet := alphabet.NewCompositeAlphabet([]alphabet.Alphabet{
		alphabet.NewRussianAlphabet(),
		alphabet.NewEnglishAlphabet(),
		alphabet.NewNumberAlphabet(),
		alphabet.NewSimpleAlphabet([]rune{'$'}),
	})

	indexer := NewIndexer(3, NewGenerator(3, alphabet), NewCleaner(alphabet.Chars(), "$", [2]string{"$", "$"}))

	return indexer.Index(dictionary)
}

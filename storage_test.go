package suggest

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func TestOnDiscWriter_Save(t *testing.T) {
	header, err := os.Create("testdata/db/test.hd")
	defer func() {
		header.Close()
		os.Remove(header.Name())
	}()
	if err != nil {
		log.Fatal(err)
	}

	docList, err := os.Create("testdata/db/test.dl")
	defer func() {
		docList.Close()
		os.Remove(docList.Name())
	}()

	if err != nil {
		log.Fatal(err)
	}

	dict, err := os.Open("testdata/cars.dict")
	defer dict.Close()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(dict)
	collection := make([]string, 0)
	for scanner.Scan() {
		collection = append(collection, scanner.Text())
	}

	index := buildIndex()
	writer := &onDiscWriter{VBEncoder(), header, docList, 0}
	writer.Save(index)

	reader := &onDiscReader{VBDecoder(), header, docList, 0}
	reader.Load()
}

func buildIndex() Index {
	dict, err := os.Open("testdata/cars.dict")
	defer dict.Close()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(dict)
	collection := make([]string, 0)
	for scanner.Scan() {
		collection = append(collection, scanner.Text())
	}

	dictionary := NewInMemoryDictionary(collection)

	alphabet := NewCompositeAlphabet([]Alphabet{
		NewRussianAlphabet(),
		NewEnglishAlphabet(),
		NewNumberAlphabet(),
		NewSimpleAlphabet([]rune{'$'}),
	})

	indexer := NewIndexer(3, NewGenerator(3, alphabet), NewCleaner(alphabet.Chars(), "$", "$"))

	return indexer.Index(dictionary)
}

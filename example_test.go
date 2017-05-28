// This example demonstrates an usage of suggest.Service
package suggest_test

import (
	"fmt"
	"github.com/alldroll/suggest"
)

// This example demonstrates how to use this package.
func Example() {
	// here we define our alphabet for given collection of words
	alphabet := suggest.NewCompositeAlphabet([]suggest.Alphabet{
		suggest.NewEnglishAlphabet(),
		suggest.NewSimpleAlphabet([]rune{'$'}), // pad wrap
	})

	// create IndexConfig with ngramSize, alphabet, wrap and pad
	wrap, pad := "$", "$"
	conf, err := suggest.NewIndexConfig(3 /*Ngram size*/, alphabet, wrap, pad)
	if err != nil {
		panic(err)
	}

	// we create InMemoryDictionary. Here we can use anything we want,
	// for example SqlDictionary
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
	dictionary := suggest.NewInMemoryDictionary(collection)

	service := suggest.NewService()
	service.AddDictionary("cars", dictionary, conf)

	topK := 5
	sim := 0.5
	query := "niss ma"
	searchConf, err := suggest.NewSearchConfig(query, topK, suggest.CosineMetric(), sim)
	if err != nil {
		panic(err)
	}

	fmt.Println(service.Suggest("cars", searchConf))
	// Output: [Nissan Maxima Nissan March]
}

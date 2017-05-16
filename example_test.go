package suggest_test

import (
	"fmt"
	"github.com/alldroll/suggest"
)

func ExampleSuggest() {
	alphabet := suggest.NewCompositeAlphabet([]suggest.Alphabet{
		suggest.NewEnglishAlphabet(),
		suggest.NewSimpleAlphabet([]rune{'$'}),
	})

	wrap, pad := "$", "$"
	conf, err := suggest.NewIndexConfig(3 /*Ngram size*/, alphabet, wrap, pad)
	if err != nil {
		panic(err)
	}

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

	service := suggest.NewSuggestService()
	service.AddDictionary("cars", dictionary, conf)

	topK := 5
	sim := 0.5
	query := "niss ma"
	searchConf, err := suggest.NewSearchConfig(query, topK, suggest.COSINE, sim)
	if err != nil {
		panic(err)
	}

	fmt.Println(service.Suggest("cars", searchConf))
	// Output: [Nissan Maxima Nissan March]
}

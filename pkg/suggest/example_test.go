// This example demonstrates an usage of suggest.Service
package suggest_test

import (
	"fmt"
	"log"

	"github.com/alldroll/suggest/pkg/dictionary"
	"github.com/alldroll/suggest/pkg/metric"
	"github.com/alldroll/suggest/pkg/suggest"
)

// This example demonstrates how to use this package.
func Example() {
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

	service := suggest.NewService()

	// create dictionary of autocomplete items
	dictionary := dictionary.NewInMemoryDictionary(collection)
	// describe index configuration
	indexDescription := suggest.IndexDescription{
		Name:      "cars",                   // name of the dictionary
		NGramSize: 3,                        // size of the nGram
		Wrap:      [2]string{"$", "$"},      // wrap symbols (front and rear)
		Pad:       "$",                      // pad to replace with forbidden chars
		Alphabet:  []string{"english", "$"}, // alphabet of allowed chars (other chars will be replaced with pad symbol)
	}

	// create runtime search index builder
	builder, err := suggest.NewRAMBuilder(dictionary, indexDescription)

	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	// asking our service for adding a new search index with given configuration
	if err := service.AddIndex(indexDescription.Name, dictionary, builder); err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	// declare a search configuration (query, topK elements, type of metric, min similarity)
	searchConf, err := suggest.NewSearchConfig("niss ma", 5, metric.CosineMetric(), 0.4)

	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	result, err := service.Suggest("cars", searchConf)

	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	values := make([]string, 0, len(result))

	for _, item := range result {
		values = append(values, item.Value)
	}

	fmt.Println(values)
	// Output: [Nissan Maxima Nissan March]
}
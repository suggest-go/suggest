# suggest

Tool for Top-k Approximate String Matching.

[![Go Report Card](https://goreportcard.com/badge/github.com/alldroll/suggest)](https://goreportcard.com/report/github.com/alldroll/suggest)
[![GoDoc](https://godoc.org/github.com/alldroll/suggest?status.svg)](https://godoc.org/github.com/alldroll/suggest)

## Usage

```go
// This example demonstrates an usage of suggest.Service
package suggest_test

import (
	"fmt"
	"github.com/alldroll/suggest/pkg/alphabet"
	"github.com/alldroll/suggest/pkg/dictionary"
	"github.com/alldroll/suggest/pkg/metric"
	"github.com/alldroll/suggest/pkg/suggest"
)

// This example demonstrates how to use this package.
func Example() {
	// here we define our alphabet for given collection of words
	alphabet := alphabet.NewCompositeAlphabet([]alphabet.Alphabet{
		alphabet.NewEnglishAlphabet(),
		alphabet.NewSimpleAlphabet([]rune{'$'}), // pad wrap
	})

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

	dictionary := dictionary.NewInMemoryDictionary(collection)

	// create IndexConfig with ngramSize, alphabet, wrap and pad
	wrap, pad := "$", "$"
	ngramSize := 3
	conf, err := suggest.NewIndexConfig(ngramSize, dictionary, alphabet, wrap, pad)
	if err != nil {
		panic(err)
	}

	service := suggest.NewService()
	service.AddRunTimeIndex("cars", conf)

	topK := 5
	sim := 0.4
	query := "niss ma"
	searchConf, err := suggest.NewSearchConfig(query, topK, metric.CosineMetric(), sim)
	if err != nil {
		panic(err)
	}

	result, err := service.Suggest("cars", searchConf)
	if err != nil {
		panic(err)
	}

	values := make([]string, 0, len(result))
	for _, item := range result {
		values = append(values, item.Value)
	}

	fmt.Println(values)
	// Output: [Nissan Maxima Nissan March]
}
```

## Demo
see https://suggest-demo.herokuapp.com/ as complete example (https://github.com/alldroll/suggest_demo)

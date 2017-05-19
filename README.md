# suggest

Tool for Top-k Approximate String Matching.

[![Go Report Card](https://goreportcard.com/badge/github.com/alldroll/suggest)](https://goreportcard.com/report/github.com/alldroll/suggest)
[![GoDoc](https://godoc.org/github.com/Lazin/go-ngram?status.png)](https://godoc.org/github.com/alldroll/suggest)

## Usage

```go
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
searchConf, err := suggest.NewSearchConfig(query, topK, suggest.Cosine, sim)
if err != nil {
    panic(err)
}

fmt.Println(service.Suggest("cars", searchConf))
// Output: [Nissan Maxima Nissan March]
```

## Demo
see https://tranquil-journey-12522.herokuapp.com/ as complete example
or run it localy by `go run cmd/web/main.go`

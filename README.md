# Suggest

Library for Top-k Approximate String Matching.

[![Go Report Card](https://goreportcard.com/badge/github.com/alldroll/suggest)](https://goreportcard.com/report/github.com/alldroll/suggest)
[![GoDoc](https://godoc.org/github.com/alldroll/suggest?status.svg)](https://godoc.org/github.com/alldroll/suggest)

The library was mostly inspired by
- http://www.chokkan.org/software/simstring/
- http://www.aaai.org/ocs/index.php/AAAI/AAAI10/paper/viewFile/1939/2234
- http://nlp.stanford.edu/IR-book/
- http://bazhenov.me/blog/2012/08/04/autocomplete.html
- http://www.aclweb.org/anthology/C10-1096


## Purpose

Let's imagine you have a website, for instance a pharmacy website.
There could be a lot of dictionaries, such as a list of medical drugs,
a list of cities (countries), where you can deliver your goods and so on.
Some of these dictionaries could be a pretty large, and it might be a
tedious for a customer to choose the correct option from the dictionary.
Have the possibility of `Top-k approximate string search` in a dictionary
is a significant these cases.

This library provides API and the simple http service for such purposes.


## Demo

Please, see [demo](http://54.183.244.111:8000/) as a complete example.
This example provides a fuzzy search in a list of dictionaries, with ability
of choosing a similarity, type of metric and topK.

## Index description format

`IndexDescription` is a crucial part of the library. It tells how to configure and to build search indexes. You can define your index descriptions in `JSON` format, as following
```
[
  ...,
  {
    "driver": "DISC", // DISC (if you have indexed data by `indexer` command) or RAM (if you would like to build a search index runtime and store all in memory)
    "name": "cars", // Name of a dictionary
    "nGramSize": 3, // The size of nGram
    "alphabet": ["russian", "english", "numbers", "$"], // Alphabet of allowed chars
    "source": "testdata/cars.dict", // Path to a source for a dictionary (the list of words, each word should start from the new line)
    "output": "testdata/db", // Path to a indexed data (if we have indexed the source by using `indexer` command)
    "pad": "$", // Pad symbol, to replace undefined chars
    "wrap": ["$", "$"] // Wrap symbols, to wrap front and rear of a word
  },
  ...
]
```

## Usage

```go
    // The dictionary, on which we expect fuzzy search
    dictionary := dictionary.NewInMemoryDictionary([]string{
        "Nissan March",
        "Nissan Juke",
        "Nissan Maxima",
        "Nissan Murano",
        "Nissan Note",
        "Toyota Mark II",
        "Toyota Corolla",
        "Toyota Corona",
    })

    // create suggest service
    service := suggest.NewService()

    // here we describe our index configuration
    indexDescription := suggest.IndexDescription{
        Name:      "cars",                   // name of the dictionary
        NGramSize: 3,                        // size of the nGram
        Wrap:      [2]string{"$", "$"},      // wrap symbols (front and rear)
        Pad:       "$",                      // pad to replace with forbidden chars
        Alphabet:  []string{"english", "$"}, // alphabet of allowed chars (other chars will be replaced with pad symbol)
    }

    // create runtime search index builder (because we don't have indexed data)
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
```

## Suggest eval

Eval command is a command line tool for approximate string search.

## Suggest indexer

Indexer command builds a search index with the given [configuration](##index-description-format).
Generated data is required by `DISC` implementation of a index driver.

## Suggest service-run

Runs a http web server with suggest methods.

### REST API

#### suggest

Returns json data about a single user.

* **/suggest/{dict}/{query}/?metric={metric}&similarity={similarity}&topK={topK}**

* **Method:**

  `GET`

*  **URL Params**

   **Required:**

   `dict=[string]`

   `query=[string]`

   `metric=[string]` `Jaccard` | `Cosine` | `Dice` | `Exact`

   `similarity=[float]` `in (0, 1]`

   `topK=[integer]` `in [1, n]`

* **Success Response:**

  * **Code:** 200 <br />
    **Content:**
    ```
    [
        {Score: 1, Value: "test"},
        {Score: 0.8, Value: "tests"},
        {Score: 0.6, Value: "tested"},
    ]
    ```

* **Error Response:**

  * **Code:** 400 BAD REQUEST <br />
    **Content:** `{ error : "" }`

  * **Code:** 500 SERVER ERROR <br />
    **Content:** `description`


#### reindex

Commands to reload (rebuild) all indexes

* **/internal/reindex/**

* **Method:**

  `POST`

* **Success Response:**

  * **Code:** 200 <br />
    **Content:**
    ```OK```

* **Error Response:**

  * **Code:** 500 SERVER ERROR <br />
    **Content:** `description`

#### dictionary list

Returns a list of managed dictionaries

* **/dict/list/**

* **Method:**

  `GET`

* **Success Response:**

  * **Code:** 200 <br />
    **Content:**
    ```
    [
  		"dict1",
  		"dict2",
  		"dict3"
	]
    ```

* **Error Response:**

  * **Code:** 500 SERVER ERROR <br />
    **Content:** `description`

## TODO

* Autocomplete (to improve initial prototype)
* NGram language model
* Spellchecker

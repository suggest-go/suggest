# Suggest

Library for Top-k Approximate String Matching, autocomplete and spell checking.

[![Build Status](https://travis-ci.com/suggest-go/suggest.svg?branch=master)](https://travis-ci.com/suggest-go/suggest)
[![Go Report Card](https://goreportcard.com/badge/github.com/suggest-go/suggest)](https://goreportcard.com/report/github.com/suggest-go/suggest)
[![GoDoc](https://godoc.org/github.com/suggest-go/suggest?status.svg)](https://godoc.org/github.com/suggest-go/suggest)

The library was mostly inspired by
- http://www.chokkan.org/software/simstring/
- http://www.aaai.org/ocs/index.php/AAAI/AAAI10/paper/viewFile/1939/2234
- http://nlp.stanford.edu/IR-book/
- http://bazhenov.me/blog/2012/08/04/autocomplete.html
- http://www.aclweb.org/anthology/C10-1096


## Purpose

Let's imagine you have a website, for instance, a pharmacy website.
There could be a lot of dictionaries, such as a list of medical drugs,
a list of cities (countries), where you can deliver your goods and so on.
Some of these dictionaries could be pretty large, and it might be a
tedious for a customer to choose the correct option from the dictionary.
Having the possibility of `Top-k approximate string search` in a dictionary
is significant in these cases.

Also, the library provides spell checking functionality, that allows you to predict the next word.

The library provides API and the simple `HTTP service` for such purposes.

## Demo

The demo shows approximate string search in a dictionary with more than 200k English words

```
$ make build
$ ./build/suggest eval -c pkg/suggest/testdata/config.json -d words -s 0.5 -k 5
```

![Suggest eval Demo](suggest-eval.gif)

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
// we create InMemoryDictionary. Here we can use anything we want,
// for example SqlDictionary, CDBDictionary and so on
dict := dictionary.NewInMemoryDictionary([]string{
    "Nissan March",
    "Nissan Juke",
    "Nissan Maxima",
    "Nissan Murano",
    "Nissan Note",
    "Toyota Mark II",
    "Toyota Corolla",
    "Toyota Corona",
})

// describe index configuration
indexDescription := suggest.IndexDescription{
    Name:      "cars",                   // name of the dictionary
    NGramSize: 3,                        // size of the nGram
    Wrap:      [2]string{"$", "$"},      // wrap symbols (front and rear)
    Pad:       "$",                      // pad to replace with forbidden chars
    Alphabet:  []string{"english", "$"}, // alphabet of allowed chars (other chars will be replaced with pad symbol)
}

// create runtime search index builder
builder, err := suggest.NewRAMBuilder(dict, indexDescription)

if err != nil {
    log.Fatalf("Unexpected error: %v", err)
}

service := suggest.NewService()

// add a new search index with the given configuration
if err := service.AddIndex(indexDescription.Name, dict, builder); err != nil {
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

Eval command is a command-line tool for approximate string search.

## Suggest indexer

Indexer command builds a search index with the given [configuration](##index-description-format).
Generated data is required by `DISC` implementation of an index driver.

## Suggest service-run

Runs HTTP webserver with suggest methods.

## Language model ngram-count

Creates Google n-grams format

## Language model build-lm

Builds a binary representation of a stupid-backoff language model and writes it to disk

## Language model eval

Eval command is a cli for lm scoring

## Spellchecker

Cli for spell checking

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

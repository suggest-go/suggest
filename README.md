# suggest

Tool for Top-k Approximate String Matching.

This code was developed only for self education, so algorithm is not memory effective and so on.

Main idea was taken from
* http://www.aaai.org/ocs/index.php/AAAI/AAAI10/paper/viewFile/1939/2234
* http://nlp.stanford.edu/IR-book/html/htmledition/k-gram-indexes-for-wildcard-queries-1.html
* http://bazhenov.me/blog/2012/08/04/autocomplete.html

## Usage

```go
service := suggest.NewSuggestService(3, suggest.JACCARD)
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
service.AddDictionary("cars", collection)

service.Suggest("cars", "niss mar", 2)
// >>> [Nissan March Nissan Maxima]

service.Suggest("cars", "guke", 2)
// >>> [Nissan Juke]

```

or see `./cmd/web/main.go` as complete example

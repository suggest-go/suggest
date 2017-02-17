package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/alldroll/suggest"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const TOP_K = 5

var (
	suggesters    map[int]*suggest.SuggestService
	indexTemplate = template.Must(template.ParseFiles("public/index.html"))
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, nil)
}

func SuggestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dict, query := vars["dict"], vars["query"]

	type candidates struct {
		Metric  string   `json:"metric"`
		Data    []string `json:"data"`
		Elapsed string   `json:"elapsed"`
	}

	ch := make(chan candidates, len(suggesters))
	var wg sync.WaitGroup
	for metric, service := range suggesters {
		wg.Add(1)
		go func(metric int, service *suggest.SuggestService) {
			start := time.Now()
			data := service.Suggest(dict, query, TOP_K)
			elapsed := time.Since(start).String()
			metricName := suggest.MetricName[metric]
			ch <- candidates{metricName, data, elapsed}
			wg.Done()
		}(metric, service)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	result := make([]candidates, 0, len(suggesters))
	for cand := range ch {
		result = append(result, cand)
	}

	response := struct {
		Status     bool         `json:"status"`
		Candidates []candidates `json:"candidates"`
	}{
		true,
		result,
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

/*
* inspired by https://github.com/jprichardson/readline-go/blob/master/readline.go
 */
func GetWordsFromFile(fileName string) []string {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	var result []string
	defer f.Close()
	buf := bufio.NewReader(f)
	line, err := buf.ReadBytes('\n')
	for err == nil {
		line = bytes.TrimRight(line, "\n")
		if len(line) > 0 {
			if line[len(line)-1] == 13 { //'\r'
				line = bytes.TrimRight(line, "\r")
			}

			result = append(result, string(line))
		}

		line, err = buf.ReadBytes('\n')
	}

	if len(line) > 0 {
		result = append(result, string(line))
	}

	return result
}

func init() {
	words := GetWordsFromFile("cars.dict")
	suggesters = map[int]*suggest.SuggestService{
		suggest.LEVENSHTEIN: suggest.NewSuggestService(3, suggest.LEVENSHTEIN),
		suggest.NGRAM:       suggest.NewSuggestService(3, suggest.NGRAM),
		suggest.JACCARD:     suggest.NewSuggestService(3, suggest.JACCARD),
	}

	for _, sug := range suggesters {
		sug.AddDictionary("cars", words)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/suggest/{dict}/{query}/", SuggestHandler)
	log.Fatal(http.ListenAndServe(":8000", r))
}

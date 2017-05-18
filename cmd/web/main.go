package main

import (
	"encoding/json"
	"fmt"
	"github.com/alldroll/suggest"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"time"
)

const TOP_K = 5

var (
	suggestService *suggest.SuggestService
	configs        []*suggest.IndexConfig
	publicPath     = "./cmd/web/public"
	dictPath       = "./cmd/web/cars.dict"
)

func SuggestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dict, query := vars["dict"], vars["query"]

	type candidates struct {
		Config  string   `json:"config"`
		Data    []string `json:"data"`
		Elapsed string   `json:"elapsed"`
	}

	lenS := len(configs)
	ch := make(chan candidates)
	for i, _ := range configs {
		go func(i int) {
			start := time.Now()
			searchConf, err := suggest.NewSearchConfig(query, 5, suggest.COSINE, 0.5)
			if err == nil {
				// TODO fixme
			}

			data := suggestService.Suggest(dict+string(i), searchConf)
			elapsed := time.Since(start).String()
			configName := fmt.Sprintf("n%d", i+2)
			ch <- candidates{configName, data, elapsed}
		}(i)
	}

	result := make([]candidates, 0, lenS)
	for i := 0; i < lenS; i++ {
		result = append(result, <-ch)
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
	//w.Header().Set("Cache-Control", "public, cache, store, s-maxage=3600, max-age=3600")
	w.Write(data)
}

func init() {
	words := suggest.GetWordsFromFile(dictPath)
	alphabet := suggest.NewCompositeAlphabet([]suggest.Alphabet{
		suggest.NewEnglishAlphabet(),
		suggest.NewNumberAlphabet(),
		suggest.NewRussianAlphabet(),
		suggest.NewSimpleAlphabet([]rune{'$'}),
	})

	for j := 2; j <= 4; j++ {
		conf, err := suggest.NewIndexConfig(j, alphabet, "$", "$")
		if err != nil {
			panic(err)
		}

		configs = append(configs, conf)
	}

	dictionary := suggest.NewInMemoryDictionary(words)
	suggestService = suggest.NewSuggestService()
	for i, config := range configs {
		suggestService.AddDictionary("cars"+string(i), dictionary, config)
	}
}

func AttachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	r := mux.NewRouter()
	AttachProfiler(r)
	r.Handle("/", http.FileServer(http.Dir(publicPath)))
	r.HandleFunc("/suggest/{dict}/{query}/", SuggestHandler)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

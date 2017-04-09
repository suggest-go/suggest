package main

import (
	"encoding/json"
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
	configs        []*suggest.Config
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
	for i, config := range configs {
		go func(i int, config *suggest.Config) {
			start := time.Now()
			data := suggestService.Suggest(dict+string(i), query)
			elapsed := time.Since(start).String()
			configName := config.GetName()
			ch <- candidates{configName, data, elapsed}
		}(i, config)
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
	w.Header().Set("Cache-Control", "public, cache, store, s-maxage=3600, max-age=3600")
	w.Write(data)
}

func init() {
	words := suggest.GetWordsFromFile(dictPath)
	configs = []*suggest.Config{
		suggest.NewConfig(2, 5, "n2"),
		suggest.NewConfig(3, 5, "n3"),
		suggest.NewConfig(4, 5, "n4"),
	}

	suggestService = suggest.NewSuggestService()
	for i, config := range configs {
		suggestService.AddDictionary("cars"+string(i), words, config)
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

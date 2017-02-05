package main

import (
	"encoding/json"
	"github.com/alldroll/suggest"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
)

const TOP_K = 5

var (
	suggester     *suggest.SuggestService
	indexTemplate = template.Must(template.ParseFiles("public/index.html"))
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, nil)
}

func SuggestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dict, query := vars["dict"], vars["query"]
	candidates := suggester.Suggest(dict, query)

	response := struct {
		Status     bool     `json:"status"`
		Candidates []string `json:"candidates"`
	}{
		true,
		candidates,
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func init() {
	suggester = suggest.NewSuggestService(TOP_K)
	f, err := os.Open("cars.dict")
	if err != nil {
		panic(err)
	}

	suggester.AddDictionary("cars", f)
	f.Close()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/suggest/{dict}/{query}/", SuggestHandler)
	log.Fatal(http.ListenAndServe(":8000", r))
}

package main

import (
	"github.com/alldroll/suggest"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Suggest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Suggest!\n"))
}

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", Suggest)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}

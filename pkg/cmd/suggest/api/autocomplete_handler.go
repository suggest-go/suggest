package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/suggest-go/suggest/pkg/suggest"
)

// autocompleteHandler is responsible for query autocomplete
type autocompleteHandler struct {
	suggestService *suggest.Service
}

// handle performs autocomplete for the given query
func (h *autocompleteHandler) handle(w http.ResponseWriter, r *http.Request) {
	var (
		vars  = mux.Vars(r)
		dict  = vars["dict"]
		query = vars["query"]
		k     = r.FormValue("topK")
	)

	i64, err := strconv.ParseInt(k, 10, 0)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultItems, err := h.suggestService.Autocomplete(dict, query, int(i64))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resultItems)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

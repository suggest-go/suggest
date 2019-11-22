package api

import (
	"encoding/json"
	httputil "github.com/suggest-go/suggest/internal/http"
	"net/http"

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
	)

	topK, err := httputil.FormTopKValue(r, "topK", defaultTopK)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultItems, err := h.suggestService.Autocomplete(dict, query, topK)

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

	if _, err := w.Write(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

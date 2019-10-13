package api

import (
	"encoding/json"
	"net/http"

	"github.com/suggest-go/suggest/pkg/suggest"
)

// dictionaryHandler handles requests with dictionaries purpose
type dictionaryHandler struct {
	suggestService *suggest.Service
}

// handle returns all managed dictionaries by the current suggestService
func (h *dictionaryHandler) handle(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(h.suggestService.GetDictionaries())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

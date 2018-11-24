package api

import (
	"encoding/json"
	"github.com/alldroll/suggest"
	"net/http"
)

//
type dictionaryHandler struct {
	suggestService *suggest.Service
}

//
func (h *dictionaryHandler) handle(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(h.suggestService.GetDictionaries())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

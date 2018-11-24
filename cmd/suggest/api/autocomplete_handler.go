package api

import (
	"encoding/json"
	"github.com/alldroll/suggest"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

//
type autocompleteHandler struct {
	suggestService *suggest.Service
}

//
func (h *autocompleteHandler) handle(w http.ResponseWriter, r *http.Request) {
	var (
		vars  = mux.Vars(r)
		dict  = vars["dict"]
		query = vars["query"]
		k     = r.FormValue("topK")
	)

	type candidates struct {
		Data    []suggest.ResultItem `json:"data"`
		Elapsed string               `json:"elapsed"`
	}

	i64, err := strconv.ParseInt(k, 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	topK := int(i64)
	start := time.Now()
	resultItems, _ := h.suggestService.AutoComplete(dict, query, topK)
	elapsed := time.Since(start).String()

	result := candidates{resultItems, elapsed}
	data, err := json.Marshal(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

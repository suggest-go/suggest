package api

import (
	"encoding/json"
	"github.com/suggest-go/suggest/pkg/spellchecker"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// predictHandler is responsible for query prediction using the spellchecker
type predictHandler struct {
	spellchecker *spellchecker.SpellChecker
}

// handle performs prediction for the provided search query
func (h *predictHandler) handle(w http.ResponseWriter, r *http.Request) {
	var (
		vars  = mux.Vars(r)
		query = vars["query"]
		k     = r.FormValue("topK")
	)

	i64, err := strconv.ParseInt(k, 10, 0)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultItems, err := h.spellchecker.Predict(query, int(i64), 0.3)

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

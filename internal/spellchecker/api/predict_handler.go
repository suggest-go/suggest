package api

import (
	"encoding/json"
	"net/http"

	httputil "github.com/suggest-go/suggest/internal/http"
	"github.com/suggest-go/suggest/pkg/spellchecker"

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
	)

	topK, err := httputil.FormTopKValue(r, "topK", 5)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	similarity, err := httputil.FormSimilarityValue(r, "similarity", 0.5)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultItems, err := h.spellchecker.Predict(query, topK, similarity)

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

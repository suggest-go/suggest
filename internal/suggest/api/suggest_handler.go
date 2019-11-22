package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	httputil "github.com/suggest-go/suggest/internal/http"
	"github.com/suggest-go/suggest/pkg/metric"
	"github.com/suggest-go/suggest/pkg/suggest"
	"net/http"
)

const (
	jaccard = "Jaccard"
	cosine  = "Cosine"
	dice    = "Dice"
	exact   = "Exact"
	overlap = "Overlap"

	defaultSimilarity = 0.5
	defaultTopK = 5
)

var metrics map[string]metric.Metric

func init() {
	metrics = map[string]metric.Metric{
		jaccard: metric.JaccardMetric(),
		cosine:  metric.CosineMetric(),
		dice:    metric.DiceMetric(),
		exact:   metric.ExactMetric(),
		overlap: metric.OverlapMetric(),
	}
}

// suggestHandler responses for handling suggest requests
type suggestHandler struct {
	suggestService *suggest.Service
}

// handle performs topK approximate string search
func (h *suggestHandler) handle(w http.ResponseWriter, r *http.Request) {
	var (
		vars       = mux.Vars(r)
		dict       = vars["dict"]
	)

	searchConf, err := buildSearchConfig(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO return 4** on dictionary not found
	resultItems, err := h.suggestService.Suggest(dict, searchConf)

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

// buildSearchConfig builds a search config for the given list of parameters
func buildSearchConfig(r *http.Request) (suggest.SearchConfig, error) {
	vars := mux.Vars(r)
	topK, err := httputil.FormTopKValue(r, "topK", defaultTopK)

	if err != nil {
		return suggest.SearchConfig{}, err
	}

	metricName := r.FormValue("metric")
	m, ok := metrics[metricName]

	if !ok {
		return suggest.SearchConfig{}, errors.New("metric is not found")
	}

	similarity, err := httputil.FormSimilarityValue(r, "similarity", defaultSimilarity)

	if err != nil {
		return suggest.SearchConfig{}, err
	}

	return suggest.NewSearchConfig(vars["query"], topK, m, similarity)
}

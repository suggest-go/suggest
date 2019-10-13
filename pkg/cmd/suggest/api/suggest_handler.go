package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/suggest-go/suggest/pkg/metric"
	"github.com/suggest-go/suggest/pkg/suggest"
	"github.com/gorilla/mux"
)

const (
	jaccard = "Jaccard"
	cosine  = "Cosine"
	dice    = "Dice"
	exact   = "Exact"
	overlap = "Overlap"
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
		query      = vars["query"]
		metricName = r.FormValue("metric")
		similarity = r.FormValue("similarity")
		k          = r.FormValue("topK")
	)

	i64, err := strconv.ParseInt(k, 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	searchConf, err := buildSearchConfig(query, metricName, similarity, int(i64))

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
	w.Write(data)
}

// buildSearchConfig builds a search config for the given list of parameters
func buildSearchConfig(query, metricName, sim string, k int) (suggest.SearchConfig, error) {
	if _, ok := metrics[metricName]; !ok {
		return suggest.SearchConfig{}, errors.New("Metric not found")
	}

	metric := metrics[metricName]
	similarity, err := strconv.ParseFloat(sim, 64)

	if err != nil {
		return suggest.SearchConfig{}, err
	}

	return suggest.NewSearchConfig(query, k, metric, similarity)
}

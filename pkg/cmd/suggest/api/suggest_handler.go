package api

import (
	"encoding/json"
	"errors"
	"github.com/alldroll/suggest/pkg/metric"
	"github.com/alldroll/suggest/pkg/suggest"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

const (
	Jaccard = "Jaccard"
	Cosine  = "Cosine"
	Dice    = "Dice"
	Exact   = "Exact"
	Overlap = "Overlap"
)

var metrics map[string]metric.Metric

//
func init() {
	metrics = map[string]metric.Metric{
		Jaccard: metric.JaccardMetric(),
		Cosine:  metric.CosineMetric(),
		Dice:    metric.DiceMetric(),
		Exact:   metric.ExactMetric(),
		Overlap: metric.OverlapMetric(),
	}
}

//
type suggestHandler struct {
	suggestService *suggest.Service
}

//
func (h *suggestHandler) handle(w http.ResponseWriter, r *http.Request) {
	var (
		vars       = mux.Vars(r)
		dict       = vars["dict"]
		query      = vars["query"]
		metricName = r.FormValue("metric")
		similarity = r.FormValue("similarity")
		k          = r.FormValue("topK")
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
	searchConf, err := buildSearchConfig(query, metricName, similarity, topK)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	start := time.Now()
	resultItems, err := h.suggestService.Suggest(dict, searchConf)
	elapsed := time.Since(start).String()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := candidates{resultItems, elapsed}
	data, err := json.Marshal(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(data)
}

//
func buildSearchConfig(query, metricName, sim string, k int) (*suggest.SearchConfig, error) {
	if _, ok := metrics[metricName]; !ok {
		return nil, errors.New("Metric not found")
	}

	metric := metrics[metricName]
	similarity, err := strconv.ParseFloat(sim, 64)
	if err != nil {
		return nil, err
	}

	return suggest.NewSearchConfig(query, k, metric, similarity)
}

package api

import (
	"encoding/json"
	"net/http"
)

// homeHandler handles requests for a home page
type homeHandler struct {
}

// handle returns home page status
func (h *homeHandler) handle(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(map[string]string{
		"project": "suggest-demo", // TODO add project name
		"version": "v1",           // TODO add version
	})

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

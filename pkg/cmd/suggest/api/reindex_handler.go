package api

import "net/http"

// reindexHandler is an entity that responsible for handling reindex requests
type reindexHandler struct {
	reindexJob func() error
}

//
func (h *reindexHandler) handle(w http.ResponseWriter, r *http.Request) {
	if err := h.reindexJob(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}

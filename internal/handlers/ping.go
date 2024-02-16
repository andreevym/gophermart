package handlers

import "net/http"

// GetPingHandler api for check is service is avalibale for a work
func (h ServiceHandlers) GetPingHandler(w http.ResponseWriter, r *http.Request) {
	err := h.dbClient.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

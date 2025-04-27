package api

import (
	"net/http"
)

func (h *Handler) GetServerStatus(w http.ResponseWriter, r *http.Request) {
	status := h.statsService.GetServerStatus()
	writeJSONResponse(w, http.StatusOK, status)
}

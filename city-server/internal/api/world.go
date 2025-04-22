package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Получение состояния мира
func (h *Handler) GetWorldState(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	worldId := params["worldId"]
	platform := params["platform"]

	state, err := h.worldService.GetWorldState(worldId, platform)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка получения состояния мира")
		return
	}

	writeJSONResponse(w, http.StatusOK, state)
}

// Сохранение состояния мира
func (h *Handler) SaveWorldState(w http.ResponseWriter, r *http.Request) {
	var inputState map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&inputState); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Ошибка декодирования запроса")
		return
	}

	params := mux.Vars(r)
	worldId := params["worldId"]
	platform := params["platform"]

	err := h.worldService.SaveWorldState(worldId, platform, inputState)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка сохранения состояния мира")
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"status": "success"})
}

// Получение патчей (delta) мира
func (h *Handler) GetWorldDelta(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	worldId := params["worldId"]
	platform := params["platform"]
	lastKnownSnapshotHash := params["lastKnownSnapshotHash"]

	delta, err := h.worldService.GetWorldDelta(worldId, platform, lastKnownSnapshotHash)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка получения патча мира")
		return
	}

	writeJSONResponse(w, http.StatusOK, delta)
}

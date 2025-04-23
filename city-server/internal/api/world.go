package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Получение состояния мира
func (h *Handler) GetWorldState(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	worldId := params["worldId"]
	platform := params["platform"]

	log.Printf("[GetWorldState] Запрос состояния для worldId=%s, platform=%s", worldId, platform)

	state, err := h.worldService.GetWorldState(worldId, platform)
	if err != nil {
		log.Printf("[GetWorldState] Ошибка получения состояния мира worldId=%s: %v", worldId, err)
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка получения состояния мира")
		return
	}

	log.Printf("[GetWorldState] Успешно получили состояние мира worldId=%s, отправляем ответ", worldId)
	writeJSONResponse(w, http.StatusOK, state)
}

// Сохранение состояния мира
func (h *Handler) SaveWorldState(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	worldId := params["worldId"]
	platform := params["platform"]

	log.Printf("[SaveWorldState] Запрос сохранения для worldId=%s, platform=%s", worldId, platform)

	var inputState map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&inputState); err != nil {
		log.Printf("[SaveWorldState] Ошибка декодирования тела запроса для worldId=%s: %v", worldId, err)
		writeErrorResponse(w, http.StatusBadRequest, "Ошибка декодирования запроса")
		return
	}

	log.Printf("[SaveWorldState] Декодирован входной state для worldId=%s: %+v", worldId, inputState)

	err := h.worldService.SaveWorldState(worldId, platform, inputState)
	if err != nil {
		log.Printf("[SaveWorldState] Ошибка сохранения state для worldId=%s: %v", worldId, err)
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка сохранения состояния мира")
		return
	}

	log.Printf("[SaveWorldState] Успешно сохранили состояние мира worldId=%s", worldId)
	writeJSONResponse(w, http.StatusOK, map[string]string{"status": "success"})
}

// Получение патчей (delta) мира
func (h *Handler) GetWorldDelta(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	worldId := params["worldId"]
	platform := params["platform"]
	lastKnownSnapshotHash := params["lastKnownSnapshotHash"]

	log.Printf("[GetWorldDelta] Запрос delta для worldId=%s, platform=%s, lastHash=%s", worldId, platform, lastKnownSnapshotHash)

	delta, err := h.worldService.GetWorldDelta(worldId, platform, lastKnownSnapshotHash)
	if err != nil {
		log.Printf("[GetWorldDelta] Ошибка получения delta для worldId=%s: %v", worldId, err)
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка получения патча мира")
		return
	}

	log.Printf("[GetWorldDelta] Успешно получили delta для worldId=%s, отправляем ответ", worldId)
	writeJSONResponse(w, http.StatusOK, delta)
}

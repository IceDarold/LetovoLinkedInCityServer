package api

import (
	"city-server/internal/services"
	"encoding/json"
	"log"
	"net/http"
)

// Handler для API
type Handler struct {
	worldService *services.WorldService
	assetService *services.AssetService
	authService  *services.AuthService
	statsService *services.StatsService
}

// Новый обработчик
func NewHandler(ws *services.WorldService, as *services.AssetService, auth *services.AuthService, stats *services.StatsService) *Handler {
	return &Handler{
		worldService: ws,
		assetService: as,
		authService:  auth,
		statsService: stats,
	}
}

// Отправка JSON-ответа
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Ошибка отправки JSON-ответа: %s", err)
	}
}

// Обработка ошибок
func writeErrorResponse(w http.ResponseWriter, status int, message string) {
	writeJSONResponse(w, status, map[string]string{"error": message})
}

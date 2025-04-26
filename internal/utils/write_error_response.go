package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

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
func WriteErrorResponse(w http.ResponseWriter, status int, message string) {
	writeJSONResponse(w, status, map[string]string{"error": message})
}

package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// HandleError - общая обработка ошибок
func HandleError(w http.ResponseWriter, err error, statusCode int) {
	log.Printf("Error: %v", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

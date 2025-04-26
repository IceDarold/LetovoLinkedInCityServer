package api

import (
	"log"
	"net/http"
)

func (h *Handler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "Нет заголовка Authorization")
		return
	}

	const prefix = "Bearer "
	if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
		writeErrorResponse(w, http.StatusUnauthorized, "Неверный формат Authorization")
		return
	}

	token := authHeader[len(prefix):]
	log.Printf("🔐 Incoming token: %s", token)

	playerId, valid := h.authService.ValidateToken(token)
	if !valid {
		writeErrorResponse(w, http.StatusUnauthorized, "Неверный токен")
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"valid":    true,
		"playerId": playerId,
	})
}

package api

import (
	"encoding/json"
	"net/http"
)

// Проверка подписи данных
func (h *Handler) ValidateSignature(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Signature string `json:"signature"`
		Payload   string `json:"payload"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Ошибка декодирования запроса")
		return
	}

	valid, err := h.authService.ValidateSignature(requestData.Signature, requestData.Payload)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка валидации подписи")
		return
	}

	if !valid {
		writeErrorResponse(w, http.StatusUnauthorized, "Неверная подпись")
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"status": "valid"})
}

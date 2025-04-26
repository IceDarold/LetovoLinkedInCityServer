package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Получение ассет‑бандла по хэшу
func (h *Handler) GetAssetBundle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	assetBundleHash := params["assetBundleHash"]

	assetBundle, err := h.assetService.GetAssetBundle(assetBundleHash)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка получения ассет‑бандла")
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(assetBundle)
}

// Загрузка нового ассет‑бандла
func (h *Handler) UploadAssetBundle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	worldId := params["worldId"]
	platform := params["platform"]
	assetBundleHash := params["assetBundleHash"]

	// Прочитаем файл из запроса
	assetBundleData := make([]byte, r.ContentLength)
	_, err := r.Body.Read(assetBundleData)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Ошибка чтения данных")
		return
	}

	err = h.assetService.SaveAssetBundle(worldId, platform, assetBundleHash, assetBundleData)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка сохранения ассет‑бандла")
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"status": "success"})
}

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/errors"
)

func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "erro ao gerar resposta JSON", http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, statusCode int, apiErr errors.APIError) {
	WriteJSON(w, statusCode, apiErr)
}

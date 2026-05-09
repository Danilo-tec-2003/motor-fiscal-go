package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/errors"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/middleware"
)

func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		statusCode = http.StatusInternalServerError
		body, _ = json.Marshal(errors.APIError{
			Code:          "INTERNAL_ERROR",
			Message:       "Erro interno inesperado.",
			CorrelationID: "",
			Details:       []errors.ValidationDetail{},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(append(body, '\n'))
}

func WriteError(w http.ResponseWriter, statusCode int, apiErr errors.APIError) {
	if apiErr.Details == nil {
		apiErr.Details = []errors.ValidationDetail{}
	}

	WriteJSON(w, statusCode, apiErr)
}

func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, http.StatusNotFound, errors.APIError{
			Code:          "NOT_FOUND",
			Message:       "Recurso nao encontrado.",
			CorrelationID: middleware.GetCorrelationID(r.Context()),
			Details:       []errors.ValidationDetail{},
		})
	}
}

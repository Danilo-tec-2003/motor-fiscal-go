package handler

import (
	"net/http"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/errors"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/middleware"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

func HealthHandler(serviceName, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			correlationID := middleware.GetCorrelationID(r.Context())
			w.Header().Set("Allow", http.MethodGet)

			WriteError(w, http.StatusMethodNotAllowed, errors.APIError{
				Code:          "METHOD_NOT_ALLOWED",
				Message:       "Metodo nao permitido. Use GET.",
				CorrelationID: correlationID,
				Details:       []errors.ValidationDetail{},
			})
			return
		}

		response := HealthResponse{
			Status:  "UP",
			Service: serviceName,
			Version: version,
		}

		WriteJSON(w, http.StatusOK, response)
	}
}

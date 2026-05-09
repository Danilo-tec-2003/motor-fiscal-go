package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/errors"
)

const APIKeyHeader = "X-API-Key"

func APIKey(requiredKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requiredKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		apiKey := r.Header.Get(APIKeyHeader)
		if apiKey != requiredKey {
			writeUnauthorized(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func writeUnauthorized(w http.ResponseWriter, r *http.Request) {
	correlationID := GetCorrelationID(r.Context())
	if correlationID != "" {
		w.Header().Set(CorrelationIDHeader, correlationID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	_ = json.NewEncoder(w).Encode(errors.APIError{
		Code:          "UNAUTHORIZED",
		Message:       "Token ausente ou invalido.",
		CorrelationID: correlationID,
		Details:       []errors.ValidationDetail{},
	})
}

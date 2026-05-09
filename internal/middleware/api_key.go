package middleware

import (
	"net/http"
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
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

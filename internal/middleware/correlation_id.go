package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const CorrelationIDHeader = "X-Correlation-ID"

type contextKey string

const correlationIDKey contextKey = "correlation_id"

func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = generateCorrelationID()
		}

		w.Header().Set(CorrelationIDHeader, correlationID)

		ctx := context.WithValue(r.Context(), correlationIDKey, correlationID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCorrelationID(ctx context.Context) string {
	correlationID, ok := ctx.Value(correlationIDKey).(string)
	if !ok {
		return ""
	}

	return correlationID
}

func generateCorrelationID() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}

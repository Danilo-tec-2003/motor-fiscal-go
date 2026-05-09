package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apierrors "github.com/Danilo-tec-2003/motor-fiscal-go/internal/errors"
)

func TestAPIKeyReturnsStandardUnauthorizedError(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	handler := CorrelationID(APIKey("dev-token", next))

	request := httptest.NewRequest(http.MethodPost, "/api/v1/tax/simulate", nil)
	request.Header.Set(CorrelationIDHeader, "req-test-401")

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected content type application/json, got %s", contentType)
	}

	var response apierrors.APIError
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("expected valid JSON error response: %v", err)
	}

	if response.Code != "UNAUTHORIZED" {
		t.Fatalf("expected code UNAUTHORIZED, got %s", response.Code)
	}

	if response.CorrelationID != "req-test-401" {
		t.Fatalf("expected correlation id req-test-401, got %s", response.CorrelationID)
	}

	if len(response.Details) != 0 {
		t.Fatalf("expected empty details, got %d item(s)", len(response.Details))
	}
}

func TestAPIKeyAllowsRequestWithValidToken(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})

	handler := CorrelationID(APIKey("dev-token", next))

	request := httptest.NewRequest(http.MethodPost, "/api/v1/tax/simulate", nil)
	request.Header.Set(APIKeyHeader, "dev-token")

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if !called {
		t.Fatal("expected next handler to be called")
	}

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}
}

package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	HealthHandler("motor-fiscal", "1.0.0").ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected content type application/json, got %s", contentType)
	}

	expectedBody := `{"status":"UP","service":"motor-fiscal","version":"1.0.0"}`
	actualBody := strings.TrimSpace(recorder.Body.String())
	if actualBody != expectedBody {
		t.Fatalf("expected body %s, got %s", expectedBody, actualBody)
	}
}

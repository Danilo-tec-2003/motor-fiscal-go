package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/config"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/middleware"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/model"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/service"
	"github.com/shopspring/decimal"
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

func TestHealtHandlerWithCorrelationID(t *testing.T) {
	cfg := config.Config{
		Service: "motor-fiscal",
		Version: "1.0.0",
	}

	router := newTestRouter(cfg)

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set(middleware.CorrelationIDHeader, "req-test-123")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	correlationID := recorder.Header().Get(middleware.CorrelationIDHeader)
	if correlationID != "req-test-123" {
		t.Fatalf("expected correlation id req-test-123, got %s", correlationID)
	}
}

func TestHealthHandlerGeneratesCorrelationID(t *testing.T) {
	cfg := config.Config{
		Service: "motor-fiscal",
		Version: "1.0.0",
	}

	router := newTestRouter(cfg)

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	correlationID := recorder.Header().Get(middleware.CorrelationIDHeader)
	if correlationID == "" {
		t.Fatal("expected generated correlation id, got empty value")
	}
}

func newTestRouter(cfg config.Config) http.Handler {
	ruleService := service.NewFiscalRuleService(fakeFiscalRuleRepository{})
	taxService := service.NewTaxService(ruleService, fakeFiscalSimulationRepository{})
	cteService := service.NewCTeService(ruleService)

	return NewRouter(
		cfg,
		NewTaxHandler(taxService),
		NewCTeHandler(cteService),
	)
}

type fakeFiscalRuleRepository struct{}

func (fakeFiscalRuleRepository) FindActiveRule(ctx context.Context, request dto.TaxSimulationRequest) (model.FiscalRule, error) {
	return model.FiscalRule{
		RuleVersion:   "2026.01",
		OriginUF:      request.OriginUF,
		DestinationUF: request.DestinationUF,
		OperationType: request.OperationType,
		CustomerType:  request.CustomerType,
		ICMSRate:      decimal.RequireFromString("12.00"),
		IBSRate:       decimal.RequireFromString("3.60"),
		CBSRate:       decimal.RequireFromString("0.90"),
		CFOP:          "6351",
		ValidFrom:     "2026-01-01",
		ValidTo:       "2026-12-31",
	}, nil
}

type fakeFiscalSimulationRepository struct{}

func (fakeFiscalSimulationRepository) Save(ctx context.Context, request dto.TaxSimulationRequest, response dto.TaxSimulationResponse) error {
	return nil
}

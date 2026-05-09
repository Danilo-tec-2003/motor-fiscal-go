package handler

import (
	"net/http"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/config"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/middleware"
)

func NewRouter(cfg config.Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", HealthHandler(cfg.Service, cfg.Version))
	mux.HandleFunc("/api/v1/tax/simulate", TaxSimulationHandler())
	mux.HandleFunc("/api/v1/cte/validate", CTeValidationHandler())
	mux.HandleFunc("/api/v1/tax/compare", TaxComparisonHandler())
	mux.HandleFunc("/api/v1/tax/batch", TaxBatchHandler())

	return middleware.CorrelationID(mux)
}

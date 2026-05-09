package handler

import (
	"net/http"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/config"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/middleware"
)

func NewRouter(cfg config.Config, taxHandler TaxHandler, cteHandler CTeHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", HealthHandler(cfg.Service, cfg.Version))
	mux.HandleFunc("/api/v1/tax/simulate", taxHandler.Simulate())
	mux.HandleFunc("/api/v1/tax/compare", taxHandler.Compare())
	mux.HandleFunc("/api/v1/tax/batch", taxHandler.Batch())
	mux.HandleFunc("/api/v1/cte/validate", cteHandler.Validate())

	return middleware.CorrelationID(mux)
}

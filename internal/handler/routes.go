package handler

import (
	"net/http"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/config"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/middleware"
)

func NewRouter(cfg config.Config, taxHandler TaxHandler, cteHandler CTeHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", HealthHandler(cfg.Service, cfg.Version))

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/v1/tax/simulate", taxHandler.Simulate())
	protectedMux.HandleFunc("/api/v1/tax/compare", taxHandler.Compare())
	protectedMux.HandleFunc("/api/v1/tax/batch", taxHandler.Batch())
	protectedMux.HandleFunc("/api/v1/cte/validate", cteHandler.Validate())

	mux.Handle("/api/v1/", middleware.APIKey(cfg.InternalAPIKey, protectedMux))

	return middleware.CorrelationID(mux)
}

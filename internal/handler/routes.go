package handler

import (
	"net/http"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/config"
)

func NewRouter(cfg config.Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthHandler(cfg.Service, cfg.Version))

	return mux
}

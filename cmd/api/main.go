package main

import (
	"fmt"
	"net/http"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/config"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/handler"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/service"
)

func main() {
	cfg := config.Load()

	taxService := service.NewTaxService()
	cteService := service.NewCTeService()

	taxHandler := handler.NewTaxHandler(taxService)
	cteHandler := handler.NewCTeHandler(cteService)

	router := handler.NewRouter(cfg, taxHandler, cteHandler)

	addr := ":" + cfg.Port

	fmt.Printf("Motor Fiscal API iniciado na porta %s\n", cfg.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
	}
}

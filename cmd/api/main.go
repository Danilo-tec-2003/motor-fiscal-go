package main

import (
	"fmt"
	"net/http"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/config"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/handler"
)

func main() {
	cfg := config.Load()
	router := handler.NewRouter(cfg)
	addr := ":" + cfg.Port

	fmt.Printf("Motor Fiscal API iniciado na porta %s\n", cfg.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
	}
}

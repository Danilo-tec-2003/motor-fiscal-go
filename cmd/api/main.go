package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/config"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/handler"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/repository"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Erro ao criar pool PostgreSQL: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Erro ao conectar no PostgreSQL: %v", err)
	}

	fiscalRuleRepository := repository.NewFiscalRuleRepository(dbPool)
	simulationRepository := repository.NewFiscalSimulationRepository(dbPool)

	fiscalRuleService := service.NewFiscalRuleService(fiscalRuleRepository)
	taxService := service.NewTaxService(fiscalRuleService, simulationRepository)
	cteService := service.NewCTeService(fiscalRuleService)

	taxHandler := handler.NewTaxHandler(taxService)
	cteHandler := handler.NewCTeHandler(cteService)

	router := handler.NewRouter(cfg, taxHandler, cteHandler)

	addr := ":" + cfg.Port

	fmt.Printf("Motor Fiscal API iniciado na porta %s\n", cfg.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
	}
}

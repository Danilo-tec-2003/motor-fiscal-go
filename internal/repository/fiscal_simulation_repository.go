package repository

import (
	"context"
	"fmt"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FiscalSimulationRepository struct {
	db *pgxpool.Pool
}

func NewFiscalSimulationRepository(db *pgxpool.Pool) FiscalSimulationRepository {
	return FiscalSimulationRepository{db: db}
}

func (r FiscalSimulationRepository) Save(ctx context.Context, request dto.TaxSimulationRequest, response dto.TaxSimulationResponse) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO fiscal_simulations (
			freight_id,
			operation_date,
			origin_uf,
			destination_uf,
			freight_value,
			icms_rate,
			icms_amount,
			ibs_rate,
			ibs_amount,
			cbs_rate,
			cbs_amount,
			total_tax,
			total_with_tax,
			cfop,
			rule_version,
			rule_id,
			rule_code,
			rule_status,
			calculation_basis,
			from_cache
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16,
			$17, $18, $19, $20
		)
		RETURNING id
	`

	var simulationID int64
	err = tx.QueryRow(
		ctx,
		query,
		request.FreightID,
		request.OperationDate,
		request.OriginUF,
		request.DestinationUF,
		request.FreightValue,
		response.ICMS.Rate,
		response.ICMS.Amount,
		response.IBS.Rate,
		response.IBS.Amount,
		response.CBS.Rate,
		response.CBS.Amount,
		response.TotalTax,
		response.TotalWithTax,
		response.CFOP,
		response.RuleVersion,
		response.RuleID,
		response.RuleCode,
		response.RuleStatus,
		response.CalculationBasis,
		response.FromCache,
	).Scan(&simulationID)
	if err != nil {
		return err
	}

	for _, detail := range response.CalculationDetails {
		if err := insertSimulationTaxDetail(ctx, tx, simulationID, detail); err != nil {
			return fmt.Errorf("insert fiscal simulation tax detail: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func insertSimulationTaxDetail(ctx context.Context, tx pgx.Tx, simulationID int64, detail dto.TaxCalculationDetail) error {
	query := `
		INSERT INTO fiscal_simulation_tax_details (
			fiscal_simulation_id,
			tax_name,
			base_value,
			base_reduction_rate,
			effective_base_value,
			rate,
			amount,
			formula
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`

	_, err := tx.Exec(
		ctx,
		query,
		simulationID,
		detail.TaxName,
		detail.BaseValue,
		detail.BaseReductionRate,
		detail.EffectiveBaseValue,
		detail.Rate,
		detail.Amount,
		detail.Formula,
	)

	return err
}

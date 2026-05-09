package repository

import (
	"context"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FiscalSimulationRepository struct {
	db *pgxpool.Pool
}

func NewFiscalSimulationRepository(db *pgxpool.Pool) FiscalSimulationRepository {
	return FiscalSimulationRepository{db: db}
}

func (r FiscalSimulationRepository) Save(ctx context.Context, request dto.TaxSimulationRequest, response dto.TaxSimulationResponse) error {
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
			from_cache
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16
		)
	`

	_, err := r.db.Exec(
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
		response.FromCache,
	)

	return err
}

package repository

import (
	"context"
	"errors"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

var ErrFiscalRuleNotFound = errors.New("fiscal rule not found")

type FiscalRuleRepository struct {
	db *pgxpool.Pool
}

func NewFiscalRuleRepository(db *pgxpool.Pool) FiscalRuleRepository {
	return FiscalRuleRepository{db: db}
}

func (r FiscalRuleRepository) FindActiveRule(ctx context.Context, request dto.TaxSimulationRequest) (model.FiscalRule, error) {
	query := `
		SELECT
			rule_version,
			origin_uf,
			destination_uf,
			operation_type,
			customer_type,
			icms_rate::text,
			ibs_rate::text,
			cbs_rate::text,
			cfop,
			to_char(valid_from, 'YYYY-MM-DD'),
			to_char(valid_to, 'YYYY-MM-DD')
		FROM fiscal_rules
		WHERE origin_uf = $1
		  AND destination_uf = $2
		  AND operation_type = $3
		  AND customer_type = $4
		  AND active = TRUE
		  AND $5::date BETWEEN valid_from AND valid_to
		ORDER BY valid_from DESC
		LIMIT 1
	`

	var rule model.FiscalRule
	var icmsRate string
	var ibsRate string
	var cbsRate string

	err := r.db.QueryRow(
		ctx,
		query,
		request.OriginUF,
		request.DestinationUF,
		request.OperationType,
		request.CustomerType,
		request.OperationDate,
	).Scan(
		&rule.RuleVersion,
		&rule.OriginUF,
		&rule.DestinationUF,
		&rule.OperationType,
		&rule.CustomerType,
		&icmsRate,
		&ibsRate,
		&cbsRate,
		&rule.CFOP,
		&rule.ValidFrom,
		&rule.ValidTo,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.FiscalRule{}, ErrFiscalRuleNotFound
		}

		return model.FiscalRule{}, err
	}

	rule.ICMSRate, err = decimal.NewFromString(icmsRate)
	if err != nil {
		return model.FiscalRule{}, err
	}

	rule.IBSRate, err = decimal.NewFromString(ibsRate)
	if err != nil {
		return model.FiscalRule{}, err
	}

	rule.CBSRate, err = decimal.NewFromString(cbsRate)
	if err != nil {
		return model.FiscalRule{}, err
	}

	return rule, nil
}

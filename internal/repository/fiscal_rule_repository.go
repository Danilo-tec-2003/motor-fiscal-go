package repository

import (
	"context"
	"errors"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/model"
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

func (r FiscalRuleRepository) FindCandidateRules(ctx context.Context, request dto.TaxSimulationRequest) ([]model.FiscalRuleEngine, error) {
	query := `
		SELECT
			id,
			rule_code,
			rule_version,
			description,
			priority,
			status,
			calculation_basis,
			cfop,
			to_char(valid_from, 'YYYY-MM-DD'),
			to_char(valid_to, 'YYYY-MM-DD'),
			active
		FROM fiscal_rules
		WHERE active = TRUE
		  AND status IN ('APPROVED', 'PENDING_REVIEW')
		  AND $1::date BETWEEN valid_from AND valid_to
		ORDER BY priority ASC, valid_from DESC, id ASC
	`

	rows, err := r.db.Query(ctx, query, request.OperationDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rules := []model.FiscalRuleEngine{}
	for rows.Next() {
		var rule model.FiscalRuleEngine

		err := rows.Scan(
			&rule.ID,
			&rule.RuleCode,
			&rule.RuleVersion,
			&rule.Description,
			&rule.Priority,
			&rule.Status,
			&rule.CalculationBasis,
			&rule.CFOP,
			&rule.ValidFrom,
			&rule.ValidTo,
			&rule.Active,
		)
		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(rules) == 0 {
		return nil, ErrFiscalRuleNotFound
	}

	ruleIDs := make([]int64, 0, len(rules))
	ruleIndex := make(map[int64]int, len(rules))
	for index, rule := range rules {
		ruleIDs = append(ruleIDs, rule.ID)
		ruleIndex[rule.ID] = index
	}

	conditions, err := r.findRuleConditions(ctx, ruleIDs)
	if err != nil {
		return nil, err
	}

	for _, condition := range conditions {
		index, ok := ruleIndex[condition.FiscalRuleID]
		if !ok {
			continue
		}

		rules[index].Conditions = append(rules[index].Conditions, condition)
	}

	taxes, err := r.findRuleTaxes(ctx, ruleIDs)
	if err != nil {
		return nil, err
	}

	for _, tax := range taxes {
		index, ok := ruleIndex[tax.FiscalRuleID]
		if !ok {
			continue
		}

		rules[index].Taxes = append(rules[index].Taxes, tax)
	}

	return rules, nil
}

func (r FiscalRuleRepository) findRuleConditions(ctx context.Context, ruleIDs []int64) ([]model.FiscalRuleCondition, error) {
	query := `
		SELECT
			id,
			fiscal_rule_id,
			field_name,
			operator,
			field_value
		FROM fiscal_rule_conditions
		WHERE fiscal_rule_id = ANY($1)
		ORDER BY fiscal_rule_id, id
	`

	rows, err := r.db.Query(ctx, query, ruleIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conditions := []model.FiscalRuleCondition{}
	for rows.Next() {
		var condition model.FiscalRuleCondition

		err := rows.Scan(
			&condition.ID,
			&condition.FiscalRuleID,
			&condition.FieldName,
			&condition.Operator,
			&condition.FieldValue,
		)
		if err != nil {
			return nil, err
		}

		conditions = append(conditions, condition)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return conditions, nil
}

func (r FiscalRuleRepository) findRuleTaxes(ctx context.Context, ruleIDs []int64) ([]model.FiscalRuleTax, error) {
	query := `
		SELECT
			id,
			fiscal_rule_id,
			tax_name,
			rate::text,
			base_reduction_rate::text,
			calculation_order
		FROM fiscal_rule_taxes
		WHERE fiscal_rule_id = ANY($1)
		ORDER BY fiscal_rule_id, calculation_order, id
	`

	rows, err := r.db.Query(ctx, query, ruleIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	taxes := []model.FiscalRuleTax{}
	for rows.Next() {
		var tax model.FiscalRuleTax
		var rate string
		var baseReductionRate string

		err := rows.Scan(
			&tax.ID,
			&tax.FiscalRuleID,
			&tax.TaxName,
			&rate,
			&baseReductionRate,
			&tax.CalculationOrder,
		)
		if err != nil {
			return nil, err
		}

		tax.Rate, err = decimal.NewFromString(rate)
		if err != nil {
			return nil, err
		}

		tax.BaseReductionRate, err = decimal.NewFromString(baseReductionRate)
		if err != nil {
			return nil, err
		}

		taxes = append(taxes, tax)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return taxes, nil
}

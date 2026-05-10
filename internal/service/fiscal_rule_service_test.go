package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/model"
	"github.com/shopspring/decimal"
)

func TestFiscalRuleServiceFindRuleMatchesConditionsAndReturnsTaxes(t *testing.T) {
	repository := fakeRuleRepository{
		rules: []model.FiscalRuleEngine{
			buildRuleEngine(1, "RULE_PE_SP", 100, "2026-01-01", []model.FiscalRuleCondition{
				condition(1, "origin_uf", model.RuleConditionOperatorEquals, "PE"),
				condition(1, "destination_uf", model.RuleConditionOperatorEquals, "SP"),
				condition(1, "operation_type", model.RuleConditionOperatorEquals, "INTERESTADUAL"),
				condition(1, "customer_type", model.RuleConditionOperatorEquals, "PJ"),
			}),
		},
	}

	service := NewFiscalRuleService(repository)

	rule, err := service.FindRule(context.Background(), defaultTaxRequest())
	if err != nil {
		t.Fatalf("expected matched rule, got error: %v", err)
	}

	if rule.RuleCode != "RULE_PE_SP" {
		t.Fatalf("expected rule code RULE_PE_SP, got %s", rule.RuleCode)
	}

	if rule.ICMSRate.StringFixed(2) != "12.00" {
		t.Fatalf("expected ICMS rate 12.00, got %s", rule.ICMSRate.StringFixed(2))
	}

	if len(rule.Taxes) != 3 {
		t.Fatalf("expected 3 taxes, got %d", len(rule.Taxes))
	}
}

func TestFiscalRuleServiceFindRuleReturnsNotFoundWhenConditionsDoNotMatch(t *testing.T) {
	repository := fakeRuleRepository{
		rules: []model.FiscalRuleEngine{
			buildRuleEngine(1, "RULE_BA_SP", 100, "2026-01-01", []model.FiscalRuleCondition{
				condition(1, "origin_uf", model.RuleConditionOperatorEquals, "BA"),
			}),
		},
	}

	service := NewFiscalRuleService(repository)

	_, err := service.FindRule(context.Background(), defaultTaxRequest())
	if !errors.Is(err, ErrFiscalRuleNotFound) {
		t.Fatalf("expected ErrFiscalRuleNotFound, got %v", err)
	}
}

func TestFiscalRuleServiceFindRuleUsesLowestPriorityNumber(t *testing.T) {
	repository := fakeRuleRepository{
		rules: []model.FiscalRuleEngine{
			buildRuleEngine(1, "RULE_GENERIC", 100, "2026-01-01", []model.FiscalRuleCondition{
				condition(1, "origin_uf", model.RuleConditionOperatorEquals, "PE"),
			}),
			buildRuleEngine(2, "RULE_SPECIFIC", 50, "2026-01-01", []model.FiscalRuleCondition{
				condition(2, "origin_uf", model.RuleConditionOperatorEquals, "PE"),
				condition(2, "destination_uf", model.RuleConditionOperatorEquals, "SP"),
			}),
		},
	}

	service := NewFiscalRuleService(repository)

	rule, err := service.FindRule(context.Background(), defaultTaxRequest())
	if err != nil {
		t.Fatalf("expected matched rule, got error: %v", err)
	}

	if rule.RuleCode != "RULE_SPECIFIC" {
		t.Fatalf("expected RULE_SPECIFIC, got %s", rule.RuleCode)
	}
}

func TestFiscalRuleServiceFindRuleReturnsConflictWhenRulesHaveSamePriorityAndVigency(t *testing.T) {
	repository := fakeRuleRepository{
		rules: []model.FiscalRuleEngine{
			buildRuleEngine(1, "RULE_A", 100, "2026-01-01", []model.FiscalRuleCondition{
				condition(1, "origin_uf", model.RuleConditionOperatorEquals, "PE"),
			}),
			buildRuleEngine(2, "RULE_B", 100, "2026-01-01", []model.FiscalRuleCondition{
				condition(2, "origin_uf", model.RuleConditionOperatorEquals, "PE"),
			}),
		},
	}

	service := NewFiscalRuleService(repository)

	_, err := service.FindRule(context.Background(), defaultTaxRequest())
	if !errors.Is(err, ErrFiscalRuleConflict) {
		t.Fatalf("expected ErrFiscalRuleConflict, got %v", err)
	}
}

func TestFiscalRuleServiceFindRuleReturnsIncompleteWhenRequiredTaxIsMissing(t *testing.T) {
	rule := buildRuleEngine(1, "RULE_INCOMPLETE", 100, "2026-01-01", []model.FiscalRuleCondition{
		condition(1, "origin_uf", model.RuleConditionOperatorEquals, "PE"),
	})
	rule.Taxes = []model.FiscalRuleTax{
		tax(1, model.TaxNameICMS, "12.00"),
		tax(1, model.TaxNameIBS, "3.60"),
	}

	repository := fakeRuleRepository{
		rules: []model.FiscalRuleEngine{rule},
	}

	service := NewFiscalRuleService(repository)

	_, err := service.FindRule(context.Background(), defaultTaxRequest())
	if !errors.Is(err, ErrFiscalRuleIncomplete) {
		t.Fatalf("expected ErrFiscalRuleIncomplete, got %v", err)
	}
}

func TestFiscalRuleServiceFindRuleSupportsInAndBetweenOperators(t *testing.T) {
	repository := fakeRuleRepository{
		rules: []model.FiscalRuleEngine{
			buildRuleEngine(1, "RULE_WITH_OPERATORS", 100, "2026-01-01", []model.FiscalRuleCondition{
				condition(1, "customer_type", model.RuleConditionOperatorIn, "PF, PJ"),
				condition(1, "operation_date", model.RuleConditionOperatorBetween, "2026-01-01..2026-12-31"),
			}),
		},
	}

	service := NewFiscalRuleService(repository)

	rule, err := service.FindRule(context.Background(), defaultTaxRequest())
	if err != nil {
		t.Fatalf("expected matched rule, got error: %v", err)
	}

	if rule.RuleCode != "RULE_WITH_OPERATORS" {
		t.Fatalf("expected RULE_WITH_OPERATORS, got %s", rule.RuleCode)
	}
}

type fakeRuleRepository struct {
	rules []model.FiscalRuleEngine
	err   error
}

func (r fakeRuleRepository) FindCandidateRules(ctx context.Context, request dto.TaxSimulationRequest) ([]model.FiscalRuleEngine, error) {
	if r.err != nil {
		return nil, r.err
	}

	return r.rules, nil
}

func defaultTaxRequest() dto.TaxSimulationRequest {
	return dto.TaxSimulationRequest{
		FreightID:     1001,
		OperationDate: "2026-03-10",
		OriginUF:      "PE",
		DestinationUF: "SP",
		FreightValue:  "3500.00",
		CustomerType:  "PJ",
		OperationType: "INTERESTADUAL",
	}
}

func buildRuleEngine(id int64, code string, priority int, validFrom string, conditions []model.FiscalRuleCondition) model.FiscalRuleEngine {
	return model.FiscalRuleEngine{
		ID:               id,
		RuleCode:         code,
		RuleVersion:      "2026.01",
		Description:      "Regra de teste",
		Priority:         priority,
		Status:           model.FiscalRuleStatusApproved,
		CalculationBasis: model.CalculationBasisFreightValue,
		CFOP:             "6351",
		ValidFrom:        validFrom,
		ValidTo:          "2026-12-31",
		Active:           true,
		Conditions:       conditions,
		Taxes: []model.FiscalRuleTax{
			tax(id, model.TaxNameICMS, "12.00"),
			tax(id, model.TaxNameIBS, "3.60"),
			tax(id, model.TaxNameCBS, "0.90"),
		},
	}
}

func condition(ruleID int64, fieldName string, operator string, fieldValue string) model.FiscalRuleCondition {
	return model.FiscalRuleCondition{
		FiscalRuleID: ruleID,
		FieldName:    fieldName,
		Operator:     operator,
		FieldValue:   fieldValue,
	}
}

func tax(ruleID int64, taxName string, rate string) model.FiscalRuleTax {
	return model.FiscalRuleTax{
		FiscalRuleID:      ruleID,
		TaxName:           taxName,
		Rate:              decimal.RequireFromString(rate),
		BaseReductionRate: decimal.Zero,
	}
}

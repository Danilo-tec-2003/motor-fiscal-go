package service

import (
	"errors"
	"testing"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/model"
	"github.com/shopspring/decimal"
)

func TestCalculateTaxesReturnsAuditableCalculationDetails(t *testing.T) {
	rule := model.FiscalRule{
		ID:               1,
		RuleCode:         "RULE_PE_SP",
		RuleVersion:      "2026.01",
		Status:           model.FiscalRuleStatusApproved,
		CalculationBasis: model.CalculationBasisFreightValue,
		CFOP:             "6351",
		Taxes: []model.FiscalRuleTax{
			{
				FiscalRuleID:      1,
				TaxName:           model.TaxNameICMS,
				Rate:              decimal.RequireFromString("12.00"),
				BaseReductionRate: decimal.RequireFromString("10.00"),
				CalculationOrder:  1,
			},
			{
				FiscalRuleID:      1,
				TaxName:           model.TaxNameIBS,
				Rate:              decimal.RequireFromString("3.60"),
				BaseReductionRate: decimal.Zero,
				CalculationOrder:  2,
			},
			{
				FiscalRuleID:      1,
				TaxName:           model.TaxNameCBS,
				Rate:              decimal.RequireFromString("0.90"),
				BaseReductionRate: decimal.Zero,
				CalculationOrder:  3,
			},
		},
	}

	calculation, err := CalculateTaxes(decimal.RequireFromString("3500.00"), rule)
	if err != nil {
		t.Fatalf("expected calculation, got error: %v", err)
	}

	if calculation.ICMS.Amount != "378.00" {
		t.Fatalf("expected ICMS amount 378.00, got %s", calculation.ICMS.Amount)
	}

	if calculation.IBS.Amount != "126.00" {
		t.Fatalf("expected IBS amount 126.00, got %s", calculation.IBS.Amount)
	}

	if calculation.CBS.Amount != "31.50" {
		t.Fatalf("expected CBS amount 31.50, got %s", calculation.CBS.Amount)
	}

	if calculation.TotalTax.StringFixedBank(2) != "535.50" {
		t.Fatalf("expected total tax 535.50, got %s", calculation.TotalTax.StringFixedBank(2))
	}

	if calculation.TotalWithTax.StringFixedBank(2) != "4035.50" {
		t.Fatalf("expected total with tax 4035.50, got %s", calculation.TotalWithTax.StringFixedBank(2))
	}

	if len(calculation.CalculationDetails) != 3 {
		t.Fatalf("expected 3 calculation details, got %d", len(calculation.CalculationDetails))
	}

	icmsDetail := calculation.CalculationDetails[0]
	if icmsDetail.TaxName != model.TaxNameICMS {
		t.Fatalf("expected first detail to be ICMS, got %s", icmsDetail.TaxName)
	}

	if icmsDetail.BaseReductionRate != "10.00" {
		t.Fatalf("expected ICMS base reduction 10.00, got %s", icmsDetail.BaseReductionRate)
	}

	if icmsDetail.EffectiveBaseValue != "3150.00" {
		t.Fatalf("expected ICMS effective base 3150.00, got %s", icmsDetail.EffectiveBaseValue)
	}

	if icmsDetail.Formula != "effective_base_value * rate / 100" {
		t.Fatalf("unexpected formula: %s", icmsDetail.Formula)
	}
}

func TestCalculateTaxesReturnsErrorForUnsupportedCalculationBasis(t *testing.T) {
	rule := model.FiscalRule{
		CalculationBasis: "UNKNOWN_BASIS",
		Taxes: []model.FiscalRuleTax{
			tax(1, model.TaxNameICMS, "12.00"),
			tax(1, model.TaxNameIBS, "3.60"),
			tax(1, model.TaxNameCBS, "0.90"),
		},
	}

	_, err := CalculateTaxes(decimal.RequireFromString("3500.00"), rule)
	if !errors.Is(err, ErrUnsupportedCalculationBasis) {
		t.Fatalf("expected ErrUnsupportedCalculationBasis, got %v", err)
	}
}

func TestCalculateTaxesReturnsIncompleteWhenRequiredTaxIsMissing(t *testing.T) {
	rule := model.FiscalRule{
		CalculationBasis: model.CalculationBasisFreightValue,
		Taxes: []model.FiscalRuleTax{
			tax(1, model.TaxNameICMS, "12.00"),
			tax(1, model.TaxNameIBS, "3.60"),
		},
	}

	_, err := CalculateTaxes(decimal.RequireFromString("3500.00"), rule)
	if !errors.Is(err, ErrFiscalRuleIncomplete) {
		t.Fatalf("expected ErrFiscalRuleIncomplete, got %v", err)
	}
}

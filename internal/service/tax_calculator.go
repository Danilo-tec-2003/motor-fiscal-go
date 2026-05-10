package service

import (
	"errors"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/model"
	"github.com/shopspring/decimal"
)

var ErrUnsupportedCalculationBasis = errors.New("unsupported calculation basis")

type TaxCalculation struct {
	BaseValue          decimal.Decimal
	ICMS               dto.TaxAmount
	IBS                dto.TaxAmount
	CBS                dto.TaxAmount
	TotalTax           decimal.Decimal
	TotalWithTax       decimal.Decimal
	CalculationDetails []dto.TaxCalculationDetail
}

func CalculateTaxes(baseValue decimal.Decimal, rule model.FiscalRule) (TaxCalculation, error) {
	if rule.CalculationBasis != "" && rule.CalculationBasis != model.CalculationBasisFreightValue {
		return TaxCalculation{}, ErrUnsupportedCalculationBasis
	}

	taxes := rule.Taxes
	if len(taxes) == 0 {
		taxes = legacyRuleTaxes(rule)
	}

	calculation := TaxCalculation{
		BaseValue:          baseValue,
		CalculationDetails: []dto.TaxCalculationDetail{},
	}

	for _, tax := range taxes {
		effectiveBaseValue := applyBaseReduction(baseValue, tax.BaseReductionRate)
		amount := effectiveBaseValue.Mul(tax.Rate).Div(decimal.NewFromInt(100)).RoundBank(2)

		taxAmount := dto.TaxAmount{
			Rate:   formatRate(tax.Rate),
			Amount: formatMoney(amount),
		}

		switch tax.TaxName {
		case model.TaxNameICMS:
			calculation.ICMS = taxAmount
		case model.TaxNameIBS:
			calculation.IBS = taxAmount
		case model.TaxNameCBS:
			calculation.CBS = taxAmount
		default:
			continue
		}

		calculation.TotalTax = calculation.TotalTax.Add(amount)
		calculation.CalculationDetails = append(calculation.CalculationDetails, dto.TaxCalculationDetail{
			TaxName:            tax.TaxName,
			BaseValue:          formatMoney(baseValue),
			BaseReductionRate:  formatRate(tax.BaseReductionRate),
			EffectiveBaseValue: formatMoney(effectiveBaseValue),
			Rate:               formatRate(tax.Rate),
			Amount:             formatMoney(amount),
			Formula:            "effective_base_value * rate / 100",
		})
	}

	if calculation.ICMS.Rate == "" || calculation.IBS.Rate == "" || calculation.CBS.Rate == "" {
		return TaxCalculation{}, ErrFiscalRuleIncomplete
	}

	calculation.TotalWithTax = baseValue.Add(calculation.TotalTax)

	return calculation, nil
}

func applyBaseReduction(baseValue decimal.Decimal, baseReductionRate decimal.Decimal) decimal.Decimal {
	if baseReductionRate.LessThanOrEqual(decimal.Zero) {
		return baseValue
	}

	reductionAmount := baseValue.Mul(baseReductionRate).Div(decimal.NewFromInt(100))
	return baseValue.Sub(reductionAmount).RoundBank(2)
}

func legacyRuleTaxes(rule model.FiscalRule) []model.FiscalRuleTax {
	return []model.FiscalRuleTax{
		{
			FiscalRuleID:      rule.ID,
			TaxName:           model.TaxNameICMS,
			Rate:              rule.ICMSRate,
			BaseReductionRate: decimal.Zero,
			CalculationOrder:  1,
		},
		{
			FiscalRuleID:      rule.ID,
			TaxName:           model.TaxNameIBS,
			Rate:              rule.IBSRate,
			BaseReductionRate: decimal.Zero,
			CalculationOrder:  2,
		},
		{
			FiscalRuleID:      rule.ID,
			TaxName:           model.TaxNameCBS,
			Rate:              rule.CBSRate,
			BaseReductionRate: decimal.Zero,
			CalculationOrder:  3,
		},
	}
}

func formatMoney(value decimal.Decimal) string {
	return value.StringFixedBank(2)
}

func formatRate(value decimal.Decimal) string {
	return value.StringFixed(2)
}

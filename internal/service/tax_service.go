package service

import (
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	"github.com/shopspring/decimal"
)

type TaxService struct{}

func NewTaxService() TaxService {
	return TaxService{}
}

func (s TaxService) Simulate(request dto.TaxSimulationRequest) (dto.TaxSimulationResponse, error) {
	freightValue, err := decimal.NewFromString(request.FreightValue)
	if err != nil {
		return dto.TaxSimulationResponse{}, err
	}

	icmsRate := decimal.NewFromFloat(12.00)
	ibsRate := decimal.NewFromFloat(3.60)
	cbsRate := decimal.NewFromFloat(0.90)

	icmsAmount := calculateTaxAmount(freightValue, icmsRate)
	ibsAmount := calculateTaxAmount(freightValue, ibsRate)
	cbsAmount := calculateTaxAmount(freightValue, cbsRate)

	totalTax := icmsAmount.Add(ibsAmount).Add(cbsAmount)
	totalWithTax := freightValue.Add(totalTax)

	return dto.TaxSimulationResponse{
		FreightID: request.FreightID,
		BaseValue: formatMoney(freightValue),
		ICMS: dto.TaxAmount{
			Rate:   formatRate(icmsRate),
			Amount: formatMoney(icmsAmount),
		},
		IBS: dto.TaxAmount{
			Rate:   formatRate(ibsRate),
			Amount: formatMoney(ibsAmount),
		},
		CBS: dto.TaxAmount{
			Rate:   formatRate(cbsRate),
			Amount: formatMoney(cbsAmount),
		},
		TotalTax:     formatMoney(totalTax),
		TotalWithTax: formatMoney(totalWithTax),
		CFOP:         "6351",
		RuleVersion:  "2026.01",
		FromCahe:     false,
	}, nil
}

func calculateTaxAmount(baseValue decimal.Decimal, rate decimal.Decimal) decimal.Decimal {
	return baseValue.Mul(rate).Div(decimal.NewFromInt(100)).RoundBank(2)
}

func formatMoney(value decimal.Decimal) string {
	return value.StringFixedBank(2)
}

func formatRate(value decimal.Decimal) string {
	return value.StringFixed(2)
}

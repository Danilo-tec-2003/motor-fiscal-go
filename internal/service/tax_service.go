package service

import (
	"context"
	"fmt"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	"github.com/shopspring/decimal"
)

type simulationRepository interface {
	Save(ctx context.Context, request dto.TaxSimulationRequest, response dto.TaxSimulationResponse) error
}

type TaxService struct {
	ruleService          FiscalRuleService
	simulationRepository simulationRepository
}

func NewTaxService(ruleService FiscalRuleService, simulationRepository simulationRepository) TaxService {
	return TaxService{
		ruleService:          ruleService,
		simulationRepository: simulationRepository,
	}
}

func (s TaxService) Simulate(ctx context.Context, request dto.TaxSimulationRequest) (dto.TaxSimulationResponse, error) {
	freightValue, err := decimal.NewFromString(request.FreightValue)
	if err != nil {
		return dto.TaxSimulationResponse{}, err
	}

	rule, err := s.ruleService.FindRule(ctx, request)
	if err != nil {
		return dto.TaxSimulationResponse{}, err
	}

	icmsRate := rule.ICMSRate
	ibsRate := rule.IBSRate
	cbsRate := rule.CBSRate

	icmsAmount := calculateTaxAmount(freightValue, icmsRate)
	ibsAmount := calculateTaxAmount(freightValue, ibsRate)
	cbsAmount := calculateTaxAmount(freightValue, cbsRate)

	totalTax := icmsAmount.Add(ibsAmount).Add(cbsAmount)
	totalWithTax := freightValue.Add(totalTax)

	response := dto.TaxSimulationResponse{
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
		CFOP:         rule.CFOP,
		RuleVersion:  rule.RuleVersion,
		FromCache:    false,
	}

	if err := s.simulationRepository.Save(ctx, request, response); err != nil {
		return dto.TaxSimulationResponse{}, fmt.Errorf("save fiscal simulation: %w", err)
	}

	return response, nil
}

func (s TaxService) Compare(ctx context.Context, request dto.TaxSimulationRequest) (dto.TaxComparisonResponse, error) {
	freightValue, err := decimal.NewFromString(request.FreightValue)
	if err != nil {
		return dto.TaxComparisonResponse{}, err
	}

	rule, err := s.ruleService.FindRule(ctx, request)
	if err != nil {
		return dto.TaxComparisonResponse{}, err
	}

	icmsAmount := calculateTaxAmount(freightValue, rule.ICMSRate)
	ibsAmount := calculateTaxAmount(freightValue, rule.IBSRate)
	cbsAmount := calculateTaxAmount(freightValue, rule.CBSRate)

	currentTotalTax := icmsAmount
	currentTotalWithTax := freightValue.Add(currentTotalTax)

	reformTotalTax := icmsAmount.Add(ibsAmount).Add(cbsAmount)
	reformTotalWithTax := freightValue.Add(reformTotalTax)

	difference := reformTotalTax.Sub(currentTotalTax)

	analysis := "Cenario da reforma manteve a mesma carga tributaria estimada."
	if difference.GreaterThan(decimal.Zero) {
		analysis = "Cenario da reforma apresentou aumento estimado de tributos."
	} else if difference.LessThan(decimal.Zero) {
		analysis = "Cenario da reforma apresentou reducao estimada de tributos."
	}

	return dto.TaxComparisonResponse{
		FreightID: request.FreightID,
		CurrentScenario: dto.TaxScenario{
			BaseValue:    formatMoney(freightValue),
			TotalTax:     formatMoney(currentTotalTax),
			TotalWithTax: formatMoney(currentTotalWithTax),
		},
		ReformScenario: dto.TaxScenario{
			BaseValue:    formatMoney(freightValue),
			TotalTax:     formatMoney(reformTotalTax),
			TotalWithTax: formatMoney(reformTotalWithTax),
		},
		Difference: formatMoney(difference),
		Analysis:   analysis,
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

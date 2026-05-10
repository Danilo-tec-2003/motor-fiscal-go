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
	response, err := s.calculate(ctx, request)
	if err != nil {
		return dto.TaxSimulationResponse{}, err
	}

	if err := s.simulationRepository.Save(ctx, request, response); err != nil {
		return dto.TaxSimulationResponse{}, fmt.Errorf("save fiscal simulation: %w", err)
	}

	return response, nil
}

func (s TaxService) Preview(ctx context.Context, request dto.TaxPreviewRequest) (dto.TaxSimulationResponse, error) {
	return s.calculate(ctx, dto.TaxSimulationRequest{
		OperationDate: request.OperationDate,
		OriginUF:      request.OriginUF,
		DestinationUF: request.DestinationUF,
		FreightValue:  request.FreightValue,
		CustomerType:  request.CustomerType,
		OperationType: request.OperationType,
	})
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

	calculation, err := CalculateTaxes(freightValue, rule)
	if err != nil {
		return dto.TaxComparisonResponse{}, err
	}

	currentTotalTax, err := decimal.NewFromString(calculation.ICMS.Amount)
	if err != nil {
		return dto.TaxComparisonResponse{}, err
	}
	currentTotalWithTax := freightValue.Add(currentTotalTax)

	reformTotalTax := calculation.TotalTax
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

func (s TaxService) calculate(ctx context.Context, request dto.TaxSimulationRequest) (dto.TaxSimulationResponse, error) {
	freightValue, err := decimal.NewFromString(request.FreightValue)
	if err != nil {
		return dto.TaxSimulationResponse{}, err
	}

	rule, err := s.ruleService.FindRule(ctx, request)
	if err != nil {
		return dto.TaxSimulationResponse{}, err
	}

	calculation, err := CalculateTaxes(freightValue, rule)
	if err != nil {
		return dto.TaxSimulationResponse{}, err
	}

	return dto.TaxSimulationResponse{
		FreightID:          request.FreightID,
		BaseValue:          formatMoney(calculation.BaseValue),
		ICMS:               calculation.ICMS,
		IBS:                calculation.IBS,
		CBS:                calculation.CBS,
		TotalTax:           formatMoney(calculation.TotalTax),
		TotalWithTax:       formatMoney(calculation.TotalWithTax),
		CFOP:               rule.CFOP,
		RuleID:             rule.ID,
		RuleCode:           rule.RuleCode,
		RuleVersion:        rule.RuleVersion,
		RuleStatus:         rule.Status,
		CalculationBasis:   rule.CalculationBasis,
		CalculationDetails: calculation.CalculationDetails,
		FromCache:          false,
	}, nil
}

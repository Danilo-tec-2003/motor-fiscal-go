package service

import "github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"

type TaxService struct{}

func NewTaxService() TaxService {
	return TaxService{}
}

func (s TaxService) Simulate(request dto.TaxSimulationRequest) dto.TaxSimulationResponse {
	return dto.TaxSimulationResponse{
		FreightID: request.FreightID,
		BaseValue: request.FreightValue,
		ICMS: dto.TaxAmount{
			Rate:   "12.00",
			Amount: "420.00",
		},
		IBS: dto.TaxAmount{
			Rate:   "3.60",
			Amount: "126.00",
		},
		CBS: dto.TaxAmount{
			Rate:   "0.98",
			Amount: "31.50",
		},
		TotalTax:     "577.50",
		TotalWithTax: "4077.50",
		CFOP:         "6351",
		RuleVersion:  "2026.01",
		FromCahe:     false,
	}
}

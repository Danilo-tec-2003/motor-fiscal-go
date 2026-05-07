package service

import (
	"errors"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/model"
	"github.com/shopspring/decimal"
)

var ErrFiscalRuleNotFound = errors.New("fiscal rule not found")

type FiscalRuleService struct{}

func NewFiscalRuleService() FiscalRuleService {
	return FiscalRuleService{}
}

func (s FiscalRuleService) FindRule(request dto.TaxSimulationRequest) (model.FiscalRule, error) {
	return model.FiscalRule{
		RuleVersion:   "2026.01",
		OriginUF:      request.OriginUF,
		DestinationUF: request.DestinationUF,
		OperationType: request.OperationType,
		CustomerType:  request.CustomerType,
		ICMSRate:      decimal.NewFromFloat(12.00),
		IBSRate:       decimal.NewFromFloat(3.60),
		CBSRate:       decimal.NewFromFloat(0.90),
		CFOP:          "6351",
	}, nil
}

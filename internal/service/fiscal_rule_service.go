package service

import (
	"errors"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/model"
	"github.com/shopspring/decimal"
)

var ErrFiscalRuleNotFound = errors.New("fiscal rule not found")

var defaultFiscalRule = []model.FiscalRule{

	{
		RuleVersion:   "2026.01",
		OriginUF:      "PE",
		DestinationUF: "SP",
		OperationType: "INTERESTADUAL",
		CustomerType:  "PJ",
		ICMSRate:      decimal.NewFromFloat(12.00),
		IBSRate:       decimal.NewFromFloat(3.60),
		CBSRate:       decimal.NewFromFloat(0.90),
		CFOP:          "6351",
		ValidFrom:     "2026-01-01",
		ValidTo:       "2026-12-31",
	},

	{
		RuleVersion:   "2026.01",
		OriginUF:      "PE",
		DestinationUF: "PE",
		OperationType: "INTERNA",
		CustomerType:  "PF",
		ICMSRate:      decimal.NewFromFloat(18.00),
		IBSRate:       decimal.NewFromFloat(3.60),
		CBSRate:       decimal.NewFromFloat(0.90),
		CFOP:          "5351",
		ValidFrom:     "2026-01-01",
		ValidTo:       "2026-12-31",
	},
}

type FiscalRuleService struct{}

func NewFiscalRuleService() FiscalRuleService {
	return FiscalRuleService{}
}

func (s FiscalRuleService) FindRule(request dto.TaxSimulationRequest) (model.FiscalRule, error) {
	for _, rule := range defaultFiscalRule {
		if ruleMatchesRequest(rule, request) {
			return rule, nil
		}
	}

	return model.FiscalRule{}, ErrFiscalRuleNotFound
}

func ruleMatchesRequest(rule model.FiscalRule, request dto.TaxSimulationRequest) bool {
	return rule.OriginUF == request.OriginUF &&
		rule.DestinationUF == request.DestinationUF &&
		rule.OperationType == request.OperationType &&
		rule.CustomerType == request.CustomerType &&
		request.OperationDate >= rule.ValidFrom &&
		request.OperationDate <= rule.ValidTo
}

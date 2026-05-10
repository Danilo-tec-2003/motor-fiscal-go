package service

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/model"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/repository"
	"github.com/shopspring/decimal"
)

var (
	ErrFiscalRuleNotFound   = errors.New("fiscal rule not found")
	ErrFiscalRuleConflict   = errors.New("fiscal rule conflict")
	ErrFiscalRuleIncomplete = errors.New("fiscal rule incomplete")
)

type fiscalRuleRepository interface {
	FindCandidateRules(ctx context.Context, request dto.TaxSimulationRequest) ([]model.FiscalRuleEngine, error)
}

type FiscalRuleService struct {
	repository fiscalRuleRepository
}

func NewFiscalRuleService(repository fiscalRuleRepository) FiscalRuleService {
	return FiscalRuleService{
		repository: repository,
	}
}

func (s FiscalRuleService) FindRule(ctx context.Context, request dto.TaxSimulationRequest) (model.FiscalRule, error) {
	candidates, err := s.repository.FindCandidateRules(ctx, request)
	if err != nil {
		if errors.Is(err, repository.ErrFiscalRuleNotFound) {
			return model.FiscalRule{}, ErrFiscalRuleNotFound
		}

		return model.FiscalRule{}, err
	}

	matchedRules := make([]model.FiscalRuleEngine, 0, len(candidates))
	for _, candidate := range candidates {
		if ruleMatchesRequest(candidate, request) {
			matchedRules = append(matchedRules, candidate)
		}
	}

	if len(matchedRules) == 0 {
		return model.FiscalRule{}, ErrFiscalRuleNotFound
	}

	sort.SliceStable(matchedRules, func(i, j int) bool {
		if matchedRules[i].Priority != matchedRules[j].Priority {
			return matchedRules[i].Priority < matchedRules[j].Priority
		}

		if matchedRules[i].ValidFrom != matchedRules[j].ValidFrom {
			return matchedRules[i].ValidFrom > matchedRules[j].ValidFrom
		}

		return matchedRules[i].ID < matchedRules[j].ID
	})

	if hasRuleConflict(matchedRules) {
		return model.FiscalRule{}, ErrFiscalRuleConflict
	}

	return fiscalRuleFromEngine(matchedRules[0], request)
}

func ruleMatchesRequest(rule model.FiscalRuleEngine, request dto.TaxSimulationRequest) bool {
	if !rule.Active {
		return false
	}

	if len(rule.Conditions) == 0 {
		return false
	}

	for _, condition := range rule.Conditions {
		if !conditionMatchesRequest(condition, request) {
			return false
		}
	}

	return true
}

func conditionMatchesRequest(condition model.FiscalRuleCondition, request dto.TaxSimulationRequest) bool {
	requestValue, ok := requestFieldValue(condition.FieldName, request)
	if !ok {
		return false
	}

	switch condition.Operator {
	case model.RuleConditionOperatorEquals:
		return strings.EqualFold(requestValue, condition.FieldValue)
	case model.RuleConditionOperatorNotEquals:
		return !strings.EqualFold(requestValue, condition.FieldValue)
	case model.RuleConditionOperatorIn:
		return valueInList(requestValue, condition.FieldValue)
	case model.RuleConditionOperatorBetween:
		return dateBetween(requestValue, condition.FieldValue)
	default:
		return false
	}
}

func requestFieldValue(fieldName string, request dto.TaxSimulationRequest) (string, bool) {
	switch fieldName {
	case "operation_date":
		return request.OperationDate, true
	case "origin_uf":
		return request.OriginUF, true
	case "destination_uf":
		return request.DestinationUF, true
	case "operation_type":
		return request.OperationType, true
	case "customer_type":
		return request.CustomerType, true
	default:
		return "", false
	}
}

func valueInList(requestValue string, listValue string) bool {
	values := strings.Split(listValue, ",")
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value), requestValue) {
			return true
		}
	}

	return false
}

func dateBetween(requestValue string, rangeValue string) bool {
	separator := ","
	if strings.Contains(rangeValue, "..") {
		separator = ".."
	}

	parts := strings.Split(rangeValue, separator)
	if len(parts) != 2 {
		return false
	}

	requestDate, err := time.Parse("2006-01-02", requestValue)
	if err != nil {
		return false
	}

	startDate, err := time.Parse("2006-01-02", strings.TrimSpace(parts[0]))
	if err != nil {
		return false
	}

	endDate, err := time.Parse("2006-01-02", strings.TrimSpace(parts[1]))
	if err != nil {
		return false
	}

	return !requestDate.Before(startDate) && !requestDate.After(endDate)
}

func hasRuleConflict(matchedRules []model.FiscalRuleEngine) bool {
	if len(matchedRules) < 2 {
		return false
	}

	return matchedRules[0].Priority == matchedRules[1].Priority &&
		matchedRules[0].ValidFrom == matchedRules[1].ValidFrom
}

func fiscalRuleFromEngine(rule model.FiscalRuleEngine, request dto.TaxSimulationRequest) (model.FiscalRule, error) {
	taxRates := map[string]decimal.Decimal{}
	for _, tax := range rule.Taxes {
		taxRates[tax.TaxName] = tax.Rate
	}

	icmsRate, ok := taxRates[model.TaxNameICMS]
	if !ok {
		return model.FiscalRule{}, ErrFiscalRuleIncomplete
	}

	ibsRate, ok := taxRates[model.TaxNameIBS]
	if !ok {
		return model.FiscalRule{}, ErrFiscalRuleIncomplete
	}

	cbsRate, ok := taxRates[model.TaxNameCBS]
	if !ok {
		return model.FiscalRule{}, ErrFiscalRuleIncomplete
	}

	return model.FiscalRule{
		ID:               rule.ID,
		RuleCode:         rule.RuleCode,
		RuleVersion:      rule.RuleVersion,
		Description:      rule.Description,
		Priority:         rule.Priority,
		Status:           rule.Status,
		OriginUF:         request.OriginUF,
		DestinationUF:    request.DestinationUF,
		OperationType:    request.OperationType,
		CustomerType:     request.CustomerType,
		CalculationBasis: rule.CalculationBasis,
		ICMSRate:         icmsRate,
		IBSRate:          ibsRate,
		CBSRate:          cbsRate,
		CFOP:             rule.CFOP,
		ValidFrom:        rule.ValidFrom,
		ValidTo:          rule.ValidTo,
		Active:           rule.Active,
	}, nil
}

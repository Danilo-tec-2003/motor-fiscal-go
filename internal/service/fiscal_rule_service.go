package service

import (
	"context"
	"errors"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/model"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/repository"
)

var ErrFiscalRuleNotFound = errors.New("fiscal rule not found")

type fiscalRuleRepository interface {
	FindActiveRule(ctx context.Context, request dto.TaxSimulationRequest) (model.FiscalRule, error)
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
	rule, err := s.repository.FindActiveRule(ctx, request)
	if err != nil {
		if errors.Is(err, repository.ErrFiscalRuleNotFound) {
			return model.FiscalRule{}, ErrFiscalRuleNotFound
		}

		return model.FiscalRule{}, err
	}

	return rule, nil
}

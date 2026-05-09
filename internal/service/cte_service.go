package service

import (
	"context"
	"errors"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	apierrors "github.com/Danilo-tec-2003/motor-fiscal-go/internal/errors"
)

type CTeService struct {
	ruleService FiscalRuleService
}

func NewCTeService(ruleService FiscalRuleService) CTeService {
	return CTeService{
		ruleService: ruleService,
	}
}

func (s CTeService) Validate(ctx context.Context, request dto.CTeValidationRequest) (dto.CTeValidationResponse, error) {
	var validationErrors []apierrors.ValidationDetail
	var warnings []apierrors.ValidationDetail

	if request.Sender.Name == "" {
		validationErrors = append(validationErrors, apierrors.ValidationDetail{
			Field:   "sender.name",
			Message: "Nome do remetente e obrigatorio.",
		})
	}

	if request.Sender.Document == "" {
		validationErrors = append(validationErrors, apierrors.ValidationDetail{
			Field:   "sender.document",
			Message: "Documento do remetente e obrigatorio.",
		})
	}

	if request.Recipient.Name == "" {
		validationErrors = append(validationErrors, apierrors.ValidationDetail{
			Field:   "recipient.name",
			Message: "Nome do destinatario e obrigatorio.",
		})
	}

	if request.Recipient.Document == "" {
		validationErrors = append(validationErrors, apierrors.ValidationDetail{
			Field:   "recipient.document",
			Message: "Documento do destinatario e obrigatorio.",
		})
	}

	cfop := ""
	rule, err := s.ruleService.FindRule(ctx, request.TaxSimulationRequest)
	if err == nil {
		cfop = rule.CFOP
	} else if errors.Is(err, ErrFiscalRuleNotFound) {
		warnings = append(warnings, apierrors.ValidationDetail{
			Field:   "tax_simulation",
			Message: "Nao foi encontrada regra fiscal para determinar CFOP.",
		})
	} else {
		return dto.CTeValidationResponse{}, err
	}

	return dto.CTeValidationResponse{
		FreightID: request.FreightID,
		Valid:     len(validationErrors) == 0,
		CFOP:      cfop,
		Errors:    validationErrors,
		Warnings:  warnings,
	}, nil
}

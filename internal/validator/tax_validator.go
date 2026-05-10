package validator

import (
	"time"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	apierrors "github.com/Danilo-tec-2003/motor-fiscal-go/internal/errors"
	"github.com/shopspring/decimal"
)

func ValidateTaxSimulationRequest(request dto.TaxSimulationRequest) []apierrors.ValidationDetail {
	details := validateTaxRequestFields(request, true)
	return details
}

func ValidateTaxPreviewRequest(request dto.TaxPreviewRequest) []apierrors.ValidationDetail {
	return validateTaxRequestFields(dto.TaxSimulationRequest{
		OperationDate: request.OperationDate,
		OriginUF:      request.OriginUF,
		DestinationUF: request.DestinationUF,
		FreightValue:  request.FreightValue,
		CustomerType:  request.CustomerType,
		OperationType: request.OperationType,
	}, false)
}

func validateTaxRequestFields(request dto.TaxSimulationRequest, requireFreightID bool) []apierrors.ValidationDetail {
	var details []apierrors.ValidationDetail

	if requireFreightID && request.FreightID <= 0 {
		details = append(details, apierrors.ValidationDetail{
			Field:   "freight_id",
			Message: "Deve ser maior que zero.",
		})
	}

	if request.OperationDate == "" {
		details = append(details, apierrors.ValidationDetail{
			Field:   "operation_date",
			Message: "Campo obrigatorio.",
		})
	} else if _, err := time.Parse("2006-01-02", request.OperationDate); err != nil {
		details = append(details, apierrors.ValidationDetail{
			Field:   "operation_date",
			Message: "Data deve estar no formato YYYY-MM-DD.",
		})
	}

	if request.OriginUF == "" {
		details = append(details, apierrors.ValidationDetail{
			Field:   "origin_uf",
			Message: "Campo obrigatorio.",
		})
	} else if !isValidUF(request.OriginUF) {
		details = append(details, apierrors.ValidationDetail{
			Field:   "origin_uf",
			Message: "UF invalida.",
		})
	}

	if request.DestinationUF == "" {
		details = append(details, apierrors.ValidationDetail{
			Field:   "destination_uf",
			Message: "Campo obrigatorio.",
		})
	} else if !isValidUF(request.DestinationUF) {
		details = append(details, apierrors.ValidationDetail{
			Field:   "destination_uf",
			Message: "UF invalida.",
		})
	}

	if isValidUF(request.OriginUF) &&
		isValidUF(request.DestinationUF) &&
		isValidOperationType(request.OperationType) {

		if request.OriginUF == request.DestinationUF && request.OperationType != "INTERNA" {
			details = append(details, apierrors.ValidationDetail{
				Field:   "operation_type",
				Message: "Operacao com mesma UF deve ser INTERNA.",
			})
		}

		if request.OriginUF != request.DestinationUF && request.OperationType != "INTERESTADUAL" {
			details = append(details, apierrors.ValidationDetail{
				Field:   "operation_type",
				Message: "Operacao entre UFs diferentes deve ser INTERESTADUAL.",
			})
		}
	}

	if request.FreightValue == "" {
		details = append(details, apierrors.ValidationDetail{
			Field:   "freight_value",
			Message: "Campo obrigatorio.",
		})
	} else {
		freightValue, err := decimal.NewFromString(request.FreightValue)
		if err != nil {
			details = append(details, apierrors.ValidationDetail{
				Field:   "freight_value",
				Message: "Valor monetario invalido.",
			})
		} else if freightValue.LessThanOrEqual(decimal.Zero) {
			details = append(details, apierrors.ValidationDetail{
				Field:   "freight_value",
				Message: "Deve ser maior que zero.",
			})
		}
	}

	if request.CustomerType == "" {
		details = append(details, apierrors.ValidationDetail{
			Field:   "customer_type",
			Message: "Campo obrigatorio.",
		})
	} else if !isValidCustomerType(request.CustomerType) {
		details = append(details, apierrors.ValidationDetail{
			Field:   "customer_type",
			Message: "Tipo de cliente invalido. Deve ser 'PF' ou 'PJ'.",
		})
	}

	if request.OperationType == "" {
		details = append(details, apierrors.ValidationDetail{
			Field:   "operation_type",
			Message: "Campo obrigatorio.",
		})
	} else if !isValidOperationType(request.OperationType) {
		details = append(details, apierrors.ValidationDetail{
			Field:   "operation_type",
			Message: "Tipo de operacao invalido. Deve ser 'INTERNA' ou 'INTERESTADUAL'.",
		})
	}

	return details
}

func isValidUF(uf string) bool {
	validUFs := map[string]bool{
		"AC": true,
		"AL": true,
		"AP": true,
		"AM": true,
		"BA": true,
		"CE": true,
		"DF": true,
		"ES": true,
		"GO": true,
		"MA": true,
		"MT": true,
		"MS": true,
		"MG": true,
		"PA": true,
		"PB": true,
		"PR": true,
		"PE": true,
		"PI": true,
		"RJ": true,
		"RN": true,
		"RS": true,
		"RO": true,
		"RR": true,
		"SC": true,
		"SP": true,
		"SE": true,
		"TO": true,
	}
	return validUFs[uf]
}

func isValidCustomerType(customerType string) bool {
	return customerType == "PF" || customerType == "PJ"
}

func isValidOperationType(operationType string) bool {
	return operationType == "INTERNA" || operationType == "INTERESTADUAL"
}

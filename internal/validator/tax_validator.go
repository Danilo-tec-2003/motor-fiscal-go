package validator

import (
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	apierrors "github.com/Danilo-tec-2003/motor-fiscal-go/internal/errors"
	"github.com/shopspring/decimal"
)

func ValidateTaxSimulationRequest(request dto.TaxSimulationRequest) []apierrors.ValidationDetail {
	var details []apierrors.ValidationDetail

	if request.FreightID <= 0 {
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
	}

	if request.OriginUF == "" {
		details = append(details, apierrors.ValidationDetail{
			Field:   "origin_uf",
			Message: "Campo obrigatorio.",
		})
	}

	if request.DestinationUF == "" {
		details = append(details, apierrors.ValidationDetail{
			Field:   "destination_uf",
			Message: "Campo obrigatorio.",
		})
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
	}

	if request.OperationType == "" {
		details = append(details, apierrors.ValidationDetail{
			Field:   "operation_type",
			Message: "Campo obrigatorio.",
		})
	}

	return details
}

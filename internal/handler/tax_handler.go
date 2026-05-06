package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	apierrors "github.com/Danilo-tec-2003/motor-fiscal-go/internal/errors"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/middleware"
)

func TaxSimulationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		correlationID := middleware.GetCorrelationID(r.Context())

		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)

			WriteError(w, http.StatusMethodNotAllowed, apierrors.APIError{
				Code:          "METHOD_NOT_ALLOWED",
				Message:       "Metodo nao permitido. Use POST.",
				CorrelationID: correlationID,
				Details:       []apierrors.ValidationDetail{},
			})
			return
		}

		var request dto.TaxSimulationRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			WriteError(w, http.StatusBadRequest, apierrors.APIError{
				Code:          "BAD_REQUEST",
				Message:       "JSON invalido ou malformado.",
				CorrelationID: correlationID,
				Details:       []apierrors.ValidationDetail{},
			})
			return
		}

		details := validateTaxSimulationRequest(request)
		if len(details) > 0 {
			WriteError(w, http.StatusUnprocessableEntity, apierrors.APIError{
				Code:          "VALIDATION_ERROR",
				Message:       "Payload invalido.",
				CorrelationID: correlationID,
				Details:       details,
			})
			return
		}

		WriteError(w, http.StatusNotImplemented, apierrors.APIError{
			Code:          "NOT_IMPLEMENTED",
			Message:       "Calculo fiscal ainda nao implementado.",
			CorrelationID: correlationID,
			Details:       []apierrors.ValidationDetail{},
		})
	}
}

func validateTaxSimulationRequest(request dto.TaxSimulationRequest) []apierrors.ValidationDetail {
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

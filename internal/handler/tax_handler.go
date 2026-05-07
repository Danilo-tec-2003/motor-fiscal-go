package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	apierrors "github.com/Danilo-tec-2003/motor-fiscal-go/internal/errors"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/middleware"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/service"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/validator"
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

		details := validator.ValidateTaxSimulationRequest(request)

		if len(details) > 0 {
			WriteError(w, http.StatusUnprocessableEntity, apierrors.APIError{
				Code:          "VALIDATION_ERROR",
				Message:       "Payload invalido.",
				CorrelationID: correlationID,
				Details:       details,
			})
			return
		}

		taxService := service.NewTaxService()
		response, err := taxService.Simulate(request)
		if err != nil {
			if errors.Is(err, service.ErrFiscalRuleNotFound) {
				WriteError(w, http.StatusNotFound, apierrors.APIError{
					Code:          "FISCAL_RULE_NOT_FOUND",
					Message:       "Nenhuma regra fiscal vigente.",
					CorrelationID: correlationID,
					Details: []apierrors.ValidationDetail{
						{
							Field:   "operation_date",
							Message: "Verifique se existe regra ativa para os parametros informados.",
						},
					},
				})
				return
			}

			WriteError(w, http.StatusUnprocessableEntity, apierrors.APIError{
				Code:          "VALIDATION_ERROR",
				Message:       "Payload invalido.",
				CorrelationID: correlationID,
				Details: []apierrors.ValidationDetail{
					{
						Field:   "freight_value",
						Message: "Valor monetario invalido.",
					},
				},
			})
			return
		}

		WriteJSON(w, http.StatusOK, response)
	}
}

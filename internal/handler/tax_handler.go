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

type taxService interface {
	Simulate(request dto.TaxSimulationRequest) (dto.TaxSimulationResponse, error)
	Compare(request dto.TaxSimulationRequest) (dto.TaxComparisonResponse, error)
}

type TaxHandler struct {
	taxService taxService
}

func NewTaxHandler(taxService taxService) TaxHandler {
	return TaxHandler{
		taxService: taxService,
	}
}

func (h TaxHandler) Simulate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		correlationID := middleware.GetCorrelationID(r.Context())

		if r.Method != http.MethodPost {
			writeTaxMethodNotAllowed(w, correlationID)
			return
		}

		var request dto.TaxSimulationRequest
		if !decodeJSON(w, r, correlationID, &request) {
			return
		}

		if !validateTaxRequest(w, request, correlationID) {
			return
		}

		response, err := h.taxService.Simulate(request)
		if err != nil {
			writeTaxServiceError(w, err, correlationID)
			return
		}

		WriteJSON(w, http.StatusOK, response)
	}
}

func (h TaxHandler) Compare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		correlationID := middleware.GetCorrelationID(r.Context())

		if r.Method != http.MethodPost {
			writeTaxMethodNotAllowed(w, correlationID)
			return
		}

		var request dto.TaxSimulationRequest
		if !decodeJSON(w, r, correlationID, &request) {
			return
		}

		if !validateTaxRequest(w, request, correlationID) {
			return
		}

		response, err := h.taxService.Compare(request)
		if err != nil {
			writeTaxServiceError(w, err, correlationID)
			return
		}

		WriteJSON(w, http.StatusOK, response)
	}
}

func (h TaxHandler) Batch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		correlationID := middleware.GetCorrelationID(r.Context())

		if r.Method != http.MethodPost {
			writeTaxMethodNotAllowed(w, correlationID)
			return
		}

		var request dto.TaxBatchRequest
		if !decodeJSON(w, r, correlationID, &request) {
			return
		}

		if len(request.Items) == 0 {
			WriteError(w, http.StatusUnprocessableEntity, apierrors.APIError{
				Code:          "VALIDATION_ERROR",
				Message:       "Payload invalido.",
				CorrelationID: correlationID,
				Details: []apierrors.ValidationDetail{
					{
						Field:   "items",
						Message: "Informe ao menos um frete para processamento em lote.",
					},
				},
			})
			return
		}

		response := dto.TaxBatchResponse{
			TotalItems: len(request.Items),
			Results:    []dto.TaxBatchItemResult{},
		}

		for _, item := range request.Items {
			details := validator.ValidateTaxSimulationRequest(item)
			if len(details) > 0 {
				apiErr := apierrors.APIError{
					Code:          "VALIDATION_ERROR",
					Message:       "Payload invalido.",
					CorrelationID: correlationID,
					Details:       details,
				}

				response.ErrorCount++
				response.Results = append(response.Results, dto.TaxBatchItemResult{
					FreightID: item.FreightID,
					Success:   false,
					Error:     &apiErr,
				})
				continue
			}

			simulation, err := h.taxService.Simulate(item)
			if err != nil {
				_, apiErr := taxServiceAPIError(err, correlationID)

				response.ErrorCount++
				response.Results = append(response.Results, dto.TaxBatchItemResult{
					FreightID: item.FreightID,
					Success:   false,
					Error:     &apiErr,
				})
				continue
			}

			response.SuccessCount++
			response.Results = append(response.Results, dto.TaxBatchItemResult{
				FreightID: item.FreightID,
				Success:   true,
				Data:      &simulation,
			})
		}

		WriteJSON(w, http.StatusOK, response)
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, correlationID string, target any) bool {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		WriteError(w, http.StatusBadRequest, apierrors.APIError{
			Code:          "BAD_REQUEST",
			Message:       "JSON invalido ou malformado.",
			CorrelationID: correlationID,
			Details:       []apierrors.ValidationDetail{},
		})
		return false
	}

	return true
}

func validateTaxRequest(w http.ResponseWriter, request dto.TaxSimulationRequest, correlationID string) bool {
	details := validator.ValidateTaxSimulationRequest(request)
	if len(details) == 0 {
		return true
	}

	WriteError(w, http.StatusUnprocessableEntity, apierrors.APIError{
		Code:          "VALIDATION_ERROR",
		Message:       "Payload invalido.",
		CorrelationID: correlationID,
		Details:       details,
	})
	return false
}

func writeTaxMethodNotAllowed(w http.ResponseWriter, correlationID string) {
	w.Header().Set("Allow", http.MethodPost)

	WriteError(w, http.StatusMethodNotAllowed, apierrors.APIError{
		Code:          "METHOD_NOT_ALLOWED",
		Message:       "Metodo nao permitido. Use POST.",
		CorrelationID: correlationID,
		Details:       []apierrors.ValidationDetail{},
	})
}

func writeTaxServiceError(w http.ResponseWriter, err error, correlationID string) {
	statusCode, apiErr := taxServiceAPIError(err, correlationID)
	WriteError(w, statusCode, apiErr)
}

func taxServiceAPIError(err error, correlationID string) (int, apierrors.APIError) {
	if errors.Is(err, service.ErrFiscalRuleNotFound) {
		return http.StatusNotFound, apierrors.APIError{
			Code:          "FISCAL_RULE_NOT_FOUND",
			Message:       "Nenhuma regra fiscal vigente.",
			CorrelationID: correlationID,
			Details: []apierrors.ValidationDetail{
				{
					Field:   "operation_date",
					Message: "Verifique se existe regra ativa para os parametros informados.",
				},
			},
		}
	}

	return http.StatusUnprocessableEntity, apierrors.APIError{
		Code:          "VALIDATION_ERROR",
		Message:       "Payload invalido.",
		CorrelationID: correlationID,
		Details: []apierrors.ValidationDetail{
			{
				Field:   "freight_value",
				Message: "Valor monetario invalido.",
			},
		},
	}
}

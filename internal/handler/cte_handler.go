package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/dto"
	apierrors "github.com/Danilo-tec-2003/motor-fiscal-go/internal/errors"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/middleware"
	"github.com/Danilo-tec-2003/motor-fiscal-go/internal/validator"
)

type cteService interface {
	Validate(request dto.CTeValidationRequest) dto.CTeValidationResponse
}

type CTeHandler struct {
	cteService cteService
}

func NewCTeHandler(cteService cteService) CTeHandler {
	return CTeHandler{
		cteService: cteService,
	}
}

func (h CTeHandler) Validate() http.HandlerFunc {
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

		var request dto.CTeValidationRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			WriteError(w, http.StatusBadRequest, apierrors.APIError{
				Code:          "BAD_REQUEST",
				Message:       "JSON invalido ou malformado.",
				CorrelationID: correlationID,
				Details:       []apierrors.ValidationDetail{},
			})
			return
		}

		details := validator.ValidateTaxSimulationRequest(request.TaxSimulationRequest)
		if len(details) > 0 {
			WriteError(w, http.StatusUnprocessableEntity, apierrors.APIError{
				Code:          "VALIDATION_ERROR",
				Message:       "Payload invalido.",
				CorrelationID: correlationID,
				Details:       details,
			})
			return
		}

		response := h.cteService.Validate(request)

		WriteJSON(w, http.StatusOK, response)
	}
}

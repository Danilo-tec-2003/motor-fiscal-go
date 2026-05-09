package dto

import apierrors "github.com/Danilo-tec-2003/motor-fiscal-go/internal/errors"

type TaxBatchResponse struct {
	TotalItems   int                  `json:"total_items"`
	SuccessCount int                  `json:"success_count"`
	ErrorCount   int                  `json:"error_count"`
	Results      []TaxBatchItemResult `json:"results"`
}

type TaxBatchItemResult struct {
	FreightID int64                  `json:"freight_id"`
	Success   bool                   `json:"success"`
	Data      *TaxSimulationResponse `json:"data,omitempty"`
	Error     *apierrors.APIError    `json:"error,omitempty"`
}

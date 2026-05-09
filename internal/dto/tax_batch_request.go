package dto

type TaxBatchRequest struct {
	Items []TaxSimulationRequest `json:"items"`
}

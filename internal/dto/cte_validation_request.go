package dto

type Participant struct {
	Name     string `json:"name"`
	Document string `json:"document"`
	UF       string `json:"uf"`
}

type CTeValidationRequest struct {
	TaxSimulationRequest
	Sender    Participant `json:"sender"`
	Recipient Participant `json:"recipient"`
}

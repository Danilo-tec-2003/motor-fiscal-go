package dto

type TaxAmount struct {
	Rate   string `json:"rate"`
	Amount string `json:"amount"`
}

type TaxSimulationResponse struct {
	FreightID    int64     `json:"freight_id"`
	BaseValue    string    `json:"base_value"`
	ICMS         TaxAmount `json:"icms"`
	IBS          TaxAmount `json:"ibs"`
	CBS          TaxAmount `json:"cbs"`
	TotalTax     string    `json:"total_tax"`
	TotalWithTax string    `json:"total_with_tax"`
	CFOP         string    `json:"cfop"`
	RuleVersion  string    `json:"rule_version"`
	FromCahe     bool      `json:"from_cache"`
}

package dto

type TaxAmount struct {
	Rate   string `json:"rate"`
	Amount string `json:"amount"`
}

type TaxCalculationDetail struct {
	TaxName            string `json:"tax_name"`
	BaseValue          string `json:"base_value"`
	BaseReductionRate  string `json:"base_reduction_rate"`
	EffectiveBaseValue string `json:"effective_base_value"`
	Rate               string `json:"rate"`
	Amount             string `json:"amount"`
	Formula            string `json:"formula"`
}

type TaxSimulationResponse struct {
	FreightID          int64                  `json:"freight_id"`
	BaseValue          string                 `json:"base_value"`
	ICMS               TaxAmount              `json:"icms"`
	IBS                TaxAmount              `json:"ibs"`
	CBS                TaxAmount              `json:"cbs"`
	TotalTax           string                 `json:"total_tax"`
	TotalWithTax       string                 `json:"total_with_tax"`
	CFOP               string                 `json:"cfop"`
	RuleID             int64                  `json:"rule_id"`
	RuleCode           string                 `json:"rule_code"`
	RuleVersion        string                 `json:"rule_version"`
	RuleStatus         string                 `json:"rule_status"`
	CalculationBasis   string                 `json:"calculation_basis"`
	CalculationDetails []TaxCalculationDetail `json:"calculation_details"`
	FromCache          bool                   `json:"from_cache"`
}

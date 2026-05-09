package dto

type TaxScenario struct {
	BaseValue    string `json:"base_value"`
	TotalTax     string `json:"total_tax"`
	TotalWithTax string `json:"total_with_tax"`
}

type TaxComparisonResponse struct {
	FreightID       int64       `json:"freight_id"`
	CurrentScenario TaxScenario `json:"current_scenario"`
	ReformScenario  TaxScenario `json:"reform_scenario"`
	Difference      string      `json:"difference"`
	Analysis        string      `json:"analysis"`
}

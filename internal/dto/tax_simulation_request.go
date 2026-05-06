package dto

type TaxSimulationRequest struct {
	FreightID     int64  `json:"freight_id"`
	OperationDate string `json:"operation_date"`
	OriginUF      string `json:"origin_uf"`
	DestinationUF string `json:"destination_uf"`
	FreightValue  string `json:"freight_value"` //Valor monetario trafegado como string para evitar perda de precisao; no service sera convertido para decimal.
	CustomerType  string `json:"customer_type"`
	OperationType string `json:"operation_type"`
}

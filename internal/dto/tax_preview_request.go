package dto

type TaxPreviewRequest struct {
	OperationDate string `json:"operation_date"`
	OriginUF      string `json:"origin_uf"`
	DestinationUF string `json:"destination_uf"`
	FreightValue  string `json:"freight_value"`
	CustomerType  string `json:"customer_type"`
	OperationType string `json:"operation_type"`
}

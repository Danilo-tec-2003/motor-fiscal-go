package errors

type ValidationDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type APIError struct {
	Code          string             `json:"code"`
	Message       string             `json:"message"`
	CorrelationID string             `json:"correlation_id"`
	Details       []ValidationDetail `json:"details"`
}

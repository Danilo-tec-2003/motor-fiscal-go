package dto

import (
	apierrors "github.com/Danilo-tec-2003/motor-fiscal-go/internal/errors"
)

type CTeValidationResponse struct {
	FreightID int64                        `json:"freight_id"`
	Valid     bool                         `json:"valid"`
	CFOP      string                       `json:"cfop,omitempty"`
	Errors    []apierrors.ValidationDetail `json:"errors"`
	Warnings  []apierrors.ValidationDetail `json:"warnings"`
}

package model

import "github.com/shopspring/decimal"

type FiscalRule struct {
	ID               int64
	RuleCode         string
	RuleVersion      string
	Description      string
	Priority         int
	Status           string
	OriginUF         string
	DestinationUF    string
	OperationType    string
	CustomerType     string
	CalculationBasis string
	ICMSRate         decimal.Decimal
	IBSRate          decimal.Decimal
	CBSRate          decimal.Decimal
	CFOP             string
	ValidFrom        string
	ValidTo          string
	Active           bool
}

package model

import "github.com/shopspring/decimal"

type FiscalRule struct {
	RuleVersion   string
	OriginUF      string
	DestinationUF string
	OperationType string
	CustomerType  string
	ICMSRate      decimal.Decimal
	IBSRate       decimal.Decimal
	CBSRate       decimal.Decimal
	CFOP          string
}

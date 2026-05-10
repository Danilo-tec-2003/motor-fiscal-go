package model

import "github.com/shopspring/decimal"

const (
	FiscalRuleStatusDraft         = "DRAFT"
	FiscalRuleStatusPendingReview = "PENDING_REVIEW"
	FiscalRuleStatusApproved      = "APPROVED"
	FiscalRuleStatusInactive      = "INACTIVE"

	RuleConditionOperatorEquals    = "EQUALS"
	RuleConditionOperatorNotEquals = "NOT_EQUALS"
	RuleConditionOperatorIn        = "IN"
	RuleConditionOperatorBetween   = "BETWEEN"

	CalculationBasisFreightValue = "FREIGHT_VALUE"

	TaxNameICMS = "ICMS"
	TaxNameIBS  = "IBS"
	TaxNameCBS  = "CBS"
)

type FiscalRuleEngine struct {
	ID               int64
	RuleCode         string
	RuleVersion      string
	Description      string
	Priority         int
	Status           string
	CalculationBasis string
	CFOP             string
	ValidFrom        string
	ValidTo          string
	Active           bool
	Conditions       []FiscalRuleCondition
	Taxes            []FiscalRuleTax
}

type FiscalRuleCondition struct {
	ID           int64
	FiscalRuleID int64
	FieldName    string
	Operator     string
	FieldValue   string
}

type FiscalRuleTax struct {
	ID                int64
	FiscalRuleID      int64
	TaxName           string
	Rate              decimal.Decimal
	BaseReductionRate decimal.Decimal
	CalculationOrder  int
}

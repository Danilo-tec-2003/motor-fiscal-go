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

	FiscalRuleSourceTypeFederalLaw         = "FEDERAL_LAW"
	FiscalRuleSourceTypeStateLaw           = "STATE_LAW"
	FiscalRuleSourceTypeCONFAZ             = "CONFAZ"
	FiscalRuleSourceTypeSEFAZ              = "SEFAZ"
	FiscalRuleSourceTypeSenateResolution   = "SENATE_RESOLUTION"
	FiscalRuleSourceTypeAccountingGuidance = "ACCOUNTING_GUIDANCE"
	FiscalRuleSourceTypeInternalNote       = "INTERNAL_NOTE"

	FiscalRuleSourceStatusPendingConfirmation = "PENDING_CONFIRMATION"
	FiscalRuleSourceStatusConfirmed           = "CONFIRMED"
	FiscalRuleSourceStatusReplaced            = "REPLACED"

	AccountingReviewStatusPendingReview  = "PENDING_REVIEW"
	AccountingReviewStatusApproved       = "APPROVED"
	AccountingReviewStatusRejected       = "REJECTED"
	AccountingReviewStatusChangesRequest = "CHANGES_REQUESTED"
)

type FiscalRuleEngine struct {
	ID                int64
	RuleCode          string
	RuleVersion       string
	Description       string
	Priority          int
	Status            string
	CalculationBasis  string
	CFOP              string
	ValidFrom         string
	ValidTo           string
	Active            bool
	Conditions        []FiscalRuleCondition
	Taxes             []FiscalRuleTax
	Sources           []FiscalRuleSource
	AccountingReviews []FiscalRuleAccountingReview
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

type FiscalRuleSource struct {
	ID           int64
	FiscalRuleID int64
	SourceType   string
	SourceStatus string
	Title        string
	Reference    string
	URL          string
	PublishedAt  string
	Notes        string
}

type FiscalRuleAccountingReview struct {
	ID           int64
	FiscalRuleID int64
	ReviewStatus string
	ReviewerName string
	ReviewerRole string
	ReviewedAt   string
	Notes        string
}

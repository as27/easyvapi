package model

// BillingAccount represents a billing/cost account in easyVerein (Buchungskonto).
type BillingAccount struct {
	// ID is the unique identifier of the billing account.
	ID int `json:"id"`
	// Name is the display name of the billing account.
	Name string `json:"name"`
	// AccountKind classifies the account type (e.g. "income", "expense").
	AccountKind string `json:"accountKind"`
	// Skr is the standard chart of accounts number (Standardkontenrahmen-Nummer).
	Skr string `json:"skr"`
	// Description is an optional free-text description.
	Description string `json:"description"`
	// Balance is the current account balance.
	Balance flexFloat64 `json:"balance"`
	// StandardFormOfAccounts is the standard chart of accounts this account belongs to.
	StandardFormOfAccounts string `json:"standardFormOfAccounts"`
}

// BillingAccountCreate holds the fields for creating or updating a billing account
// via POST / PATCH /billing-account.
type BillingAccountCreate struct {
	// Name is the display name (required for create).
	Name string `json:"name,omitempty"`
	// AccountKind classifies the account type.
	AccountKind string `json:"accountKind,omitempty"`
	// Skr is the standard chart of accounts number.
	Skr string `json:"skr,omitempty"`
	// Description is an optional free-text description.
	Description string `json:"description,omitempty"`
	// Balance is the initial or updated account balance.
	Balance float64 `json:"balance,omitempty"`
	// StandardFormOfAccounts is the standard chart of accounts reference.
	StandardFormOfAccounts string `json:"standardFormOfAccounts,omitempty"`
}

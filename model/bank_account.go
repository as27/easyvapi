package model

// BankAccount represents a bank account in easyVerein (Bankkonto).
type BankAccount struct {
	// ID is the unique identifier of the bank account.
	ID int `json:"id"`
	// Name is the display name of the bank account.
	Name string `json:"name"`
	// IBAN is the International Bank Account Number.
	IBAN string `json:"iban"`
	// BIC is the Bank Identifier Code (SWIFT code).
	BIC string `json:"bic"`
	// Balance is the current account balance.
	Balance flexFloat64 `json:"balance"`
	// Description is an optional free-text description.
	Description string `json:"description"`
	// BankAccountOwner is the name of the account owner.
	BankAccountOwner string `json:"bankAccountOwner"`
	// SepaCreditorId is the SEPA creditor identifier.
	SepaCreditorId string `json:"sepaCreditorId"`
	// SepaScheme is the SEPA scheme (e.g. "CORE", "B2B").
	SepaScheme string `json:"sepaScheme"`
}

// BankAccountCreate holds the fields for creating or updating a bank account
// via POST / PATCH /bank-account.
type BankAccountCreate struct {
	// Name is the display name (required for create).
	Name string `json:"name,omitempty"`
	// IBAN is the International Bank Account Number.
	IBAN string `json:"iban,omitempty"`
	// BIC is the Bank Identifier Code.
	BIC string `json:"bic,omitempty"`
	// Balance is the current account balance.
	Balance float64 `json:"balance,omitempty"`
	// Description is an optional free-text description.
	Description string `json:"description,omitempty"`
	// BankAccountOwner is the name of the account owner.
	BankAccountOwner string `json:"bankAccountOwner,omitempty"`
	// SepaCreditorId is the SEPA creditor identifier.
	SepaCreditorId string `json:"sepaCreditorId,omitempty"`
	// SepaScheme is the SEPA scheme.
	SepaScheme string `json:"sepaScheme,omitempty"`
}

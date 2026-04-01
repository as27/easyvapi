package model

// CustomTaxRate represents a custom tax rate in easyVerein (benutzerdefinierter Steuersatz).
type CustomTaxRate struct {
	// ID is the unique identifier of the custom tax rate.
	ID int `json:"id"`
	// TaxName is the display name of the tax rate (e.g. "Reduced VAT").
	TaxName string `json:"taxName"`
	// CustomTaxRate is the tax rate as a percentage (e.g. 7.0 for 7%).
	CustomTaxRate flexFloat64 `json:"customTaxRate"`
}

// CustomTaxRateCreate holds the fields for creating or updating a custom tax rate
// via POST / PATCH /custom-tax-rate.
type CustomTaxRateCreate struct {
	// TaxName is the display name (required for create).
	TaxName string `json:"taxName,omitempty"`
	// CustomTaxRate is the tax rate as a percentage (required for create).
	CustomTaxRate float64 `json:"customTaxRate,omitempty"`
}

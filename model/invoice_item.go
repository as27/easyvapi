package model

// InvoiceItem represents a line item within an invoice in easyVerein (Rechnungsposition).
type InvoiceItem struct {
	// ID is the unique identifier of the invoice item.
	ID int `json:"id"`
	// Title is the short description/name of the item.
	Title string `json:"title"`
	// Quantity is the number of units.
	Quantity flexFloat64 `json:"quantity"`
	// UnitPrice is the price per unit.
	UnitPrice flexFloat64 `json:"unitPrice"`
	// TaxRate is the applicable tax rate as a percentage.
	TaxRate flexFloat64 `json:"taxRate"`
	// TaxName is the name of the tax (e.g. "MwSt.", "VAT").
	TaxName string `json:"taxName"`
	// Description is an optional extended description.
	Description string `json:"description"`
	// BillingAccount is the ID of the billing account for this item.
	BillingAccount int `json:"billingAccount"`
	// Gross indicates whether prices are gross (tax-inclusive).
	Gross bool `json:"gross"`
}

// InvoiceItemCreate holds the fields for creating or updating an invoice item
// via POST / PATCH /invoice-item.
type InvoiceItemCreate struct {
	// Quantity is the number of units (required for create).
	Quantity float64 `json:"quantity,omitempty"`
	// UnitPrice is the price per unit (required for create).
	UnitPrice float64 `json:"unitPrice,omitempty"`
	// Title is the short description/name (required for create).
	Title string `json:"title,omitempty"`
	// Description is an optional extended description.
	Description string `json:"description,omitempty"`
	// TaxRate is the applicable tax rate as a percentage.
	TaxRate float64 `json:"taxRate,omitempty"`
	// Gross indicates whether prices are gross (tax-inclusive).
	Gross bool `json:"gross,omitempty"`
	// TaxName is the name of the tax.
	TaxName string `json:"taxName,omitempty"`
	// BillingAccount is the ID of the billing account for this item.
	BillingAccount int `json:"billingAccount,omitempty"`
}

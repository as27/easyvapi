package model

// Invoice represents an invoice in easyVerein.
type Invoice struct {
	// ID is the unique identifier of the invoice.
	ID int `json:"id"`
	// InvNumber is the human-readable invoice number.
	InvNumber string `json:"invNumber"`
	// Date is the invoice date in YYYY-MM-DD format.
	Date string `json:"date"`
	// Receiver is the name of the invoice recipient.
	Receiver string `json:"receiver"`
	// TotalPrice is the total invoice amount.
	TotalPrice flexFloat64 `json:"totalPrice"`
	// Kind classifies the invoice (e.g. "outgoing", "incoming").
	Kind string `json:"kind"`
	// IsDraft indicates whether the invoice is still a draft.
	IsDraft bool `json:"isDraft"`
	// IsTemplate indicates whether the invoice is a template.
	IsTemplate bool `json:"isTemplate"`
	// Description is an optional free-text description.
	Description string `json:"description"`
}

// InvoiceCreate holds the fields used when creating a new invoice via
// POST /invoice.
type InvoiceCreate struct {
	// Date is the invoice date in YYYY-MM-DD format.
	Date string `json:"date,omitempty"`
	// Receiver is the name of the invoice recipient.
	Receiver string `json:"receiver,omitempty"`
	// TotalPrice is the total invoice amount.
	TotalPrice float64 `json:"totalPrice,omitempty"`
	// Kind classifies the invoice.
	Kind string `json:"kind,omitempty"`
	// IsDraft marks the invoice as a draft.
	IsDraft bool `json:"isDraft,omitempty"`
	// Description is an optional free-text description.
	Description string `json:"description,omitempty"`
	// RelatedAddress is the ID of the contact-details record to link.
	RelatedAddress int `json:"relatedAddress,omitempty"`
}

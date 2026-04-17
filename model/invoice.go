package model

// InvoiceCharges holds the charge summary returned inside an Invoice.
type InvoiceCharges struct {
	Charge     float64 `json:"charge"`
	ChargeBack float64 `json:"chargeBack"`
	Total      float64 `json:"total"`
}

// Invoice represents an invoice in easyVerein.
type Invoice struct {
	// ID is the unique identifier of the invoice.
	ID int `json:"id"`
	// RelatedBookings holds URLs of bookings linked to this invoice.
	RelatedBookings []string `json:"relatedBookings"`
	// Org is the URL of the organisation this invoice belongs to.
	Org string `json:"org"`
	// Path is the URL to the PDF file, or empty if not yet generated.
	Path string `json:"path"`
	// RelatedAddress is the URL of the linked contact-details record, or nil.
	RelatedAddress *string `json:"relatedAddress"`
	// PayedFromUser is the URL of the user who paid, or nil.
	PayedFromUser *string `json:"payedFromUser"`
	// ApprovedFromAdmin is the URL of the approving admin, or nil.
	ApprovedFromAdmin *string `json:"approvedFromAdmin"`
	// CanceledInvoice is the URL of the invoice that was canceled by this one, or nil.
	CanceledInvoice *string `json:"canceledInvoice"`
	// CancelInvoice is the URL of the cancellation invoice for this one, or nil.
	CancelInvoice *string `json:"cancelInvoice"`
	// Charges holds charge, chargeBack and total amounts.
	Charges InvoiceCharges `json:"charges"`
	// BankAccount is the URL of the linked bank account, or nil.
	BankAccount *string `json:"bankAccount"`
	// InvoiceItems holds URLs of the line items belonging to this invoice.
	InvoiceItems []string `json:"invoiceItems"`
	// DeleteAfterDate is the scheduled deletion date, or nil.
	DeleteAfterDate *string `json:"_deleteAfterDate"`
	// DeletedBy is the URL of the user who deleted this record, or nil.
	DeletedBy *string `json:"_deletedBy"`
	// Gross indicates whether prices include tax.
	Gross bool `json:"gross"`
	// CancellationDescription is the reason for cancellation, or nil.
	CancellationDescription *string `json:"cancellationDescription"`
	// TemplateName is the name of the document template used, or nil.
	TemplateName *string `json:"templateName"`
	// Date is the invoice date in YYYY-MM-DD format, or nil if not set.
	Date *string `json:"date"`
	// DateItHappend is the date the transaction occurred, or nil.
	DateItHappend *string `json:"dateItHappend"`
	// DateSent is the date the invoice was sent, or nil.
	DateSent *string `json:"dateSent"`
	// InvNumber is the human-readable invoice number.
	InvNumber string `json:"invNumber"`
	// Receiver is the name/address of the invoice recipient, or nil.
	Receiver *string `json:"receiver"`
	// Description is an optional free-text or HTML description, or nil.
	Description *string `json:"description"`
	// TotalPrice is the total invoice amount as a string (e.g. "49.99").
	TotalPrice flexFloat64 `json:"totalPrice"`
	// Tax is the tax amount as a string (e.g. "0.00").
	Tax flexFloat64 `json:"tax"`
	// Kind classifies the invoice (e.g. "expense", "membership", "cancel").
	Kind string `json:"kind"`
	// RefNumber is the reference number.
	RefNumber string `json:"refNumber"`
	// PaymentDifference is the outstanding balance as a string (e.g. "-49.99").
	PaymentDifference flexFloat64 `json:"paymentDifference"`
	// IsDraft indicates whether the invoice is still a draft.
	IsDraft bool `json:"isDraft"`
	// IsTemplate indicates whether the invoice is a template.
	IsTemplate bool `json:"isTemplate"`
	// CreationDateForRecurringInvoices is the start date for recurring billing, or nil.
	CreationDateForRecurringInvoices *string `json:"creationDateForRecurringInvoices"`
	// RecurringInvoicesInterval is the recurrence interval in days (-1 = disabled).
	RecurringInvoicesInterval int `json:"recurringInvoicesInterval"`
	// PaymentInformation describes the payment method (e.g. "debit", "nothing").
	PaymentInformation string `json:"paymentInformation"`
	// IsRequest indicates whether the invoice is a payment request.
	IsRequest bool `json:"isRequest"`
	// TaxRate is the applied tax rate as a string, or nil.
	TaxRate *string `json:"taxRate"`
	// TaxName is the name of the tax category, or nil.
	TaxName *string `json:"taxName"`
	// ActualCallStateName is the current dunning state name.
	ActualCallStateName string `json:"actualCallStateName"`
	// CallStateDelayDays is the number of delay days for the dunning state.
	CallStateDelayDays int `json:"callStateDelayDays"`
	// Accnumber is the accounting number.
	Accnumber int `json:"accnumber"`
	// GUID is the global unique identifier string (may be "0" when unset).
	GUID flexString `json:"guid"`
	// SelectionAcc is the URL of the selected accounting plan entry, or "0"/"" when unset.
	// The API returns an integer 0 instead of null when no account is selected.
	SelectionAcc flexString `json:"selectionAcc"`
	// RemoveFileOnDelete controls whether the PDF is deleted with the record.
	RemoveFileOnDelete bool `json:"removeFileOnDelete"`
	// CustomPaymentMethod is the URL of a custom payment method, or nil.
	CustomPaymentMethod *string `json:"customPaymentMethod"`
	// IsReceipt indicates whether the invoice acts as a receipt.
	IsReceipt bool `json:"isReceipt"`
	// IsTaxRatePerInvoiceItem indicates per-item tax rates, or nil if unset.
	IsTaxRatePerInvoiceItem *bool `json:"_isTaxRatePerInvoiceItem"`
	// IsSubjectToTax indicates VAT liability, or nil if unset.
	IsSubjectToTax *bool `json:"_isSubjectToTax"`
	// Mode is the document mode (e.g. "invoice").
	Mode string `json:"mode"`
	// OfferStatus is the offer state (e.g. "open").
	OfferStatus string `json:"offerStatus"`
	// OfferValidUntil is the offer expiry date, or nil.
	OfferValidUntil *string `json:"offerValidUntil"`
	// OfferNumber is the offer number string.
	OfferNumber string `json:"offerNumber"`
	// RelatedOffer is the URL of a linked offer, or nil.
	RelatedOffer *string `json:"relatedOffer"`
	// ClosingDescription is the closing text, or nil.
	ClosingDescription *string `json:"closingDescription"`
	// UseAddressBalance controls whether the address balance is used for payment.
	UseAddressBalance bool `json:"useAddressBalance"`
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

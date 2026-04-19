package model

// Booking represents a financial booking entry in easyVerein.
type Booking struct {
	// ID is the unique identifier of the booking.
	ID int `json:"id"`
	// Amount is the monetary amount of the booking.
	Amount flexFloat64 `json:"amount"`
	// Date is the booking date in ISO 8601 format.
	Date string `json:"date"`
	// Description is an optional free-text description.
	Description string `json:"description"`
	// Receiver is the counterpart of the transaction.
	Receiver string `json:"receiver"`
	// BillingID is an optional external billing reference.
	BillingID string `json:"billingId"`
	// RelatedInvoice holds related invoice URLs/IDs.
	RelatedInvoice []interface{} `json:"relatedInvoice"`
	// Org is the URL of the owning organization.
	Org string `json:"org"`
	// BankAccount is the URL of the associated bank account.
	BankAccount string `json:"bankAccount"`
	// BillingAccount is the URL of the associated billing account (nullable).
	BillingAccount *string `json:"billingAccount"`
	// DeleteAfterDate is the scheduled deletion date (nullable).
	DeleteAfterDate *string `json:"_deleteAfterDate"`
	// DeletedBy indicates who deleted the entry (nullable).
	DeletedBy *string `json:"_deletedBy"`
	// ImportDate is the Unix timestamp of when the booking was imported.
	ImportDate int64 `json:"importDate"`
	// Blocked indicates whether the booking is blocked.
	Blocked bool `json:"blocked"`
	// PaymentDifference is the difference between booked and actual payment.
	PaymentDifference flexFloat64 `json:"paymentDifference"`
	// CounterpartIban is the IBAN of the counterpart.
	CounterpartIban string `json:"counterpartIban"`
	// CounterpartBic is the BIC of the counterpart.
	CounterpartBic string `json:"counterpartBic"`
	// TwingleDonation indicates whether this booking is a Twingle donation.
	TwingleDonation bool `json:"twingleDonation"`
	// BookingProject is the associated booking project (nullable).
	BookingProject *string `json:"bookingProject"`
	// Sphere is an internal categorization value.
	Sphere int `json:"sphere"`
}

// BookingCreate holds the fields used when creating a new booking via
// POST /booking.
type BookingCreate struct {
	// Amount is the monetary amount (required).
	Amount float64 `json:"amount"`
	// BillingAccount is the ID of the billing account to book against (required).
	BillingAccount int `json:"billingAccount"`
	// Description is an optional free-text description.
	Description string `json:"description,omitempty"`
	// Date is the booking date in YYYY-MM-DD format (required).
	Date string `json:"date"`
	// Receiver is an optional counterpart name.
	Receiver string `json:"receiver,omitempty"`
	// RelatedInvoice is an optional list of invoice URLs to link to this booking.
	RelatedInvoice []string `json:"relatedInvoice,omitempty"`
}

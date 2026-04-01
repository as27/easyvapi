package model

// Booking represents a financial booking entry in easyVerein.
type Booking struct {
	// ID is the unique identifier of the booking.
	ID int `json:"id"`
	// Amount is the monetary amount of the booking.
	Amount flexFloat64 `json:"amount"`
	// Date is the booking date in YYYY-MM-DD format.
	Date string `json:"date"`
	// Description is an optional free-text description.
	Description string `json:"description"`
	// Receiver is the counterpart of the transaction.
	Receiver string `json:"receiver"`
	// BillingID is an optional external billing reference.
	BillingID string `json:"billingId"`
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
}

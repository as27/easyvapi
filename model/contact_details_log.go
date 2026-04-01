package model

// ContactDetailsLog represents a log entry for a contact record in easyVerein
// (Kontakthistorie). Log entries track communication and changes for a contact.
type ContactDetailsLog struct {
	// ID is the unique identifier of the log entry.
	ID int `json:"id"`
	// ContactDetails is the ID of the contact this log entry belongs to.
	ContactDetails int `json:"contactDetails"`
	// Title is a short subject line for the log entry.
	Title string `json:"title"`
	// Message is the full content of the log entry.
	Message string `json:"message"`
	// Date is the date of the log entry in YYYY-MM-DD format (set by the API).
	Date string `json:"date"`
}

// ContactDetailsLogCreate holds the fields for creating or updating a log entry
// via POST / PATCH /contact-details-log.
type ContactDetailsLogCreate struct {
	// ContactDetails is the ID of the contact this log entry belongs to (required).
	ContactDetails int `json:"contactDetails,omitempty"`
	// Title is the subject line of the log entry (required).
	Title string `json:"title,omitempty"`
	// Message is the full content of the log entry.
	Message string `json:"message,omitempty"`
}

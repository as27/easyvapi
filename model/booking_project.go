package model

// BookingProject represents a booking project in easyVerein (Buchungsprojekt).
// Booking projects allow grouping bookings under a common project for reporting.
type BookingProject struct {
	// ID is the unique identifier of the booking project.
	ID int `json:"id"`
	// Name is the display name of the booking project.
	Name string `json:"name"`
	// Description is an optional free-text description.
	Description string `json:"description"`
}

// BookingProjectCreate holds the fields for creating or updating a booking project
// via POST / PATCH /booking-project.
type BookingProjectCreate struct {
	// Name is the display name (required for create).
	Name string `json:"name,omitempty"`
	// Description is an optional free-text description.
	Description string `json:"description,omitempty"`
}

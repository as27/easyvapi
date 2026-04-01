package model

// Event represents a calendar event / Veranstaltung in easyVerein.
type Event struct {
	// ID is the unique identifier of the event.
	ID int `json:"id"`
	// Name is the title of the event.
	Name string `json:"name"`
	// Start is the start date/time in ISO 8601 format.
	Start string `json:"start"`
	// End is the end date/time in ISO 8601 format.
	End string `json:"end"`
	// AllDay indicates whether the event spans the full day.
	AllDay bool `json:"allDay"`
	// IsPublic indicates whether the event is publicly visible.
	IsPublic bool `json:"isPublic"`
	// Canceled indicates whether the event has been canceled.
	Canceled bool `json:"canceled"`
	// Description is an optional free-text description.
	Description string `json:"description"`
	// LocationName is the human-readable name of the venue.
	LocationName string `json:"locationName"`
}

// EventCreate holds the fields used when creating a new event via POST /event.
type EventCreate struct {
	// Name is the title of the event (required).
	Name string `json:"name"`
	// Start is the start date/time in ISO 8601 format (required).
	Start string `json:"start"`
	// End is the end date/time in ISO 8601 format (required).
	End string `json:"end"`
	// AllDay indicates whether the event spans the full day.
	AllDay bool `json:"allDay,omitempty"`
	// IsPublic indicates whether the event is publicly visible.
	IsPublic bool `json:"isPublic,omitempty"`
	// Description is an optional free-text description.
	Description string `json:"description,omitempty"`
	// LocationName is the human-readable name of the venue.
	LocationName string `json:"locationName,omitempty"`
}

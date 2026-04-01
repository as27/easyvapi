package model

// Calendar represents a calendar in easyVerein (Kalender).
// Events are always assigned to a calendar.
type Calendar struct {
	// ID is the unique identifier of the calendar.
	ID int `json:"id"`
	// Name is the display name of the calendar.
	Name string `json:"name"`
	// Description is an optional free-text description.
	Description string `json:"description"`
	// Color is the display color (hex code, e.g. "#3498db").
	Color string `json:"color"`
	// CalendarURL is an optional external iCal feed URL.
	CalendarURL string `json:"calendarURL"`
	// CalendarKind classifies the calendar type.
	CalendarKind string `json:"calendarKind"`
	// IsPublic indicates whether the calendar is publicly visible.
	IsPublic bool `json:"isPublic"`
	// LinkedCalendars holds IDs of calendars linked to this one.
	LinkedCalendars []int `json:"linkedCalendars"`
}

// CalendarCreate holds the fields for creating or updating a calendar
// via POST / PATCH /calendar.
type CalendarCreate struct {
	// Name is the display name (required for create).
	Name string `json:"name,omitempty"`
	// Description is an optional free-text description.
	Description string `json:"description,omitempty"`
	// Color is the display color.
	Color string `json:"color,omitempty"`
	// CalendarURL is an optional external iCal feed URL.
	CalendarURL string `json:"calendarURL,omitempty"`
	// CalendarKind classifies the calendar type.
	CalendarKind string `json:"calendarKind,omitempty"`
	// IsPublic indicates whether the calendar is publicly visible.
	IsPublic bool `json:"isPublic,omitempty"`
	// LinkedCalendars holds IDs of calendars linked to this one.
	LinkedCalendars []int `json:"linkedCalendars,omitempty"`
}

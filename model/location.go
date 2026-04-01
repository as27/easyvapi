package model

// Location represents a physical venue in easyVerein (Veranstaltungsort).
// Locations can be linked to events.
type Location struct {
	// ID is the unique identifier of the location.
	ID int `json:"id"`
	// Name is the display name of the location.
	Name string `json:"name"`
	// Description is an optional free-text description.
	Description string `json:"description"`
	// Street is the street address including number.
	Street string `json:"street"`
	// City is the city name.
	City string `json:"city"`
	// State is the federal state or region.
	State string `json:"state"`
	// Zip is the postal code.
	Zip string `json:"zip"`
	// Country is the ISO 3166-1 alpha-2 country code (e.g. "DE").
	Country string `json:"country"`
	// Longlat holds the geographic coordinates as a string (e.g. "13.405,52.52").
	Longlat string `json:"longlat"`
}

// LocationCreate holds the fields for creating or updating a location
// via POST / PATCH /location.
type LocationCreate struct {
	// Name is the display name (required for create).
	Name string `json:"name,omitempty"`
	// Description is an optional free-text description.
	Description string `json:"description,omitempty"`
	// Street is the street address.
	Street string `json:"street,omitempty"`
	// City is the city name.
	City string `json:"city,omitempty"`
	// State is the federal state or region.
	State string `json:"state,omitempty"`
	// Zip is the postal code.
	Zip string `json:"zip,omitempty"`
	// Country is the ISO 3166-1 alpha-2 country code.
	Country string `json:"country,omitempty"`
	// Longlat holds the geographic coordinates.
	Longlat string `json:"longlat,omitempty"`
}

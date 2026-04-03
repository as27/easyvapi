package model

// DosbSport represents a DOSB (Deutscher Olympischer Sportbund) sport entry.
type DosbSport struct {
	// ID is the unique identifier.
	ID int `json:"id"`
	// Name is the name of the sport.
	Name string `json:"name"`
	// Code is the DOSB sport code.
	Code string `json:"code"`
}

// DosbSportCreate holds the fields for creating or updating a DOSB sport entry.
type DosbSportCreate struct {
	// Name is the name of the sport.
	Name string `json:"name,omitempty"`
	// Code is the DOSB sport code.
	Code string `json:"code,omitempty"`
}

// LsbSport represents an LSB (Landessportbund) sport entry.
type LsbSport struct {
	// ID is the unique identifier.
	ID int `json:"id"`
	// Name is the name of the sport.
	Name string `json:"name"`
	// Code is the LSB sport code.
	Code string `json:"code"`
}

// LsbSportCreate holds the fields for creating or updating an LSB sport entry.
type LsbSportCreate struct {
	// Name is the name of the sport.
	Name string `json:"name,omitempty"`
	// Code is the LSB sport code.
	Code string `json:"code,omitempty"`
}

package model

// Lending represents an item lending record in easyVerein (Ausleihe).
type Lending struct {
	// ID is the unique identifier of the lending record.
	ID int `json:"id"`
	// InventoryObject is the ID of the borrowed inventory item.
	InventoryObject int `json:"inventoryObject"`
	// LendingPerson is the ID of the contact who borrowed the item.
	LendingPerson int `json:"lendingPerson"`
	// LendingStart is the start date/time of the lending period (ISO 8601).
	LendingStart string `json:"lendingStart"`
	// LendingEnd is the end date/time of the lending period (ISO 8601).
	LendingEnd string `json:"lendingEnd"`
	// State is the current lending state (e.g. "borrowed", "returned").
	State string `json:"state"`
	// Note is an optional free-text note about the lending.
	Note string `json:"note"`
}

// LendingCreate holds the fields for creating or updating a lending record
// via POST / PATCH /lending.
type LendingCreate struct {
	// InventoryObject is the ID of the item to borrow (required for create).
	InventoryObject int `json:"inventoryObject,omitempty"`
	// LendingPerson is the ID of the borrowing contact (required for create).
	LendingPerson int `json:"lendingPerson,omitempty"`
	// LendingStart is the start date/time (ISO 8601).
	LendingStart string `json:"lendingStart,omitempty"`
	// LendingEnd is the end date/time (ISO 8601).
	LendingEnd string `json:"lendingEnd,omitempty"`
	// State is the lending state.
	State string `json:"state,omitempty"`
	// Note is an optional free-text note.
	Note string `json:"note,omitempty"`
}

package model

// CustomFieldCollection represents a collection that groups custom fields in
// easyVerein (Feldkollektion).
type CustomFieldCollection struct {
	// ID is the unique identifier of the collection.
	ID int `json:"id"`
	// Name is the display name of the collection.
	Name string `json:"name"`
	// OrderSequence controls the display order of the collection.
	OrderSequence int `json:"orderSequence"`
	// Position is the display position of the collection.
	Position int `json:"position"`
}

// CustomFieldCollectionCreate holds the fields for creating or updating a
// custom field collection via POST / PATCH /custom-field-collection.
type CustomFieldCollectionCreate struct {
	// Name is the display name (required for create).
	Name string `json:"name,omitempty"`
	// OrderSequence controls the display order.
	OrderSequence int `json:"orderSequence,omitempty"`
	// Position is the display position.
	Position int `json:"position,omitempty"`
}

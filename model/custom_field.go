package model

// CustomField represents a custom field definition in easyVerein
// (Benutzerdefiniertes Feld). Custom fields extend member, contact, event
// or inventory records with organisation-specific attributes.
type CustomField struct {
	// ID is the unique identifier of the custom field.
	ID int `json:"id"`
	// Label is the display name of the field.
	Label string `json:"label"`
	// FieldKind defines the data type (e.g. "text", "number", "date", "select").
	FieldKind string `json:"fieldKind"`
	// OrderSequence controls the display order within its collection.
	OrderSequence int `json:"orderSequence"`
	// ShowInMemberArea indicates whether the field is visible in the member area.
	ShowInMemberArea bool `json:"showInMemberArea"`
	// FieldCollection is the ID of the collection this field belongs to.
	FieldCollection int `json:"fieldCollection"`
	// MaxSelections limits the number of selectable options for select fields.
	MaxSelections int `json:"maxSelections"`
}

// CustomFieldCreate holds the fields for creating or updating a custom field
// via POST / PATCH /custom-field.
type CustomFieldCreate struct {
	// Label is the display name (required for create).
	Label string `json:"label,omitempty"`
	// FieldKind defines the data type (required for create).
	FieldKind string `json:"fieldKind,omitempty"`
	// OrderSequence controls the display order within its collection.
	OrderSequence int `json:"orderSequence,omitempty"`
	// ShowInMemberArea indicates whether the field is visible in the member area.
	ShowInMemberArea bool `json:"showInMemberArea,omitempty"`
	// FieldCollection is the ID of the collection this field belongs to.
	FieldCollection int `json:"fieldCollection,omitempty"`
	// MaxSelections limits the number of selectable options for select fields.
	MaxSelections int `json:"maxSelections,omitempty"`
}

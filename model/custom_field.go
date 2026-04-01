package model

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
)

// urlOrInt unmarshals a JSON value that is either null, a plain integer, or a
// URL string whose last path segment is the integer ID (e.g.
// "https://…/custom-field-collection/42"). The extracted integer is stored in
// Value; zero means absent/null.
type urlOrInt int

func (u *urlOrInt) UnmarshalJSON(data []byte) error {
	// null → 0
	if string(data) == "null" {
		*u = 0
		return nil
	}
	// plain integer
	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		*u = urlOrInt(n)
		return nil
	}
	// URL string: extract last path segment as ID
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("urlOrInt: %w", err)
	}
	id, err := strconv.Atoi(path.Base(s))
	if err != nil {
		return fmt.Errorf("urlOrInt: cannot parse ID from %q: %w", s, err)
	}
	*u = urlOrInt(id)
	return nil
}

// CustomField represents a custom field definition in easyVerein
// (Benutzerdefiniertes Feld). Custom fields extend member, contact, event
// or inventory records with organisation-specific attributes.
//
// Note: the API uses different field names than the spec suggests.
// The actual JSON keys are mapped here as observed from the live API.
type CustomField struct {
	// ID is the unique identifier of the custom field.
	ID int `json:"id"`
	// Label is the display name of the field (API field: "name").
	Label string `json:"name"`
	// FieldKind defines the data type (API field: "kind", e.g. "s"=text, "e"=select).
	FieldKind string `json:"kind"`
	// OrderSequence controls the display order (API field: "position").
	OrderSequence int `json:"position"`
	// ShowInMemberArea indicates whether the field is visible in the member area
	// (API field: "member_show").
	ShowInMemberArea bool `json:"member_show"`
	// FieldCollection is the ID of the collection this field belongs to. The API
	// returns this as a URL string or null; the ID is extracted automatically.
	// (API field: "collection").
	FieldCollection urlOrInt `json:"collection"`
	// Description is an optional free-text description.
	Description string `json:"description"`
	// Placeholder is the generated placeholder variable for use in templates.
	Placeholder string `json:"placeHolder"`
}

// CustomFieldCreate holds the fields for creating or updating a custom field
// via POST / PATCH /custom-field.
type CustomFieldCreate struct {
	// Label is the display name (required for create, API field: "name").
	Label string `json:"name,omitempty"`
	// FieldKind defines the data type (API field: "kind").
	FieldKind string `json:"kind,omitempty"`
	// OrderSequence controls the display order (API field: "position").
	OrderSequence int `json:"position,omitempty"`
	// ShowInMemberArea indicates visibility in the member area (API field: "member_show").
	ShowInMemberArea bool `json:"member_show,omitempty"`
	// FieldCollection is the ID of the collection (API field: "collection").
	FieldCollection int `json:"collection,omitempty"`
	// Description is an optional free-text description.
	Description string `json:"description,omitempty"`
}

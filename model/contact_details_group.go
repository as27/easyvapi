package model

// ContactDetailsGroup represents a group of contact records in easyVerein
// (Kontaktgruppe). Groups allow organizing contacts into named categories.
type ContactDetailsGroup struct {
	// ID is the unique identifier of the contact details group.
	ID int `json:"id"`
	// Name is the display name of the group.
	Name string `json:"name"`
	// Description is an optional free-text description.
	Description string `json:"description"`
	// ContactDetails holds the IDs of the contacts assigned to this group.
	ContactDetails []int `json:"contactDetails"`
}

// ContactDetailsGroupCreate holds the fields for creating or updating a contact
// details group via POST / PATCH /contact-details-group.
type ContactDetailsGroupCreate struct {
	// Name is the display name (required for create).
	Name string `json:"name,omitempty"`
	// Description is an optional free-text description.
	Description string `json:"description,omitempty"`
	// ContactDetails is the list of contact-details IDs to assign to this group.
	ContactDetails []int `json:"contactDetails,omitempty"`
}

package model

// CustomFilter represents a saved filter definition in easyVerein
// (Benutzerdefinierter Filter). Filters can be applied to member or contact lists.
type CustomFilter struct {
	// ID is the unique identifier of the custom filter.
	ID int `json:"id"`
	// Name is the display name of the filter.
	Name string `json:"name"`
	// Model is the resource type this filter applies to (e.g. "member", "contact-details").
	Model string `json:"model"`
	// Rules is the raw filter rule definition (JSON encoded by the API).
	Rules any `json:"rules"`
}

// CustomFilterCreate holds the fields for creating or updating a custom filter
// via POST / PATCH /custom-filter.
type CustomFilterCreate struct {
	// Name is the display name (required for create).
	Name string `json:"name,omitempty"`
	// Rules is the filter rule definition.
	Rules any `json:"rules,omitempty"`
}

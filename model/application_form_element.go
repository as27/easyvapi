package model

// ApplicationFormElement represents a single field/element within an
// application form in easyVerein (Formularelement).
type ApplicationFormElement struct {
	// ID is the unique identifier of the form element.
	ID int `json:"id"`
	// ApplicationForm is the ID of the form this element belongs to.
	// The API may return this as a URL string or plain integer.
	ApplicationForm urlOrInt `json:"applicationForm"`
	// Kind classifies the element type (e.g. "text", "select", "checkbox").
	Kind string `json:"kind"`
	// Label is the display label for the element.
	Label string `json:"label"`
	// Content is the element content or options.
	Content string `json:"content"`
	// DefaultValue is the pre-filled default value.
	DefaultValue string `json:"defaultValue"`
	// Required indicates whether the element must be filled out.
	Required bool `json:"required"`
	// Position controls the display order within the form.
	Position int `json:"position"`
	// MaxMemberGroupCount limits the number of selectable member groups.
	MaxMemberGroupCount int `json:"maxMemberGroupCount"`
	// AllowedMemberGroups holds IDs of member groups that may be selected.
	AllowedMemberGroups []int `json:"allowedMemberGroups"`
}

// ApplicationFormElementCreate holds the fields for creating or updating a
// form element via POST / PATCH /application-form-element.
type ApplicationFormElementCreate struct {
	// MaxMemberGroupCount limits selectable member groups.
	MaxMemberGroupCount int `json:"maxMemberGroupCount,omitempty"`
	// AllowedMemberGroups holds IDs of selectable member groups.
	AllowedMemberGroups []int `json:"allowedMemberGroups,omitempty"`
}

package model

// ApplicationForm represents a member application form in easyVerein
// (Aufnahmeformular). Forms are used for new member registrations.
type ApplicationForm struct {
	// ID is the unique identifier of the application form.
	ID int `json:"id"`
	// Title is the display name of the form.
	Title string `json:"title"`
	// Public indicates whether the form is publicly accessible.
	Public bool `json:"public"`
	// HidePrivacyNotice suppresses the privacy notice on the form.
	HidePrivacyNotice bool `json:"hidePrivacyNotice"`
	// ShowInSelect indicates whether the form appears in selection dropdowns.
	ShowInSelect bool `json:"showInSelect"`
	// Language is the form language code (e.g. "de", "en").
	Language string `json:"language"`
	// FormularKind classifies the form type.
	FormularKind string `json:"formularKind"`
}

// ApplicationFormCreate holds the fields for creating or updating an application
// form via POST / PATCH /application-form.
type ApplicationFormCreate struct {
	// Title is the display name (required for create).
	Title string `json:"title,omitempty"`
	// Public indicates whether the form is publicly accessible.
	Public bool `json:"public,omitempty"`
	// HidePrivacyNotice suppresses the privacy notice.
	HidePrivacyNotice bool `json:"hidePrivacyNotice,omitempty"`
	// ShowInSelect indicates whether the form appears in selection dropdowns.
	ShowInSelect bool `json:"showInSelect,omitempty"`
	// Language is the form language code.
	Language string `json:"language,omitempty"`
	// FormularKind classifies the form type.
	FormularKind string `json:"formularKind,omitempty"`
}

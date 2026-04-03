package model

// ChatSettings represents the chat configuration for the organization.
type ChatSettings struct {
	// ID is the unique identifier.
	ID int `json:"id"`
	// Enabled indicates whether the chat feature is active.
	Enabled bool `json:"enabled"`
	// AllowMemberToMember controls whether members can message each other.
	AllowMemberToMember bool `json:"allowMemberToMember"`
}

// ChatSettingsCreate holds the fields for updating chat settings via PATCH.
type ChatSettingsCreate struct {
	// Enabled indicates whether the chat feature is active.
	Enabled *bool `json:"enabled,omitempty"`
	// AllowMemberToMember controls whether members can message each other.
	AllowMemberToMember *bool `json:"allowMemberToMember,omitempty"`
}

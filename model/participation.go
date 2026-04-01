package model

// Participation represents a participation record for an event in easyVerein
// (Veranstaltungsteilnahme).
type Participation struct {
	// ID is the unique identifier of the participation record.
	ID int `json:"id"`
	// ParticipationAddress is the ID of the contact-details record of the participant.
	ParticipationAddress int `json:"participationAddress"`
	// Name is an optional display name for the participation.
	Name string `json:"name"`
	// ShowName indicates whether the participant's name is shown publicly.
	ShowName bool `json:"showName"`
	// State is the participation state (integer code defined by the API).
	State int `json:"state"`
	// Description is an optional free-text note about this participation.
	Description string `json:"description"`
	// IsCompanion indicates whether this participant is a companion of another.
	IsCompanion bool `json:"isCompanion"`
}

// ParticipationCreate holds the fields for creating or updating a participation
// via POST / PATCH /event/{eventPk}/participation.
type ParticipationCreate struct {
	// ParticipationAddress is the contact-details ID of the participant (required).
	ParticipationAddress int `json:"participationAddress,omitempty"`
	// Name is an optional display name.
	Name string `json:"name,omitempty"`
	// ShowName indicates whether the name is shown publicly.
	ShowName bool `json:"showName,omitempty"`
	// State is the participation state.
	State int `json:"state,omitempty"`
	// Description is an optional free-text note.
	Description string `json:"description,omitempty"`
	// IsCompanion indicates whether this participant is a companion.
	IsCompanion bool `json:"isCompanion,omitempty"`
}

// InviteGroupsRequest holds the parameters for inviting groups to an event
// via POST /event/{pk}/invite-groups.
type InviteGroupsRequest struct {
	// ContactDetailsGroupID is the ID of the contact-details group to invite.
	ContactDetailsGroupID int `json:"contactDetailsGroupId,omitempty"`
	// MemberGroupID is the ID of the member group to invite.
	MemberGroupID int `json:"memberGroupId,omitempty"`
}

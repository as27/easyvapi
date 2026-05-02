package model

// Protocol represents a protocol (Protokoll) in easyVerein.
// Protocols document meetings and are linked to events.
type Protocol struct {
	// ID is the unique identifier of the protocol.
	ID int `json:"id"`
	// Org is the organization ID.
	Org int `json:"org"`
	// LocationObject is the ID of the linked location, or nil.
	LocationObject *int `json:"locationObject"`
	// Calendar is the ID of the linked calendar, or nil.
	Calendar *int `json:"calendar"`
	// ProtocolElements holds references to protocol elements.
	ProtocolElements []interface{} `json:"protocolElements"`
	// ProtocolUploads holds references to protocol uploads.
	ProtocolUploads []interface{} `json:"protocolUploads"`
	// AllowedGroups holds IDs of member groups with access.
	AllowedGroups []int `json:"allowedGroups"`
	// DeleteAfterDate is an optional date after which the entry should be deleted.
	DeleteAfterDate *string `json:"_deleteAfterDate"`
	// DeletedBy is the user who deleted the entry (if in wastebasket).
	DeletedBy *string `json:"_deletedBy"`
	// UID is the unique calendar identifier.
	UID string `json:"uid"`
	// Name is the title of the protocol.
	Name string `json:"name"`
	// LocationName is the human-readable location name.
	LocationName string `json:"locationName"`
	// Description is the detailed content of the protocol (may contain HTML).
	Description string `json:"description"`
	// Prologue is an optional foreword.
	Prologue string `json:"prologue"`
	// MinParticipators is the minimum number of participants.
	MinParticipators int `json:"minParticipators"`
	// MaxParticipators is the maximum number of participants (0 = unlimited).
	MaxParticipators int `json:"maxParticipators"`
	// StartParticipation is the start of the registration period, or nil.
	StartParticipation *string `json:"startParticipation"`
	// EndParticipation is the end of the registration period, or nil.
	EndParticipation *string `json:"endParticipation"`
	// Access is the access level as numeric value.
	Access int `json:"access"`
	// Note is an optional internal note.
	Note string `json:"note"`
	// Start is the start date/time in ISO 8601 format.
	Start string `json:"start"`
	// End is the end date/time in ISO 8601 format.
	End string `json:"end"`
	// AllDay indicates whether the event lasts the entire day.
	AllDay bool `json:"allDay"`
	// Weekdays are the weekdays on which the event takes place.
	Weekdays *string `json:"weekdays"`
	// ConfirmationToAddresses holds IDs of addresses that receive confirmations.
	ConfirmationToAddresses []int `json:"confirmationToAddresses"`
	// SendMailCheck indicates whether email notifications are enabled.
	SendMailCheck bool `json:"sendMailCheck"`
	// ShowMemberarea indicates whether to show in the member area.
	ShowMemberarea bool `json:"showMemberarea"`
	// IsPublic indicates whether the protocol is publicly visible.
	IsPublic bool `json:"isPublic"`
	// MassParticipations indicates whether bulk registrations are allowed.
	MassParticipations bool `json:"massParticipations"`
	// Visible indicates whether the protocol is visible to members.
	Visible bool `json:"visible"`
	// MeetingLeader is the name of the meeting leader.
	MeetingLeader string `json:"meetingLeader"`
	// MeetingSecretary is the name of the meeting secretary / minute-taker.
	MeetingSecretary string `json:"meetingSecretary"`
	// IsLocked indicates whether the protocol is locked from editing.
	IsLocked bool `json:"isLocked"`
	// ParentProtocolElements holds references to parent protocol elements.
	ParentProtocolElements []interface{} `json:"parentProtocolElements"`
	// PublicProtocolUploads holds references to public protocol uploads.
	PublicProtocolUploads []interface{} `json:"publicProtocolUploads"`
}

// ProtocolCreate holds the fields for creating or updating a protocol
// via POST / PATCH /protocol.
type ProtocolCreate struct {
	// LocationObject is the ID of the linked location.
	LocationObject *int `json:"locationObject,omitempty"`
	// AllowedGroups holds IDs of member groups with access.
	AllowedGroups []int `json:"allowedGroups,omitempty"`
	// Name is the title of the protocol (required for create, max 500 chars).
	Name string `json:"name,omitempty"`
	// LocationName is the human-readable location name.
	LocationName string `json:"locationName,omitempty"`
	// Description is the detailed content (may contain HTML).
	Description string `json:"description,omitempty"`
	// Prologue is an optional foreword.
	Prologue string `json:"prologue,omitempty"`
	// MinParticipators is the minimum number of participants.
	MinParticipators *int `json:"minParticipators,omitempty"`
	// MaxParticipators is the maximum number of participants (0 = unlimited).
	MaxParticipators *int `json:"maxParticipators,omitempty"`
	// StartParticipation is the start of the registration period.
	StartParticipation *string `json:"startParticipation,omitempty"`
	// EndParticipation is the end of the registration period.
	EndParticipation *string `json:"endParticipation,omitempty"`
	// Access is the access level as numeric value.
	Access *int `json:"access,omitempty"`
	// Note is an optional internal note.
	Note string `json:"note,omitempty"`
	// Start is the start date/time in ISO 8601 format.
	Start string `json:"start,omitempty"`
	// End is the end date/time in ISO 8601 format.
	End string `json:"end,omitempty"`
	// AllDay indicates whether the event lasts the entire day.
	AllDay *bool `json:"allDay,omitempty"`
	// Weekdays are the weekdays on which the event takes place.
	Weekdays string `json:"weekdays,omitempty"`
	// ConfirmationToAddresses holds IDs of addresses that receive confirmations.
	ConfirmationToAddresses []int `json:"confirmationToAddresses,omitempty"`
	// SendMailCheck indicates whether email notifications are enabled.
	SendMailCheck *bool `json:"sendMailCheck,omitempty"`
	// ShowMemberarea indicates whether to show in the member area.
	ShowMemberarea *bool `json:"showMemberarea,omitempty"`
	// IsPublic indicates whether the protocol is publicly visible.
	IsPublic *bool `json:"isPublic,omitempty"`
	// MassParticipations indicates whether bulk registrations are allowed.
	MassParticipations *bool `json:"massParticipations,omitempty"`
	// Visible indicates whether the protocol is visible to members.
	Visible *bool `json:"visible,omitempty"`
	// MeetingLeader is the name of the meeting leader (max 128 chars).
	MeetingLeader string `json:"meetingLeader,omitempty"`
	// MeetingSecretary is the name of the meeting secretary (max 128 chars).
	MeetingSecretary string `json:"meetingSecretary,omitempty"`
	// IsLocked indicates whether the protocol is locked from editing.
	IsLocked *bool `json:"isLocked,omitempty"`
}

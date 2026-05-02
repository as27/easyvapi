package model

// Task represents a task (Aufgabe) in easyVerein.
type Task struct {
	// ID is the unique identifier of the task.
	ID int `json:"id"`
	// Org is the URL reference to the organization.
	Org string `json:"org"`
	// ParentEvent is the URL reference to the related event, or nil.
	ParentEvent *string `json:"parentEvent"`
	// Member is the URL reference to the assigned member, or nil.
	Member *string `json:"member"`
	// TaskGroup is the URL reference to the task group, or nil.
	TaskGroup *string `json:"taskGroup"`
	// ProtocolElement is the URL reference to the related protocol element, or nil.
	ProtocolElement *string `json:"protocolElement"`
	// TaskComments holds the URL references to comments on this task.
	TaskComments []string `json:"taskComments"`
	// DeleteAfterDate is an optional date after which the entry should be deleted.
	DeleteAfterDate *string `json:"_deleteAfterDate"`
	// DeletedBy is the user who deleted the entry (if in wastebasket).
	DeletedBy *string `json:"_deletedBy"`
	// Name is the title of the task.
	Name string `json:"name"`
	// Description is the detailed text of the task.
	Description string `json:"description"`
	// Due is the due date/time in ISO 8601 format, or nil if no due date.
	Due *string `json:"due"`
	// State is the current state of the task (e.g. "offen", "erledigt").
	State string `json:"state"`
	// Public indicates whether the task is visible to all members.
	Public bool `json:"public"`
	// CommentAvailable indicates whether comments are enabled (read-only).
	CommentAvailable bool `json:"commentAvailable"`
	// SendMailCheck indicates whether an email notification was sent.
	SendMailCheck bool `json:"sendMailCheck"`
	// Starttime is the start time (HH:MM:SS.ffffff), or nil.
	Starttime *string `json:"starttime"`
	// Endtime is the end time (HH:MM:SS.ffffff), or nil.
	Endtime *string `json:"endtime"`
}

// TaskCreate holds the fields for creating or updating a task
// via POST / PATCH /task.
type TaskCreate struct {
	// Name is the title of the task (required for create).
	Name string `json:"name,omitempty"`
	// Description is the detailed text of the task.
	Description string `json:"description,omitempty"`
	// Due is the due date/time in ISO 8601 format.
	Due *string `json:"due,omitempty"`
	// State is the current state (e.g. "offen", "erledigt").
	State string `json:"state,omitempty"`
	// Public indicates whether the task is visible to all members.
	Public *bool `json:"public,omitempty"`
	// SendMailCheck indicates whether to send an email notification.
	SendMailCheck *bool `json:"sendMailCheck,omitempty"`
	// Starttime is the start time (HH:MM:SS).
	Starttime *string `json:"starttime,omitempty"`
	// Endtime is the end time (HH:MM:SS).
	Endtime *string `json:"endtime,omitempty"`
	// Member is the ID of the member to assign the task to.
	Member *int `json:"member,omitempty"`
	// ParentEvent is the ID of the related event.
	ParentEvent *int `json:"parentEvent,omitempty"`
	// TaskGroup is the ID of the task group.
	TaskGroup *int `json:"taskGroup,omitempty"`
	// ProtocolElement is the ID of the related protocol element.
	ProtocolElement *int `json:"protocolElement,omitempty"`
}

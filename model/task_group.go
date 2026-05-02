package model

// TaskGroup represents a task group in easyVerein (Aufgabengruppe).
// Task groups are used to organize and categorize tasks.
type TaskGroup struct {
	// ID is the unique identifier of the task group.
	ID int `json:"id"`
	// Org is the URL reference to the organization.
	Org string `json:"org"`
	// Name is the display name of the group.
	Name string `json:"name"`
	// Color is the hex color code of the group (e.g. "#34e8eb").
	Color string `json:"color"`
	// Short is the abbreviation of the group (max 4 chars).
	Short string `json:"short"`
	// LinkedItems is the number of tasks linked to this group (read-only).
	LinkedItems int `json:"linkedItems"`
	// DeleteAfterDate is an optional date after which the entry should be deleted.
	DeleteAfterDate *string `json:"_deleteAfterDate"`
	// DeletedBy is the user who deleted the entry (if in wastebasket).
	DeletedBy *string `json:"_deletedBy"`
}

// TaskGroupCreate holds the fields for creating or updating a task group
// via POST / PATCH /task-group.
type TaskGroupCreate struct {
	// Name is the display name (required for create, max 200 chars).
	Name string `json:"name,omitempty"`
	// Color is the hex color code (required for groups, max 7 chars).
	Color string `json:"color,omitempty"`
	// Short is the abbreviation (required for groups, max 4 chars).
	Short string `json:"short,omitempty"`
}

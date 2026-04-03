package model

// Forum represents a discussion forum in easyVerein.
type Forum struct {
	// ID is the unique identifier of the forum.
	ID int `json:"id"`
	// Name is the display name of the forum.
	Name string `json:"name"`
	// Description is an optional description of the forum.
	Description string `json:"description"`
	// Public indicates whether the forum is publicly visible.
	Public bool `json:"public"`
}

// ForumCreate holds the fields for creating or updating a forum.
type ForumCreate struct {
	// Name is the display name.
	Name string `json:"name,omitempty"`
	// Description is an optional description.
	Description string `json:"description,omitempty"`
	// Public indicates whether the forum is publicly visible.
	Public bool `json:"public,omitempty"`
}

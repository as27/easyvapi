package model

// ChairmanNote represents an internal note for board members in easyVerein
// (Vorstandsnotiz). Notes are visible only to users with appropriate access.
type ChairmanNote struct {
	// ID is the unique identifier of the chairman note.
	ID int `json:"id"`
	// Text is the content of the note.
	Text string `json:"text"`
	// Date is the date of the note in YYYY-MM-DD format.
	Date string `json:"date"`
	// DeleteAfterDate is an optional date after which the note should be deleted.
	DeleteAfterDate string `json:"_deleteAfterDate"`
}

// ChairmanNoteCreate holds the fields for creating or updating a chairman note
// via POST / PATCH /chairman-note.
type ChairmanNoteCreate struct {
	// Text is the content of the note (required for create).
	Text string `json:"text,omitempty"`
	// Date is the date of the note in YYYY-MM-DD format.
	Date string `json:"date,omitempty"`
	// DeleteAfterDate sets an optional expiry date for the note.
	DeleteAfterDate string `json:"_deleteAfterDate,omitempty"`
}

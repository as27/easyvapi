package model

// AnniversaryMailing represents a scheduled anniversary notification in
// easyVerein (Jubiläumsbenachrichtigung). These are automated emails sent
// to members on their membership anniversary or birthday.
type AnniversaryMailing struct {
	// ID is the unique identifier of the anniversary mailing.
	ID int `json:"id"`
	// Name is the display name of the mailing configuration.
	Name string `json:"name"`
	// Subject is the email subject line.
	Subject string `json:"subject"`
	// Content is the email body (may contain HTML and placeholders).
	Content string `json:"content"`
	// AnniversaryYears holds the years that trigger the mailing (e.g. [1, 5, 10]).
	AnniversaryYears []int `json:"anniversaryYears"`
	// AnniversaryKind classifies the anniversary type (e.g. membership, birthday).
	AnniversaryKind int `json:"anniversaryKind"`
}

// AnniversaryMailingCreate holds the fields for creating or updating an
// anniversary mailing via POST / PATCH /anniversary-mailing.
type AnniversaryMailingCreate struct {
	// Name is the display name (required for create).
	Name string `json:"name,omitempty"`
	// Subject is the email subject line.
	Subject string `json:"subject,omitempty"`
	// Content is the email body.
	Content string `json:"content,omitempty"`
	// AnniversaryYears holds the years that trigger the mailing (e.g. [1, 5, 10]).
	AnniversaryYears []int `json:"anniversaryYears,omitempty"`
	// AnniversaryKind classifies the anniversary type.
	AnniversaryKind int `json:"anniversaryKind,omitempty"`
}

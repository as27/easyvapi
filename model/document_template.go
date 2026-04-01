package model

// DocumentTemplate represents a document template in easyVerein
// (Dokumentvorlage). Templates are used to generate letters, certificates,
// or other documents for members.
type DocumentTemplate struct {
	// ID is the unique identifier of the document template.
	ID int `json:"id"`
	// Title is the display name of the template.
	Title string `json:"title"`
	// Content is the template body (may contain HTML and placeholder variables).
	Content string `json:"content"`
	// DocumentKind classifies the template type (e.g. "letter", "certificate").
	DocumentKind string `json:"documentKind"`
	// SignatureKind defines the signature style for the document.
	SignatureKind string `json:"signatureKind"`
}

// DocumentTemplateCreate holds the fields for creating or updating a document
// template via POST / PATCH /document-template.
type DocumentTemplateCreate struct {
	// Title is the display name (required for create).
	Title string `json:"title,omitempty"`
	// Content is the template body.
	Content string `json:"content,omitempty"`
	// DocumentKind classifies the template type.
	DocumentKind string `json:"documentKind,omitempty"`
	// SignatureKind defines the signature style.
	SignatureKind string `json:"signatureKind,omitempty"`
}

package model

// WastebasketItem represents a deleted record that can be restored.
type WastebasketItem struct {
	// ID is the unique identifier of the wastebasket entry.
	ID int `json:"id"`
	// ObjectID is the ID of the deleted object.
	ObjectID int `json:"objectId"`
	// Model is the type/model of the deleted object (e.g. "member", "invoice").
	Model string `json:"model"`
	// DeletedAt is the ISO-8601 timestamp of when the object was deleted.
	DeletedAt string `json:"deletedAt"`
}

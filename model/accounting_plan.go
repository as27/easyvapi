package model

// AccountingPlan represents an accounting plan entry in easyVerein (Kontenplan).
type AccountingPlan struct {
	// ID is the unique identifier of the accounting plan entry.
	ID int `json:"id"`
	// Name is the display name of the accounting plan entry.
	Name string `json:"name"`
	// Description is an optional free-text description.
	Description string `json:"description"`
}

// AccountingPlanCreate holds the fields for creating or updating an accounting plan entry
// via POST / PUT /accounting-plan.
type AccountingPlanCreate struct {
	// Name is the display name.
	Name string `json:"name,omitempty"`
	// Description is an optional free-text description.
	Description string `json:"description,omitempty"`
}

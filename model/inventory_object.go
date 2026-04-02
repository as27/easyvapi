package model

// InventoryObject represents an inventory item in easyVerein (Inventargegenstand).
type InventoryObject struct {
	// ID is the unique identifier of the inventory object.
	ID int `json:"id"`
	// Name is the display name of the item.
	Name string `json:"name"`
	// Identifier is an optional asset tag or serial number.
	Identifier string `json:"identifier"`
	// Description is an optional free-text description.
	Description string `json:"description"`
	// Pieces is the number of available units.
	Pieces int `json:"pieces"`
	// Price is the purchase price of the item.
	Price flexFloat64 `json:"price"`
	// PurchaseDate is the date the item was purchased (YYYY-MM-DD).
	PurchaseDate string `json:"purchaseDate"`
	// LocationName is the human-readable storage location.
	LocationName string `json:"locationName"`
	// LendingAvailable indicates whether the item can be borrowed.
	LendingAvailable bool `json:"lendingAvailable"`
}

// InventoryObjectCreate holds the fields for creating or updating an inventory
// object via POST / PATCH /inventory-object.
type InventoryObjectCreate struct {
	// Name is the display name (required for create).
	Name string `json:"name,omitempty"`
	// Identifier is an optional asset tag or serial number.
	Identifier string `json:"identifier,omitempty"`
	// Description is an optional free-text description.
	Description string `json:"description,omitempty"`
	// Pieces is the number of available units.
	Pieces int `json:"pieces,omitempty"`
	// Price is the purchase price.
	Price float64 `json:"price,omitempty"`
	// PurchaseDate is the purchase date in YYYY-MM-DD format.
	PurchaseDate string `json:"purchaseDate,omitempty"`
	// LocationName is the human-readable storage location.
	LocationName string `json:"locationName,omitempty"`
	// LendingAvailable indicates whether the item can be borrowed.
	LendingAvailable bool `json:"lendingAvailable,omitempty"`
	// LendingResponsible is the ID of the contact responsible for lending.
	LendingResponsible int `json:"lendingResponsible,omitempty"`
}

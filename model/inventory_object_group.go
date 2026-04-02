package model

// InventoryObjectGroup represents a group of inventory objects in easyVerein
// (Inventargruppe).
type InventoryObjectGroup struct {
	// ID is the unique identifier of the inventory object group.
	ID int `json:"id"`
	// Name is the display name of the group.
	Name string `json:"name"`
	// Description is an optional free-text description.
	Description string `json:"description"`
	// InventoryObjects holds the IDs of inventory objects in this group.
	InventoryObjects []int `json:"inventoryObjects"`
}

// InventoryObjectGroupCreate holds the fields for creating or updating an
// inventory object group via POST / PATCH /inventory-object-group.
type InventoryObjectGroupCreate struct {
	// Name is the display name (required for create).
	Name string `json:"name,omitempty"`
	// Description is an optional free-text description.
	Description string `json:"description,omitempty"`
	// InventoryObjects holds the IDs of inventory objects to assign.
	InventoryObjects []int `json:"inventoryObjects,omitempty"`
}

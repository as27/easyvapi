package model

// ChairmanLevel represents an access/permission level for board members in
// easyVerein (Vorstandsebene / Zugriffsebene).
type ChairmanLevel struct {
	// ID is the unique identifier of the chairman level.
	ID int `json:"id"`
	// Name is the display name of the level.
	Name string `json:"name"`
	// Color is the display color (hex code, e.g. "#ff0000").
	Color string `json:"color"`
	// Short is an abbreviation for the level.
	Short string `json:"short"`
	// ModuleMembers grants access to the members module.
	ModuleMembers bool `json:"module_members"`
	// ModuleEvents grants access to the events module.
	ModuleEvents bool `json:"module_events"`
	// ModuleProtocols grants access to the protocols module.
	ModuleProtocols bool `json:"module_protocols"`
	// ModuleAddresses grants access to the addresses module.
	ModuleAddresses bool `json:"module_addresses"`
	// ModuleBookings grants access to the bookings module.
	ModuleBookings bool `json:"module_bookings"`
	// ModuleInventory grants access to the inventory module.
	ModuleInventory bool `json:"module_inventory"`
	// ModuleFiles grants access to the files module.
	ModuleFiles bool `json:"module_files"`
	// ModuleAccount grants access to the account module.
	ModuleAccount bool `json:"module_account"`
	// ModuleTodo grants access to the todo module.
	ModuleTodo bool `json:"module_todo"`
	// ModuleVotings grants access to the votings module.
	ModuleVotings bool `json:"module_votings"`
	// ModuleForum grants access to the forum module.
	ModuleForum bool `json:"module_forum"`
}

// ChairmanLevelCreate holds the fields for creating or updating a chairman level
// via POST / PATCH /chairman-level.
type ChairmanLevelCreate struct {
	// Name is the display name (required for create).
	Name string `json:"name,omitempty"`
	// Color is the display color (hex code).
	Color string `json:"color,omitempty"`
	// Short is an abbreviation for the level.
	Short string `json:"short,omitempty"`
	// ModuleMembers grants access to the members module.
	ModuleMembers bool `json:"module_members,omitempty"`
	// ModuleEvents grants access to the events module.
	ModuleEvents bool `json:"module_events,omitempty"`
	// ModuleProtocols grants access to the protocols module.
	ModuleProtocols bool `json:"module_protocols,omitempty"`
	// ModuleAddresses grants access to the addresses module.
	ModuleAddresses bool `json:"module_addresses,omitempty"`
	// ModuleBookings grants access to the bookings module.
	ModuleBookings bool `json:"module_bookings,omitempty"`
	// ModuleInventory grants access to the inventory module.
	ModuleInventory bool `json:"module_inventory,omitempty"`
	// ModuleFiles grants access to the files module.
	ModuleFiles bool `json:"module_files,omitempty"`
	// ModuleAccount grants access to the account module.
	ModuleAccount bool `json:"module_account,omitempty"`
	// ModuleTodo grants access to the todo module.
	ModuleTodo bool `json:"module_todo,omitempty"`
	// ModuleVotings grants access to the votings module.
	ModuleVotings bool `json:"module_votings,omitempty"`
	// ModuleForum grants access to the forum module.
	ModuleForum bool `json:"module_forum,omitempty"`
}

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
	// ModuleMembers grants access to the members module ("none", "read", "write").
	ModuleMembers string `json:"module_members"`
	// ModuleEvents grants access to the events module ("none", "read", "write").
	ModuleEvents string `json:"module_events"`
	// ModuleProtocols grants access to the protocols module ("none", "read", "write").
	ModuleProtocols string `json:"module_protocols"`
	// ModuleAddresses grants access to the addresses module ("none", "read", "write").
	ModuleAddresses string `json:"module_addresses"`
	// ModuleBookings grants access to the bookings module ("none", "read", "write").
	ModuleBookings string `json:"module_bookings"`
	// ModuleInventory grants access to the inventory module ("none", "read", "write").
	ModuleInventory string `json:"module_inventory"`
	// ModuleFiles grants access to the files module ("none", "read", "write").
	ModuleFiles string `json:"module_files"`
	// ModuleAccount grants access to the account module ("none", "read", "write").
	ModuleAccount string `json:"module_account"`
	// ModuleTodo grants access to the todo module ("none", "read", "write").
	ModuleTodo string `json:"module_todo"`
	// ModuleVotings grants access to the votings module ("none", "read", "write").
	ModuleVotings string `json:"module_votings"`
	// ModuleForum grants access to the forum module ("none", "read", "write").
	ModuleForum string `json:"module_forum"`
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
	// ModuleMembers grants access to the members module ("none", "read", "write").
	ModuleMembers string `json:"module_members,omitempty"`
	// ModuleEvents grants access to the events module ("none", "read", "write").
	ModuleEvents string `json:"module_events,omitempty"`
	// ModuleProtocols grants access to the protocols module ("none", "read", "write").
	ModuleProtocols string `json:"module_protocols,omitempty"`
	// ModuleAddresses grants access to the addresses module ("none", "read", "write").
	ModuleAddresses string `json:"module_addresses,omitempty"`
	// ModuleBookings grants access to the bookings module ("none", "read", "write").
	ModuleBookings string `json:"module_bookings,omitempty"`
	// ModuleInventory grants access to the inventory module ("none", "read", "write").
	ModuleInventory string `json:"module_inventory,omitempty"`
	// ModuleFiles grants access to the files module ("none", "read", "write").
	ModuleFiles string `json:"module_files,omitempty"`
	// ModuleAccount grants access to the account module ("none", "read", "write").
	ModuleAccount string `json:"module_account,omitempty"`
	// ModuleTodo grants access to the todo module ("none", "read", "write").
	ModuleTodo string `json:"module_todo,omitempty"`
	// ModuleVotings grants access to the votings module ("none", "read", "write").
	ModuleVotings string `json:"module_votings,omitempty"`
	// ModuleForum grants access to the forum module ("none", "read", "write").
	ModuleForum string `json:"module_forum,omitempty"`
}

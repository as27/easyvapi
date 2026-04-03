package model

// FileSystemPath represents a path entry in the easyVerein file system.
type FileSystemPath struct {
	// ID is the unique identifier.
	ID int `json:"id"`
	// Name is the display name of the path entry.
	Name string `json:"name"`
	// Path is the file system path string.
	Path string `json:"path"`
	// Parent is the ID of the parent path, or 0 for root.
	Parent int `json:"parent"`
}

// FileSystemPathCreate holds the fields for creating or updating a file system path.
type FileSystemPathCreate struct {
	// Name is the display name.
	Name string `json:"name,omitempty"`
	// Path is the file system path string.
	Path string `json:"path,omitempty"`
	// Parent is the ID of the parent path.
	Parent int `json:"parent,omitempty"`
}

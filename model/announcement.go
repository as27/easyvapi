package model

// Announcement represents a system announcement or banner in easyVerein
// (Ankündigung). Announcements can be shown as banners in the member area
// or admin interface.
type Announcement struct {
	// ID is the unique identifier of the announcement.
	ID int `json:"id"`
	// Title is the headline of the announcement.
	Title string `json:"title"`
	// Content is the body text of the announcement (may contain HTML).
	Content string `json:"content"`
	// ShowBanner indicates whether the announcement is displayed as a banner.
	ShowBanner bool `json:"showBanner"`
	// Platform specifies which platform shows the announcement (integer code).
	Platform int `json:"platform"`
	// StartDate is the date from which the announcement is visible (YYYY-MM-DD).
	StartDate string `json:"startDate"`
	// EndDate is the date until which the announcement is visible (YYYY-MM-DD).
	EndDate string `json:"endDate"`
}

// AnnouncementCreate holds the fields for creating or updating an announcement
// via POST / PATCH /announcement.
type AnnouncementCreate struct {
	// Title is the headline (required for create).
	Title string `json:"title,omitempty"`
	// Content is the body text.
	Content string `json:"content,omitempty"`
	// ShowBanner indicates whether to display as a banner.
	ShowBanner bool `json:"showBanner,omitempty"`
	// Platform specifies the target platform (integer code).
	Platform int `json:"platform,omitempty"`
	// StartDate is the visibility start date (YYYY-MM-DD).
	StartDate string `json:"startDate,omitempty"`
	// EndDate is the visibility end date (YYYY-MM-DD).
	EndDate string `json:"endDate,omitempty"`
}

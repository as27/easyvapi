package model

// FormerMemberData represents the archived data of a former member in easyVerein
// (Ehemalige Mitgliederdaten). This resource is read-only.
type FormerMemberData struct {
	// ID is the unique identifier of the former member record.
	ID int `json:"id"`
	// MembershipNumber is the former member's membership number.
	MembershipNumber string `json:"membershipNumber"`
	// JoinDate is when the former member joined (YYYY-MM-DD format).
	JoinDate string `json:"joinDate"`
	// ResignationDate is when the former member resigned (YYYY-MM-DD format).
	ResignationDate string `json:"resignationDate"`
	// FirstName is the former member's first name.
	FirstName string `json:"firstName"`
	// FamilyName is the former member's family name.
	FamilyName string `json:"familyName"`
	// Email is the former member's email address.
	Email string `json:"email"`
}

// Package model contains the data types returned and accepted by the
// easyVerein API v2.0.
package model

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// flexString unmarshals a JSON string or number into a Go string.
// This handles API fields that return 0 (integer) instead of null or "".
type flexString string

func (f *flexString) UnmarshalJSON(data []byte) error {
	// null → empty string
	if string(data) == "null" {
		*f = ""
		return nil
	}
	// try string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = flexString(s)
		return nil
	}
	// fall back to number → convert to string
	var n json.Number
	if err := json.Unmarshal(data, &n); err != nil {
		return fmt.Errorf("flexString: %w", err)
	}
	*f = flexString(n.String())
	return nil
}

// flexFloat64 unmarshals a JSON number or a quoted number string into float64.
type flexFloat64 float64

func (f *flexFloat64) UnmarshalJSON(data []byte) error {
	var n float64
	if err := json.Unmarshal(data, &n); err == nil {
		*f = flexFloat64(n)
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("flexFloat64: %w", err)
	}
	if s == "" {
		*f = 0
		return nil
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("flexFloat64: %w", err)
	}
	*f = flexFloat64(n)
	return nil
}

// ContactDetails contains the contact information and address details of a person.
// This struct is nested within Member and ContactDetailsService records.
// Date fields (DateOfBirth) are in YYYY-MM-DD format. Email and phone fields
// may be empty strings if not provided.
type ContactDetails struct {
	// ID is the unique identifier for this contact record.
	ID int `json:"id"`
	// FirstName is the person's first/given name.
	FirstName string `json:"firstName"`
	// FamilyName is the person's family/last name.
	FamilyName string `json:"familyName"`
	// Salutation is the formal greeting (e.g., "Mr.", "Ms.", "Dr.").
	Salutation string `json:"salutation"`
	// Street is the street address including number.
	Street string `json:"street"`
	// Zip is the postal code.
	Zip string `json:"zip"`
	// City is the city name.
	City string `json:"city"`
	// Country is the ISO 3166-1 alpha-2 country code (e.g., "DE", "CH", "AT").
	Country string `json:"country"`
	// PrivateEmail is a private email address.
	PrivateEmail string `json:"privateEmail"`
	// PrimaryEmail is the main email address for communication.
	PrimaryEmail string `json:"primaryEmail"`
	// PrivatePhone is a private phone number.
	PrivatePhone string `json:"privatePhone"`
	// MobilePhone is a mobile/cell phone number.
	MobilePhone string `json:"mobilePhone"`
	// DateOfBirth is the person's date of birth in YYYY-MM-DD format.
	DateOfBirth string `json:"dateOfBirth"`
}

// UnmarshalJSON handles the case where the API returns either a URL string or a
// full ContactDetails object for nested contact-details fields.
func (cd *ContactDetails) UnmarshalJSON(data []byte) error {
	// Try a string (URL reference) first.
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		// Only the URL is available; leave all fields at zero values.
		return nil
	}
	// Fall back to a plain struct decode using an alias to avoid recursion.
	type contactDetailsAlias ContactDetails
	var alias contactDetailsAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("model: unmarshal ContactDetails: %w", err)
	}
	*cd = ContactDetails(alias)
	return nil
}

// MemberGroup represents a membership category/group (like "Active Members", "Sponsors", etc).
// Members can belong to multiple groups through the memberGroups field.
type MemberGroup struct {
	// ID is the unique identifier for this member group.
	ID int `json:"id"`
	// Name is the full name of the member group.
	Name string `json:"name"`
	// Short is an abbreviation or short code for the group.
	Short string `json:"short"`
	// Description provides additional context about the group.
	Description string `json:"description"`
}

// UnmarshalJSON handles the case where the API returns either a URL string or a
// full MemberGroup object for nested member-group fields.
func (g *MemberGroup) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		return nil
	}
	type memberGroupAlias MemberGroup
	var alias memberGroupAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("model: unmarshal MemberGroup: %w", err)
	}
	*g = MemberGroup(alias)
	return nil
}

// MemberGroupCreate contains the fields for creating or updating a member group
// via POST / PATCH /member-group.
type MemberGroupCreate struct {
	Name        string `json:"name"`
	Short       string `json:"short,omitempty"`
	Description string `json:"description,omitempty"`
}

// Member represents a member record in easyVerein.
// A member has personal information, contact details, group memberships,
// and payment details. Fields prefixed with underscore in the API are computed
// or read-only fields. Date fields are in RFC3339 format as returned by the API.
type Member struct {
	// ID is the unique identifier for this member.
	ID int `json:"id"`
	// MembershipNumber is the human-readable membership identifier (e.g., "M-12345").
	MembershipNumber string `json:"membershipNumber"`
	// JoinDate is when the member joined (RFC3339 / ISO 8601 format).
	JoinDate string `json:"joinDate"`
	// ResignationDate is when the member resigned, or empty if still active.
	ResignationDate string `json:"resignationDate"`
	// ResignationNoticeDate is the date the resignation notice was received.
	ResignationNoticeDate string `json:"resignationNoticeDate"`
	// DeclarationOfApplication is a URL to the application declaration document.
	DeclarationOfApplication string `json:"declarationOfApplication"`
	// PaymentStartDate is the date from which the payment obligation starts.
	PaymentStartDate string `json:"_paymentStartDate"`
	// PaymentAmount is the recurring payment amount in the local currency.
	PaymentAmount float64 `json:"paymentAmount"`
	// PaymentIntervallMonths is how often the member pays (e.g., 12 for annually, 1 for monthly).
	PaymentIntervallMonths int `json:"paymentIntervallMonths"`
	// UseBalanceForMembershipFee indicates whether the member's balance is used for the fee.
	UseBalanceForMembershipFee bool `json:"useBalanceForMembershipFee"`
	// BulletinBoardNewPostNotification enables email notifications for new bulletin board posts.
	BulletinBoardNewPostNotification bool `json:"bulletinBoardNewPostNotification"`
	// IntegrationDosbSport is the list of DOSB sport codes for this member.
	IntegrationDosbSport []string `json:"integrationDosbSport"`
	// IntegrationDosbGender is the DOSB gender code (e.g. "m", "w", "d").
	IntegrationDosbGender string `json:"integrationDosbGender"`
	// IntegrationLsbSport is the list of LSB sport codes for this member.
	IntegrationLsbSport []string `json:"integrationLsbSport"`
	// IntegrationLsbGender is the LSB gender code (e.g. "m", "w", "d").
	IntegrationLsbGender string `json:"integrationLsbGender"`
	// IsApplication indicates whether this member is still in the application/approval process.
	IsApplication bool `json:"_isApplication"`
	// RelatedMembers is a list of URL references to related members (e.g. family members).
	RelatedMembers []string `json:"relatedMembers"`
	// Org is the URL reference to the organization this member belongs to.
	Org string `json:"org"`
	// DeleteAfterDate is the date after which this record will be auto-deleted.
	DeleteAfterDate string `json:"_deleteAfterDate"`
	// DeletedBy is the identifier of who deleted this record.
	DeletedBy string `json:"_deletedBy"`
	// DeclarationOfResignation is a URL to the resignation declaration document.
	DeclarationOfResignation string `json:"declarationOfResignation"`
	// DeclarationOfConsent is a URL to the consent declaration document.
	DeclarationOfConsent string `json:"declarationOfConsent"`
	// SepaMandateFile is a URL to the SEPA mandate file.
	SepaMandateFile string `json:"sepaMandateFile"`
	// MemberGroups is the list of groups this member belongs to.
	MemberGroups []MemberGroup `json:"memberGroups"`
	// CustomFields is a list of URL references to the member's custom field values.
	CustomFields []string `json:"customFields"`
	// ApplicationDate is the date the membership application was submitted.
	ApplicationDate string `json:"_applicationDate"`
	// ApplicationWasAcceptedAt is the timestamp when the application was accepted.
	ApplicationWasAcceptedAt string `json:"_applicationWasAcceptedAt"`
	// IsChairman indicates whether this member has a chairman/board role.
	IsChairman bool `json:"_isChairman"`
	// ChairmanPermissionGroup is a URL reference to the chairman permission level.
	ChairmanPermissionGroup string `json:"_chairmanPermissionGroup"`
	// ProfilePicture is a URL to the member's profile picture.
	ProfilePicture string `json:"_profilePicture"`
	// ContactDetails contains the member's personal and address information.
	ContactDetails *ContactDetails `json:"contactDetails"`
	// EmailOrUserName is the login email or username for the member's account.
	EmailOrUserName string `json:"emailOrUserName"`
	// SignatureText is the member's email signature text.
	SignatureText string `json:"signatureText"`
	// EditableByRelatedMembers indicates whether related members can edit this profile.
	EditableByRelatedMembers bool `json:"_editableByRelatedMembers"`
	// RequirePasswordChange indicates whether the member must change their password on next login.
	RequirePasswordChange bool `json:"requirePasswordChange"`
	// IsBlocked indicates whether the member account is blocked/suspended.
	IsBlocked bool `json:"_isBlocked"`
	// BlockReason is the reason the member was blocked.
	BlockReason string `json:"blockReason"`
	// ApplicationKind is the type of membership application.
	ApplicationKind string `json:"applicationKind"`
	// WantsToCancelAt is the date the member requested cancellation.
	WantsToCancelAt string `json:"wantsToCancelAt"`
	// CancelReason is the reason the member provided for cancellation.
	CancelReason string `json:"cancelReason"`
	// ShowWarningsAndNotesToAdminsInProfile shows admin notes in the member profile.
	ShowWarningsAndNotesToAdminsInProfile bool `json:"showWarningsAndNotesToAdminsInProfile"`
	// ApplicationForm is a URL reference to the application form used.
	ApplicationForm string `json:"applicationForm"`
	// IsMatrixSearchable indicates whether the member is findable in the Matrix chat directory.
	IsMatrixSearchable bool `json:"_isMatrixSearchable"`
	// MatrixBlockReason is the reason the member was blocked from Matrix.
	MatrixBlockReason string `json:"matrixBlockReason"`
	// BlockedFromMatrix indicates whether the member is blocked from the Matrix chat.
	BlockedFromMatrix bool `json:"blockedFromMatrix"`
	// MatrixCommunicationPermission is the Matrix communication permission level.
	MatrixCommunicationPermission int `json:"_matrixCommunicationPermission"`
	// UseMatrixGroupSettings indicates whether to use global Matrix group settings.
	UseMatrixGroupSettings bool `json:"useMatrixGroupSettings"`
	// RelatedMember may be returned as an integer ID or as a nested Member object,
	// typically representing a related family member or sponsor.
	RelatedMember *Member `json:"-"`

	relatedMemberRaw json.RawMessage `json:"-"`
}

// memberJSON is used internally to decode Member without triggering the custom
// UnmarshalJSON recursively.
type memberJSON struct {
	ID                                    int             `json:"id"`
	MembershipNumber                      string          `json:"membershipNumber"`
	JoinDate                              string          `json:"joinDate"`
	ResignationDate                       string          `json:"resignationDate"`
	ResignationNoticeDate                 string          `json:"resignationNoticeDate"`
	DeclarationOfApplication              string          `json:"declarationOfApplication"`
	PaymentStartDate                      string          `json:"_paymentStartDate"`
	PaymentAmount                         flexFloat64     `json:"paymentAmount"`
	PaymentIntervallMonths                int             `json:"paymentIntervallMonths"`
	UseBalanceForMembershipFee            bool            `json:"useBalanceForMembershipFee"`
	BulletinBoardNewPostNotification      bool            `json:"bulletinBoardNewPostNotification"`
	IntegrationDosbSport                  []string        `json:"integrationDosbSport"`
	IntegrationDosbGender                 string          `json:"integrationDosbGender"`
	IntegrationLsbSport                   []string        `json:"integrationLsbSport"`
	IntegrationLsbGender                  string          `json:"integrationLsbGender"`
	IsApplication                         bool            `json:"_isApplication"`
	RelatedMembers                        []string        `json:"relatedMembers"`
	Org                                   string          `json:"org"`
	DeleteAfterDate                       string          `json:"_deleteAfterDate"`
	DeletedBy                             string          `json:"_deletedBy"`
	DeclarationOfResignation              string          `json:"declarationOfResignation"`
	DeclarationOfConsent                  string          `json:"declarationOfConsent"`
	SepaMandateFile                       string          `json:"sepaMandateFile"`
	MemberGroups                          []MemberGroup   `json:"memberGroups"`
	CustomFields                          []string        `json:"customFields"`
	ApplicationDate                       string          `json:"_applicationDate"`
	ApplicationWasAcceptedAt              string          `json:"_applicationWasAcceptedAt"`
	IsChairman                            bool            `json:"_isChairman"`
	ChairmanPermissionGroup               string          `json:"_chairmanPermissionGroup"`
	ProfilePicture                        string          `json:"_profilePicture"`
	ContactDetails                        *ContactDetails `json:"contactDetails"`
	EmailOrUserName                       string          `json:"emailOrUserName"`
	SignatureText                         string          `json:"signatureText"`
	EditableByRelatedMembers              bool            `json:"_editableByRelatedMembers"`
	RequirePasswordChange                 bool            `json:"requirePasswordChange"`
	IsBlocked                             bool            `json:"_isBlocked"`
	BlockReason                           string          `json:"blockReason"`
	ApplicationKind                       string          `json:"applicationKind"`
	WantsToCancelAt                       string          `json:"wantsToCancelAt"`
	CancelReason                          string          `json:"cancelReason"`
	ShowWarningsAndNotesToAdminsInProfile bool            `json:"showWarningsAndNotesToAdminsInProfile"`
	ApplicationForm                       string          `json:"applicationForm"`
	IsMatrixSearchable                    bool            `json:"_isMatrixSearchable"`
	MatrixBlockReason                     string          `json:"matrixBlockReason"`
	BlockedFromMatrix                     bool            `json:"blockedFromMatrix"`
	MatrixCommunicationPermission         int             `json:"_matrixCommunicationPermission"`
	UseMatrixGroupSettings                bool            `json:"useMatrixGroupSettings"`
	RelatedMember                         json.RawMessage `json:"_relatedMember"`
}

// UnmarshalJSON handles the _relatedMember field, which may be either an integer
// ID or a full nested Member object.
func (m *Member) UnmarshalJSON(data []byte) error {
	var raw memberJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("model: unmarshal Member: %w", err)
	}
	m.ID = raw.ID
	m.MembershipNumber = raw.MembershipNumber
	m.JoinDate = raw.JoinDate
	m.ResignationDate = raw.ResignationDate
	m.ResignationNoticeDate = raw.ResignationNoticeDate
	m.DeclarationOfApplication = raw.DeclarationOfApplication
	m.PaymentStartDate = raw.PaymentStartDate
	m.PaymentAmount = float64(raw.PaymentAmount)
	m.PaymentIntervallMonths = raw.PaymentIntervallMonths
	m.UseBalanceForMembershipFee = raw.UseBalanceForMembershipFee
	m.BulletinBoardNewPostNotification = raw.BulletinBoardNewPostNotification
	m.IntegrationDosbSport = raw.IntegrationDosbSport
	m.IntegrationDosbGender = raw.IntegrationDosbGender
	m.IntegrationLsbSport = raw.IntegrationLsbSport
	m.IntegrationLsbGender = raw.IntegrationLsbGender
	m.IsApplication = raw.IsApplication
	m.RelatedMembers = raw.RelatedMembers
	m.Org = raw.Org
	m.DeleteAfterDate = raw.DeleteAfterDate
	m.DeletedBy = raw.DeletedBy
	m.DeclarationOfResignation = raw.DeclarationOfResignation
	m.DeclarationOfConsent = raw.DeclarationOfConsent
	m.SepaMandateFile = raw.SepaMandateFile
	m.MemberGroups = raw.MemberGroups
	m.CustomFields = raw.CustomFields
	m.ApplicationDate = raw.ApplicationDate
	m.ApplicationWasAcceptedAt = raw.ApplicationWasAcceptedAt
	m.IsChairman = raw.IsChairman
	m.ChairmanPermissionGroup = raw.ChairmanPermissionGroup
	m.ProfilePicture = raw.ProfilePicture
	m.ContactDetails = raw.ContactDetails
	m.EmailOrUserName = raw.EmailOrUserName
	m.SignatureText = raw.SignatureText
	m.EditableByRelatedMembers = raw.EditableByRelatedMembers
	m.RequirePasswordChange = raw.RequirePasswordChange
	m.IsBlocked = raw.IsBlocked
	m.BlockReason = raw.BlockReason
	m.ApplicationKind = raw.ApplicationKind
	m.WantsToCancelAt = raw.WantsToCancelAt
	m.CancelReason = raw.CancelReason
	m.ShowWarningsAndNotesToAdminsInProfile = raw.ShowWarningsAndNotesToAdminsInProfile
	m.ApplicationForm = raw.ApplicationForm
	m.IsMatrixSearchable = raw.IsMatrixSearchable
	m.MatrixBlockReason = raw.MatrixBlockReason
	m.BlockedFromMatrix = raw.BlockedFromMatrix
	m.MatrixCommunicationPermission = raw.MatrixCommunicationPermission
	m.UseMatrixGroupSettings = raw.UseMatrixGroupSettings

	if len(raw.RelatedMember) > 0 && string(raw.RelatedMember) != "null" {
		// Could be a plain integer ID.
		var id int
		if err := json.Unmarshal(raw.RelatedMember, &id); err == nil {
			m.RelatedMember = &Member{ID: id}
		} else {
			// Could be a URL string reference.
			var s string
			if err := json.Unmarshal(raw.RelatedMember, &s); err == nil {
				// Only reference URL available; leave RelatedMember as zero Member.
				m.RelatedMember = &Member{}
			} else {
				// Try a full object.
				var nested Member
				if err := json.Unmarshal(raw.RelatedMember, &nested); err != nil {
					return fmt.Errorf("model: unmarshal Member._relatedMember: %w", err)
				}
				m.RelatedMember = &nested
			}
		}
	}
	return nil
}

// MemberCreate holds the fields used when creating or updating a member.
// Use this struct with [MemberService.Create] or [MemberService.Update].
// Only non-zero fields are sent to the API, allowing partial updates.
type MemberCreate struct {
	// JoinDate is the date the member joined in YYYY-MM-DD format.
	// Omit to use the current date when creating.
	JoinDate string `json:"joinDate,omitempty"`
	// ResignationDate is the date the member resigned in YYYY-MM-DD format.
	// Set this to resign a member; omit to leave active.
	ResignationDate string `json:"resignationDate,omitempty"`
	// PaymentAmount is the recurring payment amount in the local currency.
	PaymentAmount float64 `json:"paymentAmount,omitempty"`
	// PaymentIntervallMonths is the payment interval in months (e.g., 12 for annual, 1 for monthly).
	PaymentIntervallMonths int `json:"paymentIntervallMonths,omitempty"`
}

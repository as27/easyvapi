// Package model contains the data types returned and accepted by the
// easyVerein API v2.0.
package model

import (
	"encoding/json"
	"fmt"
	"strconv"
)

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
// and payment details. The IsBlocked and IsApplication fields are computed
// by the API and prefixed with underscores in queries.
// Date fields (JoinDate, ResignationDate) are in YYYY-MM-DD format.
type Member struct {
	// ID is the unique identifier for this member.
	ID int `json:"id"`
	// MembershipNumber is the human-readable membership identifier (e.g., "M-12345").
	MembershipNumber string `json:"membershipNumber"`
	// JoinDate is when the member joined (YYYY-MM-DD format).
	JoinDate string `json:"joinDate"`
	// ResignationDate is when the member resigned, or empty if still active (YYYY-MM-DD format).
	ResignationDate string `json:"resignationDate"`
	// PaymentAmount is the recurring payment amount in the local currency.
	PaymentAmount float64 `json:"paymentAmount"`
	// PaymentIntervallMonths is how often the member pays (e.g., 12 for annually, 1 for monthly).
	PaymentIntervallMonths int `json:"paymentIntervallMonths"`
	// IsBlocked indicates whether the member account is blocked/suspended.
	IsBlocked bool `json:"isBlocked"`
	// IsApplication indicates whether this member is still in the application/approval process.
	IsApplication bool `json:"isApplication"`
	// ContactDetails contains the member's personal and address information.
	ContactDetails *ContactDetails `json:"contactDetails"`
	// MemberGroups is the list of groups this member belongs to.
	MemberGroups []MemberGroup `json:"memberGroups"`
	// RelatedMember may be returned as an integer ID or as a nested Member object,
	// typically representing a related family member or sponsor.
	RelatedMember *Member `json:"-"`

	relatedMemberRaw json.RawMessage `json:"-"`
}

// memberJSON is used internally to decode Member without triggering the custom
// UnmarshalJSON recursively.
type memberJSON struct {
	ID                     int             `json:"id"`
	MembershipNumber       string          `json:"membershipNumber"`
	JoinDate               string          `json:"joinDate"`
	ResignationDate        string          `json:"resignationDate"`
	PaymentAmount          flexFloat64     `json:"paymentAmount"`
	PaymentIntervallMonths int             `json:"paymentIntervallMonths"`
	IsBlocked              bool            `json:"isBlocked"`
	IsApplication          bool            `json:"isApplication"`
	ContactDetails         *ContactDetails `json:"contactDetails"`
	MemberGroups           []MemberGroup   `json:"memberGroups"`
	RelatedMember          json.RawMessage `json:"_relatedMember"`
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
	m.PaymentAmount = float64(raw.PaymentAmount)
	m.PaymentIntervallMonths = raw.PaymentIntervallMonths
	m.IsBlocked = raw.IsBlocked
	m.IsApplication = raw.IsApplication
	m.ContactDetails = raw.ContactDetails
	m.MemberGroups = raw.MemberGroups

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

package model

// Organization represents the easyVerein organization/club configuration.
type Organization struct {
	// ID is the unique identifier of the organization.
	ID int `json:"id"`
	// Name is the official name of the organization.
	Name string `json:"name"`
	// ShortName is an abbreviated name.
	ShortName string `json:"shortName"`
	// Email is the contact email of the organization.
	Email string `json:"email"`
	// Phone is the contact phone number.
	Phone string `json:"phone"`
	// Website is the organization's website URL.
	Website string `json:"website"`
	// Street is the street address.
	Street string `json:"street"`
	// Zip is the postal code.
	Zip string `json:"zip"`
	// City is the city name.
	City string `json:"city"`
	// Country is the ISO 3166-1 alpha-2 country code.
	Country string `json:"country"`
	// Logo is the URL to the organization's logo.
	Logo string `json:"logo"`
}

// OrganizationCreate holds the fields for updating an organization via PATCH.
type OrganizationCreate struct {
	// Name is the official name of the organization.
	Name string `json:"name,omitempty"`
	// ShortName is an abbreviated name.
	ShortName string `json:"shortName,omitempty"`
	// Email is the contact email.
	Email string `json:"email,omitempty"`
	// Phone is the contact phone number.
	Phone string `json:"phone,omitempty"`
	// Website is the organization's website URL.
	Website string `json:"website,omitempty"`
	// Street is the street address.
	Street string `json:"street,omitempty"`
	// Zip is the postal code.
	Zip string `json:"zip,omitempty"`
	// City is the city name.
	City string `json:"city,omitempty"`
	// Country is the ISO 3166-1 alpha-2 country code.
	Country string `json:"country,omitempty"`
}

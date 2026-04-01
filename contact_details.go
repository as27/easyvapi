package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/as27/easyvapi/model"
)

// ContactDetailsService manages all CRUD operations on the /contact-details endpoint.
// Use this service to manage contact information and communication details.
type ContactDetailsService struct {
	client *Client
}

// defaultContactDetailsQuery requests all fields defined in model.ContactDetails.
var defaultContactDetailsQuery = NewQuery().
	Fields("id", "firstName", "familyName", "salutation", "street", "zip",
		"city", "country", "privateEmail", "primaryEmail", "privatePhone",
		"mobilePhone", "dateOfBirth")

// ContactDetailsListOptions holds all filter and pagination options for
// ContactDetails list requests.
type ContactDetailsListOptions struct {
	ListOptions
	// Country filters by ISO 3166-1 alpha-2 country code.
	Country string
	// FirstName filters by first name (case-insensitive contains).
	FirstName string
	// FamilyName filters by family/last name (case-insensitive contains).
	FamilyName string
	// IsCompany when non-nil filters by whether the record represents a company.
	IsCompany *bool
}

// contactDetailsListParams converts opts into URL query parameters.
func contactDetailsListParams(opts *ContactDetailsListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultContactDetailsQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultContactDetailsQuery)
	if opts.Country != "" {
		params.Set("country", opts.Country)
	}
	if opts.FirstName != "" {
		params.Set("firstName__icontains", opts.FirstName)
	}
	if opts.FamilyName != "" {
		params.Set("familyName__icontains", opts.FamilyName)
	}
	if opts.IsCompany != nil {
		params.Set("isCompany", strconv.FormatBool(*opts.IsCompany))
	}
	return params
}

// List returns a lazy Iterator over all ContactDetails records matching opts.
// Pages are fetched on-demand as iteration progresses.
// Pass nil for opts to use default pagination.
//
// Example: Find all contacts in Germany named Max
//
//	opts := &easyvapi.ContactDetailsListOptions{
//		Country:   "DE",
//		FirstName: "Max",
//	}
//	iter := client.ContactDetails.List(ctx, opts)
//	for iter.Next() {
//		contact := iter.Value()
//		fmt.Printf("%s %s <%s>\n", contact.FirstName, contact.FamilyName, contact.PrimaryEmail)
//	}
func (s *ContactDetailsService) List(ctx context.Context, opts *ContactDetailsListOptions) *Iterator[model.ContactDetails] {
	startURL := s.client.buildURL("/contact-details", contactDetailsListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.ContactDetails, *string, error) {
		return fetchPage[model.ContactDetails](s.client, ctx, pageURL)
	})
}

// ListAll fetches all ContactDetails records matching opts and returns them as a slice.
// This is a convenience wrapper that collects all pages into memory.
//
// Example: Get all contacts with a phone number
//
//	opts := &easyvapi.ContactDetailsListOptions{
//		ListOptions: easyvapi.ListOptions{Search: "089"},
//	}
//	contacts, err := client.ContactDetails.ListAll(ctx, opts)
func (s *ContactDetailsService) ListAll(ctx context.Context, opts *ContactDetailsListOptions) ([]model.ContactDetails, error) {
	var all []model.ContactDetails
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single ContactDetails record by its ID.
func (s *ContactDetailsService) Get(ctx context.Context, id int, query *Query) (*model.ContactDetails, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/contact-details/%d", id), params)
	if err != nil {
		return nil, err
	}
	var cd model.ContactDetails
	if err := s.client.decodeJSON(resp, &cd); err != nil {
		return nil, err
	}
	return &cd, nil
}

// Create creates a new ContactDetails record and returns it.
func (s *ContactDetailsService) Create(ctx context.Context, cd model.ContactDetails) (*model.ContactDetails, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/contact-details", nil), cd)
	if err != nil {
		return nil, err
	}
	var created model.ContactDetails
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the ContactDetails record with the
// given ID.
func (s *ContactDetailsService) Update(ctx context.Context, id int, cd model.ContactDetails) (*model.ContactDetails, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/contact-details/%d", id), nil), cd)
	if err != nil {
		return nil, err
	}
	var updated model.ContactDetails
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the ContactDetails record with the given ID.
func (s *ContactDetailsService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/contact-details/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

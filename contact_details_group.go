package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// ContactDetailsGroupService manages all CRUD operations on the
// /contact-details-group endpoint.
type ContactDetailsGroupService struct {
	client *Client
}

// defaultContactDetailsGroupQuery is nil because the /contact-details-group
// endpoint does not support field selection via the query parameter.
var defaultContactDetailsGroupQuery *Query = nil

// ContactDetailsGroupListOptions holds all filter and pagination options for
// ContactDetailsGroup list requests.
type ContactDetailsGroupListOptions struct {
	ListOptions
	// Name filters by group name.
	Name string
}

// contactDetailsGroupListParams converts opts into URL query parameters.
func contactDetailsGroupListParams(opts *ContactDetailsGroupListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultContactDetailsGroupQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultContactDetailsGroupQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	return params
}

// List returns a lazy Iterator over all ContactDetailsGroup records matching opts.
func (s *ContactDetailsGroupService) List(ctx context.Context, opts *ContactDetailsGroupListOptions) *Iterator[model.ContactDetailsGroup] {
	startURL := s.client.buildURL("/contact-details-group", contactDetailsGroupListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.ContactDetailsGroup, *string, error) {
		return fetchPage[model.ContactDetailsGroup](s.client, ctx, pageURL)
	})
}

// ListAll fetches all ContactDetailsGroup records matching opts and returns them as a slice.
func (s *ContactDetailsGroupService) ListAll(ctx context.Context, opts *ContactDetailsGroupListOptions) ([]model.ContactDetailsGroup, error) {
	var all []model.ContactDetailsGroup
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single ContactDetailsGroup by its ID.
func (s *ContactDetailsGroupService) Get(ctx context.Context, id int, query *Query) (*model.ContactDetailsGroup, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/contact-details-group/%d", id), params)
	if err != nil {
		return nil, err
	}
	var g model.ContactDetailsGroup
	if err := s.client.decodeJSON(resp, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

// Create creates a new ContactDetailsGroup and returns the created record.
func (s *ContactDetailsGroupService) Create(ctx context.Context, g model.ContactDetailsGroupCreate) (*model.ContactDetailsGroup, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/contact-details-group", nil), g)
	if err != nil {
		return nil, err
	}
	var created model.ContactDetailsGroup
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the ContactDetailsGroup with the given ID.
func (s *ContactDetailsGroupService) Update(ctx context.Context, id int, g model.ContactDetailsGroupCreate) (*model.ContactDetailsGroup, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/contact-details-group/%d", id), nil), g)
	if err != nil {
		return nil, err
	}
	var updated model.ContactDetailsGroup
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the ContactDetailsGroup with the given ID.
func (s *ContactDetailsGroupService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/contact-details-group/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

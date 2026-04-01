package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// LocationService manages all CRUD operations on the /location endpoint.
// Locations represent physical venues that can be linked to events.
type LocationService struct {
	client *Client
}

// defaultLocationQuery is nil because the /location endpoint does not support
// field selection via the query parameter.
var defaultLocationQuery *Query = nil

// LocationListOptions holds all filter and pagination options for Location
// list requests.
type LocationListOptions struct {
	ListOptions
	// Name filters by location name.
	Name string
	// Country filters by ISO 3166-1 alpha-2 country code.
	Country string
	// Zip filters by postal code.
	Zip string
}

// locationListParams converts opts into URL query parameters.
func locationListParams(opts *LocationListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultLocationQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultLocationQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	if opts.Country != "" {
		params.Set("country", opts.Country)
	}
	if opts.Zip != "" {
		params.Set("zip", opts.Zip)
	}
	return params
}

// List returns a lazy Iterator over all Location records matching opts.
func (s *LocationService) List(ctx context.Context, opts *LocationListOptions) *Iterator[model.Location] {
	startURL := s.client.buildURL("/location", locationListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.Location, *string, error) {
		return fetchPage[model.Location](s.client, ctx, pageURL)
	})
}

// ListAll fetches all Location records matching opts and returns them as a slice.
func (s *LocationService) ListAll(ctx context.Context, opts *LocationListOptions) ([]model.Location, error) {
	var all []model.Location
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single Location by its ID.
func (s *LocationService) Get(ctx context.Context, id int, query *Query) (*model.Location, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/location/%d", id), params)
	if err != nil {
		return nil, err
	}
	var loc model.Location
	if err := s.client.decodeJSON(resp, &loc); err != nil {
		return nil, err
	}
	return &loc, nil
}

// Create creates a new Location and returns the created record.
func (s *LocationService) Create(ctx context.Context, loc model.LocationCreate) (*model.Location, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/location", nil), loc)
	if err != nil {
		return nil, err
	}
	var created model.Location
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the Location with the given ID.
func (s *LocationService) Update(ctx context.Context, id int, loc model.LocationCreate) (*model.Location, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/location/%d", id), nil), loc)
	if err != nil {
		return nil, err
	}
	var updated model.Location
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the Location with the given ID.
func (s *LocationService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/location/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/as27/easyvapi/model"
)

// ProtocolService manages all CRUD operations on the /protocol endpoint.
// Protocols document meetings and are linked to events.
type ProtocolService struct {
	client *Client
}

// defaultProtocolQuery is nil because the /protocol endpoint
// does not support field selection via the query parameter.
var defaultProtocolQuery *Query = nil

// ProtocolListOptions holds all filter and pagination options for
// Protocol list requests.
type ProtocolListOptions struct {
	ListOptions
	// Name filters by protocol name.
	Name string
	// LocationName filters by location name.
	LocationName string
	// StartGte filters protocols starting on or after this date.
	StartGte string
	// StartLte filters protocols starting on or before this date.
	StartLte string
	// EndGte filters protocols ending on or after this date.
	EndGte string
	// EndLte filters protocols ending on or before this date.
	EndLte string
	// IsPublic when non-nil filters by public visibility.
	IsPublic *bool
	// AllDay when non-nil filters by all-day flag.
	AllDay *bool
	// Calendar filters by calendar ID.
	Calendar *int
	// LocationObject filters by location ID.
	LocationObject *int
}

// protocolListParams converts opts into URL query parameters.
func protocolListParams(opts *ProtocolListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultProtocolQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultProtocolQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	if opts.LocationName != "" {
		params.Set("locationName", opts.LocationName)
	}
	if opts.StartGte != "" {
		params.Set("start__gte", opts.StartGte)
	}
	if opts.StartLte != "" {
		params.Set("start__lte", opts.StartLte)
	}
	if opts.EndGte != "" {
		params.Set("end__gte", opts.EndGte)
	}
	if opts.EndLte != "" {
		params.Set("end__lte", opts.EndLte)
	}
	if opts.IsPublic != nil {
		params.Set("isPublic", strconv.FormatBool(*opts.IsPublic))
	}
	if opts.AllDay != nil {
		params.Set("allDay", strconv.FormatBool(*opts.AllDay))
	}
	if opts.Calendar != nil {
		params.Set("calendar", strconv.Itoa(*opts.Calendar))
	}
	if opts.LocationObject != nil {
		params.Set("locationObject", strconv.Itoa(*opts.LocationObject))
	}
	return params
}

// List returns a lazy Iterator over all Protocol records matching opts.
func (s *ProtocolService) List(ctx context.Context, opts *ProtocolListOptions) *Iterator[model.Protocol] {
	startURL := s.client.buildURL("/protocol", protocolListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.Protocol, *string, error) {
		return fetchPage[model.Protocol](s.client, ctx, pageURL)
	})
}

// ListAll fetches all Protocol records matching opts and returns them as a slice.
func (s *ProtocolService) ListAll(ctx context.Context, opts *ProtocolListOptions) ([]model.Protocol, error) {
	var all []model.Protocol
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single Protocol by its ID.
func (s *ProtocolService) Get(ctx context.Context, id int, query *Query) (*model.Protocol, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/protocol/%d", id), params)
	if err != nil {
		return nil, err
	}
	var p model.Protocol
	if err := s.client.decodeJSON(resp, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Create creates a new Protocol and returns the created record.
func (s *ProtocolService) Create(ctx context.Context, p model.ProtocolCreate) (*model.Protocol, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/protocol", nil), p)
	if err != nil {
		return nil, err
	}
	var created model.Protocol
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the Protocol with the given ID.
func (s *ProtocolService) Update(ctx context.Context, id int, p model.ProtocolCreate) (*model.Protocol, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/protocol/%d", id), nil), p)
	if err != nil {
		return nil, err
	}
	var updated model.Protocol
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the Protocol with the given ID.
func (s *ProtocolService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/protocol/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// CustomFilterService manages all CRUD operations on the /custom-filter endpoint.
// Custom filters are saved filter configurations for member or contact lists.
type CustomFilterService struct {
	client *Client
}

// defaultCustomFilterQuery is nil because the /custom-filter endpoint does
// not support field selection via the query parameter.
var defaultCustomFilterQuery *Query = nil

// CustomFilterListOptions holds all filter and pagination options for
// CustomFilter list requests.
type CustomFilterListOptions struct {
	ListOptions
	// Name filters by filter name.
	Name string
	// Model filters by the resource type this filter applies to.
	Model string
}

// customFilterListParams converts opts into URL query parameters.
func customFilterListParams(opts *CustomFilterListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultCustomFilterQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultCustomFilterQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	if opts.Model != "" {
		params.Set("model", opts.Model)
	}
	return params
}

// List returns a lazy Iterator over all CustomFilter records matching opts.
func (s *CustomFilterService) List(ctx context.Context, opts *CustomFilterListOptions) *Iterator[model.CustomFilter] {
	startURL := s.client.buildURL("/custom-filter", customFilterListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.CustomFilter, *string, error) {
		return fetchPage[model.CustomFilter](s.client, ctx, pageURL)
	})
}

// ListAll fetches all CustomFilter records matching opts and returns them as a slice.
func (s *CustomFilterService) ListAll(ctx context.Context, opts *CustomFilterListOptions) ([]model.CustomFilter, error) {
	var all []model.CustomFilter
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single CustomFilter by its ID.
func (s *CustomFilterService) Get(ctx context.Context, id int, query *Query) (*model.CustomFilter, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/custom-filter/%d", id), params)
	if err != nil {
		return nil, err
	}
	var f model.CustomFilter
	if err := s.client.decodeJSON(resp, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// Create creates a new CustomFilter and returns the created record.
func (s *CustomFilterService) Create(ctx context.Context, f model.CustomFilterCreate) (*model.CustomFilter, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/custom-filter", nil), f)
	if err != nil {
		return nil, err
	}
	var created model.CustomFilter
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the CustomFilter with the given ID.
func (s *CustomFilterService) Update(ctx context.Context, id int, f model.CustomFilterCreate) (*model.CustomFilter, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/custom-filter/%d", id), nil), f)
	if err != nil {
		return nil, err
	}
	var updated model.CustomFilter
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the CustomFilter with the given ID.
func (s *CustomFilterService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/custom-filter/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

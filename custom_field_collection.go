package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// CustomFieldCollectionService manages all CRUD operations on the
// /custom-field-collection endpoint. Collections group custom fields together
// for display purposes.
type CustomFieldCollectionService struct {
	client *Client
}

// defaultCustomFieldCollectionQuery requests all fields defined in
// model.CustomFieldCollection.
var defaultCustomFieldCollectionQuery = NewQuery().
	Fields("id", "name", "orderSequence", "position")

// CustomFieldCollectionListOptions holds all filter and pagination options for
// CustomFieldCollection list requests.
type CustomFieldCollectionListOptions struct {
	ListOptions
}

// customFieldCollectionListParams converts opts into URL query parameters.
func customFieldCollectionListParams(opts *CustomFieldCollectionListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultCustomFieldCollectionQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultCustomFieldCollectionQuery)
	return params
}

// List returns a lazy Iterator over all CustomFieldCollection records.
func (s *CustomFieldCollectionService) List(ctx context.Context, opts *CustomFieldCollectionListOptions) *Iterator[model.CustomFieldCollection] {
	startURL := s.client.buildURL("/custom-field-collection", customFieldCollectionListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.CustomFieldCollection, *string, error) {
		return fetchPage[model.CustomFieldCollection](s.client, ctx, pageURL)
	})
}

// ListAll fetches all CustomFieldCollection records and returns them as a slice.
func (s *CustomFieldCollectionService) ListAll(ctx context.Context, opts *CustomFieldCollectionListOptions) ([]model.CustomFieldCollection, error) {
	var all []model.CustomFieldCollection
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single CustomFieldCollection by its ID.
func (s *CustomFieldCollectionService) Get(ctx context.Context, id int, query *Query) (*model.CustomFieldCollection, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/custom-field-collection/%d", id), params)
	if err != nil {
		return nil, err
	}
	var c model.CustomFieldCollection
	if err := s.client.decodeJSON(resp, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// Create creates a new CustomFieldCollection and returns the created record.
func (s *CustomFieldCollectionService) Create(ctx context.Context, c model.CustomFieldCollectionCreate) (*model.CustomFieldCollection, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/custom-field-collection", nil), c)
	if err != nil {
		return nil, err
	}
	var created model.CustomFieldCollection
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the CustomFieldCollection with the given ID.
func (s *CustomFieldCollectionService) Update(ctx context.Context, id int, c model.CustomFieldCollectionCreate) (*model.CustomFieldCollection, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/custom-field-collection/%d", id), nil), c)
	if err != nil {
		return nil, err
	}
	var updated model.CustomFieldCollection
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the CustomFieldCollection with the given ID.
func (s *CustomFieldCollectionService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/custom-field-collection/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

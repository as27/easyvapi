package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// InventoryObjectGroupService manages all CRUD operations on the
// /inventory-object-group endpoint.
type InventoryObjectGroupService struct {
	client *Client
}

// defaultInventoryObjectGroupQuery is nil because the /inventory-object-group
// endpoint does not support field selection via the query parameter.
var defaultInventoryObjectGroupQuery *Query = nil

// InventoryObjectGroupListOptions holds all filter and pagination options for
// InventoryObjectGroup list requests.
type InventoryObjectGroupListOptions struct {
	ListOptions
	// Name filters by group name.
	Name string
}

// inventoryObjectGroupListParams converts opts into URL query parameters.
func inventoryObjectGroupListParams(opts *InventoryObjectGroupListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultInventoryObjectGroupQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultInventoryObjectGroupQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	return params
}

// List returns a lazy Iterator over all InventoryObjectGroup records matching opts.
func (s *InventoryObjectGroupService) List(ctx context.Context, opts *InventoryObjectGroupListOptions) *Iterator[model.InventoryObjectGroup] {
	startURL := s.client.buildURL("/inventory-object-group", inventoryObjectGroupListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.InventoryObjectGroup, *string, error) {
		return fetchPage[model.InventoryObjectGroup](s.client, ctx, pageURL)
	})
}

// ListAll fetches all InventoryObjectGroup records and returns them as a slice.
func (s *InventoryObjectGroupService) ListAll(ctx context.Context, opts *InventoryObjectGroupListOptions) ([]model.InventoryObjectGroup, error) {
	var all []model.InventoryObjectGroup
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single InventoryObjectGroup by its ID.
func (s *InventoryObjectGroupService) Get(ctx context.Context, id int, query *Query) (*model.InventoryObjectGroup, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/inventory-object-group/%d", id), params)
	if err != nil {
		return nil, err
	}
	var g model.InventoryObjectGroup
	if err := s.client.decodeJSON(resp, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

// Create creates a new InventoryObjectGroup and returns the created record.
func (s *InventoryObjectGroupService) Create(ctx context.Context, g model.InventoryObjectGroupCreate) (*model.InventoryObjectGroup, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/inventory-object-group", nil), g)
	if err != nil {
		return nil, err
	}
	var created model.InventoryObjectGroup
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the InventoryObjectGroup with the given ID.
func (s *InventoryObjectGroupService) Update(ctx context.Context, id int, g model.InventoryObjectGroupCreate) (*model.InventoryObjectGroup, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/inventory-object-group/%d", id), nil), g)
	if err != nil {
		return nil, err
	}
	var updated model.InventoryObjectGroup
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the InventoryObjectGroup with the given ID.
func (s *InventoryObjectGroupService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/inventory-object-group/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

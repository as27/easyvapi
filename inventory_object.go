package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/as27/easyvapi/model"
)

// InventoryObjectService manages all CRUD operations on the /inventory-object
// endpoint. Inventory objects are physical assets tracked by the organization.
type InventoryObjectService struct {
	client *Client
}

// defaultInventoryObjectQuery is nil because the /inventory-object endpoint
// does not support field selection via the query parameter.
var defaultInventoryObjectQuery *Query = nil

// InventoryObjectListOptions holds all filter and pagination options for
// InventoryObject list requests.
type InventoryObjectListOptions struct {
	ListOptions
	// Name filters by item name.
	Name string
	// Identifier filters by asset tag or serial number.
	Identifier string
	// LendingAvailable when non-nil filters by lending availability.
	LendingAvailable *bool
}

// inventoryObjectListParams converts opts into URL query parameters.
func inventoryObjectListParams(opts *InventoryObjectListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultInventoryObjectQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultInventoryObjectQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	if opts.Identifier != "" {
		params.Set("identifier", opts.Identifier)
	}
	if opts.LendingAvailable != nil {
		params.Set("lendingAvailable", strconv.FormatBool(*opts.LendingAvailable))
	}
	return params
}

// List returns a lazy Iterator over all InventoryObject records matching opts.
func (s *InventoryObjectService) List(ctx context.Context, opts *InventoryObjectListOptions) *Iterator[model.InventoryObject] {
	startURL := s.client.buildURL("/inventory-object", inventoryObjectListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.InventoryObject, *string, error) {
		return fetchPage[model.InventoryObject](s.client, ctx, pageURL)
	})
}

// ListAll fetches all InventoryObject records matching opts and returns them as a slice.
func (s *InventoryObjectService) ListAll(ctx context.Context, opts *InventoryObjectListOptions) ([]model.InventoryObject, error) {
	var all []model.InventoryObject
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single InventoryObject by its ID.
func (s *InventoryObjectService) Get(ctx context.Context, id int, query *Query) (*model.InventoryObject, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/inventory-object/%d", id), params)
	if err != nil {
		return nil, err
	}
	var obj model.InventoryObject
	if err := s.client.decodeJSON(resp, &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// Create creates a new InventoryObject and returns the created record.
func (s *InventoryObjectService) Create(ctx context.Context, obj model.InventoryObjectCreate) (*model.InventoryObject, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/inventory-object", nil), obj)
	if err != nil {
		return nil, err
	}
	var created model.InventoryObject
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the InventoryObject with the given ID.
func (s *InventoryObjectService) Update(ctx context.Context, id int, obj model.InventoryObjectCreate) (*model.InventoryObject, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/inventory-object/%d", id), nil), obj)
	if err != nil {
		return nil, err
	}
	var updated model.InventoryObject
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the InventoryObject with the given ID.
func (s *InventoryObjectService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/inventory-object/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

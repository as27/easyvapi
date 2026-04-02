package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/as27/easyvapi/model"
)

// LendingService manages all CRUD operations on the /lending endpoint,
// including bulk operations. Lending records track borrowed inventory objects.
type LendingService struct {
	client *Client
}

// defaultLendingQuery is nil because the /lending endpoint does not support
// field selection via the query parameter.
var defaultLendingQuery *Query = nil

// LendingListOptions holds all filter and pagination options for Lending
// list requests.
type LendingListOptions struct {
	ListOptions
	// State filters by lending state (e.g. "borrowed", "returned").
	State string
	// InventoryObject filters by inventory object ID.
	InventoryObject int
	// LendingPerson filters by borrowing contact ID.
	LendingPerson int
}

// lendingListParams converts opts into URL query parameters.
func lendingListParams(opts *LendingListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultLendingQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultLendingQuery)
	if opts.State != "" {
		params.Set("state", opts.State)
	}
	if opts.InventoryObject != 0 {
		params.Set("inventoryObject", strconv.Itoa(opts.InventoryObject))
	}
	if opts.LendingPerson != 0 {
		params.Set("lendingPerson", strconv.Itoa(opts.LendingPerson))
	}
	return params
}

// List returns a lazy Iterator over all Lending records matching opts.
func (s *LendingService) List(ctx context.Context, opts *LendingListOptions) *Iterator[model.Lending] {
	startURL := s.client.buildURL("/lending", lendingListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.Lending, *string, error) {
		return fetchPage[model.Lending](s.client, ctx, pageURL)
	})
}

// ListAll fetches all Lending records matching opts and returns them as a slice.
func (s *LendingService) ListAll(ctx context.Context, opts *LendingListOptions) ([]model.Lending, error) {
	var all []model.Lending
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single Lending record by its ID.
func (s *LendingService) Get(ctx context.Context, id int, query *Query) (*model.Lending, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/lending/%d", id), params)
	if err != nil {
		return nil, err
	}
	var l model.Lending
	if err := s.client.decodeJSON(resp, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// Create creates a new Lending record and returns the created record.
func (s *LendingService) Create(ctx context.Context, l model.LendingCreate) (*model.Lending, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/lending", nil), l)
	if err != nil {
		return nil, err
	}
	var created model.Lending
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// BulkCreate creates multiple Lending records in a single request.
func (s *LendingService) BulkCreate(ctx context.Context, entries []model.LendingCreate) ([]model.Lending, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/lending/bulk-create", nil), entries)
	if err != nil {
		return nil, err
	}
	var created []model.Lending
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return created, nil
}

// BulkUpdate applies a partial update (PATCH) to multiple Lending records.
func (s *LendingService) BulkUpdate(ctx context.Context, entries []model.Lending) ([]model.Lending, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL("/lending/bulk-update", nil), entries)
	if err != nil {
		return nil, err
	}
	var updated []model.Lending
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return updated, nil
}

// Update applies a partial update (PATCH) to the Lending record with the given ID.
func (s *LendingService) Update(ctx context.Context, id int, l model.LendingCreate) (*model.Lending, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/lending/%d", id), nil), l)
	if err != nil {
		return nil, err
	}
	var updated model.Lending
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the Lending record with the given ID.
func (s *LendingService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/lending/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

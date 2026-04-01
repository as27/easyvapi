package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// ChairmanLevelService manages all CRUD operations on the /chairman-level
// endpoint. Chairman levels define access permissions for board members.
type ChairmanLevelService struct {
	client *Client
}

// defaultChairmanLevelQuery is nil because the /chairman-level endpoint does
// not support field selection via the query parameter.
var defaultChairmanLevelQuery *Query = nil

// ChairmanLevelListOptions holds all filter and pagination options for
// ChairmanLevel list requests.
type ChairmanLevelListOptions struct {
	ListOptions
}

// chairmanLevelListParams converts opts into URL query parameters.
func chairmanLevelListParams(opts *ChairmanLevelListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultChairmanLevelQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultChairmanLevelQuery)
	return params
}

// List returns a lazy Iterator over all ChairmanLevel records.
func (s *ChairmanLevelService) List(ctx context.Context, opts *ChairmanLevelListOptions) *Iterator[model.ChairmanLevel] {
	startURL := s.client.buildURL("/chairman-level", chairmanLevelListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.ChairmanLevel, *string, error) {
		return fetchPage[model.ChairmanLevel](s.client, ctx, pageURL)
	})
}

// ListAll fetches all ChairmanLevel records and returns them as a slice.
func (s *ChairmanLevelService) ListAll(ctx context.Context, opts *ChairmanLevelListOptions) ([]model.ChairmanLevel, error) {
	var all []model.ChairmanLevel
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single ChairmanLevel by its ID.
func (s *ChairmanLevelService) Get(ctx context.Context, id int, query *Query) (*model.ChairmanLevel, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/chairman-level/%d", id), params)
	if err != nil {
		return nil, err
	}
	var lvl model.ChairmanLevel
	if err := s.client.decodeJSON(resp, &lvl); err != nil {
		return nil, err
	}
	return &lvl, nil
}

// Create creates a new ChairmanLevel and returns the created record.
func (s *ChairmanLevelService) Create(ctx context.Context, lvl model.ChairmanLevelCreate) (*model.ChairmanLevel, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/chairman-level", nil), lvl)
	if err != nil {
		return nil, err
	}
	var created model.ChairmanLevel
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the ChairmanLevel with the given ID.
func (s *ChairmanLevelService) Update(ctx context.Context, id int, lvl model.ChairmanLevelCreate) (*model.ChairmanLevel, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/chairman-level/%d", id), nil), lvl)
	if err != nil {
		return nil, err
	}
	var updated model.ChairmanLevel
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the ChairmanLevel with the given ID.
func (s *ChairmanLevelService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/chairman-level/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

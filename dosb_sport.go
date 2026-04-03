package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// DosbSportService manages operations on the /dosb-sport endpoint.
// DOSB stands for Deutscher Olympischer Sportbund.
type DosbSportService struct {
	client *Client
}

// defaultDosbSportQuery is nil because the endpoint does not support
// field selection via the query parameter.
var defaultDosbSportQuery *Query = nil

// DosbSportListOptions holds filter and pagination options for list requests.
type DosbSportListOptions struct {
	ListOptions
	// Name filters by sport name.
	Name string
}

func dosbSportListParams(opts *DosbSportListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultDosbSportQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultDosbSportQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	return params
}

// List returns a lazy Iterator over all DosbSport records matching opts.
func (s *DosbSportService) List(ctx context.Context, opts *DosbSportListOptions) *Iterator[model.DosbSport] {
	startURL := s.client.buildURL("/dosb-sport", dosbSportListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.DosbSport, *string, error) {
		return fetchPage[model.DosbSport](s.client, ctx, pageURL)
	})
}

// ListAll fetches all DosbSport records matching opts and returns them as a slice.
func (s *DosbSportService) ListAll(ctx context.Context, opts *DosbSportListOptions) ([]model.DosbSport, error) {
	var all []model.DosbSport
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single DosbSport by its ID.
func (s *DosbSportService) Get(ctx context.Context, id int, query *Query) (*model.DosbSport, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/dosb-sport/%d", id), params)
	if err != nil {
		return nil, err
	}
	var sport model.DosbSport
	if err := s.client.decodeJSON(resp, &sport); err != nil {
		return nil, err
	}
	return &sport, nil
}

// Create creates a new DosbSport entry and returns the created record.
func (s *DosbSportService) Create(ctx context.Context, sport model.DosbSportCreate) (*model.DosbSport, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/dosb-sport", nil), sport)
	if err != nil {
		return nil, err
	}
	var created model.DosbSport
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the DosbSport with the given ID.
func (s *DosbSportService) Update(ctx context.Context, id int, sport model.DosbSportCreate) (*model.DosbSport, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/dosb-sport/%d", id), nil), sport)
	if err != nil {
		return nil, err
	}
	var updated model.DosbSport
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// LsbSportService manages operations on the /lsb-sport endpoint.
// LSB stands for Landessportbund.
type LsbSportService struct {
	client *Client
}

// defaultLsbSportQuery is nil because the endpoint does not support
// field selection via the query parameter.
var defaultLsbSportQuery *Query = nil

// LsbSportListOptions holds filter and pagination options for list requests.
type LsbSportListOptions struct {
	ListOptions
	// Name filters by sport name.
	Name string
}

func lsbSportListParams(opts *LsbSportListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultLsbSportQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultLsbSportQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	return params
}

// List returns a lazy Iterator over all LsbSport records matching opts.
func (s *LsbSportService) List(ctx context.Context, opts *LsbSportListOptions) *Iterator[model.LsbSport] {
	startURL := s.client.buildURL("/lsb-sport", lsbSportListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.LsbSport, *string, error) {
		return fetchPage[model.LsbSport](s.client, ctx, pageURL)
	})
}

// ListAll fetches all LsbSport records matching opts and returns them as a slice.
func (s *LsbSportService) ListAll(ctx context.Context, opts *LsbSportListOptions) ([]model.LsbSport, error) {
	var all []model.LsbSport
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single LsbSport by its ID.
func (s *LsbSportService) Get(ctx context.Context, id int, query *Query) (*model.LsbSport, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/lsb-sport/%d", id), params)
	if err != nil {
		return nil, err
	}
	var sport model.LsbSport
	if err := s.client.decodeJSON(resp, &sport); err != nil {
		return nil, err
	}
	return &sport, nil
}

// Create creates a new LsbSport entry and returns the created record.
func (s *LsbSportService) Create(ctx context.Context, sport model.LsbSportCreate) (*model.LsbSport, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/lsb-sport", nil), sport)
	if err != nil {
		return nil, err
	}
	var created model.LsbSport
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the LsbSport with the given ID.
func (s *LsbSportService) Update(ctx context.Context, id int, sport model.LsbSportCreate) (*model.LsbSport, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/lsb-sport/%d", id), nil), sport)
	if err != nil {
		return nil, err
	}
	var updated model.LsbSport
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

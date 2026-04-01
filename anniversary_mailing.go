package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// AnniversaryMailingService manages all CRUD operations on the
// /anniversary-mailing endpoint. Anniversary mailings are automated emails
// sent to members on their membership anniversary or birthday.
type AnniversaryMailingService struct {
	client *Client
}

// defaultAnniversaryMailingQuery is nil because the /anniversary-mailing
// endpoint does not support field selection via the query parameter.
var defaultAnniversaryMailingQuery *Query = nil

// AnniversaryMailingListOptions holds all filter and pagination options for
// AnniversaryMailing list requests.
type AnniversaryMailingListOptions struct {
	ListOptions
}

// anniversaryMailingListParams converts opts into URL query parameters.
func anniversaryMailingListParams(opts *AnniversaryMailingListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultAnniversaryMailingQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultAnniversaryMailingQuery)
	return params
}

// List returns a lazy Iterator over all AnniversaryMailing records.
func (s *AnniversaryMailingService) List(ctx context.Context, opts *AnniversaryMailingListOptions) *Iterator[model.AnniversaryMailing] {
	startURL := s.client.buildURL("/anniversary-mailing", anniversaryMailingListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.AnniversaryMailing, *string, error) {
		return fetchPage[model.AnniversaryMailing](s.client, ctx, pageURL)
	})
}

// ListAll fetches all AnniversaryMailing records and returns them as a slice.
func (s *AnniversaryMailingService) ListAll(ctx context.Context, opts *AnniversaryMailingListOptions) ([]model.AnniversaryMailing, error) {
	var all []model.AnniversaryMailing
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single AnniversaryMailing by its ID.
func (s *AnniversaryMailingService) Get(ctx context.Context, id int, query *Query) (*model.AnniversaryMailing, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/anniversary-mailing/%d", id), params)
	if err != nil {
		return nil, err
	}
	var m model.AnniversaryMailing
	if err := s.client.decodeJSON(resp, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// Create creates a new AnniversaryMailing and returns the created record.
func (s *AnniversaryMailingService) Create(ctx context.Context, m model.AnniversaryMailingCreate) (*model.AnniversaryMailing, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/anniversary-mailing", nil), m)
	if err != nil {
		return nil, err
	}
	var created model.AnniversaryMailing
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the AnniversaryMailing with the given ID.
func (s *AnniversaryMailingService) Update(ctx context.Context, id int, m model.AnniversaryMailingCreate) (*model.AnniversaryMailing, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/anniversary-mailing/%d", id), nil), m)
	if err != nil {
		return nil, err
	}
	var updated model.AnniversaryMailing
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the AnniversaryMailing with the given ID.
func (s *AnniversaryMailingService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/anniversary-mailing/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

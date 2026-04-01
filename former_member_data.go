package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// FormerMemberDataService provides read-only access to the /former-member-data
// endpoint. This resource contains archived data of former members and cannot
// be modified via the API.
type FormerMemberDataService struct {
	client *Client
}

// defaultFormerMemberDataQuery is nil because the /former-member-data endpoint
// does not support field selection via the query parameter.
var defaultFormerMemberDataQuery *Query = nil

// FormerMemberDataListOptions holds all filter and pagination options for
// FormerMemberData list requests.
type FormerMemberDataListOptions struct {
	ListOptions
}

// formerMemberDataListParams converts opts into URL query parameters.
func formerMemberDataListParams(opts *FormerMemberDataListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultFormerMemberDataQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultFormerMemberDataQuery)
	return params
}

// List returns a lazy Iterator over all FormerMemberData records.
// Pages are fetched on-demand as iteration progresses.
func (s *FormerMemberDataService) List(ctx context.Context, opts *FormerMemberDataListOptions) *Iterator[model.FormerMemberData] {
	startURL := s.client.buildURL("/former-member-data", formerMemberDataListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.FormerMemberData, *string, error) {
		return fetchPage[model.FormerMemberData](s.client, ctx, pageURL)
	})
}

// ListAll fetches all FormerMemberData records and returns them as a slice.
func (s *FormerMemberDataService) ListAll(ctx context.Context, opts *FormerMemberDataListOptions) ([]model.FormerMemberData, error) {
	var all []model.FormerMemberData
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single FormerMemberData record by its ID.
func (s *FormerMemberDataService) Get(ctx context.Context, id int, query *Query) (*model.FormerMemberData, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/former-member-data/%d", id), params)
	if err != nil {
		return nil, err
	}
	var fmd model.FormerMemberData
	if err := s.client.decodeJSON(resp, &fmd); err != nil {
		return nil, err
	}
	return &fmd, nil
}

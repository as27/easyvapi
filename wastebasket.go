package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// WastebasketService provides access to the /wastebasket endpoint.
// It allows listing and restoring deleted objects.
type WastebasketService struct {
	client *Client
}

// defaultWastebasketQuery is nil because the endpoint does not support
// field selection via the query parameter.
var defaultWastebasketQuery *Query = nil

// WastebasketListOptions holds filter and pagination options for list requests.
type WastebasketListOptions struct {
	ListOptions
	// Model filters by object type (e.g. "member", "invoice").
	Model string
}

func wastebasketListParams(opts *WastebasketListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultWastebasketQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultWastebasketQuery)
	if opts.Model != "" {
		params.Set("model", opts.Model)
	}
	return params
}

// List returns a lazy Iterator over all WastebasketItem records matching opts.
func (s *WastebasketService) List(ctx context.Context, opts *WastebasketListOptions) *Iterator[model.WastebasketItem] {
	startURL := s.client.buildURL("/wastebasket", wastebasketListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.WastebasketItem, *string, error) {
		return fetchPage[model.WastebasketItem](s.client, ctx, pageURL)
	})
}

// ListAll fetches all WastebasketItem records matching opts and returns them as a slice.
func (s *WastebasketService) ListAll(ctx context.Context, opts *WastebasketListOptions) ([]model.WastebasketItem, error) {
	var all []model.WastebasketItem
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Restore restores the deleted object identified by wastebasket entry ID.
func (s *WastebasketService) Restore(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL(fmt.Sprintf("/wastebasket/%d/restore", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// BillingAccountService manages all CRUD operations on the /billing-account endpoint.
// Billing accounts (Buchungskonten) are used to categorize financial bookings.
type BillingAccountService struct {
	client *Client
}

// defaultBillingAccountQuery is nil because the /billing-account endpoint does
// not support field selection via the query parameter.
var defaultBillingAccountQuery *Query = nil

// BillingAccountListOptions holds all filter and pagination options for BillingAccount
// list requests.
type BillingAccountListOptions struct {
	ListOptions
	// Name filters by the account name.
	Name string
	// AccountKind filters by account kind (e.g. "income", "expense").
	AccountKind string
}

// billingAccountListParams converts opts into URL query parameters.
func billingAccountListParams(opts *BillingAccountListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultBillingAccountQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultBillingAccountQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	if opts.AccountKind != "" {
		params.Set("accountKind", opts.AccountKind)
	}
	return params
}

// List returns a lazy Iterator over all BillingAccount records matching opts.
// Pages are fetched on-demand as iteration progresses.
//
// Example:
//
//	iter := client.BillingAccounts.List(ctx, nil)
//	for iter.Next() {
//		acc := iter.Value()
//		fmt.Printf("%s (%.2f EUR)\n", acc.Name, float64(acc.Balance))
//	}
func (s *BillingAccountService) List(ctx context.Context, opts *BillingAccountListOptions) *Iterator[model.BillingAccount] {
	startURL := s.client.buildURL("/billing-account", billingAccountListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.BillingAccount, *string, error) {
		return fetchPage[model.BillingAccount](s.client, ctx, pageURL)
	})
}

// ListAll fetches all BillingAccount records matching opts and returns them as a slice.
func (s *BillingAccountService) ListAll(ctx context.Context, opts *BillingAccountListOptions) ([]model.BillingAccount, error) {
	var all []model.BillingAccount
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single BillingAccount by its ID.
func (s *BillingAccountService) Get(ctx context.Context, id int, query *Query) (*model.BillingAccount, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/billing-account/%d", id), params)
	if err != nil {
		return nil, err
	}
	var acc model.BillingAccount
	if err := s.client.decodeJSON(resp, &acc); err != nil {
		return nil, err
	}
	return &acc, nil
}

// Create creates a new BillingAccount and returns the created record.
func (s *BillingAccountService) Create(ctx context.Context, acc model.BillingAccountCreate) (*model.BillingAccount, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/billing-account", nil), acc)
	if err != nil {
		return nil, err
	}
	var created model.BillingAccount
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the BillingAccount with the given ID.
func (s *BillingAccountService) Update(ctx context.Context, id int, acc model.BillingAccountCreate) (*model.BillingAccount, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/billing-account/%d", id), nil), acc)
	if err != nil {
		return nil, err
	}
	var updated model.BillingAccount
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the BillingAccount with the given ID.
func (s *BillingAccountService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/billing-account/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

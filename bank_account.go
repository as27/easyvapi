package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// BankAccountService manages all CRUD operations on the /bank-account endpoint.
// Bank accounts (Bankkonten) are used to track financial transactions and SEPA mandates.
type BankAccountService struct {
	client *Client
}

// defaultBankAccountQuery is nil because the /bank-account endpoint does
// not support field selection via the query parameter.
var defaultBankAccountQuery *Query = nil

// BankAccountListOptions holds all filter and pagination options for BankAccount
// list requests.
type BankAccountListOptions struct {
	ListOptions
	// Name filters by the account name.
	Name string
}

// bankAccountListParams converts opts into URL query parameters.
func bankAccountListParams(opts *BankAccountListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultBankAccountQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultBankAccountQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	return params
}

// List returns a lazy Iterator over all BankAccount records matching opts.
// Pages are fetched on-demand as iteration progresses.
//
// Example:
//
//	accounts, err := client.BankAccounts.ListAll(ctx, nil)
//	for _, acc := range accounts {
//		fmt.Printf("%s: %s\n", acc.Name, acc.IBAN)
//	}
func (s *BankAccountService) List(ctx context.Context, opts *BankAccountListOptions) *Iterator[model.BankAccount] {
	startURL := s.client.buildURL("/bank-account", bankAccountListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.BankAccount, *string, error) {
		return fetchPage[model.BankAccount](s.client, ctx, pageURL)
	})
}

// ListAll fetches all BankAccount records matching opts and returns them as a slice.
func (s *BankAccountService) ListAll(ctx context.Context, opts *BankAccountListOptions) ([]model.BankAccount, error) {
	var all []model.BankAccount
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single BankAccount by its ID.
func (s *BankAccountService) Get(ctx context.Context, id int, query *Query) (*model.BankAccount, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/bank-account/%d", id), params)
	if err != nil {
		return nil, err
	}
	var acc model.BankAccount
	if err := s.client.decodeJSON(resp, &acc); err != nil {
		return nil, err
	}
	return &acc, nil
}

// Create creates a new BankAccount and returns the created record.
func (s *BankAccountService) Create(ctx context.Context, acc model.BankAccountCreate) (*model.BankAccount, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/bank-account", nil), acc)
	if err != nil {
		return nil, err
	}
	var created model.BankAccount
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the BankAccount with the given ID.
func (s *BankAccountService) Update(ctx context.Context, id int, acc model.BankAccountCreate) (*model.BankAccount, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/bank-account/%d", id), nil), acc)
	if err != nil {
		return nil, err
	}
	var updated model.BankAccount
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the BankAccount with the given ID.
func (s *BankAccountService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/bank-account/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

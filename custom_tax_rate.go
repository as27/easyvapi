package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// CustomTaxRateService manages all CRUD operations on the /custom-tax-rate endpoint.
// Custom tax rates (benutzerdefinierte Steuersätze) allow defining organization-specific
// tax rates for use in invoices and bookings.
type CustomTaxRateService struct {
	client *Client
}

// defaultCustomTaxRateQuery requests all fields defined in model.CustomTaxRate.
var defaultCustomTaxRateQuery = NewQuery().
	Fields("id", "taxName", "customTaxRate")

// CustomTaxRateListOptions holds all filter and pagination options for CustomTaxRate
// list requests.
type CustomTaxRateListOptions struct {
	ListOptions
	// TaxName filters by the tax rate name.
	TaxName string
}

// customTaxRateListParams converts opts into URL query parameters.
func customTaxRateListParams(opts *CustomTaxRateListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultCustomTaxRateQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultCustomTaxRateQuery)
	if opts.TaxName != "" {
		params.Set("taxName", opts.TaxName)
	}
	return params
}

// List returns a lazy Iterator over all CustomTaxRate records matching opts.
// Pages are fetched on-demand as iteration progresses.
//
// Example:
//
//	rates, err := client.CustomTaxRates.ListAll(ctx, nil)
//	for _, r := range rates {
//		fmt.Printf("%s: %.1f%%\n", r.TaxName, float64(r.CustomTaxRate))
//	}
func (s *CustomTaxRateService) List(ctx context.Context, opts *CustomTaxRateListOptions) *Iterator[model.CustomTaxRate] {
	startURL := s.client.buildURL("/custom-tax-rate", customTaxRateListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.CustomTaxRate, *string, error) {
		return fetchPage[model.CustomTaxRate](s.client, ctx, pageURL)
	})
}

// ListAll fetches all CustomTaxRate records matching opts and returns them as a slice.
func (s *CustomTaxRateService) ListAll(ctx context.Context, opts *CustomTaxRateListOptions) ([]model.CustomTaxRate, error) {
	var all []model.CustomTaxRate
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single CustomTaxRate by its ID.
func (s *CustomTaxRateService) Get(ctx context.Context, id int, query *Query) (*model.CustomTaxRate, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/custom-tax-rate/%d", id), params)
	if err != nil {
		return nil, err
	}
	var r model.CustomTaxRate
	if err := s.client.decodeJSON(resp, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Create creates a new CustomTaxRate and returns the created record.
func (s *CustomTaxRateService) Create(ctx context.Context, r model.CustomTaxRateCreate) (*model.CustomTaxRate, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/custom-tax-rate", nil), r)
	if err != nil {
		return nil, err
	}
	var created model.CustomTaxRate
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the CustomTaxRate with the given ID.
func (s *CustomTaxRateService) Update(ctx context.Context, id int, r model.CustomTaxRateCreate) (*model.CustomTaxRate, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/custom-tax-rate/%d", id), nil), r)
	if err != nil {
		return nil, err
	}
	var updated model.CustomTaxRate
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the CustomTaxRate with the given ID.
func (s *CustomTaxRateService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/custom-tax-rate/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

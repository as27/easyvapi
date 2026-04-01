package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/as27/easyvapi/model"
)

// InvoiceItemService manages all CRUD operations on the /invoice-item endpoint.
// Invoice items (Rechnungspositionen) are the individual line items within an invoice.
type InvoiceItemService struct {
	client *Client
}

// defaultInvoiceItemQuery requests all fields defined in model.InvoiceItem.
var defaultInvoiceItemQuery = NewQuery().
	Fields("id", "title", "quantity", "unitPrice", "taxRate", "taxName", "description",
		"billingAccount", "gross")

// InvoiceItemListOptions holds all filter and pagination options for InvoiceItem
// list requests.
type InvoiceItemListOptions struct {
	ListOptions
	// RelatedInvoice filters items belonging to the given invoice ID.
	RelatedInvoice int
	// BillingAccount filters items by billing account ID.
	BillingAccount int
	// Title filters by item title.
	Title string
}

// invoiceItemListParams converts opts into URL query parameters.
func invoiceItemListParams(opts *InvoiceItemListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultInvoiceItemQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultInvoiceItemQuery)
	if opts.RelatedInvoice != 0 {
		params.Set("relatedInvoice", strconv.Itoa(opts.RelatedInvoice))
	}
	if opts.BillingAccount != 0 {
		params.Set("billingAccount", strconv.Itoa(opts.BillingAccount))
	}
	if opts.Title != "" {
		params.Set("title", opts.Title)
	}
	return params
}

// List returns a lazy Iterator over all InvoiceItem records matching opts.
// Pages are fetched on-demand as iteration progresses.
//
// Example: Get all items for a specific invoice
//
//	opts := &easyvapi.InvoiceItemListOptions{RelatedInvoice: 123}
//	iter := client.InvoiceItems.List(ctx, opts)
//	for iter.Next() {
//		item := iter.Value()
//		fmt.Printf("%s: %.2f x %.2f\n", item.Title, float64(item.Quantity), float64(item.UnitPrice))
//	}
func (s *InvoiceItemService) List(ctx context.Context, opts *InvoiceItemListOptions) *Iterator[model.InvoiceItem] {
	startURL := s.client.buildURL("/invoice-item", invoiceItemListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.InvoiceItem, *string, error) {
		return fetchPage[model.InvoiceItem](s.client, ctx, pageURL)
	})
}

// ListAll fetches all InvoiceItem records matching opts and returns them as a slice.
func (s *InvoiceItemService) ListAll(ctx context.Context, opts *InvoiceItemListOptions) ([]model.InvoiceItem, error) {
	var all []model.InvoiceItem
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single InvoiceItem by its ID.
func (s *InvoiceItemService) Get(ctx context.Context, id int, query *Query) (*model.InvoiceItem, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/invoice-item/%d", id), params)
	if err != nil {
		return nil, err
	}
	var item model.InvoiceItem
	if err := s.client.decodeJSON(resp, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

// Create creates a new InvoiceItem and returns the created record.
func (s *InvoiceItemService) Create(ctx context.Context, item model.InvoiceItemCreate) (*model.InvoiceItem, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/invoice-item", nil), item)
	if err != nil {
		return nil, err
	}
	var created model.InvoiceItem
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the InvoiceItem with the given ID.
func (s *InvoiceItemService) Update(ctx context.Context, id int, item model.InvoiceItemCreate) (*model.InvoiceItem, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/invoice-item/%d", id), nil), item)
	if err != nil {
		return nil, err
	}
	var updated model.InvoiceItem
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the InvoiceItem with the given ID.
func (s *InvoiceItemService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/invoice-item/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

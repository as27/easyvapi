package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/as27/easyvapi/model"
)

// InvoiceService manages all CRUD operations on the /invoice endpoint.
// Use this service to manage invoices and financial documents.
type InvoiceService struct {
	client *Client
}

// defaultInvoiceQuery requests all fields defined in model.Invoice.
var defaultInvoiceQuery = NewQuery().
	Fields(
		"id", "relatedBookings", "org", "path", "relatedAddress",
		"payedFromUser", "approvedFromAdmin", "canceledInvoice", "cancelInvoice",
		"charges", "bankAccount", "invoiceItems",
		"_deleteAfterDate", "_deletedBy",
		"gross", "cancellationDescription", "templateName",
		"date", "dateItHappend", "dateSent",
		"invNumber", "receiver", "description",
		"totalPrice", "tax", "kind", "refNumber", "paymentDifference",
		"isDraft", "isTemplate",
		"creationDateForRecurringInvoices", "recurringInvoicesInterval",
		"paymentInformation", "isRequest", "taxRate", "taxName",
		"actualCallStateName", "callStateDelayDays",
		"accnumber", "guid", "selectionAcc", "removeFileOnDelete",
		"customPaymentMethod", "isReceipt",
		"_isTaxRatePerInvoiceItem", "_isSubjectToTax",
		"mode", "offerStatus", "offerValidUntil", "offerNumber",
		"relatedOffer", "closingDescription", "useAddressBalance",
	)

// InvoiceListOptions holds all filter and pagination options for Invoice list
// requests.
type InvoiceListOptions struct {
	ListOptions
	// Kind filters by invoice kind (e.g. "outgoing", "incoming").
	Kind string
	// IsDraft when non-nil filters by draft status.
	IsDraft *bool
	// IsTemplate when non-nil filters by template status.
	IsTemplate *bool
	// DateGt filters invoices with a date strictly after this value (YYYY-MM-DD).
	DateGt string
	// DateLt filters invoices with a date strictly before this value (YYYY-MM-DD).
	DateLt string
	// PaymentDifferenceNe filters invoices where paymentDifference != given value.
	// Set to "0" to retrieve only invoices with an outstanding balance.
	PaymentDifferenceNe string
	// PaymentDifferenceGte filters invoices where paymentDifference >= given value.
	PaymentDifferenceGte string
}

// invoiceListParams converts opts into URL query parameters.
func invoiceListParams(opts *InvoiceListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultInvoiceQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultInvoiceQuery)
	if opts.Kind != "" {
		params.Set("kind", opts.Kind)
	}
	if opts.IsDraft != nil {
		params.Set("isDraft", strconv.FormatBool(*opts.IsDraft))
	}
	if opts.IsTemplate != nil {
		params.Set("isTemplate", strconv.FormatBool(*opts.IsTemplate))
	}
	if opts.DateGt != "" {
		params.Set("date__gt", opts.DateGt)
	}
	if opts.DateLt != "" {
		params.Set("date__lt", opts.DateLt)
	}
	if opts.PaymentDifferenceNe != "" {
		params.Set("paymentDifference__ne", opts.PaymentDifferenceNe)
	}
	if opts.PaymentDifferenceGte != "" {
		params.Set("paymentDifference__gte", opts.PaymentDifferenceGte)
	}
	return params
}

// List returns a lazy Iterator over all Invoice records matching opts.
// Pages are fetched on-demand as iteration progresses.
//
// Example: Find all outgoing invoices from 2026
//
//	opts := &easyvapi.InvoiceListOptions{
//		Kind:   "outgoing",
//		DateGt: "2026-01-01",
//		DateLt: "2026-12-31",
//	}
//	iter := client.Invoices.List(ctx, opts)
//	for iter.Next() {
//		inv := iter.Value()
//		fmt.Printf("%s: %.2f EUR\n", inv.InvNumber, inv.TotalPrice)
//	}
func (s *InvoiceService) List(ctx context.Context, opts *InvoiceListOptions) *Iterator[model.Invoice] {
	startURL := s.client.buildURL("/invoice", invoiceListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.Invoice, *string, error) {
		return fetchPage[model.Invoice](s.client, ctx, pageURL)
	})
}

// ListAll fetches all Invoice records matching opts and returns them as a slice.
// This is a convenience wrapper that collects all pages into memory.
//
// Example: Get all draft invoices
//
//	isDraft := true
//	opts := &easyvapi.InvoiceListOptions{
//		IsDraft: &isDraft,
//	}
//	drafts, err := client.Invoices.ListAll(ctx, opts)
func (s *InvoiceService) ListAll(ctx context.Context, opts *InvoiceListOptions) ([]model.Invoice, error) {
	var all []model.Invoice
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single Invoice by its ID.
func (s *InvoiceService) Get(ctx context.Context, id int, query *Query) (*model.Invoice, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/invoice/%d", id), params)
	if err != nil {
		return nil, err
	}
	var inv model.Invoice
	if err := s.client.decodeJSON(resp, &inv); err != nil {
		return nil, err
	}
	return &inv, nil
}

// Create creates a new Invoice and returns the created record.
func (s *InvoiceService) Create(ctx context.Context, inv model.InvoiceCreate) (*model.Invoice, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/invoice", nil), inv)
	if err != nil {
		return nil, err
	}
	var created model.Invoice
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the Invoice with the given ID.
func (s *InvoiceService) Update(ctx context.Context, id int, inv model.InvoiceCreate) (*model.Invoice, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/invoice/%d", id), nil), inv)
	if err != nil {
		return nil, err
	}
	var updated model.Invoice
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the Invoice with the given ID.
func (s *InvoiceService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/invoice/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

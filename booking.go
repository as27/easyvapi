package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/as27/easyvapi/model"
)

// BookingService manages all CRUD operations on the /booking endpoint.
// Use this service to list, retrieve, create, update, and delete financial booking records.
type BookingService struct {
	client *Client
}

// defaultBookingQuery requests all fields defined in model.Booking.
var defaultBookingQuery = NewQuery().
	Fields("id", "amount", "date", "receiver", "billingId",
		"relatedInvoice", "org", "bankAccount", "billingAccount",
		"_deleteAfterDate", "_deletedBy", "description", "importDate",
		"blocked", "paymentDifference", "counterpartIban", "counterpartBic",
		"twingleDonation", "bookingProject", "sphere")

// BookingListOptions holds all filter and pagination options for Booking list
// requests.
type BookingListOptions struct {
	ListOptions
	// BankAccount filters bookings by bank account ID.
	BankAccount int
	// DateGt filters bookings with a date strictly after this value (YYYY-MM-DD).
	DateGt string
	// DateLt filters bookings with a date strictly before this value (YYYY-MM-DD).
	DateLt string
	// BillingAccount filters bookings by billing account ID.
	BillingAccount int
}

// bookingListParams converts opts into URL query parameters.
func bookingListParams(opts *BookingListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultBookingQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultBookingQuery)
	if opts.BankAccount != 0 {
		params.Set("bankAccount", strconv.Itoa(opts.BankAccount))
	}
	if opts.DateGt != "" {
		params.Set("date__gt", opts.DateGt)
	}
	if opts.DateLt != "" {
		params.Set("date__lt", opts.DateLt)
	}
	if opts.BillingAccount != 0 {
		params.Set("billingAccount", strconv.Itoa(opts.BillingAccount))
	}
	return params
}

// List returns a lazy Iterator over all Booking records matching opts.
// Pages are fetched on-demand as iteration progresses.
// Pass nil for opts to use default filtering and pagination.
//
// Example:
//
//	iter := client.Bookings.List(ctx, nil)
//	for iter.Next() {
//		booking := iter.Value()
//		fmt.Printf("Amount: %.2f EUR, Date: %s\n", booking.Amount, booking.Date)
//	}
func (s *BookingService) List(ctx context.Context, opts *BookingListOptions) *Iterator[model.Booking] {
	startURL := s.client.buildURL("/booking", bookingListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.Booking, *string, error) {
		return fetchPage[model.Booking](s.client, ctx, pageURL)
	})
}

// ListAll fetches all Booking records matching opts and returns them as a slice.
// This is a convenience wrapper that collects all pages into memory.
// For filtering by date or billing account, use BookingListOptions.
//
// Example: Get bookings from January 2026
//
//	opts := &easyvapi.BookingListOptions{
//		DateGt: "2026-01-01",
//		DateLt: "2026-01-31",
//	}
//	bookings, err := client.Bookings.ListAll(ctx, opts)
func (s *BookingService) ListAll(ctx context.Context, opts *BookingListOptions) ([]model.Booking, error) {
	var all []model.Booking
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single Booking by its ID.
func (s *BookingService) Get(ctx context.Context, id int, query *Query) (*model.Booking, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/booking/%d", id), params)
	if err != nil {
		return nil, err
	}
	var b model.Booking
	if err := s.client.decodeJSON(resp, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// Create creates a new Booking and returns the created record.
func (s *BookingService) Create(ctx context.Context, b model.BookingCreate) (*model.Booking, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/booking", nil), b)
	if err != nil {
		return nil, err
	}
	var created model.Booking
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// BulkCreate creates multiple Booking records in a single request and returns
// the created records with IDs assigned by the server. This is more efficient
// than calling [BookingService.Create] multiple times.
//
// Example:
//
//	bookings := []model.BookingCreate{
//		{Amount: 100.00, BillingAccount: 123, Date: "2026-03-31"},
//		{Amount: 50.00, BillingAccount: 124, Date: "2026-03-31"},
//	}
//	created, err := client.Bookings.BulkCreate(ctx, bookings)
func (s *BookingService) BulkCreate(ctx context.Context, bookings []model.BookingCreate) ([]model.Booking, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/booking/bulk", nil), bookings)
	if err != nil {
		return nil, err
	}
	var created []model.Booking
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return created, nil
}

// BulkUpdate applies a partial update (PATCH) to multiple Booking records.
// Each element in bookings must include the ID field (embedded via BookingCreate
// alongside the ID supplied separately). The updated records are returned.
func (s *BookingService) BulkUpdate(ctx context.Context, bookings []model.Booking) ([]model.Booking, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL("/booking/bulk", nil), bookings)
	if err != nil {
		return nil, err
	}
	var updated []model.Booking
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return updated, nil
}

// Update applies a partial update (PATCH) to the Booking with the given ID.
func (s *BookingService) Update(ctx context.Context, id int, b model.BookingCreate) (*model.Booking, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/booking/%d", id), nil), b)
	if err != nil {
		return nil, err
	}
	var updated model.Booking
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the Booking with the given ID.
func (s *BookingService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/booking/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

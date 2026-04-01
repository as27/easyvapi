package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/as27/easyvapi/model"
)

// ContactDetailsLogService manages all CRUD operations on the
// /contact-details-log endpoint.
type ContactDetailsLogService struct {
	client *Client
}

// defaultContactDetailsLogQuery is nil because the /contact-details-log endpoint
// does not support field selection via the query parameter.
var defaultContactDetailsLogQuery *Query = nil

// ContactDetailsLogListOptions holds all filter and pagination options for
// ContactDetailsLog list requests.
type ContactDetailsLogListOptions struct {
	ListOptions
	// ContactDetails filters log entries by contact ID.
	ContactDetails int
}

// contactDetailsLogListParams converts opts into URL query parameters.
func contactDetailsLogListParams(opts *ContactDetailsLogListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultContactDetailsLogQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultContactDetailsLogQuery)
	if opts.ContactDetails != 0 {
		params.Set("contactDetails", strconv.Itoa(opts.ContactDetails))
	}
	return params
}

// List returns a lazy Iterator over all ContactDetailsLog records matching opts.
//
// Example: Get all log entries for a specific contact
//
//	opts := &easyvapi.ContactDetailsLogListOptions{ContactDetails: 123}
//	iter := client.ContactDetailsLogs.List(ctx, opts)
//	for iter.Next() {
//		entry := iter.Value()
//		fmt.Printf("[%s] %s\n", entry.Date, entry.Title)
//	}
func (s *ContactDetailsLogService) List(ctx context.Context, opts *ContactDetailsLogListOptions) *Iterator[model.ContactDetailsLog] {
	startURL := s.client.buildURL("/contact-details-log", contactDetailsLogListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.ContactDetailsLog, *string, error) {
		return fetchPage[model.ContactDetailsLog](s.client, ctx, pageURL)
	})
}

// ListAll fetches all ContactDetailsLog records matching opts and returns them as a slice.
func (s *ContactDetailsLogService) ListAll(ctx context.Context, opts *ContactDetailsLogListOptions) ([]model.ContactDetailsLog, error) {
	var all []model.ContactDetailsLog
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single ContactDetailsLog entry by its ID.
func (s *ContactDetailsLogService) Get(ctx context.Context, id int, query *Query) (*model.ContactDetailsLog, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/contact-details-log/%d", id), params)
	if err != nil {
		return nil, err
	}
	var entry model.ContactDetailsLog
	if err := s.client.decodeJSON(resp, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// Create creates a new ContactDetailsLog entry and returns the created record.
func (s *ContactDetailsLogService) Create(ctx context.Context, entry model.ContactDetailsLogCreate) (*model.ContactDetailsLog, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/contact-details-log", nil), entry)
	if err != nil {
		return nil, err
	}
	var created model.ContactDetailsLog
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the ContactDetailsLog entry with the given ID.
func (s *ContactDetailsLogService) Update(ctx context.Context, id int, entry model.ContactDetailsLogCreate) (*model.ContactDetailsLog, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/contact-details-log/%d", id), nil), entry)
	if err != nil {
		return nil, err
	}
	var updated model.ContactDetailsLog
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the ContactDetailsLog entry with the given ID.
func (s *ContactDetailsLogService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/contact-details-log/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

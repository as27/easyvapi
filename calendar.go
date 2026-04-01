package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// CalendarService manages all CRUD operations on the /calendar endpoint.
// Calendars group events and can be linked to external iCal feeds.
type CalendarService struct {
	client *Client
}

// defaultCalendarQuery is nil because the /calendar endpoint does not support
// field selection via the query parameter.
var defaultCalendarQuery *Query = nil

// CalendarListOptions holds all filter and pagination options for Calendar
// list requests.
type CalendarListOptions struct {
	ListOptions
	// Name filters by calendar name.
	Name string
}

// calendarListParams converts opts into URL query parameters.
func calendarListParams(opts *CalendarListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultCalendarQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultCalendarQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	return params
}

// List returns a lazy Iterator over all Calendar records matching opts.
func (s *CalendarService) List(ctx context.Context, opts *CalendarListOptions) *Iterator[model.Calendar] {
	startURL := s.client.buildURL("/calendar", calendarListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.Calendar, *string, error) {
		return fetchPage[model.Calendar](s.client, ctx, pageURL)
	})
}

// ListAll fetches all Calendar records matching opts and returns them as a slice.
func (s *CalendarService) ListAll(ctx context.Context, opts *CalendarListOptions) ([]model.Calendar, error) {
	var all []model.Calendar
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single Calendar by its ID.
func (s *CalendarService) Get(ctx context.Context, id int, query *Query) (*model.Calendar, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/calendar/%d", id), params)
	if err != nil {
		return nil, err
	}
	var cal model.Calendar
	if err := s.client.decodeJSON(resp, &cal); err != nil {
		return nil, err
	}
	return &cal, nil
}

// Create creates a new Calendar and returns the created record.
func (s *CalendarService) Create(ctx context.Context, cal model.CalendarCreate) (*model.Calendar, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/calendar", nil), cal)
	if err != nil {
		return nil, err
	}
	var created model.Calendar
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the Calendar with the given ID.
func (s *CalendarService) Update(ctx context.Context, id int, cal model.CalendarCreate) (*model.Calendar, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/calendar/%d", id), nil), cal)
	if err != nil {
		return nil, err
	}
	var updated model.Calendar
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the Calendar with the given ID.
func (s *CalendarService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/calendar/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

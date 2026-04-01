package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// BookingProjectService manages all CRUD operations on the /booking-project endpoint.
// Booking projects (Buchungsprojekte) allow grouping bookings under a common project.
type BookingProjectService struct {
	client *Client
}

// defaultBookingProjectQuery is nil because the /booking-project endpoint does
// not support field selection via the query parameter.
var defaultBookingProjectQuery *Query = nil

// BookingProjectListOptions holds all filter and pagination options for BookingProject
// list requests.
type BookingProjectListOptions struct {
	ListOptions
	// Name filters by project name.
	Name string
}

// bookingProjectListParams converts opts into URL query parameters.
func bookingProjectListParams(opts *BookingProjectListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultBookingProjectQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultBookingProjectQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	return params
}

// List returns a lazy Iterator over all BookingProject records matching opts.
// Pages are fetched on-demand as iteration progresses.
//
// Example:
//
//	projects, err := client.BookingProjects.ListAll(ctx, nil)
//	for _, p := range projects {
//		fmt.Printf("%d: %s\n", p.ID, p.Name)
//	}
func (s *BookingProjectService) List(ctx context.Context, opts *BookingProjectListOptions) *Iterator[model.BookingProject] {
	startURL := s.client.buildURL("/booking-project", bookingProjectListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.BookingProject, *string, error) {
		return fetchPage[model.BookingProject](s.client, ctx, pageURL)
	})
}

// ListAll fetches all BookingProject records matching opts and returns them as a slice.
func (s *BookingProjectService) ListAll(ctx context.Context, opts *BookingProjectListOptions) ([]model.BookingProject, error) {
	var all []model.BookingProject
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single BookingProject by its ID.
func (s *BookingProjectService) Get(ctx context.Context, id int, query *Query) (*model.BookingProject, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/booking-project/%d", id), params)
	if err != nil {
		return nil, err
	}
	var p model.BookingProject
	if err := s.client.decodeJSON(resp, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Create creates a new BookingProject and returns the created record.
func (s *BookingProjectService) Create(ctx context.Context, p model.BookingProjectCreate) (*model.BookingProject, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/booking-project", nil), p)
	if err != nil {
		return nil, err
	}
	var created model.BookingProject
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the BookingProject with the given ID.
func (s *BookingProjectService) Update(ctx context.Context, id int, p model.BookingProjectCreate) (*model.BookingProject, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/booking-project/%d", id), nil), p)
	if err != nil {
		return nil, err
	}
	var updated model.BookingProject
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the BookingProject with the given ID.
func (s *BookingProjectService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/booking-project/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

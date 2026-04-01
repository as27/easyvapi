package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/as27/easyvapi/model"
)

// AnnouncementService manages all CRUD operations on the /announcement endpoint.
// Announcements can be shown as banners in the member or admin area.
type AnnouncementService struct {
	client *Client
}

// defaultAnnouncementQuery is nil because the /announcement endpoint does not
// support field selection via the query parameter.
var defaultAnnouncementQuery *Query = nil

// AnnouncementListOptions holds all filter and pagination options for
// Announcement list requests.
type AnnouncementListOptions struct {
	ListOptions
	// Platform filters by target platform (integer code).
	Platform int
	// ShowBanner when non-nil filters by banner display status.
	ShowBanner *bool
}

// announcementListParams converts opts into URL query parameters.
func announcementListParams(opts *AnnouncementListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultAnnouncementQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultAnnouncementQuery)
	if opts.Platform != 0 {
		params.Set("platform", strconv.Itoa(opts.Platform))
	}
	if opts.ShowBanner != nil {
		params.Set("showBanner", strconv.FormatBool(*opts.ShowBanner))
	}
	return params
}

// List returns a lazy Iterator over all Announcement records matching opts.
func (s *AnnouncementService) List(ctx context.Context, opts *AnnouncementListOptions) *Iterator[model.Announcement] {
	startURL := s.client.buildURL("/announcement", announcementListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.Announcement, *string, error) {
		return fetchPage[model.Announcement](s.client, ctx, pageURL)
	})
}

// ListAll fetches all Announcement records matching opts and returns them as a slice.
func (s *AnnouncementService) ListAll(ctx context.Context, opts *AnnouncementListOptions) ([]model.Announcement, error) {
	var all []model.Announcement
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single Announcement by its ID.
func (s *AnnouncementService) Get(ctx context.Context, id int, query *Query) (*model.Announcement, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/announcement/%d", id), params)
	if err != nil {
		return nil, err
	}
	var a model.Announcement
	if err := s.client.decodeJSON(resp, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

// Create creates a new Announcement and returns the created record.
func (s *AnnouncementService) Create(ctx context.Context, a model.AnnouncementCreate) (*model.Announcement, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/announcement", nil), a)
	if err != nil {
		return nil, err
	}
	var created model.Announcement
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the Announcement with the given ID.
func (s *AnnouncementService) Update(ctx context.Context, id int, a model.AnnouncementCreate) (*model.Announcement, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/announcement/%d", id), nil), a)
	if err != nil {
		return nil, err
	}
	var updated model.Announcement
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the Announcement with the given ID.
func (s *AnnouncementService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/announcement/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

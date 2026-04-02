package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/as27/easyvapi/model"
)

// ApplicationFormElementService manages all CRUD operations on the
// /application-form-element endpoint, including bulk operations.
type ApplicationFormElementService struct {
	client *Client
}

// defaultApplicationFormElementQuery is nil because the
// /application-form-element endpoint does not support the query parameter.
var defaultApplicationFormElementQuery *Query = nil

// ApplicationFormElementListOptions holds all filter and pagination options
// for ApplicationFormElement list requests.
type ApplicationFormElementListOptions struct {
	ListOptions
	// ApplicationForm filters elements by their parent form ID.
	ApplicationForm int
	// Kind filters by element type.
	Kind string
	// Required when non-nil filters by required status.
	Required *bool
}

// applicationFormElementListParams converts opts into URL query parameters.
func applicationFormElementListParams(opts *ApplicationFormElementListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultApplicationFormElementQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultApplicationFormElementQuery)
	if opts.ApplicationForm != 0 {
		params.Set("applicationForm", strconv.Itoa(opts.ApplicationForm))
	}
	if opts.Kind != "" {
		params.Set("kind", opts.Kind)
	}
	if opts.Required != nil {
		params.Set("required", strconv.FormatBool(*opts.Required))
	}
	return params
}

// List returns a lazy Iterator over all ApplicationFormElement records matching opts.
func (s *ApplicationFormElementService) List(ctx context.Context, opts *ApplicationFormElementListOptions) *Iterator[model.ApplicationFormElement] {
	startURL := s.client.buildURL("/application-form-element", applicationFormElementListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.ApplicationFormElement, *string, error) {
		return fetchPage[model.ApplicationFormElement](s.client, ctx, pageURL)
	})
}

// ListAll fetches all ApplicationFormElement records matching opts and returns them as a slice.
func (s *ApplicationFormElementService) ListAll(ctx context.Context, opts *ApplicationFormElementListOptions) ([]model.ApplicationFormElement, error) {
	var all []model.ApplicationFormElement
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single ApplicationFormElement by its ID.
func (s *ApplicationFormElementService) Get(ctx context.Context, id int, query *Query) (*model.ApplicationFormElement, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/application-form-element/%d", id), params)
	if err != nil {
		return nil, err
	}
	var e model.ApplicationFormElement
	if err := s.client.decodeJSON(resp, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// Create creates a new ApplicationFormElement and returns the created record.
func (s *ApplicationFormElementService) Create(ctx context.Context, e model.ApplicationFormElementCreate) (*model.ApplicationFormElement, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/application-form-element", nil), e)
	if err != nil {
		return nil, err
	}
	var created model.ApplicationFormElement
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// BulkCreate creates multiple ApplicationFormElement records in a single request.
func (s *ApplicationFormElementService) BulkCreate(ctx context.Context, entries []model.ApplicationFormElementCreate) ([]model.ApplicationFormElement, error) {
	body := map[string]any{"entries": entries}
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/application-form-element/bulk-create", nil), body)
	if err != nil {
		return nil, err
	}
	var created []model.ApplicationFormElement
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return created, nil
}

// BulkUpdate applies a partial update (PATCH) to multiple ApplicationFormElement records.
func (s *ApplicationFormElementService) BulkUpdate(ctx context.Context, entries []model.ApplicationFormElement) ([]model.ApplicationFormElement, error) {
	body := map[string]any{"entries": entries}
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL("/application-form-element/bulk-update", nil), body)
	if err != nil {
		return nil, err
	}
	var updated []model.ApplicationFormElement
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return updated, nil
}

// Update applies a partial update (PATCH) to the ApplicationFormElement with the given ID.
func (s *ApplicationFormElementService) Update(ctx context.Context, id int, e model.ApplicationFormElementCreate) (*model.ApplicationFormElement, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/application-form-element/%d", id), nil), e)
	if err != nil {
		return nil, err
	}
	var updated model.ApplicationFormElement
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the ApplicationFormElement with the given ID.
func (s *ApplicationFormElementService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/application-form-element/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

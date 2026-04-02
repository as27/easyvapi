package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/as27/easyvapi/model"
)

// ApplicationFormService manages all CRUD operations on the /application-form
// endpoint. Application forms are used for new member registrations.
type ApplicationFormService struct {
	client *Client
}

// defaultApplicationFormQuery is nil because the /application-form endpoint
// does not support field selection via the query parameter.
var defaultApplicationFormQuery *Query = nil

// ApplicationFormListOptions holds all filter and pagination options for
// ApplicationForm list requests.
type ApplicationFormListOptions struct {
	ListOptions
	// Title filters by form title.
	Title string
	// Public when non-nil filters by public accessibility.
	Public *bool
	// Language filters by language code.
	Language string
	// FormularKind filters by form type.
	FormularKind string
}

// applicationFormListParams converts opts into URL query parameters.
func applicationFormListParams(opts *ApplicationFormListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultApplicationFormQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultApplicationFormQuery)
	if opts.Title != "" {
		params.Set("title", opts.Title)
	}
	if opts.Public != nil {
		params.Set("public", strconv.FormatBool(*opts.Public))
	}
	if opts.Language != "" {
		params.Set("language", opts.Language)
	}
	if opts.FormularKind != "" {
		params.Set("formularKind", opts.FormularKind)
	}
	return params
}

// List returns a lazy Iterator over all ApplicationForm records matching opts.
func (s *ApplicationFormService) List(ctx context.Context, opts *ApplicationFormListOptions) *Iterator[model.ApplicationForm] {
	startURL := s.client.buildURL("/application-form", applicationFormListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.ApplicationForm, *string, error) {
		return fetchPage[model.ApplicationForm](s.client, ctx, pageURL)
	})
}

// ListAll fetches all ApplicationForm records matching opts and returns them as a slice.
func (s *ApplicationFormService) ListAll(ctx context.Context, opts *ApplicationFormListOptions) ([]model.ApplicationForm, error) {
	var all []model.ApplicationForm
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single ApplicationForm by its ID.
func (s *ApplicationFormService) Get(ctx context.Context, id int, query *Query) (*model.ApplicationForm, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/application-form/%d", id), params)
	if err != nil {
		return nil, err
	}
	var f model.ApplicationForm
	if err := s.client.decodeJSON(resp, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// Create creates a new ApplicationForm and returns the created record.
func (s *ApplicationFormService) Create(ctx context.Context, f model.ApplicationFormCreate) (*model.ApplicationForm, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/application-form", nil), f)
	if err != nil {
		return nil, err
	}
	var created model.ApplicationForm
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the ApplicationForm with the given ID.
func (s *ApplicationFormService) Update(ctx context.Context, id int, f model.ApplicationFormCreate) (*model.ApplicationForm, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/application-form/%d", id), nil), f)
	if err != nil {
		return nil, err
	}
	var updated model.ApplicationForm
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the ApplicationForm with the given ID.
func (s *ApplicationFormService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/application-form/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

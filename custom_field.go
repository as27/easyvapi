package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/as27/easyvapi/model"
)

// CustomFieldService manages all CRUD operations on the /custom-field endpoint.
// Custom fields extend records with organisation-specific attributes.
type CustomFieldService struct {
	client *Client
}

// defaultCustomFieldQuery is nil because the /custom-field endpoint does
// not support field selection via the query parameter.
var defaultCustomFieldQuery *Query = nil

// CustomFieldListOptions holds all filter and pagination options for CustomField
// list requests.
type CustomFieldListOptions struct {
	ListOptions
	// Label filters by field label.
	Label string
	// FieldKind filters by field type (e.g. "text", "number", "date", "select").
	FieldKind string
	// FieldCollection filters by collection ID.
	FieldCollection int
	// ShowInMemberArea filters by visibility in the member area.
	ShowInMemberArea *bool
}

// customFieldListParams converts opts into URL query parameters.
func customFieldListParams(opts *CustomFieldListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultCustomFieldQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultCustomFieldQuery)
	if opts.Label != "" {
		params.Set("label", opts.Label)
	}
	if opts.FieldKind != "" {
		params.Set("fieldKind", opts.FieldKind)
	}
	if opts.FieldCollection != 0 {
		params.Set("fieldCollection", strconv.Itoa(opts.FieldCollection))
	}
	if opts.ShowInMemberArea != nil {
		params.Set("showInMemberArea", strconv.FormatBool(*opts.ShowInMemberArea))
	}
	return params
}

// List returns a lazy Iterator over all CustomField records matching opts.
func (s *CustomFieldService) List(ctx context.Context, opts *CustomFieldListOptions) *Iterator[model.CustomField] {
	startURL := s.client.buildURL("/custom-field", customFieldListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.CustomField, *string, error) {
		return fetchPage[model.CustomField](s.client, ctx, pageURL)
	})
}

// ListAll fetches all CustomField records matching opts and returns them as a slice.
func (s *CustomFieldService) ListAll(ctx context.Context, opts *CustomFieldListOptions) ([]model.CustomField, error) {
	var all []model.CustomField
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single CustomField by its ID.
func (s *CustomFieldService) Get(ctx context.Context, id int, query *Query) (*model.CustomField, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/custom-field/%d", id), params)
	if err != nil {
		return nil, err
	}
	var f model.CustomField
	if err := s.client.decodeJSON(resp, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// Create creates a new CustomField and returns the created record.
func (s *CustomFieldService) Create(ctx context.Context, f model.CustomFieldCreate) (*model.CustomField, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/custom-field", nil), f)
	if err != nil {
		return nil, err
	}
	var created model.CustomField
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the CustomField with the given ID.
func (s *CustomFieldService) Update(ctx context.Context, id int, f model.CustomFieldCreate) (*model.CustomField, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/custom-field/%d", id), nil), f)
	if err != nil {
		return nil, err
	}
	var updated model.CustomField
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the CustomField with the given ID.
func (s *CustomFieldService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/custom-field/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

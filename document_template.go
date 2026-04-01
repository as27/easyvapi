package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// DocumentTemplateService manages all CRUD operations on the /document-template
// endpoint. Document templates are used to generate letters, certificates, and
// other documents for members.
type DocumentTemplateService struct {
	client *Client
}

// defaultDocumentTemplateQuery is nil because the /document-template endpoint
// does not support field selection via the query parameter.
var defaultDocumentTemplateQuery *Query = nil

// DocumentTemplateListOptions holds all filter and pagination options for
// DocumentTemplate list requests.
type DocumentTemplateListOptions struct {
	ListOptions
	// Title filters by template title.
	Title string
	// DocumentKind filters by document kind.
	DocumentKind string
}

// documentTemplateListParams converts opts into URL query parameters.
func documentTemplateListParams(opts *DocumentTemplateListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultDocumentTemplateQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultDocumentTemplateQuery)
	if opts.Title != "" {
		params.Set("title", opts.Title)
	}
	if opts.DocumentKind != "" {
		params.Set("documentKind", opts.DocumentKind)
	}
	return params
}

// List returns a lazy Iterator over all DocumentTemplate records matching opts.
// Note: the content field is excluded from the default query for performance.
// Use Get to retrieve the full template content.
func (s *DocumentTemplateService) List(ctx context.Context, opts *DocumentTemplateListOptions) *Iterator[model.DocumentTemplate] {
	startURL := s.client.buildURL("/document-template", documentTemplateListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.DocumentTemplate, *string, error) {
		return fetchPage[model.DocumentTemplate](s.client, ctx, pageURL)
	})
}

// ListAll fetches all DocumentTemplate records matching opts and returns them as a slice.
func (s *DocumentTemplateService) ListAll(ctx context.Context, opts *DocumentTemplateListOptions) ([]model.DocumentTemplate, error) {
	var all []model.DocumentTemplate
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single DocumentTemplate by its ID, including its full content.
func (s *DocumentTemplateService) Get(ctx context.Context, id int, query *Query) (*model.DocumentTemplate, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/document-template/%d", id), params)
	if err != nil {
		return nil, err
	}
	var t model.DocumentTemplate
	if err := s.client.decodeJSON(resp, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// Create creates a new DocumentTemplate and returns the created record.
func (s *DocumentTemplateService) Create(ctx context.Context, t model.DocumentTemplateCreate) (*model.DocumentTemplate, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/document-template", nil), t)
	if err != nil {
		return nil, err
	}
	var created model.DocumentTemplate
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the DocumentTemplate with the given ID.
func (s *DocumentTemplateService) Update(ctx context.Context, id int, t model.DocumentTemplateCreate) (*model.DocumentTemplate, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/document-template/%d", id), nil), t)
	if err != nil {
		return nil, err
	}
	var updated model.DocumentTemplate
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the DocumentTemplate with the given ID.
func (s *DocumentTemplateService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/document-template/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

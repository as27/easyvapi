package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// ArticleObjectService manages all CRUD operations on the /article-object
// endpoint. Articles represent shop items or event tickets.
type ArticleObjectService struct {
	client *Client
}

// defaultArticleObjectQuery is nil because the /article-object endpoint does
// not support field selection via the query parameter.
var defaultArticleObjectQuery *Query = nil

// ArticleObjectListOptions holds all filter and pagination options for
// ArticleObject list requests.
type ArticleObjectListOptions struct {
	ListOptions
	// Name filters by article name.
	Name string
	// Kind filters by article type.
	Kind string
}

// articleObjectListParams converts opts into URL query parameters.
func articleObjectListParams(opts *ArticleObjectListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultArticleObjectQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultArticleObjectQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	if opts.Kind != "" {
		params.Set("kind", opts.Kind)
	}
	return params
}

// List returns a lazy Iterator over all ArticleObject records matching opts.
func (s *ArticleObjectService) List(ctx context.Context, opts *ArticleObjectListOptions) *Iterator[model.ArticleObject] {
	startURL := s.client.buildURL("/article-object", articleObjectListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.ArticleObject, *string, error) {
		return fetchPage[model.ArticleObject](s.client, ctx, pageURL)
	})
}

// ListAll fetches all ArticleObject records matching opts and returns them as a slice.
func (s *ArticleObjectService) ListAll(ctx context.Context, opts *ArticleObjectListOptions) ([]model.ArticleObject, error) {
	var all []model.ArticleObject
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single ArticleObject by its ID.
func (s *ArticleObjectService) Get(ctx context.Context, id int, query *Query) (*model.ArticleObject, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/article-object/%d", id), params)
	if err != nil {
		return nil, err
	}
	var a model.ArticleObject
	if err := s.client.decodeJSON(resp, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

// Create creates a new ArticleObject and returns the created record.
func (s *ArticleObjectService) Create(ctx context.Context, a model.ArticleObjectCreate) (*model.ArticleObject, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/article-object", nil), a)
	if err != nil {
		return nil, err
	}
	var created model.ArticleObject
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the ArticleObject with the given ID.
func (s *ArticleObjectService) Update(ctx context.Context, id int, a model.ArticleObjectCreate) (*model.ArticleObject, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/article-object/%d", id), nil), a)
	if err != nil {
		return nil, err
	}
	var updated model.ArticleObject
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the ArticleObject with the given ID.
func (s *ArticleObjectService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/article-object/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

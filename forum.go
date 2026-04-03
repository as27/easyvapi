package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// ForumService manages all CRUD operations on the /forum endpoint.
type ForumService struct {
	client *Client
}

// defaultForumQuery is nil because the endpoint does not support
// field selection via the query parameter.
var defaultForumQuery *Query = nil

// ForumListOptions holds filter and pagination options for list requests.
type ForumListOptions struct {
	ListOptions
	// Name filters by forum name.
	Name string
}

func forumListParams(opts *ForumListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultForumQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultForumQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	return params
}

// List returns a lazy Iterator over all Forum records matching opts.
func (s *ForumService) List(ctx context.Context, opts *ForumListOptions) *Iterator[model.Forum] {
	startURL := s.client.buildURL("/forum", forumListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.Forum, *string, error) {
		return fetchPage[model.Forum](s.client, ctx, pageURL)
	})
}

// ListAll fetches all Forum records matching opts and returns them as a slice.
func (s *ForumService) ListAll(ctx context.Context, opts *ForumListOptions) ([]model.Forum, error) {
	var all []model.Forum
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single Forum by its ID.
func (s *ForumService) Get(ctx context.Context, id int, query *Query) (*model.Forum, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/forum/%d", id), params)
	if err != nil {
		return nil, err
	}
	var forum model.Forum
	if err := s.client.decodeJSON(resp, &forum); err != nil {
		return nil, err
	}
	return &forum, nil
}

// Create creates a new Forum and returns the created record.
func (s *ForumService) Create(ctx context.Context, forum model.ForumCreate) (*model.Forum, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/forum", nil), forum)
	if err != nil {
		return nil, err
	}
	var created model.Forum
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the Forum with the given ID.
func (s *ForumService) Update(ctx context.Context, id int, forum model.ForumCreate) (*model.Forum, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/forum/%d", id), nil), forum)
	if err != nil {
		return nil, err
	}
	var updated model.Forum
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the Forum with the given ID.
func (s *ForumService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/forum/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

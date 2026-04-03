package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// FileSystemPathService manages all CRUD operations on the /file-system-path endpoint.
type FileSystemPathService struct {
	client *Client
}

// defaultFileSystemPathQuery is nil because the endpoint does not support
// field selection via the query parameter.
var defaultFileSystemPathQuery *Query = nil

// FileSystemPathListOptions holds filter and pagination options for list requests.
type FileSystemPathListOptions struct {
	ListOptions
	// Name filters by path name.
	Name string
}

func fileSystemPathListParams(opts *FileSystemPathListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultFileSystemPathQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultFileSystemPathQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	return params
}

// List returns a lazy Iterator over all FileSystemPath records matching opts.
func (s *FileSystemPathService) List(ctx context.Context, opts *FileSystemPathListOptions) *Iterator[model.FileSystemPath] {
	startURL := s.client.buildURL("/file-system-path", fileSystemPathListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.FileSystemPath, *string, error) {
		return fetchPage[model.FileSystemPath](s.client, ctx, pageURL)
	})
}

// ListAll fetches all FileSystemPath records matching opts and returns them as a slice.
func (s *FileSystemPathService) ListAll(ctx context.Context, opts *FileSystemPathListOptions) ([]model.FileSystemPath, error) {
	var all []model.FileSystemPath
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single FileSystemPath by its ID.
func (s *FileSystemPathService) Get(ctx context.Context, id int, query *Query) (*model.FileSystemPath, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/file-system-path/%d", id), params)
	if err != nil {
		return nil, err
	}
	var fsp model.FileSystemPath
	if err := s.client.decodeJSON(resp, &fsp); err != nil {
		return nil, err
	}
	return &fsp, nil
}

// Create creates a new FileSystemPath and returns the created record.
func (s *FileSystemPathService) Create(ctx context.Context, fsp model.FileSystemPathCreate) (*model.FileSystemPath, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/file-system-path", nil), fsp)
	if err != nil {
		return nil, err
	}
	var created model.FileSystemPath
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the FileSystemPath with the given ID.
func (s *FileSystemPathService) Update(ctx context.Context, id int, fsp model.FileSystemPathCreate) (*model.FileSystemPath, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/file-system-path/%d", id), nil), fsp)
	if err != nil {
		return nil, err
	}
	var updated model.FileSystemPath
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the FileSystemPath with the given ID.
func (s *FileSystemPathService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/file-system-path/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

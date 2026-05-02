package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// TaskGroupService manages all CRUD operations on the /task-group endpoint.
// Task groups are used to organize and categorize tasks.
type TaskGroupService struct {
	client *Client
}

// defaultTaskGroupQuery is nil because the /task-group endpoint
// does not support field selection via the query parameter.
var defaultTaskGroupQuery *Query = nil

// TaskGroupListOptions holds all filter and pagination options for
// TaskGroup list requests.
type TaskGroupListOptions struct {
	ListOptions
	// Name filters by group name.
	Name string
	// Color filters by hex color code.
	Color string
	// Short filters by group abbreviation.
	Short string
}

// taskGroupListParams converts opts into URL query parameters.
func taskGroupListParams(opts *TaskGroupListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultTaskGroupQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultTaskGroupQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	if opts.Color != "" {
		params.Set("color", opts.Color)
	}
	if opts.Short != "" {
		params.Set("short", opts.Short)
	}
	return params
}

// List returns a lazy Iterator over all TaskGroup records matching opts.
func (s *TaskGroupService) List(ctx context.Context, opts *TaskGroupListOptions) *Iterator[model.TaskGroup] {
	startURL := s.client.buildURL("/task-group", taskGroupListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.TaskGroup, *string, error) {
		return fetchPage[model.TaskGroup](s.client, ctx, pageURL)
	})
}

// ListAll fetches all TaskGroup records matching opts and returns them as a slice.
func (s *TaskGroupService) ListAll(ctx context.Context, opts *TaskGroupListOptions) ([]model.TaskGroup, error) {
	var all []model.TaskGroup
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single TaskGroup by its ID.
func (s *TaskGroupService) Get(ctx context.Context, id int, query *Query) (*model.TaskGroup, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/task-group/%d", id), params)
	if err != nil {
		return nil, err
	}
	var g model.TaskGroup
	if err := s.client.decodeJSON(resp, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

// Create creates a new TaskGroup and returns the created record.
func (s *TaskGroupService) Create(ctx context.Context, g model.TaskGroupCreate) (*model.TaskGroup, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/task-group", nil), g)
	if err != nil {
		return nil, err
	}
	var created model.TaskGroup
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the TaskGroup with the given ID.
func (s *TaskGroupService) Update(ctx context.Context, id int, g model.TaskGroupCreate) (*model.TaskGroup, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/task-group/%d", id), nil), g)
	if err != nil {
		return nil, err
	}
	var updated model.TaskGroup
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the TaskGroup with the given ID.
func (s *TaskGroupService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/task-group/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

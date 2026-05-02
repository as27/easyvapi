package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/as27/easyvapi/model"
)

// TaskService manages all CRUD operations on the /task endpoint.
// Tasks are to-do items that can be assigned to members and grouped.
type TaskService struct {
	client *Client
}

// defaultTaskQuery is nil because the /task endpoint
// does not support field selection via the query parameter.
var defaultTaskQuery *Query = nil

// TaskListOptions holds all filter and pagination options for
// Task list requests.
type TaskListOptions struct {
	ListOptions
	// Name filters by task name.
	Name string
	// State filters by task state (e.g. "offen", "erledigt").
	State string
	// StateNe excludes tasks with this state.
	StateNe string
	// Public when non-nil filters by public visibility.
	Public *bool
	// MemberIsnull when non-nil filters tasks with no assigned member.
	MemberIsnull *bool
	// DueGte filters tasks with a due date on or after this value.
	DueGte string
	// DueLte filters tasks with a due date on or before this value.
	DueLte string
	// Member filters by assigned member ID.
	Member *int
	// TaskGroup filters by task group ID.
	TaskGroup *int
	// ParentEvent filters by parent event ID.
	ParentEvent *int
}

// taskListParams converts opts into URL query parameters.
func taskListParams(opts *TaskListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultTaskQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultTaskQuery)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	if opts.State != "" {
		params.Set("state", opts.State)
	}
	if opts.StateNe != "" {
		params.Set("state__ne", opts.StateNe)
	}
	if opts.Public != nil {
		params.Set("public", strconv.FormatBool(*opts.Public))
	}
	if opts.MemberIsnull != nil {
		params.Set("member__isnull", strconv.FormatBool(*opts.MemberIsnull))
	}
	if opts.DueGte != "" {
		params.Set("due__gte", opts.DueGte)
	}
	if opts.DueLte != "" {
		params.Set("due__lte", opts.DueLte)
	}
	if opts.Member != nil {
		params.Set("member", strconv.Itoa(*opts.Member))
	}
	if opts.TaskGroup != nil {
		params.Set("taskGroup", strconv.Itoa(*opts.TaskGroup))
	}
	if opts.ParentEvent != nil {
		params.Set("parentEvent", strconv.Itoa(*opts.ParentEvent))
	}
	return params
}

// List returns a lazy Iterator over all Task records matching opts.
func (s *TaskService) List(ctx context.Context, opts *TaskListOptions) *Iterator[model.Task] {
	startURL := s.client.buildURL("/task", taskListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.Task, *string, error) {
		return fetchPage[model.Task](s.client, ctx, pageURL)
	})
}

// ListAll fetches all Task records matching opts and returns them as a slice.
func (s *TaskService) ListAll(ctx context.Context, opts *TaskListOptions) ([]model.Task, error) {
	var all []model.Task
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single Task by its ID.
func (s *TaskService) Get(ctx context.Context, id int, query *Query) (*model.Task, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/task/%d", id), params)
	if err != nil {
		return nil, err
	}
	var t model.Task
	if err := s.client.decodeJSON(resp, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// Create creates a new Task and returns the created record.
func (s *TaskService) Create(ctx context.Context, t model.TaskCreate) (*model.Task, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/task", nil), t)
	if err != nil {
		return nil, err
	}
	var created model.Task
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the Task with the given ID.
func (s *TaskService) Update(ctx context.Context, id int, t model.TaskCreate) (*model.Task, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/task/%d", id), nil), t)
	if err != nil {
		return nil, err
	}
	var updated model.Task
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the Task with the given ID.
func (s *TaskService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/task/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// AssignToMe assigns the task with the given ID to the current user.
func (s *TaskService) AssignToMe(ctx context.Context, id int) (*model.Task, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL(fmt.Sprintf("/task/%d/assign-to-me", id), nil), nil)
	if err != nil {
		return nil, err
	}
	var t model.Task
	if err := s.client.decodeJSON(resp, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// Unassign removes the member assignment from the task with the given ID.
func (s *TaskService) Unassign(ctx context.Context, id int) (*model.Task, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL(fmt.Sprintf("/task/%d/unassign", id), nil), nil)
	if err != nil {
		return nil, err
	}
	var t model.Task
	if err := s.client.decodeJSON(resp, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

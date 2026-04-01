package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// MemberGroupService manages all CRUD operations on the /member-group endpoint.
// Use this service to manage member categories and groups.
// Note: The /member-group endpoint does not support the query parameter,
// so custom field selection is not available for this service.
type MemberGroupService struct {
	client *Client
}

// MemberGroupListOptions holds filter and pagination options for MemberGroup
// list requests.
type MemberGroupListOptions struct {
	ListOptions
	// Name filters groups by name.
	Name string
}

// memberGroupListParams converts opts into URL query parameters.
func memberGroupListParams(opts *MemberGroupListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, nil)
		return params
	}
	applyListOptions(params, opts.ListOptions, nil)
	if opts.Name != "" {
		params.Set("name", opts.Name)
	}
	return params
}

// List returns a lazy Iterator over all MemberGroup records matching opts.
// Pages are fetched on-demand as iteration progresses.
//
// Example: Find all member groups with "active" in the name
//
//	opts := &easyvapi.MemberGroupListOptions{
//		Name: "active",
//	}
//	iter := client.MemberGroups.List(ctx, opts)
//	for iter.Next() {
//		group := iter.Value()
//		fmt.Printf("%s (%s): %s\n", group.Name, group.Short, group.Description)
//	}
func (s *MemberGroupService) List(ctx context.Context, opts *MemberGroupListOptions) *Iterator[model.MemberGroup] {
	startURL := s.client.buildURL("/member-group", memberGroupListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.MemberGroup, *string, error) {
		return fetchPage[model.MemberGroup](s.client, ctx, pageURL)
	})
}

// ListAll fetches all MemberGroup records matching opts and returns them as a slice.
// This is a convenience wrapper that collects all pages into memory.
//
// Example: Get all member groups
//
//	groups, err := client.MemberGroups.ListAll(ctx, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Found %d member groups\n", len(groups))
func (s *MemberGroupService) ListAll(ctx context.Context, opts *MemberGroupListOptions) ([]model.MemberGroup, error) {
	var all []model.MemberGroup
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single MemberGroup by its ID.
func (s *MemberGroupService) Get(ctx context.Context, id int, query *Query) (*model.MemberGroup, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/member-group/%d", id), params)
	if err != nil {
		return nil, err
	}
	var g model.MemberGroup
	if err := s.client.decodeJSON(resp, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

// Create creates a new MemberGroup and returns the created record.
func (s *MemberGroupService) Create(ctx context.Context, g model.MemberGroupCreate) (*model.MemberGroup, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/member-group", nil), g)
	if err != nil {
		return nil, err
	}
	var created model.MemberGroup
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the MemberGroup with the given ID.
func (s *MemberGroupService) Update(ctx context.Context, id int, g model.MemberGroupCreate) (*model.MemberGroup, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/member-group/%d", id), nil), g)
	if err != nil {
		return nil, err
	}
	var updated model.MemberGroup
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the MemberGroup with the given ID.
func (s *MemberGroupService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/member-group/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

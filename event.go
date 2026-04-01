package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/as27/easyvapi/model"
)

// EventService manages all CRUD operations on the /event endpoint.
// Use this service to list, retrieve, create, update, and delete calendar events.
type EventService struct {
	client *Client
}

// defaultEventQuery requests all fields defined in model.Event.
var defaultEventQuery = NewQuery().
	Fields("id", "name", "start", "end", "allDay", "isPublic", "canceled",
		"locationName")

// EventListOptions holds all filter and pagination options for Event list
// requests.
type EventListOptions struct {
	ListOptions
	// StartGte filters events whose start is on or after this date/time (ISO 8601).
	StartGte string
	// StartLte filters events whose start is on or before this date/time (ISO 8601).
	StartLte string
	// Calendar filters by calendar ID.
	Calendar int
	// IsPublic when non-nil filters by public visibility.
	IsPublic *bool
}

// eventListParams converts opts into URL query parameters.
func eventListParams(opts *EventListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultEventQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultEventQuery)
	if opts.StartGte != "" {
		params.Set("start__gte", opts.StartGte)
	}
	if opts.StartLte != "" {
		params.Set("start__lte", opts.StartLte)
	}
	if opts.Calendar != 0 {
		params.Set("calendar", strconv.Itoa(opts.Calendar))
	}
	if opts.IsPublic != nil {
		params.Set("isPublic", strconv.FormatBool(*opts.IsPublic))
	}
	return params
}

// List returns a lazy Iterator over all Event records matching opts.
// Pages are fetched on-demand as iteration progresses.
// Pass nil for opts to use default filtering and pagination.
//
// Example: Find events starting after a specific date
//
//	opts := &easyvapi.EventListOptions{
//		StartGte: "2026-03-01T00:00:00Z",
//		ListOptions: easyvapi.ListOptions{
//			Search: "Workshop",
//		},
//	}
//	iter := client.Events.List(ctx, opts)
//	for iter.Next() {
//		event := iter.Value()
//		fmt.Printf("%s: %s\n", event.Name, event.Start)
//	}
func (s *EventService) List(ctx context.Context, opts *EventListOptions) *Iterator[model.Event] {
	startURL := s.client.buildURL("/event", eventListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.Event, *string, error) {
		return fetchPage[model.Event](s.client, ctx, pageURL)
	})
}

// ListAll fetches all Event records matching opts and returns them as a slice.
// This is a convenience wrapper that collects all pages into memory.
// For large event databases, consider using List with Iterator instead.
//
// Example: Get all public events
//
//	isPublic := true
//	opts := &easyvapi.EventListOptions{
//		IsPublic: &isPublic,
//	}
//	events, err := client.Events.ListAll(ctx, opts)
func (s *EventService) ListAll(ctx context.Context, opts *EventListOptions) ([]model.Event, error) {
	var all []model.Event
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single Event by its ID.
func (s *EventService) Get(ctx context.Context, id int, query *Query) (*model.Event, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/event/%d", id), params)
	if err != nil {
		return nil, err
	}
	var e model.Event
	if err := s.client.decodeJSON(resp, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// Create creates a new Event and returns the created record.
func (s *EventService) Create(ctx context.Context, e model.EventCreate) (*model.Event, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/event", nil), e)
	if err != nil {
		return nil, err
	}
	var created model.Event
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the Event with the given ID.
func (s *EventService) Update(ctx context.Context, id int, e model.EventCreate) (*model.Event, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/event/%d", id), nil), e)
	if err != nil {
		return nil, err
	}
	var updated model.Event
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the Event with the given ID.
func (s *EventService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/event/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// Copy creates a copy of the event with the given ID and returns the new event.
//
// Example:
//
//	copied, err := client.Events.Copy(ctx, 123)
func (s *EventService) Copy(ctx context.Context, id int) (*model.Event, error) {
	resp, err := s.client.do(ctx, "GET", s.client.buildURL(fmt.Sprintf("/event/%d/copy", id), nil), nil)
	if err != nil {
		return nil, err
	}
	var e model.Event
	if err := s.client.decodeJSON(resp, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// GenerateInvoices triggers invoice generation for all participants of the
// event with the given ID.
func (s *EventService) GenerateInvoices(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL(fmt.Sprintf("/event/%d/generate-invoices", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// InviteGroups invites a contact-details group or member group to the event.
//
// Example:
//
//	err := client.Events.InviteGroups(ctx, 123, model.InviteGroupsRequest{MemberGroupID: 42})
func (s *EventService) InviteGroups(ctx context.Context, id int, req model.InviteGroupsRequest) error {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL(fmt.Sprintf("/event/%d/invite-groups", id), nil), req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListParticipations returns a lazy Iterator over all participation records for
// the event with the given ID.
//
// Example:
//
//	iter := client.Events.ListParticipations(ctx, 123, nil)
//	for iter.Next() {
//		p := iter.Value()
//		fmt.Printf("Participant %d, state %d\n", p.ParticipationAddress, p.State)
//	}
func (s *EventService) ListParticipations(ctx context.Context, eventID int, opts *ListOptions) *Iterator[model.Participation] {
	params := url.Values{}
	if opts != nil {
		applyListOptions(params, *opts, nil)
	} else {
		applyListOptions(params, ListOptions{}, nil)
	}
	startURL := s.client.buildURL(fmt.Sprintf("/event/%d/participation", eventID), params)
	return newIterator(startURL, func(pageURL string) ([]model.Participation, *string, error) {
		return fetchPage[model.Participation](s.client, ctx, pageURL)
	})
}

// ListAllParticipations fetches all participation records for the event and
// returns them as a slice.
func (s *EventService) ListAllParticipations(ctx context.Context, eventID int, opts *ListOptions) ([]model.Participation, error) {
	var all []model.Participation
	iter := s.ListParticipations(ctx, eventID, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// GetParticipation retrieves a single participation record by event ID and
// participation ID.
func (s *EventService) GetParticipation(ctx context.Context, eventID, participationID int) (*model.Participation, error) {
	resp, err := s.client.get(ctx, fmt.Sprintf("/event/%d/participation/%d", eventID, participationID), url.Values{})
	if err != nil {
		return nil, err
	}
	var p model.Participation
	if err := s.client.decodeJSON(resp, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// CreateParticipation creates a new participation record for the given event.
func (s *EventService) CreateParticipation(ctx context.Context, eventID int, p model.ParticipationCreate) (*model.Participation, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL(fmt.Sprintf("/event/%d/participation", eventID), nil), p)
	if err != nil {
		return nil, err
	}
	var created model.Participation
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// UpdateParticipation applies a partial update (PATCH) to a participation record.
func (s *EventService) UpdateParticipation(ctx context.Context, eventID, participationID int, p model.ParticipationCreate) (*model.Participation, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/event/%d/participation/%d", eventID, participationID), nil), p)
	if err != nil {
		return nil, err
	}
	var updated model.Participation
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// DeleteParticipation removes a participation record from the given event.
func (s *EventService) DeleteParticipation(ctx context.Context, eventID, participationID int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/event/%d/participation/%d", eventID, participationID), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

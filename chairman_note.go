package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// ChairmanNoteService manages all CRUD operations on the /chairman-note
// endpoint. Chairman notes are internal notes visible only to board members.
type ChairmanNoteService struct {
	client *Client
}

// defaultChairmanNoteQuery requests all fields defined in model.ChairmanNote.
var defaultChairmanNoteQuery = NewQuery().
	Fields("id", "text", "date", "_deleteAfterDate")

// ChairmanNoteListOptions holds all filter and pagination options for
// ChairmanNote list requests.
type ChairmanNoteListOptions struct {
	ListOptions
	// DateGte filters notes with a date on or after this value (YYYY-MM-DD).
	DateGte string
	// DateLte filters notes with a date on or before this value (YYYY-MM-DD).
	DateLte string
}

// chairmanNoteListParams converts opts into URL query parameters.
func chairmanNoteListParams(opts *ChairmanNoteListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultChairmanNoteQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultChairmanNoteQuery)
	if opts.DateGte != "" {
		params.Set("date__gte", opts.DateGte)
	}
	if opts.DateLte != "" {
		params.Set("date__lte", opts.DateLte)
	}
	return params
}

// List returns a lazy Iterator over all ChairmanNote records matching opts.
//
// Example:
//
//	iter := client.ChairmanNotes.List(ctx, nil)
//	for iter.Next() {
//		note := iter.Value()
//		fmt.Printf("[%s] %s\n", note.Date, note.Text)
//	}
func (s *ChairmanNoteService) List(ctx context.Context, opts *ChairmanNoteListOptions) *Iterator[model.ChairmanNote] {
	startURL := s.client.buildURL("/chairman-note", chairmanNoteListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.ChairmanNote, *string, error) {
		return fetchPage[model.ChairmanNote](s.client, ctx, pageURL)
	})
}

// ListAll fetches all ChairmanNote records matching opts and returns them as a slice.
func (s *ChairmanNoteService) ListAll(ctx context.Context, opts *ChairmanNoteListOptions) ([]model.ChairmanNote, error) {
	var all []model.ChairmanNote
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single ChairmanNote by its ID.
func (s *ChairmanNoteService) Get(ctx context.Context, id int, query *Query) (*model.ChairmanNote, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/chairman-note/%d", id), params)
	if err != nil {
		return nil, err
	}
	var note model.ChairmanNote
	if err := s.client.decodeJSON(resp, &note); err != nil {
		return nil, err
	}
	return &note, nil
}

// Create creates a new ChairmanNote and returns the created record.
func (s *ChairmanNoteService) Create(ctx context.Context, note model.ChairmanNoteCreate) (*model.ChairmanNote, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/chairman-note", nil), note)
	if err != nil {
		return nil, err
	}
	var created model.ChairmanNote
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the ChairmanNote with the given ID.
func (s *ChairmanNoteService) Update(ctx context.Context, id int, note model.ChairmanNoteCreate) (*model.ChairmanNote, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/chairman-note/%d", id), nil), note)
	if err != nil {
		return nil, err
	}
	var updated model.ChairmanNote
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the ChairmanNote with the given ID.
func (s *ChairmanNoteService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/chairman-note/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

package easyvapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// AccountingPlanService manages all CRUD operations on the /accounting-plan endpoint.
// Accounting plans (Kontenpläne) define the structure of the chart of accounts.
type AccountingPlanService struct {
	client *Client
}

// defaultAccountingPlanQuery requests all fields defined in model.AccountingPlan.
var defaultAccountingPlanQuery = NewQuery().
	Fields("id", "name", "description")

// AccountingPlanListOptions holds all filter and pagination options for AccountingPlan
// list requests.
type AccountingPlanListOptions struct {
	ListOptions
}

// accountingPlanListParams converts opts into URL query parameters.
func accountingPlanListParams(opts *AccountingPlanListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultAccountingPlanQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultAccountingPlanQuery)
	return params
}

// List returns a lazy Iterator over all AccountingPlan records matching opts.
// Pages are fetched on-demand as iteration progresses.
func (s *AccountingPlanService) List(ctx context.Context, opts *AccountingPlanListOptions) *Iterator[model.AccountingPlan] {
	startURL := s.client.buildURL("/accounting-plan", accountingPlanListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.AccountingPlan, *string, error) {
		return fetchPage[model.AccountingPlan](s.client, ctx, pageURL)
	})
}

// ListAll fetches all AccountingPlan records and returns them as a slice.
func (s *AccountingPlanService) ListAll(ctx context.Context, opts *AccountingPlanListOptions) ([]model.AccountingPlan, error) {
	var all []model.AccountingPlan
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single AccountingPlan entry by its ID.
func (s *AccountingPlanService) Get(ctx context.Context, id int, query *Query) (*model.AccountingPlan, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/accounting-plan/%d", id), params)
	if err != nil {
		return nil, err
	}
	var ap model.AccountingPlan
	if err := s.client.decodeJSON(resp, &ap); err != nil {
		return nil, err
	}
	return &ap, nil
}

// Create creates a new AccountingPlan entry and returns the created record.
func (s *AccountingPlanService) Create(ctx context.Context, ap model.AccountingPlanCreate) (*model.AccountingPlan, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/accounting-plan", nil), ap)
	if err != nil {
		return nil, err
	}
	var created model.AccountingPlan
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a full update (PUT) to the AccountingPlan entry with the given ID.
func (s *AccountingPlanService) Update(ctx context.Context, id int, ap model.AccountingPlanCreate) (*model.AccountingPlan, error) {
	resp, err := s.client.do(ctx, "PUT", s.client.buildURL(fmt.Sprintf("/accounting-plan/%d", id), nil), ap)
	if err != nil {
		return nil, err
	}
	var updated model.AccountingPlan
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the AccountingPlan entry with the given ID.
func (s *AccountingPlanService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/accounting-plan/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

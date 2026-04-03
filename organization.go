package easyvapi

import (
	"context"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// OrganizationService provides access to the /organization endpoint.
// The organization is a singleton resource; there is no list or create.
type OrganizationService struct {
	client *Client
}

// Get retrieves the organization record.
func (s *OrganizationService) Get(ctx context.Context, query *Query) (*model.Organization, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, "/organization", params)
	if err != nil {
		return nil, err
	}
	var org model.Organization
	if err := s.client.decodeJSON(resp, &org); err != nil {
		return nil, err
	}
	return &org, nil
}

// Update applies a partial update (PATCH) to the organization.
func (s *OrganizationService) Update(ctx context.Context, org model.OrganizationCreate) (*model.Organization, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL("/organization", nil), org)
	if err != nil {
		return nil, err
	}
	var updated model.Organization
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

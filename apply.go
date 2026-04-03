package easyvapi

import (
	"context"
)

// ApplyService provides access to the /apply endpoint.
// Use this service to submit a membership application (Mitgliedsantrag).
type ApplyService struct {
	client *Client
}

// ApplyRequest holds the fields for a membership application submission.
type ApplyRequest struct {
	// ApplicationForm is the ID of the application form to use.
	ApplicationForm int `json:"applicationForm"`
	// Data holds the form field values keyed by element ID or field name.
	Data map[string]any `json:"data,omitempty"`
}

// Submit sends a membership application to the API.
// The response body is discarded; check the returned error for failure.
func (s *ApplyService) Submit(ctx context.Context, req ApplyRequest) error {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/apply", nil), req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

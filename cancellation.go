package easyvapi

import (
	"context"
)

// CancellationService provides access to the /cancellation endpoint.
// Use this service to submit a cancellation request to the easyVerein API.
type CancellationService struct {
	client *Client
}

// CancellationRequest holds the fields for a cancellation submission.
// The exact fields depend on the cancellation type; use a map or a custom
// struct depending on context.
type CancellationRequest struct {
	// InvoiceID is the ID of the invoice to cancel, if applicable.
	InvoiceID int `json:"invoice,omitempty"`
	// MemberID is the ID of the member whose membership is to be cancelled, if applicable.
	MemberID int `json:"member,omitempty"`
	// Reason is an optional free-text reason for the cancellation.
	Reason string `json:"reason,omitempty"`
}

// Submit sends a cancellation request to the API.
// The response body is discarded; check the returned error for failure.
//
// Example:
//
//	err := client.Cancellations.Submit(ctx, easyvapi.CancellationRequest{
//		InvoiceID: 12345,
//		Reason:    "Duplicate invoice",
//	})
func (s *CancellationService) Submit(ctx context.Context, req CancellationRequest) error {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/cancellation", nil), req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

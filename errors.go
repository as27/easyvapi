package easyvapi

import (
	"fmt"
	"time"
)

// APIError represents a non-2xx response returned by the easyVerein API.
// It contains the HTTP status code, a message, and optional detail information.
// Use [errors.As] to detect and handle this error type:
//
//	var apiErr *easyvapi.APIError
//	if errors.As(err, &apiErr) {
//		fmt.Printf("API error %d: %s\n", apiErr.StatusCode, apiErr.Message)
//	}
type APIError struct {
	// StatusCode is the HTTP status code returned by the server (e.g., 400, 404, 500).
	StatusCode int
	// Message is a short human-readable error description (e.g., "Bad Request", "Not Found").
	Message string
	// Detail contains additional detail provided by the API response body, if available.
	// May include validation errors or specific error information.
	Detail string
}

// Error implements the error interface for APIError.
func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("easyvapi: HTTP %d: %s – %s", e.StatusCode, e.Message, e.Detail)
	}
	return fmt.Sprintf("easyvapi: HTTP %d: %s", e.StatusCode, e.Message)
}

// RateLimitError is returned when the API rate limit (100 requests/minute) has been exceeded.
// The client automatically throttles when X-RateLimit-Remaining < 5, but this error
// is returned when a 429 (Too Many Requests) response is received.
// Callers should wait at least RetryAfter before retrying.
// Use [errors.As] to detect and handle this error type:
//
//	var rateLimitErr *easyvapi.RateLimitError
//	if errors.As(err, &rateLimitErr) {
//		fmt.Printf("Rate limited. Retry after %v\n", rateLimitErr.RetryAfter)
//		time.Sleep(rateLimitErr.RetryAfter)
//		// retry the operation
//	}
type RateLimitError struct {
	// RetryAfter is the recommended wait duration before the next request,
	// extracted from the Retry-After header or defaulting to 60 seconds.
	RetryAfter time.Duration
}

// Error implements the error interface for RateLimitError.
func (e *RateLimitError) Error() string {
	return fmt.Sprintf("easyvapi: rate limit exceeded, retry after %s", e.RetryAfter)
}

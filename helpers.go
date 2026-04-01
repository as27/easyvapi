package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const defaultLimit = 100

// applyListOptions encodes the common ListOptions fields (Limit, Query, Ordering, Search)
// into URL query parameters. If opts.Query is nil, defaultQuery is used instead to ensure
// only known fields are requested. Pass nil for defaultQuery to skip the query parameter
// entirely (for endpoints that do not support queries, like /member-group).
//
// This function is used internally by service List methods to build API requests.
// Callers should not use this directly.
func applyListOptions(params url.Values, opts ListOptions, defaultQuery *Query) {
	limit := opts.Limit
	if limit <= 0 {
		limit = defaultLimit
	}
	params.Set("limit", strconv.Itoa(limit))

	q := opts.Query
	if q == nil {
		q = defaultQuery
	}
	if qs := q.String(); qs != "" {
		params.Set("query", qs)
	}

	if opts.Ordering != "" {
		params.Set("ordering", opts.Ordering)
	}
	if opts.Search != "" {
		params.Set("search", opts.Search)
	}
}

// fetchPage performs a GET request against the given full URL (which already
// includes query parameters) and decodes the paginated response envelope.
// Returns the slice of results from this page, the URL for the next page
// (or nil if this is the last page), and any error that occurred.
//
// This function is used internally as the fetch function for Iterator instances.
// Callers should not use this directly.
//
// Example error handling:
//
//	results, nextURL, err := fetchPage[model.Member](client, ctx, url)
//	if err != nil {
//		// Handle rate limiting or API errors
//		var rateLimitErr *easyvapi.RateLimitError
//		if errors.As(err, &rateLimitErr) {
//			time.Sleep(rateLimitErr.RetryAfter)
//			// retry
//		}
//	}
func fetchPage[T any](c *Client, ctx context.Context, pageURL string) ([]T, *string, error) {
	// pageURL already contains all query parameters.
	resp, err := c.do(ctx, "GET", pageURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("easyvapi: fetch page: %w", err)
	}
	var paged pagedResponse[T]
	if err := c.decodeJSON(resp, &paged); err != nil {
		return nil, nil, err
	}
	return paged.Results, paged.Next, nil
}

package easyvapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// buildURL constructs the full request URL for the given path and optional
// query parameters. If params is nil an empty query string is used.
func (c *Client) buildURL(path string, params url.Values) string {
	base := strings.TrimRight(c.baseURL, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := base + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	return u
}

// get sends a GET request to path with the given query parameters and returns
// the raw *http.Response. The caller is responsible for closing the body.
func (c *Client) get(ctx context.Context, path string, params url.Values) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, c.buildURL(path, params), nil)
}

// do executes an HTTP request with the given method against rawURL. If body is
// non-nil it is JSON-encoded and sent as the request body.
//
// Behaviour:
//   - Sets the Authorization: Bearer header.
//   - If the response contains the "X-Token-Refresh-Needed" header the token is
//     refreshed automatically and the request is retried once.
//   - If X-RateLimit-Remaining drops below 5, the function sleeps briefly to
//     avoid exhausting the budget.
//   - Non-2xx responses are returned as *APIError.
func (c *Client) do(ctx context.Context, method, rawURL string, body any) (*http.Response, error) {
	resp, err := c.doOnce(ctx, method, rawURL, body)
	if err != nil {
		return nil, err
	}
	// If the server signals that the token should be refreshed, do it and retry.
	if val := resp.Header.Get("tokenRefreshNeeded"); val == "true" || val == "1" {
		resp.Body.Close()
		if err := c.refreshToken(ctx); err != nil {
			return nil, err
		}
		return c.doOnce(ctx, method, rawURL, body)
	}
	return resp, nil
}

// doOnce performs a single HTTP round-trip without any retry logic.
func (c *Client) doOnce(ctx context.Context, method, rawURL string, body any) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("easyvapi: marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("easyvapi: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("easyvapi: http request: %w", err)
	}

	// Throttle proactively if we are close to the rate limit.
	if remaining := resp.Header.Get("X-RateLimit-Remaining"); remaining != "" {
		if n, err := strconv.Atoi(remaining); err == nil && n < 5 {
			time.Sleep(10 * time.Second)
		}
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		resp.Body.Close()
		retryAfter := 60 * time.Second
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if secs, err := strconv.Atoi(ra); err == nil {
				retryAfter = time.Duration(secs) * time.Second
			}
		}
		return nil, &RateLimitError{RetryAfter: retryAfter}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		apiErr := &APIError{StatusCode: resp.StatusCode}
		// Read the raw body so we can try structured parsing and fall back
		// to the raw text if the structure doesn't match.
		rawBody, _ := io.ReadAll(resp.Body)
		var errBody struct {
			Detail  string `json:"detail"`
			Message string `json:"message"`
		}
		if json.Unmarshal(rawBody, &errBody) == nil {
			apiErr.Detail = errBody.Detail
			apiErr.Message = errBody.Message
		}
		// If neither field was populated, use the raw body as detail so
		// callers can see the actual server response.
		if apiErr.Message == "" && apiErr.Detail == "" && len(rawBody) > 0 {
			apiErr.Detail = string(rawBody)
		}
		if apiErr.Message == "" {
			apiErr.Message = http.StatusText(resp.StatusCode)
		}
		return nil, apiErr
	}

	return resp, nil
}

// refreshToken calls GET /refresh-token and stores the new token on the client.
func (c *Client) refreshToken(ctx context.Context) error {
	resp, err := c.doOnce(ctx, http.MethodGet, c.buildURL("/refresh-token", nil), nil)
	if err != nil {
		return fmt.Errorf("easyvapi: refresh token: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("easyvapi: decode refresh-token response: %w", err)
	}
	if result.Token != "" {
		c.token = result.Token
		if c.onTokenRefresh != nil {
			c.onTokenRefresh(result.Token)
		}
	}
	return nil
}

// decodeJSON decodes the JSON body of resp into v. When the client is in debug
// mode unknown fields cause a decode error, which helps catch API changes early.
func (c *Client) decodeJSON(resp *http.Response, v any) error {
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if c.debug {
		dec.DisallowUnknownFields()
	}
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("easyvapi: decode response: %w", err)
	}
	return nil
}

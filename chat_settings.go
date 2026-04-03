package easyvapi

import (
	"context"
	"net/url"

	"github.com/as27/easyvapi/model"
)

// ChatSettingsService provides access to the /chat-settings endpoint.
// Chat settings are a singleton resource; there is no list or create.
type ChatSettingsService struct {
	client *Client
}

// Get retrieves the current chat settings.
func (s *ChatSettingsService) Get(ctx context.Context, query *Query) (*model.ChatSettings, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, "/chat-settings", params)
	if err != nil {
		return nil, err
	}
	var cs model.ChatSettings
	if err := s.client.decodeJSON(resp, &cs); err != nil {
		return nil, err
	}
	return &cs, nil
}

// Update applies a partial update (PATCH) to the chat settings.
func (s *ChatSettingsService) Update(ctx context.Context, cs model.ChatSettingsCreate) (*model.ChatSettings, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL("/chat-settings", nil), cs)
	if err != nil {
		return nil, err
	}
	var updated model.ChatSettings
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

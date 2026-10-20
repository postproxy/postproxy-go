package postproxy

import (
	"context"
	"encoding/json"
	"net/http"
)

// APIKeysService manages profile-group-scoped API keys. Requires an
// account-wide API key.
type APIKeysService struct {
	client *Client
}

// List returns the account's API keys.
func (s *APIKeysService) List(ctx context.Context) (*ListResponse[APIKey], error) {
	data, err := s.client.request(ctx, http.MethodGet, "/api_keys")
	if err != nil {
		return nil, err
	}

	var result ListResponse[APIKey]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates an API key scoped to a profile group. The returned Token is
// shown only once.
func (s *APIKeysService) Create(ctx context.Context, name, profileGroupID string) (*APIKey, error) {
	body := map[string]string{"name": name, "profile_group_id": profileGroupID}

	data, err := s.client.request(ctx, http.MethodPost, "/api_keys", withJSON(body))
	if err != nil {
		return nil, err
	}

	var result APIKey
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes an API key by ID.
func (s *APIKeysService) Delete(ctx context.Context, id string) (*DeleteResponse, error) {
	data, err := s.client.request(ctx, http.MethodDelete, "/api_keys/"+id)
	if err != nil {
		return nil, err
	}

	var result DeleteResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

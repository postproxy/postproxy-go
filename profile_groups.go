package postproxy

import (
	"context"
	"encoding/json"
	"net/http"
)

// ProfileGroupsService handles communication with the profile groups related methods of the PostProxy API.
type ProfileGroupsService struct {
	client *Client
}

// List returns all profile groups.
func (s *ProfileGroupsService) List(ctx context.Context) (*ListResponse[ProfileGroup], error) {
	data, err := s.client.request(ctx, http.MethodGet, "/profile_groups")
	if err != nil {
		return nil, err
	}

	var result ListResponse[ProfileGroup]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a single profile group by ID.
func (s *ProfileGroupsService) Get(ctx context.Context, id string) (*ProfileGroup, error) {
	data, err := s.client.request(ctx, http.MethodGet, "/profile_groups/"+id)
	if err != nil {
		return nil, err
	}

	var result ProfileGroup
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new profile group.
func (s *ProfileGroupsService) Create(ctx context.Context, name string) (*ProfileGroup, error) {
	body := map[string]string{"name": name}

	data, err := s.client.request(ctx, http.MethodPost, "/profile_groups", withJSON(body))
	if err != nil {
		return nil, err
	}

	var result ProfileGroup
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a profile group by ID.
func (s *ProfileGroupsService) Delete(ctx context.Context, id string) (*DeleteResponse, error) {
	data, err := s.client.request(ctx, http.MethodDelete, "/profile_groups/"+id)
	if err != nil {
		return nil, err
	}

	var result DeleteResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// InitializeConnection starts an OAuth connection flow for a platform.
// BlueSky and Telegram use dedicated helpers (ConnectBluesky, ConnectTelegram).
func (s *ProfileGroupsService) InitializeConnection(ctx context.Context, id string, platform Platform, redirectURL string) (*ConnectionResponse, error) {
	body := map[string]string{
		"platform":     string(platform),
		"redirect_url": redirectURL,
	}

	data, err := s.client.request(ctx, http.MethodPost, "/profile_groups/"+id+"/initialize_connection", withJSON(body))
	if err != nil {
		return nil, err
	}

	var result ConnectionResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BlueskyConnectOptions are the credentials used by ConnectBluesky.
type BlueskyConnectOptions struct {
	Identifier  string
	AppPassword string
}

// ConnectBluesky creates a Bluesky profile in the group synchronously using
// the user's handle and an app password (https://bsky.app/settings/app-passwords).
func (s *ProfileGroupsService) ConnectBluesky(ctx context.Context, id string, opts BlueskyConnectOptions) (*BlueskyConnectionResponse, error) {
	body := map[string]string{
		"platform":     string(PlatformBluesky),
		"identifier":   opts.Identifier,
		"app_password": opts.AppPassword,
	}

	data, err := s.client.request(ctx, http.MethodPost, "/profile_groups/"+id+"/initialize_connection", withJSON(body))
	if err != nil {
		return nil, err
	}

	var result BlueskyConnectionResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TelegramConnectOptions hold the bot token used by ConnectTelegram.
type TelegramConnectOptions struct {
	BotToken string
}

// ConnectTelegram registers a Telegram bot for the profile group. After it returns,
// poll ProfilesService.Placements until non-empty — the bot must be added as
// administrator to a Telegram channel before channels appear.
func (s *ProfileGroupsService) ConnectTelegram(ctx context.Context, id string, opts TelegramConnectOptions) (*TelegramConnectionResponse, error) {
	body := map[string]string{
		"platform":  string(PlatformTelegram),
		"bot_token": opts.BotToken,
	}

	data, err := s.client.request(ctx, http.MethodPost, "/profile_groups/"+id+"/initialize_connection", withJSON(body))
	if err != nil {
		return nil, err
	}

	var result TelegramConnectionResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

package postproxy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// WebhooksService handles communication with the webhooks related methods of the PostProxy API.
type WebhooksService struct {
	client *Client
}

// List returns all webhooks.
func (s *WebhooksService) List(ctx context.Context) (*ListResponse[Webhook], error) {
	data, err := s.client.request(ctx, http.MethodGet, "/webhooks")
	if err != nil {
		return nil, err
	}

	var result ListResponse[Webhook]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a single webhook by ID.
func (s *WebhooksService) Get(ctx context.Context, id string) (*Webhook, error) {
	data, err := s.client.request(ctx, http.MethodGet, "/webhooks/"+id)
	if err != nil {
		return nil, err
	}

	var result Webhook
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new webhook.
func (s *WebhooksService) Create(ctx context.Context, webhookURL string, events []string, opts *WebhookCreateOptions) (*Webhook, error) {
	payload := map[string]any{
		"url":    webhookURL,
		"events": events,
	}
	if opts != nil && opts.Description != nil {
		payload["description"] = *opts.Description
	}

	data, err := s.client.request(ctx, http.MethodPost, "/webhooks", withJSON(payload))
	if err != nil {
		return nil, err
	}

	var result Webhook
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates an existing webhook.
func (s *WebhooksService) Update(ctx context.Context, id string, opts *WebhookUpdateOptions) (*Webhook, error) {
	payload := map[string]any{}
	if opts != nil {
		if opts.URL != nil {
			payload["url"] = *opts.URL
		}
		if opts.Events != nil {
			payload["events"] = opts.Events
		}
		if opts.Enabled != nil {
			payload["enabled"] = *opts.Enabled
		}
		if opts.Description != nil {
			payload["description"] = *opts.Description
		}
	}

	data, err := s.client.request(ctx, "PATCH", "/webhooks/"+id, withJSON(payload))
	if err != nil {
		return nil, err
	}

	var result Webhook
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a webhook by ID.
func (s *WebhooksService) Delete(ctx context.Context, id string) (*DeleteResponse, error) {
	data, err := s.client.request(ctx, http.MethodDelete, "/webhooks/"+id)
	if err != nil {
		return nil, err
	}

	var result DeleteResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Deliveries returns delivery attempts for a webhook.
func (s *WebhooksService) Deliveries(ctx context.Context, id string, opts *WebhookDeliveryListOptions) (*PaginatedResponse[WebhookDelivery], error) {
	var reqOpts []requestOption
	if opts != nil {
		params := url.Values{}
		if opts.Page != nil {
			params.Set("page", strconv.Itoa(*opts.Page))
		}
		if opts.PerPage != nil {
			params.Set("per_page", strconv.Itoa(*opts.PerPage))
		}
		if len(params) > 0 {
			reqOpts = append(reqOpts, withParams(params))
		}
	}

	data, err := s.client.request(ctx, http.MethodGet, "/webhooks/"+id+"/deliveries", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result PaginatedResponse[WebhookDelivery]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

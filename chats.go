package postproxy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// ChatsService handles communication with the direct-message chat related methods of the PostProxy API.
type ChatsService struct {
	client *Client
}

// List returns a paginated list of chats for a profile.
func (s *ChatsService) List(ctx context.Context, profileID string, opts *ChatListOptions) (*PaginatedResponse[Chat], error) {
	var reqOpts []requestOption
	if opts != nil {
		params := url.Values{}
		if opts.Page != nil {
			params.Set("page", strconv.Itoa(*opts.Page))
		}
		if opts.PerPage != nil {
			params.Set("per_page", strconv.Itoa(*opts.PerPage))
		}
		if opts.Before != nil {
			params.Set("before", *opts.Before)
		}
		if opts.After != nil {
			params.Set("after", *opts.After)
		}
		if len(params) > 0 {
			reqOpts = append(reqOpts, withParams(params))
		}
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+profileID+"/chats", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result PaginatedResponse[Chat]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates (or finds) a chat with a participant for a profile.
func (s *ChatsService) Create(ctx context.Context, profileID, participantExternalID string, opts *ChatCreateOptions) (*Chat, error) {
	body := map[string]any{"participant_external_id": participantExternalID}
	var reqOpts []requestOption
	if opts != nil {
		if opts.ParticipantUsername != nil {
			body["participant_username"] = *opts.ParticipantUsername
		}
		if opts.ParticipantName != nil {
			body["participant_name"] = *opts.ParticipantName
		}
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}
	reqOpts = append(reqOpts, withJSON(body))

	data, err := s.client.request(ctx, http.MethodPost, "/profiles/"+profileID+"/chats", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Chat
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a single chat by ID.
func (s *ChatsService) Get(ctx context.Context, chatID string, opts *RequestOptions) (*Chat, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/chats/"+chatID, reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Chat
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Archive archives a chat (Bluesky only).
func (s *ChatsService) Archive(ctx context.Context, chatID string, opts *RequestOptions) (*Chat, error) {
	return s.archiveAction(ctx, http.MethodPost, chatID, opts)
}

// Unarchive unarchives a chat (Bluesky only).
func (s *ChatsService) Unarchive(ctx context.Context, chatID string, opts *RequestOptions) (*Chat, error) {
	return s.archiveAction(ctx, http.MethodDelete, chatID, opts)
}

func (s *ChatsService) archiveAction(ctx context.Context, method, chatID string, opts *RequestOptions) (*Chat, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, method, "/chats/"+chatID+"/archive", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Chat
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

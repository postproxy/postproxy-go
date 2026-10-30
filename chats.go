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
		if opts.PlacementID != nil {
			body["placement_id"] = *opts.PlacementID
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
	return s.chatAction(ctx, http.MethodPost, chatID, "archive", opts)
}

// Unarchive unarchives a chat (Bluesky only).
func (s *ChatsService) Unarchive(ctx context.Context, chatID string, opts *RequestOptions) (*Chat, error) {
	return s.chatAction(ctx, http.MethodDelete, chatID, "archive", opts)
}

// MarkRead marks the latest inbound message as read (shows the read receipt).
func (s *ChatsService) MarkRead(ctx context.Context, chatID string, opts *RequestOptions) (*Chat, error) {
	return s.chatAction(ctx, http.MethodPost, chatID, "mark_read", opts)
}

// StartTyping shows the typing indicator. Facebook, Instagram and WhatsApp only.
func (s *ChatsService) StartTyping(ctx context.Context, chatID string, opts *RequestOptions) (*AcceptedResponse, error) {
	return s.typingAction(ctx, chatID, "start_typing", opts)
}

// StopTyping hides the typing indicator. Facebook, Instagram and WhatsApp only.
func (s *ChatsService) StopTyping(ctx context.Context, chatID string, opts *RequestOptions) (*AcceptedResponse, error) {
	return s.typingAction(ctx, chatID, "stop_typing", opts)
}

// --- Thread control (Facebook & Instagram) ---
//
// Meta gives one app at a time control of a conversation when several share a
// Page or Instagram inbox. These return the chat, whose
// Metadata["thread_control"] holds the owner Meta reports back.

// PassThreadControl hands the conversation to another Meta app, e.g. Page
// Inbox "263902037430900" or Instagram Inbox "1217981644879628".
func (s *ChatsService) PassThreadControl(ctx context.Context, chatID, targetAppID string, opts *ThreadControlOptions) (*Chat, error) {
	return s.threadControl(ctx, chatID, "pass_thread_control", map[string]any{"target_app_id": targetAppID}, opts)
}

// TakeThreadControl takes the conversation from whichever app holds it. Needs
// Conversation Routing; Instagram must be connected via Facebook.
func (s *ChatsService) TakeThreadControl(ctx context.Context, chatID string, opts *ThreadControlOptions) (*Chat, error) {
	return s.threadControl(ctx, chatID, "take_thread_control", map[string]any{}, opts)
}

// ReleaseThreadControl gives up control; the conversation goes idle.
func (s *ChatsService) ReleaseThreadControl(ctx context.Context, chatID string, opts *ThreadControlOptions) (*Chat, error) {
	return s.threadControl(ctx, chatID, "release_thread_control", map[string]any{}, opts)
}

// RequestThreadControl asks the current owner to pass the conversation.
// Nothing changes until it does — watch for EventChatThreadControlChanged.
func (s *ChatsService) RequestThreadControl(ctx context.Context, chatID string, opts *ThreadControlOptions) (*Chat, error) {
	return s.threadControl(ctx, chatID, "request_thread_control", map[string]any{}, opts)
}

// ThreadOwner asks Meta who controls the conversation right now.
func (s *ChatsService) ThreadOwner(ctx context.Context, chatID string, opts *RequestOptions) (*ThreadOwner, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/chats/"+chatID+"/thread_owner", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result ThreadOwner
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ChatsService) threadControl(ctx context.Context, chatID, action string, body map[string]any, opts *ThreadControlOptions) (*Chat, error) {
	var reqOpts []requestOption
	if opts != nil {
		if opts.Metadata != nil {
			body["metadata"] = *opts.Metadata
		}
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}
	reqOpts = append(reqOpts, withJSON(body))

	data, err := s.client.request(ctx, http.MethodPost, "/chats/"+chatID+"/"+action, reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Chat
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ChatsService) typingAction(ctx context.Context, chatID, action string, opts *RequestOptions) (*AcceptedResponse, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodPost, "/chats/"+chatID+"/"+action, reqOpts...)
	if err != nil {
		return nil, err
	}

	var result AcceptedResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ChatsService) chatAction(ctx context.Context, method, chatID, action string, opts *RequestOptions) (*Chat, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, method, "/chats/"+chatID+"/"+action, reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Chat
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

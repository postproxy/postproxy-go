package postproxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

// MessagesService handles communication with the direct-message related methods of the PostProxy API.
type MessagesService struct {
	client *Client
}

// List returns a paginated list of messages in a chat.
func (s *MessagesService) List(ctx context.Context, chatID string, opts *MessageListOptions) (*PaginatedResponse[Message], error) {
	var reqOpts []requestOption
	if opts != nil {
		params := url.Values{}
		if opts.Page != nil {
			params.Set("page", strconv.Itoa(*opts.Page))
		}
		if opts.PerPage != nil {
			params.Set("per_page", strconv.Itoa(*opts.PerPage))
		}
		if opts.Direction != nil {
			params.Set("direction", string(*opts.Direction))
		}
		if opts.Status != nil {
			params.Set("status", string(*opts.Status))
		}
		if len(params) > 0 {
			reqOpts = append(reqOpts, withParams(params))
		}
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/chats/"+chatID+"/messages", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result PaginatedResponse[Message]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Send sends a message in a chat. If MediaFiles is set in opts, the request uses multipart/form-data.
func (s *MessagesService) Send(ctx context.Context, chatID string, opts *MessageSendOptions) (*Message, error) {
	var reqOpts []requestOption

	if opts != nil && len(opts.MediaFiles) > 0 {
		formOpts, err := s.buildFormData(opts)
		if err != nil {
			return nil, err
		}
		reqOpts = append(reqOpts, formOpts...)
	} else if opts != nil {
		reqOpts = append(reqOpts, s.buildJSON(opts))
	}

	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodPost, "/chats/"+chatID+"/messages", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Message
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *MessagesService) buildJSON(opts *MessageSendOptions) requestOption {
	payload := map[string]any{}
	if opts.Body != nil {
		payload["body"] = *opts.Body
	}
	if opts.Media != nil {
		payload["media"] = opts.Media
	}
	if opts.Tag != nil {
		payload["tag"] = *opts.Tag
	}
	if opts.ReplyToExternalID != nil {
		payload["reply_to_external_id"] = *opts.ReplyToExternalID
	}
	if opts.ReplyMarkup != nil {
		payload["reply_markup"] = opts.ReplyMarkup
	}
	return withJSON(payload)
}

func (s *MessagesService) buildFormData(opts *MessageSendOptions) ([]requestOption, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	if opts.Body != nil {
		_ = w.WriteField("body", *opts.Body)
	}
	if opts.Tag != nil {
		_ = w.WriteField("tag", *opts.Tag)
	}
	if opts.ReplyToExternalID != nil {
		_ = w.WriteField("reply_to_external_id", *opts.ReplyToExternalID)
	}

	for _, u := range opts.Media {
		_ = w.WriteField("media[]", u)
	}

	for _, filePath := range opts.MediaFiles {
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read media file %s: %w", filePath, err)
		}

		contentType := mimeTypeFromPath(filePath)
		filename := filepath.Base(filePath)

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="media[]"; filename=%q`, filename))
		h.Set("Content-Type", contentType)
		part, err := w.CreatePart(h)
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}
		if _, err := part.Write(fileData); err != nil {
			return nil, fmt.Errorf("failed to write file data: %w", err)
		}
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return []requestOption{withFormData(&buf, w.FormDataContentType())}, nil
}

// Get returns a single message by ID.
func (s *MessagesService) Get(ctx context.Context, messageID string, opts *RequestOptions) (*Message, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/messages/"+messageID, reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Message
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Edit edits an outbound message (Telegram only).
func (s *MessagesService) Edit(ctx context.Context, messageID string, opts *MessageEditOptions) (*Message, error) {
	payload := map[string]any{}
	var reqOpts []requestOption
	if opts != nil {
		if opts.Body != nil {
			payload["body"] = *opts.Body
		}
		if opts.ReplyMarkup != nil {
			payload["reply_markup"] = opts.ReplyMarkup
		}
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}
	reqOpts = append(reqOpts, withJSON(payload))

	data, err := s.client.request(ctx, http.MethodPatch, "/messages/"+messageID, reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Message
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// React adds a reaction to a message. Leave Reaction unset to let the server default to "love".
func (s *MessagesService) React(ctx context.Context, messageID string, opts *MessageReactOptions) (*Message, error) {
	payload := map[string]any{}
	var reqOpts []requestOption
	if opts != nil {
		if opts.Reaction != nil {
			payload["reaction"] = *opts.Reaction
		}
		if opts.Emoji != nil {
			payload["emoji"] = *opts.Emoji
		}
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}
	if len(payload) > 0 {
		reqOpts = append(reqOpts, withJSON(payload))
	}

	data, err := s.client.request(ctx, http.MethodPost, "/messages/"+messageID+"/react", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Message
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Unreact removes a reaction from a message.
func (s *MessagesService) Unreact(ctx context.Context, messageID string, opts *RequestOptions) (*Message, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodDelete, "/messages/"+messageID+"/unreact", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Message
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

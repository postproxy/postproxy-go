package postproxy

import (
	"context"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var mockInbound = map[string]any{
	"id":                    "msg_111",
	"chat_id":               "chat_xyz789",
	"external_id":           "mid.abc123",
	"direction":             "inbound",
	"body":                  "Hey, do you ship internationally?",
	"status":                "received",
	"tag":                   nil,
	"error_message":         nil,
	"platform_data":         nil,
	"external_posted_at":    "2026-05-31T14:02:00.000Z",
	"external_delivered_at": nil,
	"external_read_at":      nil,
	"external_edited_at":    nil,
	"reactions": []map[string]any{
		{"sender_external_id": "psid_123", "emoji": "❤️", "reaction": "love", "at": "2026-05-31T14:04:00.000Z"},
	},
	"attachments":    []any{},
	"is_unsupported": false,
	"created_at":     "2026-05-31T14:02:01.000Z",
}

var mockOutbound = map[string]any{
	"id":                 "msg_222",
	"chat_id":            "chat_xyz789",
	"external_id":        nil,
	"direction":          "outbound",
	"body":               "Yes, we ship worldwide!",
	"status":             "pending",
	"tag":                nil,
	"error_message":      nil,
	"platform_data":      nil,
	"external_posted_at": nil,
	"attachments":        []any{},
	"reactions":          []any{},
	"is_unsupported":     false,
	"created_at":         "2026-05-31T15:30:05.000Z",
}

func TestMessagesList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chats/chat_xyz789/messages" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("direction") != "inbound" {
			t.Errorf("expected direction=inbound, got %s", r.URL.Query().Get("direction"))
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"total": 1, "page": 0, "per_page": 20,
			"data": []map[string]any{mockInbound},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	dir := MessageDirectionInbound
	result, err := c.Messages.List(context.Background(), "chat_xyz789", &MessageListOptions{Direction: &dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 message, got %d", len(result.Data))
	}
	msg := result.Data[0]
	if msg.Direction != MessageDirectionInbound {
		t.Errorf("expected inbound, got %s", msg.Direction)
	}
	if len(msg.Reactions) != 1 || msg.Reactions[0].Reaction == nil || *msg.Reactions[0].Reaction != "love" {
		t.Errorf("unexpected reactions: %+v", msg.Reactions)
	}
}

func TestMessagesSendText(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/chats/chat_xyz789/messages" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		_ = json.Unmarshal(body, &req)
		if req["body"] != "Yes, we ship worldwide!" {
			t.Errorf("unexpected body: %v", req["body"])
		}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(mockOutbound)
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	text := "Yes, we ship worldwide!"
	msg, err := c.Messages.Send(context.Background(), "chat_xyz789", &MessageSendOptions{Body: &text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.ID != "msg_222" {
		t.Errorf("expected id msg_222, got %s", msg.ID)
	}
	if msg.Status != MessageStatusPending {
		t.Errorf("expected status pending, got %s", msg.Status)
	}
}

func TestMessagesSendWithTag(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		_ = json.Unmarshal(body, &req)
		if req["tag"] != "HUMAN_AGENT" {
			t.Errorf("expected tag HUMAN_AGENT, got %v", req["tag"])
		}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(mockOutbound)
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	text := "Following up."
	tag := "HUMAN_AGENT"
	_, err := c.Messages.Send(context.Background(), "chat_xyz789", &MessageSendOptions{Body: &text, Tag: &tag})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMessagesSendQuickReplies(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		_ = json.Unmarshal(body, &req)

		qrs, ok := req["quick_replies"].([]any)
		if !ok || len(qrs) != 2 {
			t.Fatalf("expected 2 quick_replies, got %v", req["quick_replies"])
		}
		first := qrs[0].(map[string]any)
		if first["title"] != "Track order" || first["payload"] != "TRACK" {
			t.Errorf("unexpected quick reply: %v", first)
		}
		// ContentType is omitempty, so an unset one must not be serialized.
		if _, present := first["content_type"]; present {
			t.Errorf("expected content_type to be omitted, got %v", first["content_type"])
		}

		w.WriteHeader(http.StatusAccepted)
		out := map[string]any{}
		for k, v := range mockOutbound {
			out[k] = v
		}
		out["quick_replies"] = []map[string]any{
			{"content_type": "text", "title": "Track order", "payload": "TRACK"},
		}
		_ = json.NewEncoder(w).Encode(out)
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	text := "What can I help with?"
	msg, err := c.Messages.Send(context.Background(), "chat_xyz789", &MessageSendOptions{
		Body: &text,
		QuickReplies: []QuickReply{
			{Title: "Track order", Payload: "TRACK"},
			{Title: "Talk to support", Payload: "HELP"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(msg.QuickReplies) != 1 || msg.QuickReplies[0].ContentType != "text" {
		t.Errorf("expected echoed quick reply with content_type text, got %+v", msg.QuickReplies)
	}
}

func TestMessagesSendButtonsWithCard(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		_ = json.Unmarshal(body, &req)

		btns, ok := req["buttons"].([]any)
		if !ok || len(btns) != 2 {
			t.Fatalf("expected 2 buttons, got %v", req["buttons"])
		}
		web := btns[0].(map[string]any)
		if web["type"] != "web_url" || web["url"] != "https://shop.example.com" {
			t.Errorf("unexpected web_url button: %v", web)
		}
		// Payload is omitempty and unset on a web_url button.
		if _, present := web["payload"]; present {
			t.Errorf("expected payload to be omitted on web_url button, got %v", web["payload"])
		}
		post := btns[1].(map[string]any)
		if post["type"] != "postback" || post["payload"] != "CANCEL:123" {
			t.Errorf("unexpected postback button: %v", post)
		}

		card, ok := req["card"].(map[string]any)
		if !ok || card["subtitle"] != "Arriving Friday" {
			t.Fatalf("unexpected card: %v", req["card"])
		}
		action, ok := card["default_action"].(map[string]any)
		if !ok || action["url"] != "https://shop.example.com" {
			t.Errorf("unexpected default_action: %v", card["default_action"])
		}

		w.WriteHeader(http.StatusAccepted)
		out := map[string]any{}
		for k, v := range mockOutbound {
			out[k] = v
		}
		out["buttons"] = []map[string]any{
			{"type": "web_url", "title": "Track", "url": "https://shop.example.com"},
		}
		out["card"] = map[string]any{"subtitle": "Arriving Friday"}
		_ = json.NewEncoder(w).Encode(out)
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	text := "Your order shipped"
	msg, err := c.Messages.Send(context.Background(), "chat_xyz789", &MessageSendOptions{
		Body: &text,
		Buttons: []MessageButton{
			{Type: "web_url", Title: "Track", URL: "https://shop.example.com"},
			{Type: "postback", Title: "Cancel", Payload: "CANCEL:123"},
		},
		Card: &MessageCard{
			Subtitle:      "Arriving Friday",
			DefaultAction: &CardDefaultAction{Type: "web_url", URL: "https://shop.example.com"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(msg.Buttons) != 1 || msg.Buttons[0].URL != "https://shop.example.com" {
		t.Errorf("unexpected echoed buttons: %+v", msg.Buttons)
	}
	if msg.Card == nil || msg.Card.Subtitle != "Arriving Friday" {
		t.Errorf("unexpected echoed card: %+v", msg.Card)
	}
}

func TestMessagesTappedAction(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out := map[string]any{}
		for k, v := range mockInbound {
			out[k] = v
		}
		out["id"] = "msg_333"
		out["body"] = "Track order"
		out["tapped_action"] = map[string]any{
			"kind": "quick_reply", "payload": "TRACK", "title": "Track order",
		}
		_ = json.NewEncoder(w).Encode(out)
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	msg, err := c.Messages.Get(context.Background(), "msg_333", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.TappedAction == nil {
		t.Fatal("expected tapped_action to be set")
	}
	if msg.TappedAction.Kind != TappedActionQuickReply {
		t.Errorf("expected kind quick_reply, got %s", msg.TappedAction.Kind)
	}
	if msg.TappedAction.Payload != "TRACK" {
		t.Errorf("expected payload TRACK, got %s", msg.TappedAction.Payload)
	}
	if msg.TappedAction.Title == nil || *msg.TappedAction.Title != "Track order" {
		t.Errorf("unexpected title: %v", msg.TappedAction.Title)
	}
}

func TestMessagesSendMediaURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		_ = json.Unmarshal(body, &req)
		media, _ := req["media"].([]any)
		if len(media) != 1 || media[0] != "https://cdn.example.com/photo.png" {
			t.Errorf("unexpected media: %v", req["media"])
		}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(mockOutbound)
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	_, err := c.Messages.Send(context.Background(), "chat_xyz789", &MessageSendOptions{Media: []string{"https://cdn.example.com/photo.png"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMessagesSendMediaFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "photo.png")
	if err := os.WriteFile(tmp, []byte("fakepngdata"), 0o644); err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("expected multipart, got %s", ct)
		}
		_, params, _ := mime.ParseMediaType(ct)
		mr := multipart.NewReader(r.Body, params["boundary"])
		foundFile := false
		for {
			part, err := mr.NextPart()
			if err != nil {
				break
			}
			if part.FormName() == "media[]" && part.FileName() != "" {
				foundFile = true
			}
		}
		if !foundFile {
			t.Errorf("expected a media[] file part")
		}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(mockOutbound)
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	_, err := c.Messages.Send(context.Background(), "chat_xyz789", &MessageSendOptions{MediaFiles: []string{tmp}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMessagesGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/messages/msg_111" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockInbound)
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	msg, err := c.Messages.Get(context.Background(), "msg_111", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.ID != "msg_111" {
		t.Errorf("expected id msg_111, got %s", msg.ID)
	}
}

func TestMessagesEdit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/messages/msg_222" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		_ = json.Unmarshal(body, &req)
		if req["body"] != "Updated" {
			t.Errorf("unexpected body: %v", req["body"])
		}
		w.WriteHeader(http.StatusOK)
		edited := map[string]any{}
		for k, v := range mockOutbound {
			edited[k] = v
		}
		edited["body"] = "Updated"
		_ = json.NewEncoder(w).Encode(edited)
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	text := "Updated"
	msg, err := c.Messages.Edit(context.Background(), "msg_222", &MessageEditOptions{Body: &text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Body == nil || *msg.Body != "Updated" {
		t.Errorf("expected body Updated, got %v", msg.Body)
	}
}

func TestMessagesReact(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/messages/msg_111/react" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		_ = json.Unmarshal(body, &req)
		if req["reaction"] != "love" {
			t.Errorf("expected reaction love, got %v", req["reaction"])
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockInbound)
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	reaction := "love"
	emoji := "❤️"
	msg, err := c.Messages.React(context.Background(), "msg_111", &MessageReactOptions{Reaction: &reaction, Emoji: &emoji})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.ID != "msg_111" {
		t.Errorf("expected id msg_111, got %s", msg.ID)
	}
}

func TestMessagesUnreact(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/messages/msg_111/unreact" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockInbound)
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	msg, err := c.Messages.Unreact(context.Background(), "msg_111", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.ID != "msg_111" {
		t.Errorf("expected id msg_111, got %s", msg.ID)
	}
}

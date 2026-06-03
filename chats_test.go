package postproxy

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mockChat = map[string]any{
	"id":                       "chat_xyz789",
	"profile_id":               "prof_abc123",
	"platform":                 "instagram",
	"participant_external_id":  "igsid_8675309",
	"participant_username":     "jane_doe",
	"participant_name":         "Jane Doe",
	"participant_avatar_url":   "https://storage.postproxy.dev/x.jpg",
	"external_conversation_id": nil,
	"last_inbound_at":          "2026-05-31T14:02:00.000Z",
	"last_outbound_at":         "2026-05-31T15:10:00.000Z",
	"last_message_at":          "2026-05-31T15:10:00.000Z",
	"metadata":                 map[string]any{"is_verified_user": false, "follower_count": 482},
	"created_at":               "2026-04-12T08:00:00.000Z",
}

func TestChatsList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/profiles/prof_abc123/chats" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("per_page") != "20" {
			t.Errorf("expected per_page=20, got %s", r.URL.Query().Get("per_page"))
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"total": 1, "page": 0, "per_page": 20,
			"data": []map[string]any{mockChat},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	perPage := 20
	result, err := c.Chats.List(context.Background(), "prof_abc123", &ChatListOptions{PerPage: &perPage})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Total)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 chat, got %d", len(result.Data))
	}
	chat := result.Data[0]
	if chat.ID != "chat_xyz789" {
		t.Errorf("expected id chat_xyz789, got %s", chat.ID)
	}
	if chat.ParticipantUsername == nil || *chat.ParticipantUsername != "jane_doe" {
		t.Errorf("unexpected participant_username: %v", chat.ParticipantUsername)
	}
	if got, _ := chat.Metadata["follower_count"].(float64); got != 482 {
		t.Errorf("expected follower_count 482, got %v", chat.Metadata["follower_count"])
	}
}

func TestChatsCreate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/profiles/prof_abc123/chats" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		_ = json.Unmarshal(body, &req)
		if req["participant_external_id"] != "igsid_8675309" {
			t.Errorf("expected participant_external_id igsid_8675309, got %v", req["participant_external_id"])
		}
		if req["participant_username"] != "jane_doe" {
			t.Errorf("expected participant_username jane_doe, got %v", req["participant_username"])
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockChat)
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	username := "jane_doe"
	chat, err := c.Chats.Create(context.Background(), "prof_abc123", "igsid_8675309", &ChatCreateOptions{ParticipantUsername: &username})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.ID != "chat_xyz789" {
		t.Errorf("expected id chat_xyz789, got %s", chat.ID)
	}
}

func TestChatsGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chats/chat_xyz789" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockChat)
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	chat, err := c.Chats.Get(context.Background(), "chat_xyz789", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.ID != "chat_xyz789" {
		t.Errorf("expected id chat_xyz789, got %s", chat.ID)
	}
}

func TestChatsArchive(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/chats/chat_xyz789/archive" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		archived := map[string]any{}
		for k, v := range mockChat {
			archived[k] = v
		}
		archived["archived"] = true
		_ = json.NewEncoder(w).Encode(archived)
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	chat, err := c.Chats.Archive(context.Background(), "chat_xyz789", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.Archived == nil || !*chat.Archived {
		t.Errorf("expected archived true, got %v", chat.Archived)
	}
}

func TestChatsUnarchive(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/chats/chat_xyz789/archive" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		archived := map[string]any{}
		for k, v := range mockChat {
			archived[k] = v
		}
		archived["archived"] = false
		_ = json.NewEncoder(w).Encode(archived)
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	chat, err := c.Chats.Unarchive(context.Background(), "chat_xyz789", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.Archived == nil || *chat.Archived {
		t.Errorf("expected archived false, got %v", chat.Archived)
	}
}

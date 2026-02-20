package postproxy

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProfileGroupsList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/profile_groups" {
			t.Errorf("expected path /api/profile_groups, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ListResponse[ProfileGroup]{
			Data: []ProfileGroup{
				{ID: "pg-1", Name: "Default", ProfilesCount: 3},
			},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.ProfileGroups.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 group, got %d", len(result.Data))
	}
	if result.Data[0].Name != "Default" {
		t.Errorf("expected name %q, got %q", "Default", result.Data[0].Name)
	}
}

func TestProfileGroupsGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/profile_groups/pg-123" {
			t.Errorf("expected path /api/profile_groups/pg-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ProfileGroup{ID: "pg-123", Name: "My Group", ProfilesCount: 5})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	group, err := c.ProfileGroups.Get(context.Background(), "pg-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.ID != "pg-123" {
		t.Errorf("expected ID %q, got %q", "pg-123", group.ID)
	}
	if group.ProfilesCount != 5 {
		t.Errorf("expected profiles count 5, got %d", group.ProfilesCount)
	}
}

func TestProfileGroupsCreate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/profile_groups" {
			t.Errorf("expected path /api/profile_groups, got %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]string
		_ = json.Unmarshal(body, &payload)
		if payload["name"] != "New Group" {
			t.Errorf("expected name %q, got %q", "New Group", payload["name"])
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ProfileGroup{ID: "pg-new", Name: "New Group", ProfilesCount: 0})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	group, err := c.ProfileGroups.Create(context.Background(), "New Group")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "New Group" {
		t.Errorf("expected name %q, got %q", "New Group", group.Name)
	}
}

func TestProfileGroupsDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/profile_groups/pg-del" {
			t.Errorf("expected path /api/profile_groups/pg-del, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(DeleteResponse{Deleted: true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.ProfileGroups.Delete(context.Background(), "pg-del")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Deleted {
		t.Error("expected deleted=true")
	}
}

func TestProfileGroupsInitializeConnection(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/profile_groups/pg-1/initialize_connection" {
			t.Errorf("expected path /api/profile_groups/pg-1/initialize_connection, got %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]string
		_ = json.Unmarshal(body, &payload)
		if payload["platform"] != "instagram" {
			t.Errorf("expected platform %q, got %q", "instagram", payload["platform"])
		}
		if payload["redirect_url"] != "https://example.com/callback" {
			t.Errorf("expected redirect_url %q, got %q", "https://example.com/callback", payload["redirect_url"])
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ConnectionResponse{
			URL:     "https://auth.postproxy.dev/connect/abc",
			Success: true,
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.ProfileGroups.InitializeConnection(
		context.Background(),
		"pg-1",
		PlatformInstagram,
		"https://example.com/callback",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("expected success=true")
	}
	if result.URL != "https://auth.postproxy.dev/connect/abc" {
		t.Errorf("expected URL %q, got %q", "https://auth.postproxy.dev/connect/abc", result.URL)
	}
}

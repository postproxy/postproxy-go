package postproxy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProfilesList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/profiles" {
			t.Errorf("expected path /api/profiles, got %s", r.URL.Path)
		}
		if pgID := r.URL.Query().Get("profile_group_id"); pgID != "pg-1" {
			t.Errorf("expected profile_group_id=pg-1, got %s", pgID)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ListResponse[Profile]{
			Data: []Profile{
				{ID: "prof-1", Name: "Test Profile", Status: ProfileStatusActive, Platform: PlatformInstagram},
			},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	pgID := "pg-1"
	result, err := c.Profiles.List(context.Background(), &RequestOptions{ProfileGroupID: &pgID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(result.Data))
	}
	if result.Data[0].ID != "prof-1" {
		t.Errorf("expected profile ID %q, got %q", "prof-1", result.Data[0].ID)
	}
	if result.Data[0].Platform != PlatformInstagram {
		t.Errorf("expected platform %q, got %q", PlatformInstagram, result.Data[0].Platform)
	}
}

func TestProfilesGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/profiles/prof-123" {
			t.Errorf("expected path /api/profiles/prof-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Profile{
			ID:       "prof-123",
			Name:     "My Profile",
			Status:   ProfileStatusActive,
			Platform: PlatformFacebook,
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	profile, err := c.Profiles.Get(context.Background(), "prof-123", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile.ID != "prof-123" {
		t.Errorf("expected profile ID %q, got %q", "prof-123", profile.ID)
	}
	if profile.Name != "My Profile" {
		t.Errorf("expected name %q, got %q", "My Profile", profile.Name)
	}
}

func TestProfilesPlacements(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/profiles/prof-1/placements" {
			t.Errorf("expected path /api/profiles/prof-1/placements, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ListResponse[Placement]{
			Data: []Placement{
				{ID: "pl-1", Name: "Feed"},
				{ID: "pl-2", Name: "Story"},
			},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Profiles.Placements(context.Background(), "prof-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 2 {
		t.Fatalf("expected 2 placements, got %d", len(result.Data))
	}
	if result.Data[0].Name != "Feed" {
		t.Errorf("expected placement name %q, got %q", "Feed", result.Data[0].Name)
	}
}

func TestProfilesAssignPlacementToGroup(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/profiles/prof-1/assign_placement_to_group" {
			t.Errorf("expected path /api/profiles/prof-1/assign_placement_to_group, got %s", r.URL.Path)
		}

		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body["placement_id"] != "pl-1" {
			t.Errorf("expected placement_id=pl-1, got %s", body["placement_id"])
		}
		if body["target_profile_group_id"] != "pg-2" {
			t.Errorf("expected target_profile_group_id=pg-2, got %s", body["target_profile_group_id"])
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Placement{ID: "pl-1", Name: "Feed", ProfileGroupID: "pg-2"})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Profiles.AssignPlacementToGroup(context.Background(), "prof-1", "pl-1", "pg-2", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ProfileGroupID != "pg-2" {
		t.Errorf("expected profile_group_id %q, got %q", "pg-2", result.ProfileGroupID)
	}
}

func TestProfilesIceBreakers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/profiles/prof-1/ice_breakers" {
			t.Errorf("expected path /api/profiles/prof-1/ice_breakers, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(IceBreakersResponse{
			IceBreakers: []IceBreaker{{Question: "What do you do?", Payload: "services"}},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Profiles.IceBreakers(context.Background(), "prof-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.IceBreakers) != 1 {
		t.Fatalf("expected 1 ice breaker, got %d", len(result.IceBreakers))
	}
	if result.IceBreakers[0].Question != "What do you do?" {
		t.Errorf("expected question %q, got %q", "What do you do?", result.IceBreakers[0].Question)
	}
}

func TestProfilesSetIceBreakers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/profiles/prof-1/ice_breakers" {
			t.Errorf("expected path /api/profiles/prof-1/ice_breakers, got %s", r.URL.Path)
		}

		var body struct {
			IceBreakers []IceBreaker `json:"ice_breakers"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if len(body.IceBreakers) != 1 || body.IceBreakers[0].Payload != "services" {
			t.Errorf("unexpected body: %+v", body)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(SuccessResponse{Success: true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Profiles.SetIceBreakers(context.Background(), "prof-1", []IceBreaker{
		{Question: "What do you do?", Payload: "services"},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("expected success=true")
	}
}

func TestProfilesDeleteIceBreakers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/profiles/prof-1/ice_breakers" {
			t.Errorf("expected path /api/profiles/prof-1/ice_breakers, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(SuccessResponse{Success: true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Profiles.DeleteIceBreakers(context.Background(), "prof-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("expected success=true")
	}
}

func TestProfilesDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/profiles/prof-del" {
			t.Errorf("expected path /api/profiles/prof-del, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(SuccessResponse{Success: true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Profiles.Delete(context.Background(), "prof-del", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("expected success=true")
	}
}

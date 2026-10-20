package postproxy

import (
	"context"
	"testing"
)

var mockAccountPlacement = map[string]any{
	"id":              "pl_1",
	"platform":        "facebook",
	"external_id":     "1234567890",
	"name":            "Acme Page",
	"avatar_url":      nil,
	"status":          "active",
	"active_profile":  map[string]any{"id": "prof_1", "name": "Jane", "connected": true},
	"linked_profiles": []any{map[string]any{"id": "prof_2", "name": "John", "connected": true}},
}

func TestPlacementsList(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{
		"total": 1, "page": 1, "per_page": 25, "data": []any{mockAccountPlacement},
	})

	page, perPage := 2, 10
	result, err := c.Placements.List(context.Background(), &PlacementListOptions{Page: &page, PerPage: &perPage})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Data[0].ActiveProfile.ID != "prof_1" {
		t.Errorf("unexpected active profile: %+v", result.Data[0].ActiveProfile)
	}
	if len(result.Data[0].LinkedProfiles) != 1 {
		t.Errorf("expected 1 linked profile, got %d", len(result.Data[0].LinkedProfiles))
	}

	req := requests()[0]
	if req.Method != "GET" || req.Path != "/api/placements" {
		t.Errorf("unexpected request: %s %s", req.Method, req.Path)
	}
	if req.Query.Get("page") != "2" || req.Query.Get("per_page") != "10" {
		t.Errorf("unexpected query: %v", req.Query)
	}
}

func TestPlacementsSetActiveProfile(t *testing.T) {
	c, requests := newRecordingClient(t, mockAccountPlacement)

	if _, err := c.Placements.SetActiveProfile(context.Background(), "pl_1", "prof_2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := requests()[0]
	if req.Method != "PATCH" || req.Path != "/api/placements/pl_1/set_active_profile" {
		t.Errorf("unexpected request: %s %s", req.Method, req.Path)
	}
	assertJSON(t, req.JSON(t), `{"profile_id":"prof_2"}`)
}

func TestAPIKeysList(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"data": []any{map[string]any{
		"id": "k1", "name": "CI", "profile_group_id": "pg_1", "sandbox": false, "created_at": "2026-09-09T00:00:00Z",
	}}})

	result, err := c.APIKeys.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Data[0].Name != "CI" || result.Data[0].Token != nil {
		t.Errorf("unexpected key: %+v", result.Data[0])
	}
	if requests()[0].Path != "/api/api_keys" {
		t.Errorf("unexpected path: %s", requests()[0].Path)
	}
}

func TestAPIKeysCreate(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{
		"id": "k1", "name": "CI", "profile_group_id": "pg_1", "sandbox": false, "token": "secret", "created_at": "2026-09-09T00:00:00Z",
	})

	ctx := WithIdempotencyKey(context.Background(), "idem")
	key, err := c.APIKeys.Create(ctx, "CI", "pg_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key.Token == nil || *key.Token != "secret" {
		t.Errorf("unexpected token: %v", key.Token)
	}

	req := requests()[0]
	if req.Method != "POST" {
		t.Errorf("expected POST, got %s", req.Method)
	}
	assertJSON(t, req.JSON(t), `{"name":"CI","profile_group_id":"pg_1"}`)
	if req.Header.Get("Idempotency-Key") != "idem" {
		t.Errorf("expected Idempotency-Key idem, got %q", req.Header.Get("Idempotency-Key"))
	}
}

func TestAPIKeysDelete(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"deleted": true})

	result, err := c.APIKeys.Delete(context.Background(), "k1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Deleted {
		t.Error("expected deleted=true")
	}
	if req := requests()[0]; req.Method != "DELETE" || req.Path != "/api/api_keys/k1" {
		t.Errorf("unexpected request: %s %s", req.Method, req.Path)
	}
}

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

func TestProfilesBackfillPosts(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/profiles/prof-1/backfill_posts" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["from"] != "2025-01-01" {
			t.Errorf("expected from=2025-01-01, got %q", body["from"])
		}

		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(PostSync{
			ID:        "sync456def",
			ProfileID: "prof-1",
			Kind:      "posts",
			Trigger:   PostSyncTriggerBackfill,
			Status:    PostSyncStatusPending,
			CreatedAt: "2026-08-06T09:15:00.000Z",
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	sync, err := c.Profiles.BackfillPosts(context.Background(), "prof-1", "2025-01-01", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sync.ID != "sync456def" {
		t.Errorf("expected sync ID sync456def, got %q", sync.ID)
	}
	if sync.Trigger != PostSyncTriggerBackfill {
		t.Errorf("expected trigger backfill, got %q", sync.Trigger)
	}
	if sync.Status != PostSyncStatusPending {
		t.Errorf("expected status pending, got %q", sync.Status)
	}
}

func TestProfilesBackfillPostsConflict(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error":           "A posts backfill is already running for this profile",
			"profile_sync_id": "sync456def",
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	_, err := c.Profiles.BackfillPosts(context.Background(), "prof-1", "2025-01-01", nil)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !IsConflictError(err) {
		t.Fatalf("expected a conflict error, got %v", err)
	}

	apiErr, ok := err.(*PostProxyError)
	if !ok {
		t.Fatalf("expected *PostProxyError, got %T", err)
	}
	if apiErr.Response["profile_sync_id"] != "sync456def" {
		t.Errorf("expected the running sync id in the response, got %v", apiErr.Response)
	}
}

func TestProfilesPostSyncs(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/profiles/prof-1/post_syncs" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("trigger") != "backfill" {
			t.Errorf("expected trigger=backfill, got %q", q.Get("trigger"))
		}
		if q.Get("status") != "running" {
			t.Errorf("expected status=running, got %q", q.Get("status"))
		}
		if q.Get("per_page") != "25" {
			t.Errorf("expected per_page=25, got %q", q.Get("per_page"))
		}

		w.WriteHeader(http.StatusOK)
		oldest := "2025-11-04T18:22:00.000Z"
		_ = json.NewEncoder(w).Encode(PaginatedResponse[PostSync]{
			Data: []PostSync{{
				ID:             "sync456def",
				ProfileID:      "prof-1",
				Trigger:        PostSyncTriggerBackfill,
				Status:         PostSyncStatusRunning,
				PostsSeen:      150,
				PostsImported:  143,
				OldestPostedAt: &oldest,
			}},
			Total:   1,
			Page:    0,
			PerPage: 25,
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	trigger := PostSyncTriggerBackfill
	status := PostSyncStatusRunning
	perPage := 25
	result, err := c.Profiles.PostSyncs(context.Background(), "prof-1", &PostSyncListOptions{
		Trigger: &trigger,
		Status:  &status,
		PerPage: &perPage,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Total)
	}
	if result.Data[0].PostsImported != 143 {
		t.Errorf("expected 143 imported, got %d", result.Data[0].PostsImported)
	}
	if result.Data[0].OldestPostedAt == nil {
		t.Error("expected oldest_posted_at to be set")
	}
}

func TestProfilesPostSync(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/profiles/prof-1/post_syncs/sync456def" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(PostSync{ID: "sync456def", Status: PostSyncStatusCompleted})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	sync, err := c.Profiles.PostSync(context.Background(), "prof-1", "sync456def", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sync.Status != PostSyncStatusCompleted {
		t.Errorf("expected status completed, got %q", sync.Status)
	}
}

func TestProfilesNewFields(t *testing.T) {
	c, _ := newRecordingClient(t, map[string]any{
		"id": "p1", "name": "Acme", "username": "acme", "status": "inactive", "inactive_reason": "network_taken",
		"platform": "instagram", "platform_id": "1789", "profile_group_id": "pg_1", "expires_at": nil, "post_count": 3,
		"avatar_url": "https://cdn.test/a.jpg", "platform_url": "https://instagram.com/acme", "placements_sync_status": "synced",
	})

	p, err := c.Profiles.Get(context.Background(), "p1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.InactiveReason == nil || *p.InactiveReason != ProfileInactiveReasonNetworkTaken {
		t.Errorf("unexpected inactive_reason: %v", p.InactiveReason)
	}
	if p.Username == nil || p.PlatformID == nil || p.AvatarURL == nil || p.PlatformURL == nil || p.PlacementsSyncStatus == nil {
		t.Errorf("expected new fields to be set: %+v", p)
	}
}

func TestProfilesDailyStats(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"data": map[string]any{
		"profile_id": "p1", "attribution": "publish",
		"records": []any{map[string]any{"date": "2026-09-01", "impressions": 10, "platform_stats": nil}},
		"metrics": map[string]any{},
	}})

	attribution := DailyStatsAttributionPublish
	result, err := c.Profiles.DailyStats(context.Background(), "p1", &DailyStatsOptions{
		PlacementID: ptr("pl1"), From: ptr("2026-09-01"), To: ptr("2026-09-30"), Attribution: &attribution,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	record := result.Data.Records[0]
	if record.Date != "2026-09-01" || record.Metrics["impressions"] != float64(10) {
		t.Errorf("unexpected record: %+v", record)
	}
	if _, ok := record.Metrics["date"]; ok {
		t.Error("date should not be in Metrics")
	}

	req := requests()[0]
	if req.Path != "/api/profiles/p1/daily_stats" {
		t.Errorf("unexpected path: %s", req.Path)
	}
	if req.Query.Get("attribution") != "publish" || req.Query.Get("placement_id") != "pl1" {
		t.Errorf("unexpected query: %v", req.Query)
	}
}

func TestProfilesFollowerHistory(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"data": map[string]any{"records": []any{
		map[string]any{"date": "2026-09-01", "followers_gained": 5, "followers_lost": 1, "followers_change": 4, "followers_total": 100},
	}}})

	result, err := c.Profiles.FollowerHistory(context.Background(), "p1", &FollowerHistoryOptions{From: ptr("2026-09-01")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r := result.Data.Records[0]; r.FollowersTotal == nil || *r.FollowersTotal != 100 {
		t.Errorf("unexpected record: %+v", r)
	}
	req := requests()[0]
	if req.Path != "/api/profiles/p1/follower_history" || req.Query.Get("from") != "2026-09-01" {
		t.Errorf("unexpected request: %s %v", req.Path, req.Query)
	}
}

func TestProfilesParticipantStatus(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"is_verified_user": true, "fetched_at": "x"})

	result, err := c.Profiles.ParticipantStatus(context.Background(), "p1", "igsid_1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsVerifiedUser == nil || !*result.IsVerifiedUser {
		t.Errorf("unexpected status: %+v", result)
	}
	if requests()[0].Path != "/api/profiles/p1/participants/igsid_1" {
		t.Errorf("unexpected path: %s", requests()[0].Path)
	}
}

func TestProfilesResolveMention(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"urn": "urn:li:organization:1", "type": "organization", "name": "Acme", "vanity_name": "acme"})

	result, err := c.Profiles.ResolveMention(context.Background(), "p1", "https://www.linkedin.com/company/acme", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.URN != "urn:li:organization:1" {
		t.Errorf("unexpected urn: %s", result.URN)
	}
	if got := requests()[0].Query.Get("url"); got != "https://www.linkedin.com/company/acme" {
		t.Errorf("unexpected url param: %s", got)
	}
}

func TestProfilesInstagramAvailableAudio(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"data": []any{map[string]any{"id": "a1"}}})

	audioType := InstagramAudioTypeMusic
	result, err := c.Profiles.InstagramAvailableAudio(context.Background(), "p1", &InstagramAvailableAudioOptions{AudioType: &audioType, SearchQuery: ptr("summer")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Data[0]["id"] != "a1" {
		t.Errorf("unexpected data: %v", result.Data)
	}
	req := requests()[0]
	if req.Path != "/api/profiles/p1/instagram/available_audio" {
		t.Errorf("unexpected path: %s", req.Path)
	}
	if req.Query.Get("audio_type") != "music" || req.Query.Get("search_query") != "summer" {
		t.Errorf("unexpected query: %v", req.Query)
	}
}

func TestProfilesFacebookPage(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"about": "Hi"})
	ctx := context.Background()

	page, err := c.Profiles.FacebookPage(ctx, "p1", "page_1", &FacebookPageOptions{Fields: ptr("about,website")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page["about"] != "Hi" {
		t.Errorf("unexpected page: %v", page)
	}
	_, _ = c.Profiles.UpdateFacebookPage(ctx, "p1", "page_1", map[string]any{"about": "New about"}, nil)

	reqs := requests()
	if reqs[0].Query.Get("page_id") != "page_1" || reqs[0].Query.Get("fields") != "about,website" {
		t.Errorf("unexpected query: %v", reqs[0].Query)
	}
	if reqs[1].Method != "PATCH" || reqs[1].Path != "/api/profiles/p1/facebook/update_page" {
		t.Errorf("unexpected request: %s %s", reqs[1].Method, reqs[1].Path)
	}
	assertJSON(t, reqs[1].JSON(t), `{"about":"New about","page_id":"page_1"}`)
}

func TestProfilesAssignToGroup(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"id": "p1"})

	if _, err := c.Profiles.AssignToGroup(context.Background(), "p1", "pg_2", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req := requests()[0]
	if req.Method != "PATCH" || req.Path != "/api/profiles/p1/assign_to_group" {
		t.Errorf("unexpected request: %s %s", req.Method, req.Path)
	}
	assertJSON(t, req.JSON(t), `{"target_profile_group_id":"pg_2"}`)
}

func TestProfilesChatSyncs(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"id": "cs1", "status": "pending", "chats_seen": 0, "messages_imported": 0})
	ctx := context.Background()

	sync, err := c.Profiles.BackfillChats(ctx, "p1", "2026-01-01T00:00:00Z", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sync.ID != "cs1" || sync.Status != PostSyncStatusPending {
		t.Errorf("unexpected sync: %+v", sync)
	}
	status := PostSyncStatusRunning
	_, _ = c.Profiles.ChatSyncs(ctx, "p1", &ChatSyncListOptions{Status: &status, PerPage: ptr(5)})
	_, _ = c.Profiles.ChatSync(ctx, "p1", "cs1", nil)

	reqs := requests()
	if reqs[0].Method != "POST" || reqs[0].Path != "/api/profiles/p1/backfill_chats" {
		t.Errorf("unexpected request: %s %s", reqs[0].Method, reqs[0].Path)
	}
	assertJSON(t, reqs[0].JSON(t), `{"from":"2026-01-01T00:00:00Z"}`)
	if reqs[1].Query.Get("status") != "running" || reqs[1].Query.Get("per_page") != "5" {
		t.Errorf("unexpected query: %v", reqs[1].Query)
	}
	if reqs[2].Path != "/api/profiles/p1/chat_syncs/cs1" {
		t.Errorf("unexpected path: %s", reqs[2].Path)
	}
}

func TestProfilesInstagramCollaborationInvites(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{
		"collaboration_invites": []any{map[string]any{
			"media_id": "1789", "media_owner_username": "partner_brand", "caption": "Collab", "media_url": "https://x",
		}},
		"next_cursor": nil,
	})
	ctx := context.Background()

	list, err := c.Profiles.InstagramCollaborationInvites(ctx, "p1", &InstagramCollaborationInvitesOptions{After: ptr("cur_1")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if list.CollaborationInvites[0].MediaID != "1789" || list.NextCursor != nil {
		t.Errorf("unexpected list: %+v", list)
	}
	if _, err := c.Profiles.AcceptInstagramCollaborationInvite(ctx, "p1", "1789", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := c.Profiles.DeclineInstagramCollaborationInvite(ctx, "p1", "1790", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reqs := requests()
	if reqs[0].Path != "/api/profiles/p1/instagram/collaboration_invites" || reqs[0].Query.Get("after") != "cur_1" {
		t.Errorf("unexpected list request: %s %v", reqs[0].Path, reqs[0].Query)
	}
	if reqs[1].Method != "POST" || reqs[1].Path != "/api/profiles/p1/instagram/accept_collaboration_invite" {
		t.Errorf("unexpected accept request: %s %s", reqs[1].Method, reqs[1].Path)
	}
	assertJSON(t, reqs[1].JSON(t), `{"media_id":"1789"}`)
	if reqs[2].Path != "/api/profiles/p1/instagram/decline_collaboration_invite" {
		t.Errorf("unexpected decline path: %s", reqs[2].Path)
	}
	assertJSON(t, reqs[2].JSON(t), `{"media_id":"1790"}`)
}

func TestProfilesInstagramCollaborationInviteResult(t *testing.T) {
	c, _ := newRecordingClient(t, map[string]any{"accepted": true, "media_id": "1789"})

	result, err := c.Profiles.AcceptInstagramCollaborationInvite(context.Background(), "p1", "1789", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Accepted == nil || !*result.Accepted || result.MediaID != "1789" || result.Declined != nil {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestProfilesPinterestBoards(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"id": "1022", "name": "Summer Recipes", "avatar_url": nil, "metadata": map[string]any{}})
	ctx := context.Background()

	privacy := PinterestBoardPrivacyProtected
	board, err := c.Profiles.CreatePinterestBoard(ctx, "p1", "Summer Recipes", &PinterestBoardOptions{
		Description: ptr("Light dinners"), Privacy: &privacy,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if board.ID != "1022" {
		t.Errorf("unexpected board: %+v", board)
	}
	if _, err := c.Profiles.DeletePinterestBoard(ctx, "p1", "1022", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reqs := requests()
	if reqs[0].Method != "POST" || reqs[0].Path != "/api/profiles/p1/pinterest/create_board" {
		t.Errorf("unexpected create request: %s %s", reqs[0].Method, reqs[0].Path)
	}
	assertJSON(t, reqs[0].JSON(t), `{"name":"Summer Recipes","description":"Light dinners","privacy":"PROTECTED"}`)
	if reqs[1].Method != "DELETE" || reqs[1].Path != "/api/profiles/p1/pinterest/delete_board" || reqs[1].Query.Get("board_id") != "1022" {
		t.Errorf("unexpected delete request: %s %s %v", reqs[1].Method, reqs[1].Path, reqs[1].Query)
	}
}

func TestProfilesYouTubePlaylists(t *testing.T) {
	playlist := map[string]any{
		"id": "PL1", "title": "Tutorials", "description": "", "privacy_status": "public",
		"item_count": 3, "thumbnail_url": nil, "published_at": "2026-03-14T09:21:00.000Z",
	}
	c, requests := newRecordingClient(t, map[string]any{"data": []any{playlist}, "next_page_token": nil})
	ctx := context.Background()

	list, err := c.Profiles.YouTubePlaylists(ctx, "p1", &YouTubePlaylistsOptions{PageToken: ptr("tok")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if list.Data[0].ID != "PL1" || list.Data[0].PrivacyStatus != YouTubePrivacyPublic || list.Data[0].ItemCount != 3 || list.NextPageToken != nil {
		t.Errorf("unexpected list: %+v", list)
	}
	unlisted := YouTubePrivacyUnlisted
	if _, err := c.Profiles.CreateYouTubePlaylist(ctx, "p1", "Tutorials", &YouTubePlaylistCreateOptions{PrivacyStatus: &unlisted}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := c.Profiles.UpdateYouTubePlaylist(ctx, "p1", "PL1", &YouTubePlaylistUpdateOptions{Description: ptr("")}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := c.Profiles.DeleteYouTubePlaylist(ctx, "p1", "PL1", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reqs := requests()
	if reqs[0].Path != "/api/profiles/p1/youtube/playlists" || reqs[0].Query.Get("page_token") != "tok" {
		t.Errorf("unexpected list request: %s %v", reqs[0].Path, reqs[0].Query)
	}
	if reqs[1].Method != "POST" || reqs[1].Path != "/api/profiles/p1/youtube/create_playlist" {
		t.Errorf("unexpected create request: %s %s", reqs[1].Method, reqs[1].Path)
	}
	assertJSON(t, reqs[1].JSON(t), `{"title":"Tutorials","privacy_status":"unlisted"}`)
	if reqs[2].Method != "PATCH" || reqs[2].Path != "/api/profiles/p1/youtube/update_playlist" {
		t.Errorf("unexpected update request: %s %s", reqs[2].Method, reqs[2].Path)
	}
	assertJSON(t, reqs[2].JSON(t), `{"playlist_id":"PL1","description":""}`)
	if reqs[3].Method != "DELETE" || reqs[3].Path != "/api/profiles/p1/youtube/delete_playlist" || reqs[3].Query.Get("playlist_id") != "PL1" {
		t.Errorf("unexpected delete request: %s %s %v", reqs[3].Method, reqs[3].Path, reqs[3].Query)
	}
}

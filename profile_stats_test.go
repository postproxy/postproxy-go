package postproxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func ptr[T any](v T) *T { return &v }

func TestProfilesService_GetProfileStats_WithPlacement(t *testing.T) {
	var gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"profile_id":"pf1","platform":"linkedin","placement_id":"org_1","records":[{"stats":{"followerCount":100},"recorded_at":"2026-05-12T00:00:00Z"}]}}`))
	}))
	defer srv.Close()

	c := NewClient("k", WithBaseURL(srv.URL))
	result, err := c.Profiles.GetProfileStats(context.Background(), "pf1", &ProfileStatsOptions{
		PlacementID: ptr("org_1"),
		From:        ptr("2026-04-01T00:00:00Z"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Data.ProfileID != "pf1" {
		t.Errorf("profile_id mismatch: %q", result.Data.ProfileID)
	}
	if got, _ := result.Data.Records[0].Stats["followerCount"].(float64); got != 100 {
		t.Errorf("followerCount mismatch: %+v", result.Data.Records[0].Stats["followerCount"])
	}
	if gotPath != "/api/profiles/pf1/stats" {
		t.Errorf("path mismatch: %q", gotPath)
	}
	if !strings.Contains(gotQuery, "placement_id=org_1") {
		t.Errorf("query missing placement_id: %q", gotQuery)
	}
}

func TestProfilesService_GetProfileStats_NoPlacement(t *testing.T) {
	var gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"profile_id":"bsky1","platform":"bluesky","placement_id":null,"records":[]}}`))
	}))
	defer srv.Close()

	c := NewClient("k", WithBaseURL(srv.URL))
	if _, err := c.Profiles.GetProfileStats(context.Background(), "bsky1", nil); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(gotQuery, "placement_id") {
		t.Errorf("unexpected placement_id in query: %q", gotQuery)
	}
}

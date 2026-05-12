package postproxy

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProfileGroupsService_ConnectBluesky(t *testing.T) {
	var gotBody map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true,"profile":{"id":"pf_bsky_1","network":"bluesky","name":"Jane","external_username":"jane.bsky.social"}}`))
	}))
	defer srv.Close()

	c := NewClient("k", WithBaseURL(srv.URL))
	result, err := c.ProfileGroups.ConnectBluesky(context.Background(), "pg-1", BlueskyConnectOptions{
		Identifier:  "jane.bsky.social",
		AppPassword: "xxxx",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Success || result.Profile.ID != "pf_bsky_1" {
		t.Errorf("unexpected result: %+v", result)
	}
	expected := map[string]string{
		"platform":     "bluesky",
		"identifier":   "jane.bsky.social",
		"app_password": "xxxx",
	}
	for k, v := range expected {
		if gotBody[k] != v {
			t.Errorf("body[%s]=%q, want %q", k, gotBody[k], v)
		}
	}
}

func TestProfileGroupsService_ConnectTelegram(t *testing.T) {
	var gotBody map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true,"profile":{"id":"pf_tg_1","network":"telegram","name":"My Bot","external_username":"my_bot"},"next_step":"Add bot as admin"}`))
	}))
	defer srv.Close()

	c := NewClient("k", WithBaseURL(srv.URL))
	result, err := c.ProfileGroups.ConnectTelegram(context.Background(), "pg-1", TelegramConnectOptions{
		BotToken: "123:ABC",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.NextStep == nil || *result.NextStep != "Add bot as admin" {
		t.Errorf("next_step mismatch: %+v", result.NextStep)
	}
	if gotBody["platform"] != "telegram" || gotBody["bot_token"] != "123:ABC" {
		t.Errorf("body mismatch: %+v", gotBody)
	}
}

package postproxy

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQueuesList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/post_queues" {
			t.Errorf("expected path /api/post_queues, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id": "q1abc", "name": "Morning Posts", "timezone": "America/New_York",
					"enabled": true, "jitter": 10, "profile_group_id": "pg123", "posts_count": 5,
					"timeslots": []map[string]any{
						{"id": 1, "day": 1, "time": "09:00"},
						{"id": 2, "day": 3, "time": "09:00"},
					},
				},
			},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Queues.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 queue, got %d", len(result.Data))
	}
	if result.Data[0].ID != "q1abc" {
		t.Errorf("expected id q1abc, got %s", result.Data[0].ID)
	}
	if len(result.Data[0].Timeslots) != 2 {
		t.Errorf("expected 2 timeslots, got %d", len(result.Data[0].Timeslots))
	}
}

func TestQueuesGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/post_queues/q1abc" {
			t.Errorf("expected path /api/post_queues/q1abc, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "q1abc", "name": "Morning Posts", "timezone": "America/New_York",
			"enabled": true, "jitter": 10, "profile_group_id": "pg123", "posts_count": 5,
			"timeslots": []map[string]any{},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	queue, err := c.Queues.Get(context.Background(), "q1abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if queue.ID != "q1abc" {
		t.Errorf("expected id q1abc, got %s", queue.ID)
	}
	if queue.Name != "Morning Posts" {
		t.Errorf("expected name Morning Posts, got %s", queue.Name)
	}
}

func TestQueuesNextSlot(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/post_queues/q1abc/next_slot" {
			t.Errorf("expected path /api/post_queues/q1abc/next_slot, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"next_slot": "2026-03-11T14:00:00Z"})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Queues.NextSlot(context.Background(), "q1abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.NextSlot != "2026-03-11T14:00:00Z" {
		t.Errorf("expected next_slot 2026-03-11T14:00:00Z, got %s", result.NextSlot)
	}
}

func TestQueuesCreate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		_ = json.Unmarshal(body, &payload)

		if payload["profile_group_id"] != "pg123" {
			t.Errorf("expected profile_group_id pg123, got %v", payload["profile_group_id"])
		}
		pq := payload["post_queue"].(map[string]any)
		if pq["name"] != "Morning Posts" {
			t.Errorf("expected name Morning Posts, got %v", pq["name"])
		}
		if pq["timezone"] != "America/New_York" {
			t.Errorf("expected timezone America/New_York, got %v", pq["timezone"])
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "q1abc", "name": "Morning Posts", "timezone": "America/New_York",
			"enabled": true, "jitter": 10, "profile_group_id": "pg123", "posts_count": 0,
			"timeslots": []map[string]any{{"id": 1, "day": 1, "time": "09:00"}},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	tz := "America/New_York"
	jitter := 10
	queue, err := c.Queues.Create(context.Background(), "Morning Posts", "pg123", &QueueCreateOptions{
		Timezone: &tz,
		Jitter:   &jitter,
		Timeslots: []TimeslotInput{
			{Day: 1, Time: "09:00"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if queue.ID != "q1abc" {
		t.Errorf("expected id q1abc, got %s", queue.ID)
	}
}

func TestQueuesUpdate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		_ = json.Unmarshal(body, &payload)

		pq := payload["post_queue"].(map[string]any)
		if pq["enabled"] != false {
			t.Errorf("expected enabled=false, got %v", pq["enabled"])
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "q1abc", "name": "Morning Posts", "timezone": "America/New_York",
			"enabled": false, "jitter": 10, "profile_group_id": "pg123", "posts_count": 0,
			"timeslots": []map[string]any{},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	enabled := false
	queue, err := c.Queues.Update(context.Background(), "q1abc", &QueueUpdateOptions{
		Enabled: &enabled,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if queue.Enabled != false {
		t.Errorf("expected enabled=false, got %v", queue.Enabled)
	}
}

func TestQueuesDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/post_queues/q1abc" {
			t.Errorf("expected path /api/post_queues/q1abc, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(DeleteResponse{Deleted: true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Queues.Delete(context.Background(), "q1abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Deleted {
		t.Error("expected deleted=true")
	}
}

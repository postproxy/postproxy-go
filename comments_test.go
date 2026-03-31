package postproxy

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommentsList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/posts/post1/comments" {
			t.Errorf("expected path /api/posts/post1/comments, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("profile_id") != "prof1" {
			t.Errorf("expected profile_id=prof1, got %s", r.URL.Query().Get("profile_id"))
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"total": 1, "page": 0, "per_page": 20,
			"data": []map[string]any{
				{
					"id": "cmt_abc123", "external_id": "123456", "body": "Great post!",
					"status": "synced", "author_username": "someuser", "like_count": 3,
					"is_hidden": false, "created_at": "2026-03-25T10:01:00Z",
					"replies": []map[string]any{
						{"id": "cmt_def456", "body": "Thanks!", "status": "synced", "created_at": "2026-03-25T10:05:00Z", "replies": []any{}},
					},
				},
			},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Comments.List(context.Background(), "post1", "prof1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Total)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(result.Data))
	}
	if result.Data[0].ID != "cmt_abc123" {
		t.Errorf("expected id cmt_abc123, got %s", result.Data[0].ID)
	}
	if len(result.Data[0].Replies) != 1 {
		t.Errorf("expected 1 reply, got %d", len(result.Data[0].Replies))
	}
}

func TestCommentsListWithPagination(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "2" {
			t.Errorf("expected page=2, got %s", r.URL.Query().Get("page"))
		}
		if r.URL.Query().Get("per_page") != "10" {
			t.Errorf("expected per_page=10, got %s", r.URL.Query().Get("per_page"))
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"total": 42, "page": 2, "per_page": 10, "data": []any{},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	page, perPage := 2, 10
	result, err := c.Comments.List(context.Background(), "post1", "prof1", &CommentListOptions{Page: &page, PerPage: &perPage})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 42 {
		t.Errorf("expected total 42, got %d", result.Total)
	}
}

func TestCommentsGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/posts/post1/comments/cmt_abc123" {
			t.Errorf("expected path /api/posts/post1/comments/cmt_abc123, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("profile_id") != "prof1" {
			t.Errorf("expected profile_id=prof1")
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "cmt_abc123", "body": "Great post!", "status": "synced",
			"like_count": 3, "created_at": "2026-03-25T10:01:00Z",
			"replies": []any{},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	comment, err := c.Comments.Get(context.Background(), "post1", "cmt_abc123", "prof1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID != "cmt_abc123" {
		t.Errorf("expected id cmt_abc123, got %s", comment.ID)
	}
	if comment.Body != "Great post!" {
		t.Errorf("expected body 'Great post!', got %s", comment.Body)
	}
}

func TestCommentsCreate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/posts/post1/comments" {
			t.Errorf("expected path /api/posts/post1/comments, got %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var req map[string]string
		_ = json.Unmarshal(body, &req)
		if req["text"] != "Nice!" {
			t.Errorf("expected text 'Nice!', got %s", req["text"])
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "cmt_new", "body": "Nice!", "status": "pending",
			"created_at": "2026-03-25T12:00:00Z", "replies": []any{},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	comment, err := c.Comments.Create(context.Background(), "post1", "prof1", "Nice!", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID != "cmt_new" {
		t.Errorf("expected id cmt_new, got %s", comment.ID)
	}
	if comment.Status != "pending" {
		t.Errorf("expected status pending, got %s", comment.Status)
	}
}

func TestCommentsCreateReply(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req map[string]string
		_ = json.Unmarshal(body, &req)
		if req["parent_id"] != "cmt_abc123" {
			t.Errorf("expected parent_id cmt_abc123, got %s", req["parent_id"])
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "cmt_reply", "body": "Thanks!", "status": "pending",
			"created_at": "2026-03-25T12:00:00Z", "replies": []any{},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	parentID := "cmt_abc123"
	comment, err := c.Comments.Create(context.Background(), "post1", "prof1", "Thanks!", &CommentCreateOptions{ParentID: &parentID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID != "cmt_reply" {
		t.Errorf("expected id cmt_reply, got %s", comment.ID)
	}
}

func TestCommentsDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/posts/post1/comments/cmt_abc123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"accepted": true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Comments.Delete(context.Background(), "post1", "cmt_abc123", "prof1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Accepted {
		t.Errorf("expected accepted true")
	}
}

func TestCommentsHide(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/posts/post1/comments/cmt_abc123/hide" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"accepted": true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Comments.Hide(context.Background(), "post1", "cmt_abc123", "prof1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Accepted {
		t.Errorf("expected accepted true")
	}
}

func TestCommentsUnhide(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/posts/post1/comments/cmt_abc123/unhide" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"accepted": true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Comments.Unhide(context.Background(), "post1", "cmt_abc123", "prof1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Accepted {
		t.Errorf("expected accepted true")
	}
}

func TestCommentsLike(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/posts/post1/comments/cmt_abc123/like" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"accepted": true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Comments.Like(context.Background(), "post1", "cmt_abc123", "prof1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Accepted {
		t.Errorf("expected accepted true")
	}
}

func TestCommentsUnlike(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/posts/post1/comments/cmt_abc123/unlike" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"accepted": true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Comments.Unlike(context.Background(), "post1", "cmt_abc123", "prof1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Accepted {
		t.Errorf("expected accepted true")
	}
}

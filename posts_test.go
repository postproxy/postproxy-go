package postproxy

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPostsList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/posts" {
			t.Errorf("expected path /api/posts, got %s", r.URL.Path)
		}

		q := r.URL.Query()
		if q.Get("page") != "2" {
			t.Errorf("expected page=2, got %s", q.Get("page"))
		}
		if q.Get("per_page") != "10" {
			t.Errorf("expected per_page=10, got %s", q.Get("per_page"))
		}
		if q.Get("status") != "draft" {
			t.Errorf("expected status=draft, got %s", q.Get("status"))
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(PaginatedResponse[Post]{
			Data:    []Post{{ID: "post-1", Body: "hello", Status: PostStatusDraft}},
			Total:   1,
			Page:    2,
			PerPage: 10,
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	page := 2
	perPage := 10
	status := PostStatusDraft
	result, err := c.Posts.List(context.Background(), &PostListOptions{
		Page:    &page,
		PerPage: &perPage,
		Status:  &status,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 post, got %d", len(result.Data))
	}
	if result.Data[0].ID != "post-1" {
		t.Errorf("expected post ID %q, got %q", "post-1", result.Data[0].ID)
	}
	if result.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Total)
	}
}

func TestPostsGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/posts/post-123" {
			t.Errorf("expected path /api/posts/post-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Post{ID: "post-123", Body: "test", Status: PostStatusPending})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	post, err := c.Posts.Get(context.Background(), "post-123", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.ID != "post-123" {
		t.Errorf("expected post ID %q, got %q", "post-123", post.ID)
	}
}

func TestPostsCreate_JSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/posts" {
			t.Errorf("expected path /api/posts, got %s", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		_ = json.Unmarshal(body, &payload)

		post, _ := payload["post"].(map[string]any)
		if post["body"] != "Hello world" {
			t.Errorf("expected body %q, got %v", "Hello world", post["body"])
		}

		profiles, _ := payload["profiles"].([]any)
		if len(profiles) != 2 {
			t.Errorf("expected 2 profiles, got %d", len(profiles))
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Post{ID: "new-post", Body: "Hello world", Status: PostStatusPending})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	post, err := c.Posts.Create(context.Background(), "Hello world", []string{"prof-1", "prof-2"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.ID != "new-post" {
		t.Errorf("expected post ID %q, got %q", "new-post", post.ID)
	}
}

func TestPostsCreate_WithOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		_ = json.Unmarshal(body, &payload)

		post, _ := payload["post"].(map[string]any)
		if post["draft"] != true {
			t.Errorf("expected draft=true, got %v", post["draft"])
		}
		if post["scheduled_at"] != "2025-01-01T00:00:00Z" {
			t.Errorf("expected scheduled_at, got %v", post["scheduled_at"])
		}

		media, _ := payload["media"].([]any)
		if len(media) != 1 {
			t.Errorf("expected 1 media URL, got %d", len(media))
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Post{ID: "new-post", Status: PostStatusDraft})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	draft := true
	scheduled := "2025-01-01T00:00:00Z"
	_, err := c.Posts.Create(context.Background(), "Post", []string{"prof-1"}, &PostCreateOptions{
		Media:       []string{"https://example.com/image.jpg"},
		Draft:       &draft,
		ScheduledAt: &scheduled,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostsCreate_FormData(t *testing.T) {
	// Create a temp file to upload.
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.jpg")
	if err := os.WriteFile(tmpFile, []byte("fake image data"), 0644); err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("expected multipart/form-data, got %s", ct)
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatalf("failed to parse form: %v", err)
		}

		if r.FormValue("post[body]") != "Upload test" {
			t.Errorf("expected body %q, got %q", "Upload test", r.FormValue("post[body]"))
		}

		profiles := r.MultipartForm.Value["profiles[]"]
		if len(profiles) != 1 || profiles[0] != "prof-1" {
			t.Errorf("expected profiles [prof-1], got %v", profiles)
		}

		files := r.MultipartForm.File["media[]"]
		if len(files) != 1 {
			t.Errorf("expected 1 file, got %d", len(files))
		} else if files[0].Filename != "test.jpg" {
			t.Errorf("expected filename test.jpg, got %s", files[0].Filename)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Post{ID: "uploaded-post", Status: PostStatusPending})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	post, err := c.Posts.Create(context.Background(), "Upload test", []string{"prof-1"}, &PostCreateOptions{
		MediaFiles: []string{tmpFile},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.ID != "uploaded-post" {
		t.Errorf("expected post ID %q, got %q", "uploaded-post", post.ID)
	}
}

func TestPostsCreate_PlatformParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		_ = json.Unmarshal(body, &payload)

		platforms, _ := payload["platforms"].(map[string]any)
		ig, _ := platforms["instagram"].(map[string]any)
		if ig["format"] != "reel" {
			t.Errorf("expected instagram format %q, got %v", "reel", ig["format"])
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Post{ID: "post-plat"})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	igFmt := InstagramFormatReel
	_, err := c.Posts.Create(context.Background(), "Test", []string{"p1"}, &PostCreateOptions{
		Platforms: &PlatformParams{
			Instagram: &InstagramParams{Format: &igFmt},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostsPublishDraft(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/posts/draft-1/publish" {
			t.Errorf("expected path /api/posts/draft-1/publish, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Post{ID: "draft-1", Status: PostStatusProcessing})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	post, err := c.Posts.PublishDraft(context.Background(), "draft-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.Status != PostStatusProcessing {
		t.Errorf("expected status %q, got %q", PostStatusProcessing, post.Status)
	}
}

func TestPostsDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/posts/post-del" {
			t.Errorf("expected path /api/posts/post-del, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(DeleteResponse{Deleted: true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Posts.Delete(context.Background(), "post-del", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Deleted {
		t.Error("expected deleted=true")
	}
}

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

func TestPostsStats(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/posts/stats" {
			t.Errorf("expected path /api/posts/stats, got %s", r.URL.Path)
		}

		q := r.URL.Query()
		if q.Get("post_ids") != "abc123,def456" {
			t.Errorf("expected post_ids=abc123,def456, got %s", q.Get("post_ids"))
		}
		if q.Get("profiles") != "instagram,twitter" {
			t.Errorf("expected profiles=instagram,twitter, got %s", q.Get("profiles"))
		}
		if q.Get("from") != "2026-02-01T00:00:00Z" {
			t.Errorf("expected from=2026-02-01T00:00:00Z, got %s", q.Get("from"))
		}
		if q.Get("to") != "2026-02-24T00:00:00Z" {
			t.Errorf("expected to=2026-02-24T00:00:00Z, got %s", q.Get("to"))
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"abc123": map[string]any{
					"platforms": []map[string]any{
						{
							"profile_id": "prof_abc",
							"platform":   "instagram",
							"records": []map[string]any{
								{
									"stats":       map[string]any{"impressions": 1200, "likes": 85},
									"recorded_at": "2026-02-20T12:00:00Z",
								},
								{
									"stats":       map[string]any{"impressions": 1523, "likes": 102},
									"recorded_at": "2026-02-21T04:00:00Z",
								},
							},
						},
					},
				},
			},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	from := "2026-02-01T00:00:00Z"
	to := "2026-02-24T00:00:00Z"
	result, err := c.Posts.Stats(context.Background(), []string{"abc123", "def456"}, &PostStatsOptions{
		Profiles: []string{"instagram", "twitter"},
		From:     &from,
		To:       &to,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	postStats, ok := result.Data["abc123"]
	if !ok {
		t.Fatal("expected data for abc123")
	}
	if len(postStats.Platforms) != 1 {
		t.Fatalf("expected 1 platform, got %d", len(postStats.Platforms))
	}
	plat := postStats.Platforms[0]
	if plat.ProfileID != "prof_abc" {
		t.Errorf("expected profile_id %q, got %q", "prof_abc", plat.ProfileID)
	}
	if plat.Platform != PlatformInstagram {
		t.Errorf("expected platform %q, got %q", PlatformInstagram, plat.Platform)
	}
	if len(plat.Records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(plat.Records))
	}
	if plat.Records[0].RecordedAt != "2026-02-20T12:00:00Z" {
		t.Errorf("expected recorded_at %q, got %q", "2026-02-20T12:00:00Z", plat.Records[0].RecordedAt)
	}
}

func TestPostsStats_NoOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("post_ids") != "abc123" {
			t.Errorf("expected post_ids=abc123, got %s", q.Get("post_ids"))
		}
		if q.Get("profiles") != "" {
			t.Errorf("expected no profiles param, got %s", q.Get("profiles"))
		}
		if q.Get("from") != "" {
			t.Errorf("expected no from param, got %s", q.Get("from"))
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{}})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Posts.Stats(context.Background(), []string{"abc123"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Data == nil {
		t.Fatal("expected non-nil data")
	}
}

func TestPostsStats_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  400,
			"error":   "Bad Request",
			"message": "param is missing or the value is empty: post_ids",
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	_, err := c.Posts.Stats(context.Background(), []string{}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsBadRequestError(err) {
		t.Errorf("expected bad request error, got %v", err)
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

func TestPostsDelete_WithDeleteOnPlatform(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Query().Get("delete_on_platform") != "true" {
			t.Errorf("expected delete_on_platform=true, got %q", r.URL.Query().Get("delete_on_platform"))
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(DeleteResponse{Deleted: true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	truthy := true
	result, err := c.Posts.Delete(context.Background(), "post-del", &PostDeleteOptions{DeleteOnPlatform: &truthy})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Deleted {
		t.Error("expected deleted=true")
	}
}

func TestPostsDeleteOnPlatform_AllPlatforms(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/posts/post-1/delete_on_platform" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		if len(body) != 0 {
			t.Errorf("expected empty body, got %q", string(body))
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(DeleteOnPlatformResponse{
			Success:  true,
			Deleting: []DeletingPlatform{{PostProfileID: "pp-1", Platform: PlatformTwitter}},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Posts.DeleteOnPlatform(context.Background(), "post-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("expected success=true")
	}
	if len(result.Deleting) != 1 || result.Deleting[0].Platform != PlatformTwitter {
		t.Errorf("unexpected deleting: %+v", result.Deleting)
	}
}

func TestPostsDeleteOnPlatform_ByNetwork(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		_ = json.Unmarshal(body, &payload)
		if payload["network"] != "twitter" {
			t.Errorf("expected network=twitter, got %v", payload["network"])
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(DeleteOnPlatformResponse{Success: true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	network := "twitter"
	_, err := c.Posts.DeleteOnPlatform(context.Background(), "post-1", &PostDeleteOnPlatformOptions{Network: &network})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostsDeleteOnPlatform_ByProfileID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		_ = json.Unmarshal(body, &payload)
		if payload["profile_id"] != "prof-1" {
			t.Errorf("expected profile_id=prof-1, got %v", payload["profile_id"])
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(DeleteOnPlatformResponse{Success: true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	profileID := "prof-1"
	_, err := c.Posts.DeleteOnPlatform(context.Background(), "post-1", &PostDeleteOnPlatformOptions{ProfileID: &profileID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostsDeleteOnPlatform_ByPostProfileID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		_ = json.Unmarshal(body, &payload)
		if payload["post_profile_id"] != "pp-1" {
			t.Errorf("expected post_profile_id=pp-1, got %v", payload["post_profile_id"])
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(DeleteOnPlatformResponse{Success: true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	ppID := "pp-1"
	_, err := c.Posts.DeleteOnPlatform(context.Background(), "post-1", &PostDeleteOnPlatformOptions{PostProfileID: &ppID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostsCreate_WithThread(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		_ = json.Unmarshal(body, &payload)

		thread, ok := payload["thread"].([]any)
		if !ok || len(thread) != 2 {
			t.Errorf("expected thread with 2 items, got %v", payload["thread"])
		}
		first := thread[0].(map[string]any)
		if first["body"] != "Reply 1" {
			t.Errorf("expected first thread body 'Reply 1', got %v", first["body"])
		}
		second := thread[1].(map[string]any)
		media := second["media"].([]any)
		if len(media) != 1 {
			t.Errorf("expected 1 media in second thread, got %d", len(media))
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Post{
			ID:   "post-thread",
			Body: "Main",
			Thread: []ThreadChild{
				{ID: "t-1", Body: "Reply 1"},
				{ID: "t-2", Body: "Reply 2"},
			},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	post, err := c.Posts.Create(context.Background(), "Main", []string{"prof-1"}, &PostCreateOptions{
		Thread: []ThreadChildInput{
			{Body: "Reply 1"},
			{Body: "Reply 2", Media: []string{"https://example.com/img.jpg"}},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(post.Thread) != 2 {
		t.Fatalf("expected 2 thread children, got %d", len(post.Thread))
	}
	if post.Thread[0].Body != "Reply 1" {
		t.Errorf("expected thread body 'Reply 1', got %q", post.Thread[0].Body)
	}
}

func TestPostsCreate_ThreadWithMediaFiles(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "thread-img.png")
	if err := os.WriteFile(tmpFile, []byte("fake png data"), 0644); err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("expected multipart/form-data, got %s", ct)
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatalf("failed to parse form: %v", err)
		}

		if r.FormValue("post[body]") != "Thread main" {
			t.Errorf("expected body 'Thread main', got %q", r.FormValue("post[body]"))
		}

		if r.FormValue("thread[0][body]") != "Reply 1" {
			t.Errorf("expected thread[0][body] 'Reply 1', got %q", r.FormValue("thread[0][body]"))
		}
		if r.FormValue("thread[1][body]") != "Reply 2" {
			t.Errorf("expected thread[1][body] 'Reply 2', got %q", r.FormValue("thread[1][body]"))
		}

		// Thread child 0 has a URL media
		urlMedia := r.MultipartForm.Value["thread[0][media][]"]
		if len(urlMedia) != 1 || urlMedia[0] != "https://example.com/img.jpg" {
			t.Errorf("expected thread[0] URL media, got %v", urlMedia)
		}

		// Thread child 1 has a file upload
		files := r.MultipartForm.File["thread[1][media][]"]
		if len(files) != 1 {
			t.Errorf("expected 1 file in thread[1], got %d", len(files))
		} else if files[0].Filename != "thread-img.png" {
			t.Errorf("expected filename thread-img.png, got %s", files[0].Filename)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Post{
			ID:   "post-thread-files",
			Body: "Thread main",
			Thread: []ThreadChild{
				{ID: "t-1", Body: "Reply 1"},
				{ID: "t-2", Body: "Reply 2"},
			},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	post, err := c.Posts.Create(context.Background(), "Thread main", []string{"prof-1"}, &PostCreateOptions{
		Thread: []ThreadChildInput{
			{Body: "Reply 1", Media: []string{"https://example.com/img.jpg"}},
			{Body: "Reply 2", MediaFiles: []string{tmpFile}},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.ID != "post-thread-files" {
		t.Errorf("expected post ID %q, got %q", "post-thread-files", post.ID)
	}
	if len(post.Thread) != 2 {
		t.Fatalf("expected 2 thread children, got %d", len(post.Thread))
	}
}

func TestPostsGet_WithMediaAndThread(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		imgURL := "https://cdn.example.com/img.jpg"
		_ = json.NewEncoder(w).Encode(Post{
			ID:     "post-1",
			Body:   "Hello",
			Status: PostStatusMediaProcessingFailed,
			Media: []Media{
				{ID: "m-1", Status: MediaStatusProcessed, ContentType: "image/jpeg", URL: &imgURL},
			},
			Thread: []ThreadChild{
				{ID: "t-1", Body: "Reply"},
			},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	post, err := c.Posts.Get(context.Background(), "post-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.Status != PostStatusMediaProcessingFailed {
		t.Errorf("expected status media_processing_failed, got %s", post.Status)
	}
	if len(post.Media) != 1 {
		t.Fatalf("expected 1 media, got %d", len(post.Media))
	}
	if post.Media[0].Status != MediaStatusProcessed {
		t.Errorf("expected media status processed, got %s", post.Media[0].Status)
	}
	if len(post.Thread) != 1 {
		t.Fatalf("expected 1 thread child, got %d", len(post.Thread))
	}
	if post.Thread[0].Body != "Reply" {
		t.Errorf("expected thread body 'Reply', got %q", post.Thread[0].Body)
	}
}

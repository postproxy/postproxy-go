package postproxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("test-key")
	if c.apiKey != "test-key" {
		t.Errorf("expected apiKey %q, got %q", "test-key", c.apiKey)
	}
	if c.baseURL != defaultBaseURL {
		t.Errorf("expected baseURL %q, got %q", defaultBaseURL, c.baseURL)
	}
	if c.Posts == nil || c.Profiles == nil || c.ProfileGroups == nil {
		t.Error("expected services to be initialized")
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	hc := &http.Client{}
	c := NewClient("key",
		WithBaseURL("https://custom.api"),
		WithHTTPClient(hc),
		WithProfileGroupID("pg-123"),
	)
	if c.baseURL != "https://custom.api" {
		t.Errorf("expected baseURL %q, got %q", "https://custom.api", c.baseURL)
	}
	if c.httpClient != hc {
		t.Error("expected custom HTTP client")
	}
	if c.profileGroupID != "pg-123" {
		t.Errorf("expected profileGroupID %q, got %q", "pg-123", c.profileGroupID)
	}
}

func TestClient_AuthHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := NewClient("my-api-key", WithBaseURL(srv.URL))
	_, _ = c.request(context.Background(), http.MethodGet, "/test")

	expected := "Bearer my-api-key"
	if gotAuth != expected {
		t.Errorf("expected Authorization %q, got %q", expected, gotAuth)
	}
}

func TestClient_ProfileGroupIDQueryParam(t *testing.T) {
	var gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("profile_group_id")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL), WithProfileGroupID("pg-default"))
	_, _ = c.request(context.Background(), http.MethodGet, "/test")

	if gotQuery != "pg-default" {
		t.Errorf("expected profile_group_id %q, got %q", "pg-default", gotQuery)
	}
}

func TestClient_ProfileGroupIDOverride(t *testing.T) {
	var gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("profile_group_id")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL), WithProfileGroupID("pg-default"))
	override := "pg-override"
	_, _ = c.request(context.Background(), http.MethodGet, "/test", withProfileGroupID(&override))

	if gotQuery != "pg-override" {
		t.Errorf("expected profile_group_id %q, got %q", "pg-override", gotQuery)
	}
}

func TestClient_APIPath(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	_, _ = c.request(context.Background(), http.MethodGet, "/posts")

	if gotPath != "/api/posts" {
		t.Errorf("expected path %q, got %q", "/api/posts", gotPath)
	}
}

func TestClient_ErrorParsing(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantMsg    string
		checkFn    func(error) bool
	}{
		{
			name:       "401 authentication error",
			statusCode: 401,
			body:       `{"error":"Invalid API key"}`,
			wantMsg:    "Invalid API key",
			checkFn:    IsAuthenticationError,
		},
		{
			name:       "404 not found",
			statusCode: 404,
			body:       `{"error":"Not found"}`,
			wantMsg:    "Not found",
			checkFn:    IsNotFoundError,
		},
		{
			name:       "422 validation error",
			statusCode: 422,
			body:       `{"message":"Validation failed"}`,
			wantMsg:    "Validation failed",
			checkFn:    IsValidationError,
		},
		{
			name:       "400 bad request",
			statusCode: 400,
			body:       `{"error":"Bad request"}`,
			wantMsg:    "Bad request",
			checkFn:    IsBadRequestError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			c := NewClient("key", WithBaseURL(srv.URL))
			_, err := c.request(context.Background(), http.MethodGet, "/test")

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			apiErr, ok := err.(*PostProxyError)
			if !ok {
				t.Fatalf("expected *PostProxyError, got %T", err)
			}
			if apiErr.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, apiErr.StatusCode)
			}
			if apiErr.Message != tt.wantMsg {
				t.Errorf("expected message %q, got %q", tt.wantMsg, apiErr.Message)
			}
			if !tt.checkFn(err) {
				t.Errorf("check function returned false for status %d", tt.statusCode)
			}
		})
	}
}

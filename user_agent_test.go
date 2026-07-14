package postproxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_UserAgentHeader(t *testing.T) {
	var gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := NewClient("k", WithBaseURL(srv.URL))
	_, _ = c.request(context.Background(), http.MethodGet, "/test")

	if !strings.HasPrefix(gotUA, "postproxy-go/"+Version) {
		t.Errorf("User-Agent does not have expected prefix: %q", gotUA)
	}
	if !strings.Contains(gotUA, "go/") {
		t.Errorf("User-Agent missing go runtime: %q", gotUA)
	}
}

func TestVersionConstant(t *testing.T) {
	if Version != "1.11.0" {
		t.Errorf("expected Version 1.11.0, got %q", Version)
	}
}

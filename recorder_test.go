package postproxy

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sync"
	"testing"
)

// recordedRequest is one request captured by newRecordingClient.
type recordedRequest struct {
	Method string
	Path   string
	Query  url.Values
	Header http.Header
	Body   []byte
}

// JSON decodes the request body, or returns nil when there is none.
func (r recordedRequest) JSON(t *testing.T) any {
	t.Helper()
	if len(r.Body) == 0 {
		return nil
	}
	var out any
	if err := json.Unmarshal(r.Body, &out); err != nil {
		t.Fatalf("request body is not JSON: %v (%s)", err, r.Body)
	}
	return out
}

// newRecordingClient starts a server that answers every request with
// response (encoded as JSON) and records what it received.
func newRecordingClient(t *testing.T, response any) (*Client, func() []recordedRequest) {
	t.Helper()
	var (
		mu       sync.Mutex
		requests []recordedRequest
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		requests = append(requests, recordedRequest{
			Method: r.Method,
			Path:   r.URL.Path,
			Query:  r.URL.Query(),
			Header: r.Header.Clone(),
			Body:   body,
		})
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	t.Cleanup(srv.Close)

	return NewClient("key", WithBaseURL(srv.URL)), func() []recordedRequest {
		mu.Lock()
		defer mu.Unlock()
		return append([]recordedRequest(nil), requests...)
	}
}

// assertJSON compares a decoded JSON value with the JSON literal want.
func assertJSON(t *testing.T, got any, want string) {
	t.Helper()
	var wantValue any
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatalf("bad want JSON: %v", err)
	}
	if !reflect.DeepEqual(got, wantValue) {
		gotJSON, _ := json.Marshal(got)
		t.Errorf("JSON mismatch\n got: %s\nwant: %s", gotJSON, want)
	}
}

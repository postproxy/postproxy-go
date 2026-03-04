package postproxy

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhooksList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/webhooks" {
			t.Errorf("expected path /api/webhooks, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ListResponse[Webhook]{
			Data: []Webhook{{
				ID:      "wh-1",
				URL:     "https://example.com/webhook",
				Events:  []string{"post.published", "post.failed"},
				Enabled: true,
			}},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Webhooks.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 webhook, got %d", len(result.Data))
	}
	if result.Data[0].ID != "wh-1" {
		t.Errorf("expected id wh-1, got %s", result.Data[0].ID)
	}
}

func TestWebhooksGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/webhooks/wh-1" {
			t.Errorf("expected path /api/webhooks/wh-1, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Webhook{
			ID:      "wh-1",
			URL:     "https://example.com/webhook",
			Events:  []string{"post.published"},
			Enabled: true,
			Secret:  strPtr("whsec_test123"),
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	webhook, err := c.Webhooks.Get(context.Background(), "wh-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if webhook.ID != "wh-1" {
		t.Errorf("expected id wh-1, got %s", webhook.ID)
	}
	if *webhook.Secret != "whsec_test123" {
		t.Errorf("expected secret whsec_test123, got %s", *webhook.Secret)
	}
}

func TestWebhooksCreate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/webhooks" {
			t.Errorf("expected path /api/webhooks, got %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		_ = json.Unmarshal(body, &payload)

		if payload["url"] != "https://example.com/webhook" {
			t.Errorf("expected url in body, got %v", payload["url"])
		}
		if payload["description"] != "Test webhook" {
			t.Errorf("expected description in body, got %v", payload["description"])
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Webhook{
			ID:     "wh-new",
			URL:    "https://example.com/webhook",
			Events: []string{"post.published"},
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	desc := "Test webhook"
	webhook, err := c.Webhooks.Create(context.Background(), "https://example.com/webhook", []string{"post.published"}, &WebhookCreateOptions{
		Description: &desc,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if webhook.ID != "wh-new" {
		t.Errorf("expected id wh-new, got %s", webhook.ID)
	}
}

func TestWebhooksUpdate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/webhooks/wh-1" {
			t.Errorf("expected path /api/webhooks/wh-1, got %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		_ = json.Unmarshal(body, &payload)

		if payload["enabled"] != false {
			t.Errorf("expected enabled=false in body, got %v", payload["enabled"])
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Webhook{
			ID:      "wh-1",
			Enabled: false,
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	enabled := false
	webhook, err := c.Webhooks.Update(context.Background(), "wh-1", &WebhookUpdateOptions{
		Enabled: &enabled,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if webhook.Enabled != false {
		t.Errorf("expected enabled=false, got %v", webhook.Enabled)
	}
}

func TestWebhooksDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/webhooks/wh-1" {
			t.Errorf("expected path /api/webhooks/wh-1, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(DeleteResponse{Deleted: true})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	result, err := c.Webhooks.Delete(context.Background(), "wh-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Deleted {
		t.Error("expected deleted=true")
	}
}

func TestWebhooksDeliveries(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/webhooks/wh-1/deliveries" {
			t.Errorf("expected path /api/webhooks/wh-1/deliveries, got %s", r.URL.Path)
		}

		q := r.URL.Query()
		if q.Get("page") != "1" {
			t.Errorf("expected page=1, got %s", q.Get("page"))
		}
		if q.Get("per_page") != "10" {
			t.Errorf("expected per_page=10, got %s", q.Get("per_page"))
		}

		status := 200
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(PaginatedResponse[WebhookDelivery]{
			Data: []WebhookDelivery{{
				ID:             "del-1",
				EventType:      "post.published",
				ResponseStatus: &status,
				Success:        true,
			}},
			Total:   1,
			Page:    1,
			PerPage: 10,
		})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	page := 1
	perPage := 10
	result, err := c.Webhooks.Deliveries(context.Background(), "wh-1", &WebhookDeliveryListOptions{
		Page:    &page,
		PerPage: &perPage,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 delivery, got %d", len(result.Data))
	}
	if result.Data[0].ID != "del-1" {
		t.Errorf("expected id del-1, got %s", result.Data[0].ID)
	}
	if !result.Data[0].Success {
		t.Error("expected success=true")
	}
}

func TestVerifyWebhookSignature_Valid(t *testing.T) {
	payload := `{"event":"post.published","data":{"id":"post-1"}}`
	secret := "whsec_test123"
	signature := "t=1234567890,v1=c8e99efbb07ac8e3152c02dd8d83e8ddb803ae8fb001d9e1ab42fb0b1f405ef2"
	if !VerifyWebhookSignature(payload, signature, secret) {
		t.Error("expected signature to be valid")
	}
}

func TestVerifyWebhookSignature_Invalid(t *testing.T) {
	payload := `{"event":"post.published","data":{"id":"post-1"}}`
	secret := "whsec_test123"
	signature := "t=1234567890,v1=invalidsignature"
	if VerifyWebhookSignature(payload, signature, secret) {
		t.Error("expected signature to be invalid")
	}
}

func TestVerifyWebhookSignature_WrongSecret(t *testing.T) {
	payload := `{"event":"post.published","data":{"id":"post-1"}}`
	secret := "wrong_secret"
	signature := "t=1234567890,v1=c8e99efbb07ac8e3152c02dd8d83e8ddb803ae8fb001d9e1ab42fb0b1f405ef2"
	if VerifyWebhookSignature(payload, signature, secret) {
		t.Error("expected signature to be invalid with wrong secret")
	}
}

func strPtr(s string) *string { return &s }

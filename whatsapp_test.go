package postproxy

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

const waBase = "/api/profiles/prof_wa/whatsapp"

var mockTemplate = map[string]any{
	"id":               "tpl_1",
	"external_id":      "998877",
	"profile_id":       "prof_wa",
	"name":             "order_update",
	"language":         "en_US",
	"category":         "UTILITY",
	"status":           "APPROVED",
	"parameter_format": "POSITIONAL",
	"components":       []any{map[string]any{"type": "BODY", "text": "Order {{1}} shipped"}},
	"variable_count":   1,
	"created_at":       "2026-09-28T00:00:00Z",
	"updated_at":       "2026-09-28T00:00:00Z",
}

func TestWhatsAppTemplatesList(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"data": []any{mockTemplate}})

	status := "APPROVED"
	refresh := true
	result, err := c.WhatsAppTemplates.List(context.Background(), "prof_wa", &WhatsAppTemplateListOptions{Status: &status, Refresh: &refresh})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Data[0].Name != "order_update" || result.Data[0].VariableCount != 1 {
		t.Errorf("unexpected template: %+v", result.Data[0])
	}

	req := requests()[0]
	if req.Path != "/api/profiles/prof_wa/templates" {
		t.Errorf("unexpected path: %s", req.Path)
	}
	if req.Query.Get("status") != "APPROVED" || req.Query.Get("refresh") != "true" {
		t.Errorf("unexpected query: %v", req.Query)
	}
}

func TestWhatsAppTemplatesGetByName(t *testing.T) {
	var rawPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawPath = r.URL.EscapedPath()
		if r.URL.Query().Get("language") != "en_US" {
			t.Errorf("expected language=en_US, got %v", r.URL.Query())
		}
		_, _ = w.Write([]byte(`{"id":"tpl_1","name":"order update"}`))
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	lang := "en_US"
	if _, err := c.WhatsAppTemplates.Get(context.Background(), "prof_wa", "order update", &WhatsAppTemplateOptions{Language: &lang}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rawPath != "/api/profiles/prof_wa/templates/order%20update" {
		t.Errorf("unexpected path: %s", rawPath)
	}
}

func TestWhatsAppTemplatesCreateUpdateDelete(t *testing.T) {
	c, requests := newRecordingClient(t, mockTemplate)
	ctx := context.Background()

	ttl := 3600
	_, _ = c.WhatsAppTemplates.Create(ctx, "prof_wa", WhatsAppTemplateInput{
		Name:       "order_update",
		Category:   "UTILITY",
		Language:   "en_US",
		Components: []map[string]any{{"type": "BODY", "text": "Order {{1}} shipped"}},
	}, nil)
	_, _ = c.WhatsAppTemplates.Update(ctx, "prof_wa", "tpl_1", WhatsAppTemplateUpdate{MessageSendTTLSeconds: &ttl}, nil)
	_, _ = c.WhatsAppTemplates.Delete(ctx, "prof_wa", "order_update", nil)

	reqs := requests()
	if reqs[0].Method != "POST" {
		t.Errorf("expected POST, got %s", reqs[0].Method)
	}
	assertJSON(t, reqs[0].JSON(t), `{"name":"order_update","category":"UTILITY","language":"en_US","components":[{"type":"BODY","text":"Order {{1}} shipped"}]}`)
	if reqs[1].Method != "PATCH" || reqs[1].Path != "/api/profiles/prof_wa/templates/tpl_1" {
		t.Errorf("unexpected request: %s %s", reqs[1].Method, reqs[1].Path)
	}
	assertJSON(t, reqs[1].JSON(t), `{"message_send_ttl_seconds":3600}`)
	if reqs[2].Method != "DELETE" || reqs[2].Query.Has("language") {
		t.Errorf("unexpected request: %s %v", reqs[2].Method, reqs[2].Query)
	}
}

func TestWhatsAppTemplatesLibrary(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"data": []any{}})

	name := "delivery"
	if _, err := c.WhatsAppTemplates.Library(context.Background(), "prof_wa", &WhatsAppTemplateLibraryOptions{Name: &name}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req := requests()[0]
	if req.Path != "/api/profiles/prof_wa/templates/library" || req.Query.Get("name") != "delivery" {
		t.Errorf("unexpected request: %s %v", req.Path, req.Query)
	}
}

func TestWhatsAppReadsUseQueryParams(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"verified_name": "Acme"})
	ctx := context.Background()

	limit := 10
	_, _ = c.WhatsApp.NumberInfo(ctx, "prof_wa", "pn_1", nil)
	_, _ = c.WhatsApp.BlockedUsers(ctx, "prof_wa", "pn_1", &WhatsAppPageOptions{Limit: &limit})

	reqs := requests()
	if reqs[0].Path != waBase+"/number_info" || reqs[0].Query.Get("phone_number_id") != "pn_1" {
		t.Errorf("unexpected request: %s %v", reqs[0].Path, reqs[0].Query)
	}
	if reqs[1].Query.Get("limit") != "10" || reqs[1].Query.Has("after") {
		t.Errorf("unexpected query: %v", reqs[1].Query)
	}
}

func TestWhatsAppWritesUseJSON(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"success": true})
	ctx := context.Background()

	method, lang := "SMS", "en"
	about := "Hi"
	transfer := "force_transfer"
	_, _ = c.WhatsApp.RegisterPhoneNumber(ctx, "prof_wa", "pn_1", "123456", nil)
	_, _ = c.WhatsApp.RequestVerificationCode(ctx, "prof_wa", "pn_1", &WhatsAppVerificationCodeOptions{Method: &method, Language: &lang})
	_, _ = c.WhatsApp.VerifyPhoneNumber(ctx, "prof_wa", "pn_1", "000111", nil)
	_, _ = c.WhatsApp.UpdateBusinessProfile(ctx, "prof_wa", "pn_1", WhatsAppBusinessProfile{About: &about, Websites: []string{"https://acme.test"}}, nil)
	_, _ = c.WhatsApp.UpdateBusinessProfilePhoto(ctx, "prof_wa", "pn_1", WhatsAppProfilePhoto{URL: "https://img.test/a.png"}, nil)
	_, _ = c.WhatsApp.RequestDisplayNameChange(ctx, "prof_wa", "pn_1", "Acme Co", nil)
	_, _ = c.WhatsApp.SetUsername(ctx, "prof_wa", "pn_1", "acme", &WhatsAppSetUsernameOptions{TransferAction: &transfer})
	_, _ = c.WhatsApp.BlockUsers(ctx, "prof_wa", "pn_1", []string{"15551234567"}, nil)

	reqs := requests()
	wantMethods := []string{"POST", "POST", "POST", "PATCH", "POST", "POST", "PUT", "POST"}
	for i, m := range wantMethods {
		if reqs[i].Method != m {
			t.Errorf("request %d: expected %s, got %s", i, m, reqs[i].Method)
		}
	}
	assertJSON(t, reqs[0].JSON(t), `{"phone_number_id":"pn_1","pin":"123456"}`)
	assertJSON(t, reqs[1].JSON(t), `{"phone_number_id":"pn_1","method":"SMS","language":"en"}`)
	assertJSON(t, reqs[3].JSON(t), `{"phone_number_id":"pn_1","about":"Hi","websites":["https://acme.test"]}`)
	assertJSON(t, reqs[4].JSON(t), `{"phone_number_id":"pn_1","url":"https://img.test/a.png"}`)
	assertJSON(t, reqs[5].JSON(t), `{"phone_number_id":"pn_1","display_name":"Acme Co"}`)
	assertJSON(t, reqs[6].JSON(t), `{"phone_number_id":"pn_1","username":"acme","transfer_action":"force_transfer"}`)
	assertJSON(t, reqs[7].JSON(t), `{"phone_number_id":"pn_1","users":["15551234567"]}`)
}

func TestWhatsAppDeleteUsername(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"success": true})

	if _, err := c.WhatsApp.DeleteUsername(context.Background(), "prof_wa", "pn_1", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req := requests()[0]
	if req.Method != "DELETE" || len(req.Body) != 0 {
		t.Errorf("unexpected request: %s body=%q", req.Method, req.Body)
	}
	if req.Query.Get("phone_number_id") != "pn_1" {
		t.Errorf("unexpected query: %v", req.Query)
	}
}

func TestWhatsAppGroups(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{})
	ctx := context.Background()

	mode := "approval_required"
	subject := "VIPs"
	_, _ = c.WhatsApp.CreateGroup(ctx, "prof_wa", "pn_1", "Customers", &WhatsAppGroupCreateOptions{JoinApprovalMode: &mode})
	_, _ = c.WhatsApp.Group(ctx, "prof_wa", "pn_1", "g_1", nil)
	_, _ = c.WhatsApp.UpdateGroup(ctx, "prof_wa", "pn_1", "g_1", WhatsAppGroupUpdate{Subject: &subject}, nil)
	_, _ = c.WhatsApp.AddGroupParticipants(ctx, "prof_wa", "pn_1", "g_1", []string{"15551234567"}, nil)
	_, _ = c.WhatsApp.CreateGroupInviteLink(ctx, "prof_wa", "pn_1", "g_1", nil)
	_, _ = c.WhatsApp.ApproveGroupJoinRequests(ctx, "prof_wa", "pn_1", "g_1", []string{"15557654321"}, nil)
	_, _ = c.WhatsApp.DeleteGroup(ctx, "prof_wa", "pn_1", "g_1", nil)

	reqs := requests()
	assertJSON(t, reqs[0].JSON(t), `{"phone_number_id":"pn_1","subject":"Customers","join_approval_mode":"approval_required"}`)
	if reqs[1].Query.Get("group_id") != "g_1" {
		t.Errorf("unexpected query: %v", reqs[1].Query)
	}
	assertJSON(t, reqs[2].JSON(t), `{"phone_number_id":"pn_1","group_id":"g_1","subject":"VIPs"}`)
	if reqs[3].Path != waBase+"/add_group_participants" {
		t.Errorf("unexpected path: %s", reqs[3].Path)
	}
	assertJSON(t, reqs[3].JSON(t), `{"phone_number_id":"pn_1","group_id":"g_1","phone_numbers":["15551234567"]}`)
	if reqs[4].Path != waBase+"/create_group_invite_link" {
		t.Errorf("unexpected path: %s", reqs[4].Path)
	}
	if reqs[5].Path != waBase+"/approve_group_join_requests" {
		t.Errorf("unexpected path: %s", reqs[5].Path)
	}
	if reqs[6].Method != "DELETE" || reqs[6].Query.Get("group_id") != "g_1" {
		t.Errorf("unexpected request: %s %v", reqs[6].Method, reqs[6].Query)
	}
}

func TestWhatsAppConversionEvents(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"events_received": 1})
	ctx := context.Background()

	value := 19.99
	_, _ = c.WhatsApp.Dataset(ctx, "prof_wa", nil)
	_, _ = c.WhatsApp.SendConversionEvent(ctx, "prof_wa", WhatsAppConversionEvent{
		EventName: "Purchase",
		EventID:   "order-42",
		ChatID:    "chat_1",
		Currency:  "USD",
		Value:     &value,
	}, nil)

	reqs := requests()
	if reqs[0].Query.Has("phone_number_id") {
		t.Errorf("unexpected phone_number_id: %v", reqs[0].Query)
	}
	if reqs[1].Path != waBase+"/send_conversion_event" {
		t.Errorf("unexpected path: %s", reqs[1].Path)
	}
	assertJSON(t, reqs[1].JSON(t), `{"event_name":"Purchase","event_id":"order-42","chat_id":"chat_1","currency":"USD","value":19.99}`)
}

func TestWhatsAppMediaReturnsRawBytes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != waBase+"/media" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("phone_number_id") != "pn_1" || r.URL.Query().Get("media_id") != "media-1" {
			t.Errorf("unexpected query: %v", r.URL.Query())
		}
		w.Header().Set("Content-Type", "image/jpeg")
		_, _ = w.Write([]byte{1, 2, 3})
	}))
	defer srv.Close()

	c := NewClient("key", WithBaseURL(srv.URL))
	media, err := c.WhatsApp.Media(context.Background(), "prof_wa", "pn_1", "media-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if media.ContentType != "image/jpeg" {
		t.Errorf("expected image/jpeg, got %q", media.ContentType)
	}
	if !bytes.Equal(media.Data, []byte{1, 2, 3}) {
		t.Errorf("unexpected data: %v", media.Data)
	}
}

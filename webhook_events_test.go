package postproxy

import (
	"encoding/json"
	"errors"
	"testing"
)

func envelope(t *testing.T, eventType string, data any) []byte {
	t.Helper()
	body, err := json.Marshal(map[string]any{
		"id":         "evt_1",
		"type":       eventType,
		"created_at": "2026-05-12T00:00:00Z",
		"data":       data,
	})
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func TestParseWebhookEvent_PostProcessed(t *testing.T) {
	body := envelope(t, "post.processed", map[string]any{
		"id": "p1", "body": "hi", "status": "processed",
		"scheduled_at": nil, "created_at": "2026-05-12T00:00:00Z",
		"platforms": []any{
			map[string]any{"id": "pf1", "platform": "twitter", "name": "X"},
		},
	})

	event, err := ParseWebhookEvent(body)
	if err != nil {
		t.Fatal(err)
	}
	if event.Type != EventPostProcessed {
		t.Errorf("expected post.processed, got %q", event.Type)
	}
	data, err := event.AsPostProcessed()
	if err != nil {
		t.Fatal(err)
	}
	if data.Body != "hi" || len(data.Platforms) != 1 {
		t.Errorf("unexpected data: %+v", data)
	}
	if data.Platforms[0].Platform != PlatformTwitter {
		t.Errorf("expected twitter, got %q", data.Platforms[0].Platform)
	}
}

func TestParseWebhookEvent_PlatformPostVariants(t *testing.T) {
	for _, eventType := range []WebhookEventType{
		EventPlatformPostPublished,
		EventPlatformPostFailed,
		EventPlatformPostFailedWaitingRetry,
		EventPlatformPostInsights,
	} {
		body := envelope(t, string(eventType), map[string]any{
			"id": "pp1", "post_id": "p1", "platform": "twitter",
			"profile_id": "pf1", "profile_name": "X", "status": "published",
			"error": nil, "error_details": nil, "platform_id": "tw_999",
		})
		event, err := ParseWebhookEvent(body)
		if err != nil {
			t.Fatalf("type %s: %v", eventType, err)
		}
		if event.Type != eventType {
			t.Errorf("expected %q, got %q", eventType, event.Type)
		}
		data, err := event.AsPlatformPost()
		if err != nil {
			t.Fatal(err)
		}
		if data.ID != "pp1" {
			t.Errorf("unexpected data: %+v", data)
		}
	}
}

func TestParseWebhookEvent_ProfileStats(t *testing.T) {
	body := envelope(t, "profile.stats", map[string]any{
		"profile_id": "pf1", "platform": "linkedin",
		"placement_id": "org_1",
		"stats":        map[string]any{"followerCount": 4567},
		"recorded_at":  "2026-05-12T00:00:00Z",
	})

	event, err := ParseWebhookEvent(body)
	if err != nil {
		t.Fatal(err)
	}
	data, err := event.AsProfileStats()
	if err != nil {
		t.Fatal(err)
	}
	if data.PlacementID == nil || *data.PlacementID != "org_1" {
		t.Errorf("placement_id mismatch: %+v", data.PlacementID)
	}
	if got, _ := data.Stats["followerCount"].(float64); got != 4567 {
		t.Errorf("followerCount mismatch: %+v", data.Stats["followerCount"])
	}
}

func TestParseWebhookEvent_CommentCreated(t *testing.T) {
	body := envelope(t, "comment.created", map[string]any{
		"id": "c1", "post_id": "p1", "platform_post_id": "pp1",
		"platform": "instagram", "external_id": "ig_c",
		"parent_external_id": nil, "body": "great", "status": "published",
		"author_external_id": "u", "author_name": "Jane",
		"author_username": "jane", "author_avatar_url": nil,
		"like_count": 0, "reply_count": 0, "is_hidden": false,
		"permalink": nil, "platform_data": nil, "posted_at": nil,
		"created_at": "2026-05-12T00:00:00Z",
	})
	event, err := ParseWebhookEvent(body)
	if err != nil {
		t.Fatal(err)
	}
	data, err := event.AsCommentCreated()
	if err != nil {
		t.Fatal(err)
	}
	if data.AuthorName == nil || *data.AuthorName != "Jane" {
		t.Errorf("author_name mismatch: %+v", data.AuthorName)
	}
}

func mockMessage() map[string]any {
	return map[string]any{
		"id": "msg_1", "chat_id": "chat_1", "external_id": "mid.1",
		"direction": "inbound", "body": "hi", "status": "received",
		"reactions": []any{}, "attachments": []any{},
		"is_unsupported": false, "created_at": "2026-06-01T00:00:00Z",
	}
}

func TestParseWebhookEvent_MessageVariants(t *testing.T) {
	for _, eventType := range []WebhookEventType{
		EventMessageReceived,
		EventMessageSent,
		EventMessageDelivered,
		EventMessageRead,
		EventMessageEdited,
		EventMessageDeleted,
		EventMessageFailedWaitingForRetry,
		EventMessageFailed,
	} {
		body := envelope(t, string(eventType), map[string]any{"message": mockMessage()})
		event, err := ParseWebhookEvent(body)
		if err != nil {
			t.Fatalf("type %s: %v", eventType, err)
		}
		if event.Type != eventType {
			t.Errorf("expected %q, got %q", eventType, event.Type)
		}
		data, err := event.AsMessage()
		if err != nil {
			t.Fatal(err)
		}
		if data.Message.ID != "msg_1" {
			t.Errorf("unexpected message: %+v", data.Message)
		}
	}
}

func TestParseWebhookEvent_ReactionReceived(t *testing.T) {
	msg := mockMessage()
	msg["reactions"] = []map[string]any{
		{"sender_external_id": "psid_123", "emoji": "❤️", "reaction": "love", "at": "2026-06-01T15:02:00Z"},
	}
	body := envelope(t, "reaction.received", map[string]any{
		"message":            msg,
		"sender_external_id": "psid_123",
		"action":             "react",
		"reaction":           "love",
		"emoji":              "❤️",
		"occurred_at":        "2026-06-01T15:02:00Z",
	})
	event, err := ParseWebhookEvent(body)
	if err != nil {
		t.Fatal(err)
	}
	if event.Type != EventReactionReceived {
		t.Errorf("expected reaction.received, got %q", event.Type)
	}
	data, err := event.AsReaction()
	if err != nil {
		t.Fatal(err)
	}
	if data.Action != "react" {
		t.Errorf("expected action react, got %q", data.Action)
	}
	if len(data.Message.Reactions) != 1 || data.Message.Reactions[0].Reaction == nil || *data.Message.Reactions[0].Reaction != "love" {
		t.Errorf("unexpected reactions: %+v", data.Message.Reactions)
	}
}

func TestParseWebhookEvent_ProfileCommentCreated(t *testing.T) {
	body := envelope(t, "profile_comment.created", map[string]any{
		"id": "abc123", "profile_id": "prof123", "platform": "google_business",
		"placement_id":       "accounts/1/locations/2",
		"external_id":        "accounts/1/locations/2/reviews/A",
		"parent_external_id": nil, "body": "Great coffee!", "status": "synced",
		"author_username": "Jane D.", "author_avatar_url": nil,
		"platform_data": map[string]any{"star_rating": 5},
		"posted_at":     "2026-05-10T11:55:00Z",
		"created_at":    "2026-05-13T18:00:00Z",
	})
	event, err := ParseWebhookEvent(body)
	if err != nil {
		t.Fatal(err)
	}
	if event.Type != EventProfileCommentCreated {
		t.Errorf("expected profile_comment.created, got %q", event.Type)
	}
	data, err := event.AsProfileCommentCreated()
	if err != nil {
		t.Fatal(err)
	}
	if data.ID != "abc123" {
		t.Errorf("expected id abc123, got %s", data.ID)
	}
	if got, _ := data.PlatformData["star_rating"].(float64); got != 5 {
		t.Errorf("expected star_rating 5, got %v", data.PlatformData["star_rating"])
	}
}

func TestParseWebhookEvent_UnknownType(t *testing.T) {
	body := envelope(t, "foo.bar", map[string]any{})
	_, err := ParseWebhookEvent(body)
	if !errors.Is(err, ErrUnknownWebhookEvent) {
		t.Errorf("expected ErrUnknownWebhookEvent, got %v", err)
	}
}

func TestParseWebhookEvent_InvalidJSON(t *testing.T) {
	_, err := ParseWebhookEvent([]byte("not json{"))
	if err == nil {
		t.Error("expected error on invalid JSON")
	}
}

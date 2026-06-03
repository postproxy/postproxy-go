package postproxy

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrUnknownWebhookEvent is returned when ParseWebhookEvent sees a `type` field
// that doesn't match any known WebhookEventType.
var ErrUnknownWebhookEvent = errors.New("unknown webhook event type")

// WebhookEvent is the envelope returned by ParseWebhookEvent. Use Type to choose
// which typed payload accessor to call.
type WebhookEvent struct {
	ID        string           `json:"id"`
	Type      WebhookEventType `json:"type"`
	CreatedAt string           `json:"created_at"`
	Data      json.RawMessage  `json:"data"`
}

// AsPostProcessed decodes Data as a PostProcessedData. Use when Type == EventPostProcessed.
func (e *WebhookEvent) AsPostProcessed() (*PostProcessedData, error) {
	return decodeAs[PostProcessedData](e.Data)
}

// AsPostImported decodes Data as a PostImportedData. Use when Type == EventPostImported.
func (e *WebhookEvent) AsPostImported() (*PostImportedData, error) {
	return decodeAs[PostImportedData](e.Data)
}

// AsPlatformPost decodes Data as a PlatformPostData. Use for platform_post.* events.
func (e *WebhookEvent) AsPlatformPost() (*PlatformPostData, error) {
	return decodeAs[PlatformPostData](e.Data)
}

// AsProfileEvent decodes Data as a ProfileEventData. Use for profile.connected/disconnected.
func (e *WebhookEvent) AsProfileEvent() (*ProfileEventData, error) {
	return decodeAs[ProfileEventData](e.Data)
}

// AsProfileStats decodes Data as a ProfileStatsData. Use when Type == EventProfileStats.
func (e *WebhookEvent) AsProfileStats() (*ProfileStatsData, error) {
	return decodeAs[ProfileStatsData](e.Data)
}

// AsMediaFailed decodes Data as a MediaFailedData. Use when Type == EventMediaFailed.
func (e *WebhookEvent) AsMediaFailed() (*MediaFailedData, error) {
	return decodeAs[MediaFailedData](e.Data)
}

// AsCommentCreated decodes Data as a CommentCreatedData. Use when Type == EventCommentCreated.
func (e *WebhookEvent) AsCommentCreated() (*CommentCreatedData, error) {
	return decodeAs[CommentCreatedData](e.Data)
}

// AsMessage decodes Data as a MessageEventData. Use for message.* events.
func (e *WebhookEvent) AsMessage() (*MessageEventData, error) {
	return decodeAs[MessageEventData](e.Data)
}

// AsReaction decodes Data as a ReactionEventData. Use when Type == EventReactionReceived.
func (e *WebhookEvent) AsReaction() (*ReactionEventData, error) {
	return decodeAs[ReactionEventData](e.Data)
}

// AsProfileCommentCreated decodes Data as a ProfileCommentCreatedData. Use when Type == EventProfileCommentCreated.
func (e *WebhookEvent) AsProfileCommentCreated() (*ProfileCommentCreatedData, error) {
	return decodeAs[ProfileCommentCreatedData](e.Data)
}

func decodeAs[T any](raw json.RawMessage) (*T, error) {
	var out T
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PostProcessedPlatform is the per-platform entry inside a post.processed event.
type PostProcessedPlatform struct {
	ID       string   `json:"id"`
	Platform Platform `json:"platform"`
	Name     string   `json:"name"`
}

// PostProcessedData is the payload of a post.processed event.
type PostProcessedData struct {
	ID          string                  `json:"id"`
	Body        string                  `json:"body"`
	Status      PostStatus              `json:"status"`
	ScheduledAt *string                 `json:"scheduled_at"`
	CreatedAt   string                  `json:"created_at"`
	Platforms   []PostProcessedPlatform `json:"platforms"`
}

// ImportedProfile is the profile field of a post.imported event.
type ImportedProfile struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Platform Platform `json:"platform"`
}

// PostImportedData is the payload of a post.imported event.
type PostImportedData struct {
	ID             string          `json:"id"`
	Body           string          `json:"body"`
	Source         string          `json:"source"`
	PostedAt       *string         `json:"posted_at"`
	CreatedAt      string          `json:"created_at"`
	Platform       Platform        `json:"platform"`
	Profile        ImportedProfile `json:"profile"`
	PlatformPostID string          `json:"platform_post_id"`
	PublicID       *string         `json:"public_id"`
}

// PlatformPostData is the payload shared by platform_post.* events.
type PlatformPostData struct {
	ID           string             `json:"id"`
	PostID       string             `json:"post_id"`
	Platform     Platform           `json:"platform"`
	ProfileID    string             `json:"profile_id"`
	ProfileName  string             `json:"profile_name"`
	Status       PlatformPostStatus `json:"status"`
	Error        *string            `json:"error"`
	ErrorDetails *ErrorDetails      `json:"error_details"`
	PlatformID   *string            `json:"platform_id"`
	Insights     map[string]any     `json:"insights,omitempty"`
}

// ProfileEventData is the payload of profile.connected and profile.disconnected events.
type ProfileEventData struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Platform       Platform `json:"platform"`
	ProfileGroupID string   `json:"profile_group_id"`
	Status         string   `json:"status"`
	UID            string   `json:"uid"`
	Username       *string  `json:"username"`
}

// ProfileStatsData is the payload of a profile.stats event.
type ProfileStatsData struct {
	ProfileID   string         `json:"profile_id"`
	Platform    Platform       `json:"platform"`
	PlacementID *string        `json:"placement_id"`
	Stats       map[string]any `json:"stats"`
	RecordedAt  string         `json:"recorded_at"`
}

// MediaFailedData is the payload of a media.failed event.
type MediaFailedData struct {
	ID           string `json:"id"`
	PostID       string `json:"post_id"`
	ContentType  string `json:"content_type"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}

// CommentCreatedData is the payload of a comment.created event.
type CommentCreatedData struct {
	ID               string         `json:"id"`
	PostID           string         `json:"post_id"`
	PlatformPostID   string         `json:"platform_post_id"`
	Platform         Platform       `json:"platform"`
	ExternalID       *string        `json:"external_id"`
	ParentExternalID *string        `json:"parent_external_id"`
	Body             string         `json:"body"`
	Status           string         `json:"status"`
	AuthorExternalID *string        `json:"author_external_id"`
	AuthorName       *string        `json:"author_name"`
	AuthorUsername   *string        `json:"author_username"`
	AuthorAvatarURL  *string        `json:"author_avatar_url"`
	LikeCount        int            `json:"like_count"`
	ReplyCount       int            `json:"reply_count"`
	IsHidden         bool           `json:"is_hidden"`
	Permalink        *string        `json:"permalink"`
	PlatformData     map[string]any `json:"platform_data"`
	PostedAt         *string        `json:"posted_at"`
	CreatedAt        string         `json:"created_at"`
}

// MessageEventData is the payload shared by all message.* events.
type MessageEventData struct {
	Message Message `json:"message"`
}

// ReactionEventData is the payload of a reaction.received event.
type ReactionEventData struct {
	Message          Message `json:"message"`
	SenderExternalID string  `json:"sender_external_id"`
	Action           string  `json:"action"`
	Reaction         *string `json:"reaction"`
	Emoji            *string `json:"emoji"`
	OccurredAt       *string `json:"occurred_at"`
}

// ProfileCommentCreatedData is the payload of a profile_comment.created event.
type ProfileCommentCreatedData struct {
	ID               string         `json:"id"`
	ProfileID        string         `json:"profile_id"`
	Platform         Platform       `json:"platform"`
	PlacementID      *string        `json:"placement_id"`
	ExternalID       *string        `json:"external_id"`
	ParentExternalID *string        `json:"parent_external_id"`
	Body             string         `json:"body"`
	Status           string         `json:"status"`
	AuthorUsername   *string        `json:"author_username"`
	AuthorAvatarURL  *string        `json:"author_avatar_url"`
	PlatformData     map[string]any `json:"platform_data"`
	PostedAt         *string        `json:"posted_at"`
	CreatedAt        string         `json:"created_at"`
}

// knownEventTypes is the lookup set used by ParseWebhookEvent.
var knownEventTypes = func() map[WebhookEventType]struct{} {
	m := make(map[WebhookEventType]struct{}, len(WebhookEventTypes))
	for _, t := range WebhookEventTypes {
		m[t] = struct{}{}
	}
	return m
}()

// ParseWebhookEvent parses a raw webhook body and returns a WebhookEvent. The
// returned event's Type is guaranteed to be one of WebhookEventTypes; use the
// As* helpers to decode Data into a typed payload.
func ParseWebhookEvent(body []byte) (*WebhookEvent, error) {
	var event WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("invalid webhook body: %w", err)
	}
	if _, ok := knownEventTypes[event.Type]; !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownWebhookEvent, event.Type)
	}
	return &event, nil
}

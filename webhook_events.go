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

// AsPlatformPost decodes Data as a PlatformPostData. Use for platform_post.*
// events, including platform_post.edited and platform_post.edit_failed.
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

// AsCommentCreated decodes Data as a CommentCreatedData. Use for comment.*
// events (created, failed, failed_waiting_for_retry, delivery_unconfirmed,
// edited, edit_failed).
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

// AsProfileCommentCreated decodes Data as a ProfileCommentCreatedData. Use for
// profile_comment.* events (created, failed, failed_waiting_for_retry, deleted).
func (e *WebhookEvent) AsProfileCommentCreated() (*ProfileCommentCreatedData, error) {
	return decodeAs[ProfileCommentCreatedData](e.Data)
}

// AsPlacementsSynced decodes Data as a PlacementsSyncedData. Use when Type == EventProfilePlacementsSynced.
func (e *WebhookEvent) AsPlacementsSynced() (*PlacementsSyncedData, error) {
	return decodeAs[PlacementsSyncedData](e.Data)
}

// AsPlacement decodes Data as a PlacementEventData. Use for placement.* events
// (orphaned, restored, name_status_updated, quality_updated).
func (e *WebhookEvent) AsPlacement() (*PlacementEventData, error) {
	return decodeAs[PlacementEventData](e.Data)
}

// AsReferral decodes Data as a ReferralEventData. Use when Type == EventReferralReceived.
func (e *WebhookEvent) AsReferral() (*ReferralEventData, error) {
	return decodeAs[ReferralEventData](e.Data)
}

// AsMessageTemplate decodes Data as a MessageTemplateEventData. Use for
// message_template.* events.
func (e *WebhookEvent) AsMessageTemplate() (*MessageTemplateEventData, error) {
	return decodeAs[MessageTemplateEventData](e.Data)
}

// AsWhatsAppAutomaticEvent decodes Data as a WhatsAppAutomaticEventData. Use
// when Type == EventWhatsAppAutomaticEvent.
func (e *WebhookEvent) AsWhatsAppAutomaticEvent() (*WhatsAppAutomaticEventData, error) {
	return decodeAs[WhatsAppAutomaticEventData](e.Data)
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
	PlatformID     string   `json:"platform_id"`
	Username       *string  `json:"username"`
	// Origin is set on profile.disconnected: "user" when a user disconnected
	// it, "automated" when PostProxy did.
	Origin *string `json:"origin,omitempty"`
}

// PlacementEventData is the payload of placement.* events. Name-status and
// quality events add the fields below their names.
type PlacementEventData struct {
	// ID is the placement's id on the platform (page id, phone number id, ...).
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Platform       Platform       `json:"platform"`
	ProfileGroupID *string        `json:"profile_group_id"`
	ProfileID      string         `json:"profile_id"`
	ProfileName    string         `json:"profile_name"`
	Metadata       map[string]any `json:"metadata"`
	// placement.name_status_updated
	Decision        *string `json:"decision,omitempty"`
	RequestedName   *string `json:"requested_name,omitempty"`
	RejectionReason *string `json:"rejection_reason,omitempty"`
	// placement.quality_updated
	Event        *string `json:"event,omitempty"`
	CurrentLimit any     `json:"current_limit,omitempty"`
	OldLimit     any     `json:"old_limit,omitempty"`
}

// PlacementsSyncedData is the payload of a profile.placements_synced event.
type PlacementsSyncedData struct {
	ProfileEventData
	PlacementsSyncStatus *string              `json:"placements_sync_status"`
	Placements           []PlacementEventData `json:"placements"`
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
	ProfileID        *string        `json:"profile_id,omitempty"`
	Error            *string        `json:"error,omitempty"`
	ErrorDetails     *ErrorDetails  `json:"error_details,omitempty"`
	ExternalEditedAt *string        `json:"external_edited_at,omitempty"`
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
	Error            *string        `json:"error,omitempty"`
	ErrorDetails     *ErrorDetails  `json:"error_details,omitempty"`
	Permalink        *string        `json:"permalink,omitempty"`
	PostedAt         *string        `json:"posted_at"`
	CreatedAt        string         `json:"created_at"`
}

// WebhookChatData is the chat carried by a referral.received event.
type WebhookChatData struct {
	ID                    string   `json:"id"`
	ProfileID             string   `json:"profile_id"`
	Platform              Platform `json:"platform"`
	ExternalPlacementID   *string  `json:"external_placement_id"`
	Group                 bool     `json:"group"`
	ParticipantExternalID string   `json:"participant_external_id"`
	ParticipantUsername   *string  `json:"participant_username"`
	ParticipantName       *string  `json:"participant_name"`
	LastInboundAt         *string  `json:"last_inbound_at"`
	LastOutboundAt        *string  `json:"last_outbound_at"`
	LastMessageAt         *string  `json:"last_message_at"`
}

// ReferralEventData is the payload of a referral.received event.
type ReferralEventData struct {
	Chat   WebhookChatData `json:"chat"`
	Ref    *string         `json:"ref,omitempty"`
	Source *string         `json:"source,omitempty"`
	Type   *string         `json:"type,omitempty"`
	// CtwaClid is the WhatsApp click-to-WhatsApp ad click id.
	CtwaClid   *string `json:"ctwa_clid,omitempty"`
	OccurredAt *string `json:"occurred_at,omitempty"`
}

// MessageTemplateEventData is the payload of message_template.* events.
type MessageTemplateEventData struct {
	ID             string   `json:"id"`
	ExternalID     *string  `json:"external_id"`
	ProfileID      string   `json:"profile_id"`
	Platform       Platform `json:"platform"`
	Name           string   `json:"name"`
	Language       string   `json:"language"`
	Category       string   `json:"category"`
	Status         string   `json:"status"`
	RejectedReason *string  `json:"rejected_reason"`
	QualityScore   *string  `json:"quality_score"`
}

// WhatsAppAutomaticEventData is the payload of a whatsapp.automatic_event
// event — a conversion Meta detected in a chat.
type WhatsAppAutomaticEventData struct {
	ProfileID             string         `json:"profile_id"`
	ChatID                *string        `json:"chat_id"`
	ParticipantExternalID *string        `json:"participant_external_id"`
	ExternalMessageID     *string        `json:"external_message_id"`
	EventName             *string        `json:"event_name"`
	CtwaClid              *string        `json:"ctwa_clid"`
	CustomData            map[string]any `json:"custom_data"`
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

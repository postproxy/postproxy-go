package postproxy

// Profile represents a social media profile.
type Profile struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Status         ProfileStatus `json:"status"`
	Platform       Platform      `json:"platform"`
	ProfileGroupID string        `json:"profile_group_id"`
	ExpiresAt      *string       `json:"expires_at"`
	PostCount      int           `json:"post_count"`
}

// ProfileGroup represents a group of profiles.
type ProfileGroup struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ProfilesCount int    `json:"profiles_count"`
}

// Insights represents engagement insights for a platform result.
type Insights struct {
	Impressions *int    `json:"impressions"`
	On          *string `json:"on"`
}

// ErrorDetails represents structured error details from the platform.
type ErrorDetails struct {
	PlatformErrorCode    *string `json:"platform_error_code"`
	PlatformErrorSubcode *string `json:"platform_error_subcode"`
	PlatformErrorMessage *string `json:"platform_error_message"`
	PostproxyNote        *string `json:"postproxy_note"`
}

// PlatformResult represents the result of posting to a specific platform.
type PlatformResult struct {
	Platform     Platform           `json:"platform"`
	Status       PlatformPostStatus `json:"status"`
	Params       map[string]any     `json:"params"`
	Error        *string            `json:"error"`
	ErrorDetails *ErrorDetails      `json:"error_details"`
	AttemptedAt  *string            `json:"attempted_at"`
	Insights     *Insights          `json:"insights"`
}

// MediaPlatformError represents a per-platform error encountered while uploading a media item.
type MediaPlatformError struct {
	Platform     Platform           `json:"platform"`
	Status       PlatformPostStatus `json:"status"`
	Error        *string            `json:"error"`
	ErrorDetails *ErrorDetails      `json:"error_details"`
}

// Media represents a media attachment on a post.
type Media struct {
	ID           string               `json:"id"`
	Status       MediaStatus          `json:"status"`
	ErrorMessage *string              `json:"error_message"`
	ContentType  string               `json:"content_type"`
	SourceURL    *string              `json:"source_url"`
	URL          *string              `json:"url"`
	Platforms    []MediaPlatformError `json:"platforms,omitempty"`
}

// ThreadChild represents a child post in a thread.
type ThreadChild struct {
	ID    string  `json:"id"`
	Body  string  `json:"body"`
	Media []Media `json:"media"`
}

// ThreadChildInput represents a thread child in a create request.
type ThreadChildInput struct {
	Body       string   `json:"body"`
	Media      []string `json:"media,omitempty"`
	MediaFiles []string `json:"-"`
}

// Post represents a social media post.
type Post struct {
	ID            string           `json:"id"`
	Body          string           `json:"body"`
	Status        PostStatus       `json:"status"`
	ScheduledAt   *string          `json:"scheduled_at"`
	CreatedAt     string           `json:"created_at"`
	Media         []Media          `json:"media"`
	Platforms     []PlatformResult `json:"platforms"`
	Thread        []ThreadChild    `json:"thread"`
	QueueID       *string          `json:"queue_id"`
	QueuePriority *string          `json:"queue_priority"`
}

// Placement represents a placement option for a profile.
type Placement struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Metadata map[string]any `json:"metadata,omitempty"`
	// ProfileGroupID is set in the response of AssignPlacementToGroup.
	ProfileGroupID string `json:"profile_group_id,omitempty"`
}

// IceBreaker represents an Instagram DM ice breaker (FAQ prompt).
type IceBreaker struct {
	Question string `json:"question"`
	Payload  string `json:"payload"`
}

// IceBreakersResponse represents the response from listing ice breakers.
type IceBreakersResponse struct {
	IceBreakers []IceBreaker `json:"ice_breakers"`
}

// ListResponse is a generic response containing a list of items.
type ListResponse[T any] struct {
	Data []T `json:"data"`
}

// PaginatedResponse is a generic paginated response containing a list of items.
type PaginatedResponse[T any] struct {
	Data    []T `json:"data"`
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// DeleteResponse represents the response from a delete operation.
type DeleteResponse struct {
	Deleted bool `json:"deleted"`
}

// DeletingPlatform represents a post profile scheduled for platform deletion.
type DeletingPlatform struct {
	PostProfileID string   `json:"post_profile_id"`
	Platform      Platform `json:"platform"`
}

// DeleteOnPlatformResponse represents the response from deleting a post from platforms.
type DeleteOnPlatformResponse struct {
	Success  bool               `json:"success"`
	Deleting []DeletingPlatform `json:"deleting"`
}

// SuccessResponse represents a generic success response.
type SuccessResponse struct {
	Success bool `json:"success"`
}

// ConnectionResponse represents the response from an OAuth initialize connection operation.
type ConnectionResponse struct {
	URL     string `json:"url"`
	Success bool   `json:"success"`
}

// SyncProfile is the profile returned by non-OAuth connection flows (BlueSky, Telegram).
type SyncProfile struct {
	ID               string   `json:"id"`
	Network          Platform `json:"network"`
	Name             string   `json:"name"`
	ExternalUsername *string  `json:"external_username"`
}

// BlueskyConnectionResponse is returned by ConnectBluesky.
type BlueskyConnectionResponse struct {
	Success bool        `json:"success"`
	Profile SyncProfile `json:"profile"`
}

// TelegramConnectionResponse is returned by ConnectTelegram. Channels populate
// asynchronously — poll ProfilesService.Placements until non-empty.
type TelegramConnectionResponse struct {
	Success  bool        `json:"success"`
	Profile  SyncProfile `json:"profile"`
	NextStep *string     `json:"next_step,omitempty"`
}

// ProfileStats holds a profile's stats timeseries returned by GetProfileStats.
type ProfileStats struct {
	ProfileID   string        `json:"profile_id"`
	Platform    Platform      `json:"platform"`
	PlacementID *string       `json:"placement_id"`
	Records     []StatsRecord `json:"records"`
}

// ProfileStatsResponse wraps the data field returned by GET /profiles/:id/stats.
type ProfileStatsResponse struct {
	Data ProfileStats `json:"data"`
}

// ProfileStatsOptions are the optional query parameters for GetProfileStats.
type ProfileStatsOptions struct {
	PlacementID    *string
	From           *string
	To             *string
	ProfileGroupID *string
}

// FacebookParams contains Facebook-specific post parameters.
type FacebookParams struct {
	Format       *FacebookFormat `json:"format,omitempty"`
	Title        *string         `json:"title,omitempty"`
	FirstComment *string         `json:"first_comment,omitempty"`
	PageID       *string         `json:"page_id,omitempty"`
}

// InstagramParams contains Instagram-specific post parameters.
type InstagramParams struct {
	Format        *InstagramFormat   `json:"format,omitempty"`
	FirstComment  *string            `json:"first_comment,omitempty"`
	Collaborators []string           `json:"collaborators,omitempty"`
	CoverURL      *string            `json:"cover_url,omitempty"`
	AudioName     *string            `json:"audio_name,omitempty"`
	TrialStrategy *bool              `json:"trial_strategy,omitempty"`
	ThumbOffset   *int               `json:"thumb_offset,omitempty"`
	UserTags      []InstagramUserTag `json:"user_tags,omitempty"`
}

// InstagramUserTag is an Instagram account to tag in a post. Images require X
// and Y; reels and video slides are tagged by username only — Instagram
// ignores coordinates there.
type InstagramUserTag struct {
	Username string   `json:"username"`
	X        *float64 `json:"x,omitempty"`
	Y        *float64 `json:"y,omitempty"`
	// MediaIndex picks which media item to tag, 0-based. Defaults to 0.
	MediaIndex *int `json:"media_index,omitempty"`
}

// TikTokParams contains TikTok-specific post parameters.
type TikTokParams struct {
	Format             *TikTokFormat  `json:"format,omitempty"`
	PrivacyStatus      *TikTokPrivacy `json:"privacy_status,omitempty"`
	PhotoCoverIndex    *int           `json:"photo_cover_index,omitempty"`
	AutoAddMusic       *bool          `json:"auto_add_music,omitempty"`
	MadeWithAI         *bool          `json:"made_with_ai,omitempty"`
	DisableComment     *bool          `json:"disable_comment,omitempty"`
	DisableDuet        *bool          `json:"disable_duet,omitempty"`
	DisableStitch      *bool          `json:"disable_stitch,omitempty"`
	BrandContentToggle *bool          `json:"brand_content_toggle,omitempty"`
	BrandOrganicToggle *bool          `json:"brand_organic_toggle,omitempty"`
}

// LinkedInParams contains LinkedIn-specific post parameters.
type LinkedInParams struct {
	Format         *LinkedInFormat `json:"format,omitempty"`
	OrganizationID *string         `json:"organization_id,omitempty"`
}

// YouTubeParams contains YouTube-specific post parameters.
type YouTubeParams struct {
	Format                 *YouTubeFormat  `json:"format,omitempty"`
	Title                  *string         `json:"title,omitempty"`
	PrivacyStatus          *YouTubePrivacy `json:"privacy_status,omitempty"`
	CoverURL               *string         `json:"cover_url,omitempty"`
	MadeForKids            *bool           `json:"made_for_kids,omitempty"`
	Tags                   []string        `json:"tags,omitempty"`
	CategoryID             *string         `json:"category_id,omitempty"`
	ContainsSyntheticMedia *bool           `json:"contains_synthetic_media,omitempty"`
}

// PinterestParams contains Pinterest-specific post parameters.
type PinterestParams struct {
	Format          *PinterestFormat `json:"format,omitempty"`
	Title           *string          `json:"title,omitempty"`
	BoardID         *string          `json:"board_id,omitempty"`
	DestinationLink *string          `json:"destination_link,omitempty"`
	CoverURL        *string          `json:"cover_url,omitempty"`
	ThumbOffset     *int             `json:"thumb_offset,omitempty"`
}

// ThreadsParams contains Threads-specific post parameters.
type ThreadsParams struct {
	Format *ThreadsFormat `json:"format,omitempty"`
}

// TwitterParams contains Twitter-specific post parameters.
type TwitterParams struct {
	Format *TwitterFormat `json:"format,omitempty"`
	// PollOptions is required when Format is "poll": 2-4 options, max 25 characters each.
	PollOptions []string `json:"poll_options,omitempty"`
	// PollDurationMinutes is required when Format is "poll": 5 to 10080 minutes (7 days).
	PollDurationMinutes *int `json:"poll_duration_minutes,omitempty"`
}

// BlueskyParams contains Bluesky-specific post parameters.
type BlueskyParams struct {
	Format *BlueskyFormat `json:"format,omitempty"`
}

// TelegramParams contains Telegram-specific post parameters. ChatID is required.
type TelegramParams struct {
	Format              *TelegramFormat    `json:"format,omitempty"`
	ChatID              string             `json:"chat_id"`
	ParseMode           *TelegramParseMode `json:"parse_mode,omitempty"`
	DisableLinkPreview  *bool              `json:"disable_link_preview,omitempty"`
	DisableNotification *bool              `json:"disable_notification,omitempty"`
}

// PlatformParams contains platform-specific parameters for a post.
type PlatformParams struct {
	Facebook       *FacebookParams  `json:"facebook,omitempty"`
	Instagram      *InstagramParams `json:"instagram,omitempty"`
	TikTok         *TikTokParams    `json:"tiktok,omitempty"`
	LinkedIn       *LinkedInParams  `json:"linkedin,omitempty"`
	YouTube        *YouTubeParams   `json:"youtube,omitempty"`
	Pinterest      *PinterestParams `json:"pinterest,omitempty"`
	Threads        *ThreadsParams   `json:"threads,omitempty"`
	Twitter        *TwitterParams   `json:"twitter,omitempty"`
	Bluesky        *BlueskyParams   `json:"bluesky,omitempty"`
	Telegram       *TelegramParams  `json:"telegram,omitempty"`
	GoogleBusiness map[string]any   `json:"google_business,omitempty"`
}

// Webhook represents a webhook subscription.
type Webhook struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Events      []string `json:"events"`
	Enabled     bool     `json:"enabled"`
	Description *string  `json:"description"`
	Secret      *string  `json:"secret,omitempty"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// WebhookDelivery represents a webhook delivery attempt.
type WebhookDelivery struct {
	ID             string `json:"id"`
	EventID        string `json:"event_id"`
	EventType      string `json:"event_type"`
	ResponseStatus *int   `json:"response_status"`
	AttemptNumber  int    `json:"attempt_number"`
	Success        bool   `json:"success"`
	AttemptedAt    string `json:"attempted_at"`
	CreatedAt      string `json:"created_at"`
}

// WebhookCreateOptions contains options for creating a webhook.
type WebhookCreateOptions struct {
	Description *string
}

// WebhookUpdateOptions contains options for updating a webhook.
type WebhookUpdateOptions struct {
	URL         *string
	Events      []string
	Enabled     *bool
	Description *string
}

// WebhookDeliveryListOptions contains options for listing webhook deliveries.
type WebhookDeliveryListOptions struct {
	Page    *int
	PerPage *int
}

// PostStatsOptions contains options for retrieving post stats.
type PostStatsOptions struct {
	Profiles []string
	From     *string
	To       *string
}

// StatsResponse is the response from the post stats endpoint.
type StatsResponse struct {
	Data map[string]PostStats `json:"data"`
}

// PostStats contains stats for a single post across platforms.
type PostStats struct {
	Platforms []PlatformStats `json:"platforms"`
}

// PlatformStats contains stats records for a single platform.
type PlatformStats struct {
	ProfileID string        `json:"profile_id"`
	Platform  Platform      `json:"platform"`
	Records   []StatsRecord `json:"records"`
}

// StatsRecord is a single stats snapshot at a point in time.
type StatsRecord struct {
	Stats map[string]any `json:"stats"`
	// RawStats carries every metric under its original platform name, e.g.
	// "views" for Instagram or "impression_count" for Twitter/X.
	RawStats   map[string]any `json:"raw_stats"`
	RecordedAt string         `json:"recorded_at"`
}

// Attachment represents a media attachment on a comment or message.
type Attachment struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	URL        *string     `json:"url"`
	Status     MediaStatus `json:"status"`
	ExternalID *string     `json:"external_id"`
}

// Comment represents a comment on a post.
type Comment struct {
	ID               string         `json:"id"`
	ExternalID       *string        `json:"external_id"`
	Body             string         `json:"body"`
	Status           string         `json:"status"`
	AuthorUsername   *string        `json:"author_username"`
	AuthorAvatarURL  *string        `json:"author_avatar_url"`
	AuthorExternalID *string        `json:"author_external_id"`
	Metadata         map[string]any `json:"metadata"`
	ParentExternalID *string        `json:"parent_external_id"`
	LikeCount        int            `json:"like_count"`
	IsHidden         bool           `json:"is_hidden"`
	Permalink        *string        `json:"permalink"`
	PlatformData     any            `json:"platform_data"`
	Attachments      []Attachment   `json:"attachments"`
	PostedAt         *string        `json:"posted_at"`
	CreatedAt        string         `json:"created_at"`
	Replies          []Comment      `json:"replies"`
}

// BulkComment is a comment from CommentsService.ListAll. Flat: replies are
// their own entries linked to their parent by ParentExternalID rather than
// nested under Replies.
type BulkComment struct {
	PostID           string         `json:"post_id"`
	ProfileID        string         `json:"profile_id"`
	Platform         Platform       `json:"platform"`
	ID               string         `json:"id"`
	ExternalID       *string        `json:"external_id"`
	Body             string         `json:"body"`
	Status           string         `json:"status"`
	AuthorUsername   *string        `json:"author_username"`
	AuthorAvatarURL  *string        `json:"author_avatar_url"`
	AuthorExternalID *string        `json:"author_external_id"`
	Metadata         map[string]any `json:"metadata"`
	ParentExternalID *string        `json:"parent_external_id"`
	LikeCount        int            `json:"like_count"`
	IsHidden         bool           `json:"is_hidden"`
	Permalink        *string        `json:"permalink"`
	PlatformData     any            `json:"platform_data"`
	Attachments      []Attachment   `json:"attachments"`
	PostedAt         *string        `json:"posted_at"`
	CreatedAt        string         `json:"created_at"`
}

// PostSync records one post pull for a profile — the sync fired when the
// profile connects, the recurring poll, or a backfill.
type PostSync struct {
	ID          string          `json:"id"`
	ProfileID   string          `json:"profile_id"`
	Kind        string          `json:"kind"`
	Trigger     PostSyncTrigger `json:"trigger"`
	Status      PostSyncStatus  `json:"status"`
	StartedAt   *string         `json:"started_at"`
	CompletedAt *string         `json:"completed_at"`
	PostsSeen   int             `json:"posts_seen"`
	// PostsImported counts posts that were new and got created — lower than
	// PostsSeen whenever the run re-read posts you already have.
	PostsImported int     `json:"posts_imported"`
	BackfillFrom  *string `json:"backfill_from"`
	// OldestPostedAt is the publish date of the oldest post the run reached.
	OldestPostedAt *string `json:"oldest_posted_at"`
	Error          *string `json:"error"`
	CreatedAt      string  `json:"created_at"`
}

// Reaction represents an emoji reaction on a message.
type Reaction struct {
	SenderExternalID string  `json:"sender_external_id"`
	Emoji            *string `json:"emoji"`
	Reaction         *string `json:"reaction"`
	At               *string `json:"at"`
}

// Chat represents a direct-message conversation with a participant.
type Chat struct {
	ID                     string         `json:"id"`
	ProfileID              string         `json:"profile_id"`
	Platform               Platform       `json:"platform"`
	ParticipantExternalID  string         `json:"participant_external_id"`
	ParticipantUsername    *string        `json:"participant_username"`
	ParticipantName        *string        `json:"participant_name"`
	ParticipantAvatarURL   *string        `json:"participant_avatar_url"`
	ExternalConversationID *string        `json:"external_conversation_id"`
	LastInboundAt          *string        `json:"last_inbound_at"`
	LastOutboundAt         *string        `json:"last_outbound_at"`
	LastMessageAt          *string        `json:"last_message_at"`
	Metadata               map[string]any `json:"metadata"`
	Archived               *bool          `json:"archived,omitempty"`
	CreatedAt              string         `json:"created_at"`
}

// Message represents a direct message within a chat.
type Message struct {
	ID                  string           `json:"id"`
	ChatID              string           `json:"chat_id"`
	ExternalID          *string          `json:"external_id"`
	Direction           MessageDirection `json:"direction"`
	Body                *string          `json:"body"`
	Status              MessageStatus    `json:"status"`
	Tag                 *string          `json:"tag"`
	ExternalCommentID   *string          `json:"external_comment_id"`
	ErrorMessage        *string          `json:"error_message"`
	PlatformData        map[string]any   `json:"platform_data"`
	ExternalPostedAt    *string          `json:"external_posted_at"`
	ExternalDeliveredAt *string          `json:"external_delivered_at"`
	ExternalReadAt      *string          `json:"external_read_at"`
	ExternalEditedAt    *string          `json:"external_edited_at"`
	ReplyToExternalID   *string          `json:"reply_to_external_id"`
	ReplyMarkup         map[string]any   `json:"reply_markup"`
	ExternalDeletedAt   *string          `json:"external_deleted_at"`
	Reactions           []Reaction       `json:"reactions"`
	Attachments         []Attachment     `json:"attachments"`
	IsUnsupported       bool             `json:"is_unsupported"`
	CreatedAt           string           `json:"created_at"`
}

// ChatListOptions contains options for listing chats.
type ChatListOptions struct {
	Page           *int
	PerPage        *int
	Before         *string
	After          *string
	ProfileGroupID *string
}

// ChatCreateOptions contains options for creating a chat.
type ChatCreateOptions struct {
	ParticipantUsername *string
	ParticipantName     *string
	ProfileGroupID      *string
}

// MessageListOptions contains options for listing messages.
type MessageListOptions struct {
	Page           *int
	PerPage        *int
	Direction      *MessageDirection
	Status         *MessageStatus
	ProfileGroupID *string
}

// MessageSendOptions contains options for sending a message. If MediaFiles is set,
// the request uses multipart/form-data (mirrors PostsService.Create).
type MessageSendOptions struct {
	Body              *string
	Media             []string
	MediaFiles        []string
	Tag               *string
	ReplyToExternalID *string
	ReplyMarkup       map[string]any
	ProfileGroupID    *string
}

// MessageEditOptions contains options for editing a message.
type MessageEditOptions struct {
	Body           *string
	ReplyMarkup    map[string]any
	ProfileGroupID *string
}

// MessageReactOptions contains options for reacting to a message.
type MessageReactOptions struct {
	Reaction       *string
	Emoji          *string
	ProfileGroupID *string
}

// ProfileComment represents a profile-level comment (Google Business review or reply).
type ProfileComment struct {
	ID               string           `json:"id"`
	ExternalID       string           `json:"external_id"`
	ParentExternalID *string          `json:"parent_external_id"`
	PlacementID      string           `json:"placement_id"`
	Body             string           `json:"body"`
	Status           string           `json:"status"`
	AuthorUsername   *string          `json:"author_username"`
	AuthorAvatarURL  *string          `json:"author_avatar_url"`
	PlatformData     map[string]any   `json:"platform_data"`
	PostedAt         string           `json:"posted_at"`
	CreatedAt        string           `json:"created_at"`
	Replies          []ProfileComment `json:"replies"`
}

// ProfileCommentListOptions contains options for listing profile comments.
type ProfileCommentListOptions struct {
	PlacementID    *string
	Page           *int
	PerPage        *int
	ProfileGroupID *string
}

// CommentListOptions contains options for listing comments.
//
// From and To filter on when PostProxy received the comment (created_at), not
// the platform's posted_at. They apply to top-level comments — one in range
// brings its full Replies slice with it.
type CommentListOptions struct {
	Page    *int
	PerPage *int
	From    *string
	To      *string
}

// BulkCommentListOptions contains options for CommentsService.ListAll. Every
// field is optional. Profiles takes profile IDs or network names, mixed.
type BulkCommentListOptions struct {
	PostIDs        []string
	Profiles       []string
	From           *string
	To             *string
	Page           *int
	PerPage        *int
	ProfileGroupID *string
}

// PostSyncListOptions contains options for ProfilesService.PostSyncs.
type PostSyncListOptions struct {
	Trigger        *PostSyncTrigger
	Status         *PostSyncStatus
	Page           *int
	PerPage        *int
	ProfileGroupID *string
}

// CommentCreateOptions contains options for creating a comment.
type CommentCreateOptions struct {
	ParentID *string
}

// AcceptedResponse represents a response indicating an action was accepted.
type AcceptedResponse struct {
	Accepted bool `json:"accepted"`
}

// RequestOptions contains common options for API requests.
type RequestOptions struct {
	ProfileGroupID *string
}

// PostListOptions contains options for listing posts.
type PostListOptions struct {
	Page           *int
	PerPage        *int
	Status         *PostStatus
	Platforms      []Platform
	ScheduledAfter *string
	ProfileGroupID *string
}

// PostCreateOptions contains options for creating a post.
type PostCreateOptions struct {
	Media          []string
	MediaFiles     []string
	Platforms      *PlatformParams
	Thread         []ThreadChildInput
	ScheduledAt    *string
	Draft          *bool
	QueueID        *string
	QueuePriority  *string
	ProfileGroupID *string
}

// PostDeleteOptions contains options for deleting a post.
type PostDeleteOptions struct {
	// DeleteOnPlatform, when true, also deletes the post from all published platforms before removing it from the database.
	DeleteOnPlatform *bool
	ProfileGroupID   *string
}

// PostDeleteOnPlatformOptions contains options for deleting a post from social platforms.
type PostDeleteOnPlatformOptions struct {
	// PostProfileID targets a specific post profile. Covers the entire thread for the underlying profile.
	PostProfileID *string
	// ProfileID targets all post profiles for a profile on the post.
	ProfileID *string
	// Network targets all post profiles for a network (e.g. "twitter").
	Network        *string
	ProfileGroupID *string
}

// PostUpdateOptions contains options for updating a post.
type PostUpdateOptions struct {
	Body           *string
	Profiles       []string
	Media          []string
	MediaFiles     []string
	Platforms      *PlatformParams
	Thread         []ThreadChildInput
	ScheduledAt    *string
	Draft          *bool
	QueueID        *string
	QueuePriority  *string
	ProfileGroupID *string
}

// Timeslot represents a weekly recurring timeslot in a queue.
type Timeslot struct {
	ID   int    `json:"id"`
	Day  int    `json:"day"`
	Time string `json:"time"`
}

// Queue represents a post queue.
type Queue struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Description    *string    `json:"description"`
	Timezone       string     `json:"timezone"`
	Enabled        bool       `json:"enabled"`
	Jitter         int        `json:"jitter"`
	ProfileGroupID string     `json:"profile_group_id"`
	Timeslots      []Timeslot `json:"timeslots"`
	PostsCount     int        `json:"posts_count"`
}

// NextSlotResponse represents the response from the next slot endpoint.
type NextSlotResponse struct {
	NextSlot string `json:"next_slot"`
}

// TimeslotInput represents a timeslot in a create/update request.
type TimeslotInput struct {
	Day     int    `json:"day,omitempty"`
	Time    string `json:"time,omitempty"`
	ID      int    `json:"id,omitempty"`
	Destroy bool   `json:"_destroy,omitempty"`
}

// QueueCreateOptions contains options for creating a queue.
type QueueCreateOptions struct {
	Description *string
	Timezone    *string
	Jitter      *int
	Timeslots   []TimeslotInput
}

// QueueUpdateOptions contains options for updating a queue.
type QueueUpdateOptions struct {
	Name        *string
	Description *string
	Timezone    *string
	Enabled     *bool
	Jitter      *int
	Timeslots   []TimeslotInput
}

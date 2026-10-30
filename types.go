package postproxy

import "encoding/json"

// Profile represents a social media profile.
type Profile struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Username *string       `json:"username"`
	Status   ProfileStatus `json:"status"`
	// InactiveReason says why an inactive profile was deactivated;
	// ProfileInactiveReasonNetworkTaken means the same account was connected elsewhere.
	InactiveReason *ProfileInactiveReason `json:"inactive_reason"`
	Platform       Platform               `json:"platform"`
	// PlatformID is the account's id on the platform.
	PlatformID           *string `json:"platform_id"`
	ProfileGroupID       string  `json:"profile_group_id"`
	ExpiresAt            *string `json:"expires_at"`
	PostCount            int     `json:"post_count"`
	AvatarURL            *string `json:"avatar_url"`
	PlatformURL          *string `json:"platform_url"`
	PlacementsSyncStatus *string `json:"placements_sync_status"`
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

// PlatformAttachment is a file PostProxy generated for one platform, e.g. a TikTok cover.
type PlatformAttachment struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Status      string  `json:"status"`
	ContentType *string `json:"content_type"`
	URL         *string `json:"url"`
}

// PlatformResult represents the result of posting to a specific platform.
type PlatformResult struct {
	Platform     Platform             `json:"platform"`
	ProfileID    string               `json:"profile_id"`
	ProfileName  string               `json:"profile_name"`
	Status       PlatformPostStatus   `json:"status"`
	Permalink    *string              `json:"permalink"`
	Params       map[string]any       `json:"params"`
	Error        *string              `json:"error"`
	ErrorDetails *ErrorDetails        `json:"error_details"`
	AttemptedAt  *string              `json:"attempted_at"`
	Insights     *Insights            `json:"insights"`
	Attachments  []PlatformAttachment `json:"attachments,omitempty"`
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
	// Platforms shows whether this thread item went out — each one is
	// delivered separately.
	Platforms []PlatformResult `json:"platforms"`
}

// ThreadChildInput represents a thread child in a create request.
type ThreadChildInput struct {
	Body       string   `json:"body"`
	Media      []string `json:"media,omitempty"`
	MediaFiles []string `json:"-"`
}

// Post represents a social media post.
type Post struct {
	ID          string     `json:"id"`
	Body        string     `json:"body"`
	Status      PostStatus `json:"status"`
	ScheduledAt *string    `json:"scheduled_at"`
	CreatedAt   string     `json:"created_at"`
	Draft       *bool      `json:"draft,omitempty"`
	// Source is "postproxy" or "imported".
	Source        *string          `json:"source,omitempty"`
	Media         []Media          `json:"media"`
	Platforms     []PlatformResult `json:"platforms"`
	Thread        []ThreadChild    `json:"thread"`
	QueueID       *string          `json:"queue_id"`
	QueuePriority *string          `json:"queue_priority"`
}

// Placement represents a placement option for a profile.
type Placement struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	AvatarURL *string        `json:"avatar_url,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	// ProfileGroupID is set in the response of AssignPlacementToGroup.
	ProfileGroupID string `json:"profile_group_id,omitempty"`
}

// AccountPlacementProfile is a connected profile a placement is reachable through.
type AccountPlacementProfile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Connected bool   `json:"connected"`
}

// AccountPlacement is an account-level placement, returned by
// PlacementsService.List. A placement (a Facebook Page, LinkedIn organization,
// WhatsApp number, ...) can be reachable through several connected profiles;
// one of them is active.
type AccountPlacement struct {
	ID             string                    `json:"id"`
	Platform       Platform                  `json:"platform"`
	ExternalID     string                    `json:"external_id"`
	Name           string                    `json:"name"`
	AvatarURL      *string                   `json:"avatar_url"`
	Status         string                    `json:"status"`
	ActiveProfile  AccountPlacementProfile   `json:"active_profile"`
	LinkedProfiles []AccountPlacementProfile `json:"linked_profiles"`
}

// APIKey is a profile-group-scoped API key.
type APIKey struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	ProfileGroupID *string `json:"profile_group_id"`
	Sandbox        bool    `json:"sandbox"`
	// Token is only returned by APIKeysService.Create — store it, it can't be
	// read again.
	Token     *string `json:"token,omitempty"`
	CreatedAt string  `json:"created_at"`
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

// AltText is alt text for a post's media, one entry per media item in order.
// A single entry is sent as a plain string, which the API applies to the first
// media item.
type AltText []string

// MarshalJSON sends a single entry as a string and several as an array.
func (a AltText) MarshalJSON() ([]byte, error) {
	if len(a) == 1 {
		return json.Marshal(a[0])
	}
	return json.Marshal([]string(a))
}

// UnmarshalJSON accepts a string or an array of strings.
func (a *AltText) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*a = AltText{single}
		return nil
	}
	var list []string
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	*a = list
	return nil
}

// FacebookParams contains Facebook-specific post parameters.
type FacebookParams struct {
	Format       *FacebookFormat `json:"format,omitempty"`
	Title        *string         `json:"title,omitempty"`
	FirstComment *string         `json:"first_comment,omitempty"`
	PageID       *string         `json:"page_id,omitempty"`
	// Link makes a link post: Facebook renders the preview card. Text-only, no media.
	Link    *string `json:"link,omitempty"`
	AltText AltText `json:"alt_text,omitempty"`
	// MadeWithAI applies to photo posts only; videos and reels ignore it.
	MadeWithAI *bool `json:"made_with_ai,omitempty"`
	// Countries are ISO 3166-1 alpha-2 codes; viewers elsewhere won't see the post.
	Countries []string `json:"countries,omitempty"`
}

// InstagramParams contains Instagram-specific post parameters.
type InstagramParams struct {
	Format        *InstagramFormat `json:"format,omitempty"`
	FirstComment  *string          `json:"first_comment,omitempty"`
	Collaborators []string         `json:"collaborators,omitempty"`
	CoverURL      *string          `json:"cover_url,omitempty"`
	// CoverFile is a local path to a cover image, uploaded with the post. A
	// base64 data URL ("data:image/jpeg;base64,...") is sent as-is instead.
	CoverFile *string `json:"cover_file,omitempty"`
	AudioName *string `json:"audio_name,omitempty"`
	// AudioID is a music or sound id from ProfilesService.InstagramAvailableAudio.
	// Profiles connected via Facebook login only.
	AudioID *string `json:"audio_id,omitempty"`
	// AudioVolume is 1-100, default 100. Facebook-login profiles only.
	AudioVolume *int `json:"audio_volume,omitempty"`
	VideoVolume *int `json:"video_volume,omitempty"`
	// StripAudio removes the video's own audio track (reel, story).
	StripAudio    *bool                   `json:"strip_audio,omitempty"`
	TrialStrategy *InstagramTrialStrategy `json:"trial_strategy,omitempty"`
	ThumbOffset   *int                    `json:"thumb_offset,omitempty"`
	UserTags      []InstagramUserTag      `json:"user_tags,omitempty"`
	AltText       AltText                 `json:"alt_text,omitempty"`
	MadeWithAI    *bool                   `json:"made_with_ai,omitempty"`
	// DisableComment turns comments off on the post, applied right after
	// publishing. Post, reel.
	DisableComment *bool `json:"disable_comment,omitempty"`
	// LocationID is the numeric id of a Facebook Page with an address. Post,
	// reel; not stories.
	LocationID *string `json:"location_id,omitempty"`
	// ShareToFeed is reels only. false publishes to the Reels tab only, not
	// the feed/grid.
	ShareToFeed *bool `json:"share_to_feed,omitempty"`
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
	// UploadAsDraft sends the video to the creator's TikTok inbox as a draft
	// instead of publishing it.
	UploadAsDraft *bool `json:"upload_as_draft,omitempty"`
	// CoverURL is a custom video cover. Profiles connected before 2026-09-30
	// must reconnect.
	CoverURL *string `json:"cover_url,omitempty"`
	// CoverFile is a local path to a cover image, uploaded with the post. A
	// base64 data URL ("data:image/jpeg;base64,...") is sent as-is instead.
	CoverFile *string `json:"cover_file,omitempty"`
	// ThumbOffset is the frame to use as cover, in milliseconds. CoverURL wins
	// if both are set.
	ThumbOffset *int `json:"thumb_offset,omitempty"`
}

// LinkedInParams contains LinkedIn-specific post parameters.
type LinkedInParams struct {
	Format         *LinkedInFormat `json:"format,omitempty"`
	OrganizationID *string         `json:"organization_id,omitempty"`
	// ReshareURN is the URN of a post to reshare, e.g. "urn:li:share:123".
	ReshareURN *string `json:"reshare_urn,omitempty"`
	AltText    AltText `json:"alt_text,omitempty"`
	// Countries applies to organization pages only; up to 25 ISO 3166-1 alpha-2 codes.
	Countries []string `json:"countries,omitempty"`
	// PollQuestion is required when Format is "poll": 140 characters or fewer.
	PollQuestion *string `json:"poll_question,omitempty"`
	// PollOptions is required when Format is "poll": 2-4 options, max 30 characters each.
	PollOptions []string `json:"poll_options,omitempty"`
	// PollDurationDays is 1, 3, 7 (default) or 14.
	PollDurationDays *int `json:"poll_duration_days,omitempty"`
	// DocumentTitle is shown above an attached document. Defaults to the file name.
	DocumentTitle *string `json:"document_title,omitempty"`
	// FirstComment is posted right after publishing, as the same author (1,250
	// characters or fewer). Not on polls; profiles connected before 2026-10-28
	// must reconnect.
	FirstComment *string `json:"first_comment,omitempty"`
}

// YouTubeParams contains YouTube-specific post parameters.
type YouTubeParams struct {
	Format        *YouTubeFormat  `json:"format,omitempty"`
	Title         *string         `json:"title,omitempty"`
	PrivacyStatus *YouTubePrivacy `json:"privacy_status,omitempty"`
	CoverURL      *string         `json:"cover_url,omitempty"`
	// CoverFile is a local path to a cover image, uploaded with the post. A
	// base64 data URL ("data:image/jpeg;base64,...") is sent as-is instead.
	CoverFile              *string  `json:"cover_file,omitempty"`
	MadeForKids            *bool    `json:"made_for_kids,omitempty"`
	Tags                   []string `json:"tags,omitempty"`
	CategoryID             *string  `json:"category_id,omitempty"`
	ContainsSyntheticMedia *bool    `json:"contains_synthetic_media,omitempty"`
	// PlaylistID adds the video to this playlist, from
	// ProfilesService.YouTubePlaylists.
	PlaylistID *string `json:"playlist_id,omitempty"`
}

// PinterestParams contains Pinterest-specific post parameters.
type PinterestParams struct {
	Format          *PinterestFormat `json:"format,omitempty"`
	Title           *string          `json:"title,omitempty"`
	BoardID         *string          `json:"board_id,omitempty"`
	DestinationLink *string          `json:"destination_link,omitempty"`
	CoverURL        *string          `json:"cover_url,omitempty"`
	// CoverFile is a local path to a cover image, uploaded with the post. A
	// base64 data URL ("data:image/jpeg;base64,...") is sent as-is instead.
	CoverFile   *string `json:"cover_file,omitempty"`
	ThumbOffset *int    `json:"thumb_offset,omitempty"`
	AltText     AltText `json:"alt_text,omitempty"`
}

// ThreadsParams contains Threads-specific post parameters.
type ThreadsParams struct {
	Format  *ThreadsFormat `json:"format,omitempty"`
	AltText AltText        `json:"alt_text,omitempty"`
	// TopicTag is one tag, 50 characters or fewer, without the leading # — no
	// "." or "&".
	TopicTag *string `json:"topic_tag,omitempty"`
}

// TwitterParams contains Twitter-specific post parameters.
type TwitterParams struct {
	Format  *TwitterFormat `json:"format,omitempty"`
	AltText AltText        `json:"alt_text,omitempty"`
	// PollOptions is required when Format is "poll": 2-4 options, max 25 characters each.
	PollOptions []string `json:"poll_options,omitempty"`
	// PollDurationMinutes is required when Format is "poll": 5 to 10080 minutes (7 days).
	PollDurationMinutes *int `json:"poll_duration_minutes,omitempty"`
}

// BlueskyParams contains Bluesky-specific post parameters.
type BlueskyParams struct {
	Format  *BlueskyFormat `json:"format,omitempty"`
	AltText AltText        `json:"alt_text,omitempty"`
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
	ExternalEditedAt *string        `json:"external_edited_at"`
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
	ExternalEditedAt *string        `json:"external_edited_at"`
	PostedAt         *string        `json:"posted_at"`
	CreatedAt        string         `json:"created_at"`
}

// PostSync records one post pull for a profile — the sync fired when the
// profile connects, the recurring poll, or a backfill.
type PostSync struct {
	ID          string          `json:"id"`
	ProfileID   string          `json:"profile_id"`
	PlacementID *string         `json:"placement_id"`
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

// ChatSync records one DM import for a profile — started by
// ProfilesService.BackfillChats or on connect.
type ChatSync struct {
	ID               string          `json:"id"`
	ProfileID        string          `json:"profile_id"`
	PlacementID      *string         `json:"placement_id"`
	Kind             string          `json:"kind"`
	Trigger          PostSyncTrigger `json:"trigger"`
	Status           PostSyncStatus  `json:"status"`
	StartedAt        *string         `json:"started_at"`
	CompletedAt      *string         `json:"completed_at"`
	ChatsSeen        int             `json:"chats_seen"`
	MessagesImported int             `json:"messages_imported"`
	BackfillFrom     *string         `json:"backfill_from"`
	OldestUpdatedAt  *string         `json:"oldest_updated_at"`
	Error            *string         `json:"error"`
	CreatedAt        string          `json:"created_at"`
}

// ChatSyncListOptions contains options for ProfilesService.ChatSyncs.
type ChatSyncListOptions struct {
	Trigger        *PostSyncTrigger
	Status         *PostSyncStatus
	Page           *int
	PerPage        *int
	ProfileGroupID *string
}

// DailyStatsRecord is one day of a DailyStats series. Normalized metrics
// (e.g. "impressions", "viewers") are top-level keys of the record; they are
// collected in Metrics.
type DailyStatsRecord struct {
	Date          string         `json:"date"`
	PlatformStats map[string]any `json:"platform_stats"`
	Metrics       map[string]any `json:"-"`
}

// UnmarshalJSON collects every key other than date and platform_stats into Metrics.
func (r *DailyStatsRecord) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if date, ok := raw["date"].(string); ok {
		r.Date = date
	}
	if ps, ok := raw["platform_stats"].(map[string]any); ok {
		r.PlatformStats = ps
	}
	delete(raw, "date")
	delete(raw, "platform_stats")
	r.Metrics = raw
	return nil
}

// DailyStats is a profile's per-day series returned by ProfilesService.DailyStats.
type DailyStats struct {
	ProfileID    string                `json:"profile_id"`
	Platform     Platform              `json:"platform"`
	PlacementID  *string               `json:"placement_id"`
	Attribution  DailyStatsAttribution `json:"attribution"`
	From         string                `json:"from"`
	To           string                `json:"to"`
	DataTimezone *string               `json:"data_timezone"`
	Records      []DailyStatsRecord    `json:"records"`
	// Metrics describes each metric in Records.
	Metrics map[string]any `json:"metrics"`
}

// DailyStatsResponse wraps the data field returned by GET /profiles/:id/daily_stats.
type DailyStatsResponse struct {
	Data DailyStats `json:"data"`
}

// DailyStatsOptions are the optional query parameters for ProfilesService.DailyStats.
type DailyStatsOptions struct {
	PlacementID *string
	From        *string
	To          *string
	// Attribution DailyStatsAttributionPublish counts post metrics toward the
	// day the post was published instead of the day they were received.
	Attribution    *DailyStatsAttribution
	ProfileGroupID *string
}

// FollowerHistoryRecord is one day of follower changes.
type FollowerHistoryRecord struct {
	Date            string `json:"date"`
	FollowersGained *int   `json:"followers_gained"`
	FollowersLost   *int   `json:"followers_lost"`
	FollowersChange *int   `json:"followers_change"`
	FollowersTotal  *int   `json:"followers_total"`
}

// FollowerHistory is returned by ProfilesService.FollowerHistory.
type FollowerHistory struct {
	ProfileID   string                  `json:"profile_id"`
	Platform    Platform                `json:"platform"`
	PlacementID *string                 `json:"placement_id"`
	From        string                  `json:"from"`
	To          string                  `json:"to"`
	Records     []FollowerHistoryRecord `json:"records"`
}

// FollowerHistoryResponse wraps the data field returned by GET /profiles/:id/follower_history.
type FollowerHistoryResponse struct {
	Data FollowerHistory `json:"data"`
}

// FollowerHistoryOptions are the optional query parameters for ProfilesService.FollowerHistory.
type FollowerHistoryOptions struct {
	PlacementID    *string
	From           *string
	To             *string
	ProfileGroupID *string
}

// ParticipantStatus says whether an Instagram DM participant is verified and
// follows you. Cached for a short while server-side.
type ParticipantStatus struct {
	IsVerifiedUser       *bool  `json:"is_verified_user,omitempty"`
	IsUserFollowBusiness *bool  `json:"is_user_follow_business,omitempty"`
	IsBusinessFollowUser *bool  `json:"is_business_follow_user,omitempty"`
	FetchedAt            string `json:"fetched_at"`
}

// ResolvedMention turns a LinkedIn company page URL (or its vanity name) into
// a URN you can @-mention in a post body.
type ResolvedMention struct {
	URN        string `json:"urn"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	VanityName string `json:"vanity_name"`
}

// InstagramAvailableAudioOptions are the optional query parameters for
// ProfilesService.InstagramAvailableAudio.
type InstagramAvailableAudioOptions struct {
	AudioType      *InstagramAudioType
	SearchQuery    *string
	ProfileGroupID *string
}

// InstagramAvailableAudioResponse lists music and original sounds a reel can
// attach via InstagramParams.AudioID.
type InstagramAvailableAudioResponse struct {
	Data []map[string]any `json:"data"`
}

// InstagramCollaborationInvitesOptions are the optional parameters for
// ProfilesService.InstagramCollaborationInvites.
type InstagramCollaborationInvitesOptions struct {
	// After is the NextCursor of the previous page.
	After          *string
	ProfileGroupID *string
}

// InstagramCollaborationInvite is a post from another account that invited
// this Instagram profile as a collaborator.
type InstagramCollaborationInvite struct {
	MediaID            string  `json:"media_id"`
	MediaOwnerUsername *string `json:"media_owner_username"`
	Caption            *string `json:"caption"`
	// MediaURL is left out by Instagram for videos with copyrighted or
	// licensed audio.
	MediaURL *string `json:"media_url,omitempty"`
}

// InstagramCollaborationInvitesResponse is a page of collaboration invites.
type InstagramCollaborationInvitesResponse struct {
	CollaborationInvites []InstagramCollaborationInvite `json:"collaboration_invites"`
	// NextCursor is nil on the last page.
	NextCursor *string `json:"next_cursor"`
}

// InstagramCollaborationInviteResult is returned by
// ProfilesService.AcceptInstagramCollaborationInvite (Accepted) and
// DeclineInstagramCollaborationInvite (Declined).
type InstagramCollaborationInviteResult struct {
	MediaID  string `json:"media_id"`
	Accepted *bool  `json:"accepted,omitempty"`
	Declined *bool  `json:"declined,omitempty"`
}

// PinterestBoardOptions are the optional parameters for
// ProfilesService.CreatePinterestBoard.
type PinterestBoardOptions struct {
	Description    *string
	Privacy        *PinterestBoardPrivacy
	ProfileGroupID *string
}

// PinterestBoardDeleted is returned by ProfilesService.DeletePinterestBoard.
type PinterestBoardDeleted struct {
	Deleted bool   `json:"deleted"`
	BoardID string `json:"board_id"`
}

// YouTubePlaylist is a playlist on a YouTube channel.
type YouTubePlaylist struct {
	ID            string         `json:"id"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	PrivacyStatus YouTubePrivacy `json:"privacy_status"`
	ItemCount     int            `json:"item_count"`
	ThumbnailURL  *string        `json:"thumbnail_url"`
	PublishedAt   string         `json:"published_at"`
}

// YouTubePlaylistsOptions are the optional parameters for
// ProfilesService.YouTubePlaylists.
type YouTubePlaylistsOptions struct {
	// PageToken is the NextPageToken of the previous page.
	PageToken      *string
	ProfileGroupID *string
}

// YouTubePlaylistsResponse is a page of up to 50 playlists.
type YouTubePlaylistsResponse struct {
	Data []YouTubePlaylist `json:"data"`
	// NextPageToken is nil on the last page.
	NextPageToken *string `json:"next_page_token"`
}

// YouTubePlaylistCreateOptions are the optional parameters for
// ProfilesService.CreateYouTubePlaylist.
type YouTubePlaylistCreateOptions struct {
	Description    *string
	PrivacyStatus  *YouTubePrivacy
	ProfileGroupID *string
}

// YouTubePlaylistUpdateOptions are the fields to change in
// ProfilesService.UpdateYouTubePlaylist. Nil fields keep their current value;
// a pointer to "" for Description clears it.
type YouTubePlaylistUpdateOptions struct {
	Title          *string
	Description    *string
	PrivacyStatus  *YouTubePrivacy
	ProfileGroupID *string
}

// YouTubePlaylistDeleted is returned by ProfilesService.DeleteYouTubePlaylist.
type YouTubePlaylistDeleted struct {
	Deleted    bool   `json:"deleted"`
	PlaylistID string `json:"playlist_id"`
}

// ThreadControlOptions are the optional parameters for the ChatsService
// thread control calls.
type ThreadControlOptions struct {
	// Metadata is a free-form string passed along to the other app.
	Metadata       *string
	ProfileGroupID *string
}

// ThreadOwner is who controls a Facebook / Instagram conversation. Both
// fields are nil when the conversation is idle.
type ThreadOwner struct {
	AppID     *string `json:"app_id"`
	ExpiresAt *string `json:"expires_at"`
}

// FacebookPageOptions are the optional parameters for ProfilesService.FacebookPage.
type FacebookPageOptions struct {
	// Fields is a comma-separated list of Page fields to read.
	Fields         *string
	ProfileGroupID *string
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
	ID        string   `json:"id"`
	ProfileID string   `json:"profile_id"`
	Platform  Platform `json:"platform"`
	// ExternalPlacementID is the page / phone number the chat runs through.
	ExternalPlacementID *string `json:"external_placement_id"`
	// Group is true for a WhatsApp group chat.
	Group                  bool    `json:"group"`
	ParticipantExternalID  string  `json:"participant_external_id"`
	ParticipantUsername    *string `json:"participant_username"`
	ParticipantName        *string `json:"participant_name"`
	ParticipantAvatarURL   *string `json:"participant_avatar_url"`
	ExternalConversationID *string `json:"external_conversation_id"`
	LastInboundAt          *string `json:"last_inbound_at"`
	LastOutboundAt         *string `json:"last_outbound_at"`
	LastMessageAt          *string `json:"last_message_at"`
	// Metadata carries "thread_control" on Facebook and Instagram chats once
	// known ("owner_app_id", "expires_at", "updated_at").
	Metadata map[string]any `json:"metadata"`
	Archived *bool          `json:"archived,omitempty"`
	// WithinMessagingWindow is WhatsApp-only: whether the 24-hour customer
	// service window is open. When it's closed only a template can be sent.
	WithinMessagingWindow *bool  `json:"within_messaging_window,omitempty"`
	CreatedAt             string `json:"created_at"`
}

// QuickReply is a tappable chip rendered above the participant's composer,
// gone once tapped. Facebook and Instagram only; up to 13 per send.
// ContentType is optional on send (only "text" is accepted) and always present
// on responses.
type QuickReply struct {
	ContentType string `json:"content_type,omitempty"`
	// 20 characters or fewer.
	Title string `json:"title"`
	// 1000 characters or fewer. Comes back on the tap as TappedAction.Payload.
	Payload string `json:"payload"`
}

// MessageButton is a button attached to the message itself, delivered as a Meta
// generic template. Facebook and Instagram only; up to 3 per send.
//
// URL is required and must be https when Type is "web_url"; Payload is required
// when Type is "postback". Type is a plain string rather than a named enum so a
// new Meta button type needs no SDK release.
type MessageButton struct {
	Type string `json:"type"`
	// 20 characters or fewer.
	Title   string `json:"title"`
	URL     string `json:"url,omitempty"`
	Payload string `json:"payload,omitempty"`
}

// CardDefaultAction is the tap-through on a MessageCard.
type CardDefaultAction struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// MessageCard fills in the rest of the generic-template element that carries
// Buttons, for a richer product-style card. Requires Buttons.
type MessageCard struct {
	// 80 characters or fewer.
	Subtitle      string             `json:"subtitle,omitempty"`
	ImageURL      string             `json:"image_url,omitempty"`
	DefaultAction *CardDefaultAction `json:"default_action,omitempty"`
}

// MessageElement is one card of a horizontally scrolling carousel. Facebook
// and Instagram only; up to 10 per send. Title plus at least one other field
// is required.
type MessageElement struct {
	// 80 characters or fewer.
	Title         string             `json:"title"`
	Subtitle      string             `json:"subtitle,omitempty"`
	ImageURL      string             `json:"image_url,omitempty"`
	DefaultAction *CardDefaultAction `json:"default_action,omitempty"`
	// Up to 3.
	Buttons []MessageButton `json:"buttons,omitempty"`
}

// TemplateButtonParam fills in a dynamic template button. SubType is "url",
// "copy_code", "flow", "quick_reply" or "voice_call"; the remaining
// sub-type-specific fields (e.g. "text", "coupon_code") go in Extra.
type TemplateButtonParam struct {
	Index   int
	SubType string
	Extra   map[string]any
}

// MarshalJSON flattens Extra next to index and sub_type.
func (p TemplateButtonParam) MarshalJSON() ([]byte, error) {
	out := make(map[string]any, len(p.Extra)+2)
	for k, v := range p.Extra {
		out[k] = v
	}
	out["index"] = p.Index
	out["sub_type"] = p.SubType
	return json.Marshal(out)
}

// TemplateHeaderMedia is an image, video or document template header, given
// by Link or by uploaded media ID.
type TemplateHeaderMedia struct {
	Link string `json:"link,omitempty"`
	ID   string `json:"id,omitempty"`
}

// MessageTemplateInput sends a WhatsApp template. Address the template by ID
// (PostProxy or Meta id) or by Name (plus Language when it exists in several).
type MessageTemplateInput struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Language string `json:"language,omitempty"`
	// Variables are positional values: header, then body, then dynamic URL buttons.
	Variables      []string              `json:"variables,omitempty"`
	ButtonParams   []TemplateButtonParam `json:"button_params,omitempty"`
	HeaderMedia    *TemplateHeaderMedia  `json:"header_media,omitempty"`
	HeaderLocation *MessageLocation      `json:"header_location,omitempty"`
}

// MessageLocation is a WhatsApp location pin.
type MessageLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

// MessageTemplateRef identifies the template an outbound WhatsApp message was sent from.
type MessageTemplateRef struct {
	ID         string   `json:"id"`
	ExternalID *string  `json:"external_id"`
	Name       string   `json:"name"`
	Language   string   `json:"language"`
	Variables  []string `json:"variables"`
}

// TappedActionKind identifies which interactive element a participant tapped.
type TappedActionKind = string

const (
	TappedActionQuickReply    TappedActionKind = "quick_reply"
	TappedActionPostback      TappedActionKind = "postback"
	TappedActionCallbackQuery TappedActionKind = "callback_query"
)

// TappedAction is set on inbound messages created by a tap on an interactive
// element you sent. It is derived from PlatformData rather than stored, so it
// also resolves for taps ingested before PostProxy exposed this field.
// TappedActionCallbackQuery is Telegram, so this is not Meta-only even though
// the send options are.
type TappedAction struct {
	Kind TappedActionKind `json:"kind"`
	// The payload you set on the quick reply, button, or ice breaker.
	Payload string `json:"payload"`
	// The label the participant tapped.
	Title *string `json:"title"`
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
	ErrorDetails        *ErrorDetails    `json:"error_details,omitempty"`
	PlatformData        map[string]any   `json:"platform_data"`
	ExternalPostedAt    *string          `json:"external_posted_at"`
	ExternalDeliveredAt *string          `json:"external_delivered_at"`
	ExternalReadAt      *string          `json:"external_read_at"`
	ExternalEditedAt    *string          `json:"external_edited_at"`
	ReplyToExternalID   *string          `json:"reply_to_external_id"`
	ReplyMarkup         map[string]any   `json:"reply_markup"`
	QuickReplies        []QuickReply     `json:"quick_replies"`
	Buttons             []MessageButton  `json:"buttons"`
	Card                *MessageCard     `json:"card"`
	Elements            []MessageElement `json:"elements"`
	TappedAction        *TappedAction    `json:"tapped_action"`
	StoryReply          map[string]any   `json:"story_reply"`
	// Interactive, Location, Contacts and Template are WhatsApp-only.
	Interactive       map[string]any      `json:"interactive"`
	Location          *MessageLocation    `json:"location"`
	Contacts          []map[string]any    `json:"contacts"`
	Template          *MessageTemplateRef `json:"template,omitempty"`
	ExternalDeletedAt *string             `json:"external_deleted_at"`
	Reactions         []Reaction          `json:"reactions"`
	Attachments       []Attachment        `json:"attachments"`
	IsUnsupported     bool                `json:"is_unsupported"`
	CreatedAt         string              `json:"created_at"`
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
	// PlacementID is required for WhatsApp: the phone number id from
	// ProfilesService.Placements.
	PlacementID    *string
	ProfileGroupID *string
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
//
// QuickReplies, Buttons, and Card are Facebook and Instagram only — they return
// 422 on Telegram and Bluesky, where ReplyMarkup is the equivalent. They are
// sent on the JSON path only, so pass Media as hosted URLs rather than
// MediaFiles when combining with an attachment.
type MessageSendOptions struct {
	Body       *string
	Media      []string
	MediaFiles []string
	Tag        *string
	// ReplyToExternalID quotes the message with this external_id. Telegram,
	// WhatsApp and Instagram.
	ReplyToExternalID *string
	ReplyMarkup       map[string]any
	QuickReplies      []QuickReply
	Buttons           []MessageButton
	// Requires Buttons.
	Card *MessageCard
	// Elements is a carousel of up to 10 cards (Facebook and Instagram). Can't
	// be combined with Body, Buttons or Card.
	Elements []MessageElement
	// WhatsApp only. Send exactly one of Body, Media, Template, Interactive,
	// Location or Contacts. Outside the 24-hour window only a Template gets
	// through.
	Template    *MessageTemplateInput
	Interactive map[string]any
	Location    *MessageLocation
	Contacts    []map[string]any
	// Category "utility" sends a text message under the utility category.
	Category    *string
	LinkPreview *bool
	// VoiceNote delivers an audio attachment as a voice note.
	VoiceNote *bool
	// Filename is the file name shown for a document attachment.
	Filename       *string
	ProfileGroupID *string
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
	PostIDs  []string
	Profiles []string
	// PlacementID narrows to one page / organization.
	PlacementID    *string
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

// PrivateReplyOptions contains the interactive extras for
// CommentsService.PrivateReply. They work as in MessageSendOptions; the text
// may be empty when sending Elements.
type PrivateReplyOptions struct {
	QuickReplies []QuickReply
	Buttons      []MessageButton
	Card         *MessageCard
	Elements     []MessageElement
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
	// SkipDuplicateCheck publishes even if an identical post was just
	// submitted, instead of failing with a 409.
	SkipDuplicateCheck bool
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
//
// On a published post only Body, Platforms.YouTube.Title,
// Platforms.Pinterest.Title and Platforms.Pinterest.DestinationLink can
// change; the edit is pushed to the platforms in the background (the
// platform_post.edited / platform_post.edit_failed webhooks).
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

// WhatsAppTemplate is a WhatsApp message template, which lives on the
// WhatsApp Business account (the profile).
type WhatsAppTemplate struct {
	ID         string  `json:"id"`
	ExternalID *string `json:"external_id"`
	ProfileID  string  `json:"profile_id"`
	Name       string  `json:"name"`
	Language   string  `json:"language"`
	Category   string  `json:"category"`
	Status     string  `json:"status"`
	// ParameterFormat is "POSITIONAL" or "NAMED".
	ParameterFormat       *string          `json:"parameter_format"`
	Components            []map[string]any `json:"components"`
	MessageSendTTLSeconds *int             `json:"message_send_ttl_seconds,omitempty"`
	RejectedReason        *string          `json:"rejected_reason,omitempty"`
	QualityScore          *string          `json:"quality_score,omitempty"`
	VariableCount         int              `json:"variable_count"`
	CreatedAt             string           `json:"created_at"`
	UpdatedAt             string           `json:"updated_at"`
}

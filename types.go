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

// Media represents a media attachment on a post.
type Media struct {
	ID           string      `json:"id"`
	Status       MediaStatus `json:"status"`
	ErrorMessage *string     `json:"error_message"`
	ContentType  string      `json:"content_type"`
	SourceURL    *string     `json:"source_url"`
	URL          *string     `json:"url"`
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
	ID             string           `json:"id"`
	Body           string           `json:"body"`
	Status         PostStatus       `json:"status"`
	ScheduledAt    *string          `json:"scheduled_at"`
	CreatedAt      string           `json:"created_at"`
	Media          []Media          `json:"media"`
	Platforms      []PlatformResult `json:"platforms"`
	Thread         []ThreadChild    `json:"thread"`
	QueueID        *string          `json:"queue_id"`
	QueuePriority  *string          `json:"queue_priority"`
}

// Placement represents a placement option for a profile.
type Placement struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

// ConnectionResponse represents the response from an initialize connection operation.
type ConnectionResponse struct {
	URL     string `json:"url"`
	Success bool   `json:"success"`
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
	Format        *InstagramFormat `json:"format,omitempty"`
	FirstComment  *string          `json:"first_comment,omitempty"`
	Collaborators []string         `json:"collaborators,omitempty"`
	CoverURL      *string          `json:"cover_url,omitempty"`
	AudioName     *string          `json:"audio_name,omitempty"`
	TrialStrategy *bool            `json:"trial_strategy,omitempty"`
	ThumbOffset   *int             `json:"thumb_offset,omitempty"`
}

// TikTokParams contains TikTok-specific post parameters.
type TikTokParams struct {
	Format              *TikTokFormat  `json:"format,omitempty"`
	PrivacyStatus       *TikTokPrivacy `json:"privacy_status,omitempty"`
	PhotoCoverIndex     *int           `json:"photo_cover_index,omitempty"`
	AutoAddMusic        *bool          `json:"auto_add_music,omitempty"`
	MadeWithAI          *bool          `json:"made_with_ai,omitempty"`
	DisableComment      *bool          `json:"disable_comment,omitempty"`
	DisableDuet         *bool          `json:"disable_duet,omitempty"`
	DisableStitch       *bool          `json:"disable_stitch,omitempty"`
	BrandContentToggle  *bool          `json:"brand_content_toggle,omitempty"`
	BrandOrganicToggle  *bool          `json:"brand_organic_toggle,omitempty"`
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
}

// PlatformParams contains platform-specific parameters for a post.
type PlatformParams struct {
	Facebook  *FacebookParams  `json:"facebook,omitempty"`
	Instagram *InstagramParams `json:"instagram,omitempty"`
	TikTok    *TikTokParams    `json:"tiktok,omitempty"`
	LinkedIn  *LinkedInParams  `json:"linkedin,omitempty"`
	YouTube   *YouTubeParams   `json:"youtube,omitempty"`
	Pinterest *PinterestParams `json:"pinterest,omitempty"`
	Threads   *ThreadsParams   `json:"threads,omitempty"`
	Twitter   *TwitterParams   `json:"twitter,omitempty"`
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
	ID             string  `json:"id"`
	EventID        string  `json:"event_id"`
	EventType      string  `json:"event_type"`
	ResponseStatus *int    `json:"response_status"`
	AttemptNumber  int     `json:"attempt_number"`
	Success        bool    `json:"success"`
	AttemptedAt    string  `json:"attempted_at"`
	CreatedAt      string  `json:"created_at"`
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
	Stats      map[string]any `json:"stats"`
	RecordedAt string         `json:"recorded_at"`
}

// Comment represents a comment on a post.
type Comment struct {
	ID               string    `json:"id"`
	ExternalID       *string   `json:"external_id"`
	Body             string    `json:"body"`
	Status           string    `json:"status"`
	AuthorUsername   *string   `json:"author_username"`
	AuthorAvatarURL  *string   `json:"author_avatar_url"`
	AuthorExternalID *string   `json:"author_external_id"`
	ParentExternalID *string   `json:"parent_external_id"`
	LikeCount        int       `json:"like_count"`
	IsHidden         bool      `json:"is_hidden"`
	Permalink        *string   `json:"permalink"`
	PlatformData     any       `json:"platform_data"`
	PostedAt         *string   `json:"posted_at"`
	CreatedAt        string    `json:"created_at"`
	Replies          []Comment `json:"replies"`
}

// CommentListOptions contains options for listing comments.
type CommentListOptions struct {
	Page    *int
	PerPage *int
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

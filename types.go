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

// PlatformResult represents the result of posting to a specific platform.
type PlatformResult struct {
	Platform    Platform           `json:"platform"`
	Status      PlatformPostStatus `json:"status"`
	Params      map[string]any     `json:"params"`
	Error       *string            `json:"error"`
	AttemptedAt *string            `json:"attempted_at"`
	Insights    *Insights          `json:"insights"`
}

// Post represents a social media post.
type Post struct {
	ID          string           `json:"id"`
	Body        string           `json:"body"`
	Status      PostStatus       `json:"status"`
	ScheduledAt *string          `json:"scheduled_at"`
	CreatedAt   string           `json:"created_at"`
	Platforms   []PlatformResult `json:"platforms"`
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
	Format        *YouTubeFormat  `json:"format,omitempty"`
	Title         *string         `json:"title,omitempty"`
	PrivacyStatus *YouTubePrivacy `json:"privacy_status,omitempty"`
	CoverURL      *string         `json:"cover_url,omitempty"`
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
	ScheduledAt    *string
	Draft          *bool
	ProfileGroupID *string
}

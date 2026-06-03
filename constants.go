package postproxy

// Version is the SDK version, sent as part of the User-Agent header.
const Version = "1.10.0"

// Platform represents a social media platform.
type Platform string

const (
	PlatformFacebook       Platform = "facebook"
	PlatformInstagram      Platform = "instagram"
	PlatformTikTok         Platform = "tiktok"
	PlatformLinkedIn       Platform = "linkedin"
	PlatformYouTube        Platform = "youtube"
	PlatformTwitter        Platform = "twitter"
	PlatformThreads        Platform = "threads"
	PlatformPinterest      Platform = "pinterest"
	PlatformBluesky        Platform = "bluesky"
	PlatformTelegram       Platform = "telegram"
	PlatformGoogleBusiness Platform = "google_business"
)

// PostStatus represents the status of a post.
type PostStatus string

const (
	PostStatusPending               PostStatus = "pending"
	PostStatusDraft                 PostStatus = "draft"
	PostStatusProcessing            PostStatus = "processing"
	PostStatusProcessed             PostStatus = "processed"
	PostStatusScheduled             PostStatus = "scheduled"
	PostStatusMediaProcessingFailed PostStatus = "media_processing_failed"
)

// MediaStatus represents the processing status of a media attachment.
type MediaStatus string

const (
	MediaStatusPending   MediaStatus = "pending"
	MediaStatusProcessed MediaStatus = "processed"
	MediaStatusFailed    MediaStatus = "failed"
)

// PlatformPostStatus represents the status of a post on a specific platform.
type PlatformPostStatus string

const (
	PlatformPostStatusPending    PlatformPostStatus = "pending"
	PlatformPostStatusProcessing PlatformPostStatus = "processing"
	PlatformPostStatusPublished  PlatformPostStatus = "published"
	PlatformPostStatusFailed     PlatformPostStatus = "failed"
	PlatformPostStatusDeleted    PlatformPostStatus = "deleted"
)

// ProfileStatus represents the status of a profile.
type ProfileStatus string

const (
	ProfileStatusActive   ProfileStatus = "active"
	ProfileStatusExpired  ProfileStatus = "expired"
	ProfileStatusInactive ProfileStatus = "inactive"
)

// InstagramFormat represents the format of an Instagram post.
type InstagramFormat string

const (
	InstagramFormatPost  InstagramFormat = "post"
	InstagramFormatReel  InstagramFormat = "reel"
	InstagramFormatStory InstagramFormat = "story"
)

// FacebookFormat represents the format of a Facebook post.
type FacebookFormat string

const (
	FacebookFormatPost  FacebookFormat = "post"
	FacebookFormatStory FacebookFormat = "story"
	FacebookFormatReel  FacebookFormat = "reel"
)

// TikTokFormat represents the format of a TikTok post.
type TikTokFormat string

const (
	TikTokFormatVideo TikTokFormat = "video"
	TikTokFormatImage TikTokFormat = "image"
)

// LinkedInFormat represents the format of a LinkedIn post.
type LinkedInFormat string

const (
	LinkedInFormatPost LinkedInFormat = "post"
)

// YouTubeFormat represents the format of a YouTube post.
type YouTubeFormat string

const (
	YouTubeFormatPost YouTubeFormat = "post"
)

// PinterestFormat represents the format of a Pinterest post.
type PinterestFormat string

const (
	PinterestFormatPin PinterestFormat = "pin"
)

// ThreadsFormat represents the format of a Threads post.
type ThreadsFormat string

const (
	ThreadsFormatPost ThreadsFormat = "post"
)

// TwitterFormat represents the format of a Twitter post.
type TwitterFormat string

const (
	TwitterFormatPost TwitterFormat = "post"
)

// TikTokPrivacy represents the privacy setting for a TikTok post.
type TikTokPrivacy string

const (
	TikTokPrivacyPublic        TikTokPrivacy = "PUBLIC_TO_EVERYONE"
	TikTokPrivacyMutualFollows TikTokPrivacy = "MUTUAL_FOLLOW_FRIENDS"
	TikTokPrivacyFollowers     TikTokPrivacy = "FOLLOWER_OF_CREATOR"
	TikTokPrivacySelfOnly      TikTokPrivacy = "SELF_ONLY"
)

// YouTubePrivacy represents the privacy setting for a YouTube post.
type YouTubePrivacy string

const (
	YouTubePrivacyPublic   YouTubePrivacy = "public"
	YouTubePrivacyUnlisted YouTubePrivacy = "unlisted"
	YouTubePrivacyPrivate  YouTubePrivacy = "private"
)

// BlueskyFormat represents the format of a Bluesky post.
type BlueskyFormat string

const (
	BlueskyFormatPost BlueskyFormat = "post"
)

// TelegramFormat represents the format of a Telegram post.
type TelegramFormat string

const (
	TelegramFormatPost TelegramFormat = "post"
)

// TelegramParseMode controls Telegram message body formatting.
type TelegramParseMode string

const (
	TelegramParseModeHTML       TelegramParseMode = "HTML"
	TelegramParseModeMarkdownV2 TelegramParseMode = "MarkdownV2"
)

// MessageDirection represents the direction of a direct message.
type MessageDirection string

const (
	MessageDirectionInbound  MessageDirection = "inbound"
	MessageDirectionOutbound MessageDirection = "outbound"
)

// MessageStatus represents the delivery status of a direct message.
type MessageStatus string

const (
	MessageStatusPending               MessageStatus = "pending"
	MessageStatusPublished             MessageStatus = "published"
	MessageStatusFailedWaitingForRetry MessageStatus = "failed_waiting_for_retry"
	MessageStatusFailed                MessageStatus = "failed"
	MessageStatusReceived              MessageStatus = "received"
)

// WebhookEventType is the discriminator for webhook events.
type WebhookEventType string

const (
	EventPostProcessed                  WebhookEventType = "post.processed"
	EventPostImported                   WebhookEventType = "post.imported"
	EventPlatformPostPublished          WebhookEventType = "platform_post.published"
	EventPlatformPostFailed             WebhookEventType = "platform_post.failed"
	EventPlatformPostFailedWaitingRetry WebhookEventType = "platform_post.failed_waiting_for_retry"
	EventPlatformPostInsights           WebhookEventType = "platform_post.insights"
	EventProfileConnected               WebhookEventType = "profile.connected"
	EventProfileDisconnected            WebhookEventType = "profile.disconnected"
	EventProfileStats                   WebhookEventType = "profile.stats"
	EventMediaFailed                    WebhookEventType = "media.failed"
	EventCommentCreated                 WebhookEventType = "comment.created"
	EventProfileCommentCreated          WebhookEventType = "profile_comment.created"
	EventMessageReceived                WebhookEventType = "message.received"
	EventMessageSent                    WebhookEventType = "message.sent"
	EventMessageDelivered               WebhookEventType = "message.delivered"
	EventMessageRead                    WebhookEventType = "message.read"
	EventMessageEdited                  WebhookEventType = "message.edited"
	EventMessageDeleted                 WebhookEventType = "message.deleted"
	EventMessageFailedWaitingForRetry   WebhookEventType = "message.failed_waiting_for_retry"
	EventMessageFailed                  WebhookEventType = "message.failed"
	EventReactionReceived               WebhookEventType = "reaction.received"
)

// WebhookEventTypes is the set of all known webhook event types.
var WebhookEventTypes = []WebhookEventType{
	EventPostProcessed,
	EventPostImported,
	EventPlatformPostPublished,
	EventPlatformPostFailed,
	EventPlatformPostFailedWaitingRetry,
	EventPlatformPostInsights,
	EventProfileConnected,
	EventProfileDisconnected,
	EventProfileStats,
	EventMediaFailed,
	EventCommentCreated,
	EventProfileCommentCreated,
	EventMessageReceived,
	EventMessageSent,
	EventMessageDelivered,
	EventMessageRead,
	EventMessageEdited,
	EventMessageDeleted,
	EventMessageFailedWaitingForRetry,
	EventMessageFailed,
	EventReactionReceived,
}

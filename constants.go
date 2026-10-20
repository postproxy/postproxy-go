package postproxy

// Version is the SDK version, sent as part of the User-Agent header.
const Version = "1.14.0"

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
	PlatformWhatsApp       Platform = "whatsapp"
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
	LinkedInFormatPoll LinkedInFormat = "poll"
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
	TwitterFormatPoll TwitterFormat = "poll"
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

// PostSyncTrigger says what started a post sync run.
type PostSyncTrigger string

const (
	PostSyncTriggerConnect   PostSyncTrigger = "connect"
	PostSyncTriggerScheduled PostSyncTrigger = "scheduled"
	PostSyncTriggerBackfill  PostSyncTrigger = "backfill"
)

// PostSyncStatus represents the state of a post sync run.
type PostSyncStatus string

const (
	PostSyncStatusPending   PostSyncStatus = "pending"
	PostSyncStatusRunning   PostSyncStatus = "running"
	PostSyncStatusCompleted PostSyncStatus = "completed"
	PostSyncStatusFailed    PostSyncStatus = "failed"
)

// InstagramTrialStrategy controls how an Instagram trial reel graduates to
// your followers.
type InstagramTrialStrategy string

const (
	InstagramTrialStrategyManual        InstagramTrialStrategy = "MANUAL"
	InstagramTrialStrategySSPerformance InstagramTrialStrategy = "SS_PERFORMANCE"
)

// InstagramAudioType filters ProfilesService.InstagramAvailableAudio.
type InstagramAudioType string

const (
	InstagramAudioTypeMusic         InstagramAudioType = "music"
	InstagramAudioTypeOriginalSound InstagramAudioType = "original_sound"
)

// ConnectionOnboarding picks the connection flow. Facebook and Instagram
// profiles connect through Facebook login; WhatsApp through the Cloud API or
// the WhatsApp Business app.
type ConnectionOnboarding string

const (
	ConnectionOnboardingInstagram   ConnectionOnboarding = "instagram"
	ConnectionOnboardingFacebook    ConnectionOnboarding = "facebook"
	ConnectionOnboardingAPI         ConnectionOnboarding = "api"
	ConnectionOnboardingBusinessApp ConnectionOnboarding = "business_app"
)

// DailyStatsAttribution says which day post metrics count toward in
// ProfilesService.DailyStats.
type DailyStatsAttribution string

const (
	DailyStatsAttributionReceived DailyStatsAttribution = "received"
	DailyStatsAttributionPublish  DailyStatsAttribution = "publish"
)

// ProfileInactiveReason says why an inactive profile was deactivated.
type ProfileInactiveReason string

const (
	// ProfileInactiveReasonNetworkTaken means the same account was connected elsewhere.
	ProfileInactiveReasonNetworkTaken ProfileInactiveReason = "network_taken"
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
	EventPostProcessed                       WebhookEventType = "post.processed"
	EventPostImported                        WebhookEventType = "post.imported"
	EventPlatformPostPublished               WebhookEventType = "platform_post.published"
	EventPlatformPostFailed                  WebhookEventType = "platform_post.failed"
	EventPlatformPostFailedWaitingRetry      WebhookEventType = "platform_post.failed_waiting_for_retry"
	EventPlatformPostEdited                  WebhookEventType = "platform_post.edited"
	EventPlatformPostEditFailed              WebhookEventType = "platform_post.edit_failed"
	EventPlatformPostInsights                WebhookEventType = "platform_post.insights"
	EventProfileConnected                    WebhookEventType = "profile.connected"
	EventProfileDisconnected                 WebhookEventType = "profile.disconnected"
	EventProfileStats                        WebhookEventType = "profile.stats"
	EventProfilePlacementsSynced             WebhookEventType = "profile.placements_synced"
	EventPlacementOrphaned                   WebhookEventType = "placement.orphaned"
	EventPlacementRestored                   WebhookEventType = "placement.restored"
	EventMediaFailed                         WebhookEventType = "media.failed"
	EventCommentCreated                      WebhookEventType = "comment.created"
	EventCommentFailed                       WebhookEventType = "comment.failed"
	EventCommentFailedWaitingForRetry        WebhookEventType = "comment.failed_waiting_for_retry"
	EventCommentDeliveryUnconfirmed          WebhookEventType = "comment.delivery_unconfirmed"
	EventCommentEdited                       WebhookEventType = "comment.edited"
	EventCommentEditFailed                   WebhookEventType = "comment.edit_failed"
	EventProfileCommentCreated               WebhookEventType = "profile_comment.created"
	EventProfileCommentFailed                WebhookEventType = "profile_comment.failed"
	EventProfileCommentFailedWaitingForRetry WebhookEventType = "profile_comment.failed_waiting_for_retry"
	EventProfileCommentDeleted               WebhookEventType = "profile_comment.deleted"
	EventMessageReceived                     WebhookEventType = "message.received"
	EventMessageSent                         WebhookEventType = "message.sent"
	EventMessageDelivered                    WebhookEventType = "message.delivered"
	EventMessageRead                         WebhookEventType = "message.read"
	EventMessageEdited                       WebhookEventType = "message.edited"
	EventMessageDeleted                      WebhookEventType = "message.deleted"
	EventMessageFailedWaitingForRetry        WebhookEventType = "message.failed_waiting_for_retry"
	EventMessageFailed                       WebhookEventType = "message.failed"
	EventReactionReceived                    WebhookEventType = "reaction.received"
	EventReferralReceived                    WebhookEventType = "referral.received"
	EventMessageTemplateStatusUpdated        WebhookEventType = "message_template.status_updated"
	EventMessageTemplateCategoryUpdated      WebhookEventType = "message_template.category_updated"
	EventPlacementNameStatusUpdated          WebhookEventType = "placement.name_status_updated"
	EventPlacementQualityUpdated             WebhookEventType = "placement.quality_updated"
	EventWhatsAppAutomaticEvent              WebhookEventType = "whatsapp.automatic_event"
)

// WebhookEventTypes is the set of all known webhook event types.
var WebhookEventTypes = []WebhookEventType{
	EventPostProcessed,
	EventPostImported,
	EventPlatformPostPublished,
	EventPlatformPostFailed,
	EventPlatformPostFailedWaitingRetry,
	EventPlatformPostEdited,
	EventPlatformPostEditFailed,
	EventPlatformPostInsights,
	EventProfileConnected,
	EventProfileDisconnected,
	EventProfileStats,
	EventProfilePlacementsSynced,
	EventPlacementOrphaned,
	EventPlacementRestored,
	EventMediaFailed,
	EventCommentCreated,
	EventCommentFailed,
	EventCommentFailedWaitingForRetry,
	EventCommentDeliveryUnconfirmed,
	EventCommentEdited,
	EventCommentEditFailed,
	EventProfileCommentCreated,
	EventProfileCommentFailed,
	EventProfileCommentFailedWaitingForRetry,
	EventProfileCommentDeleted,
	EventMessageReceived,
	EventMessageSent,
	EventMessageDelivered,
	EventMessageRead,
	EventMessageEdited,
	EventMessageDeleted,
	EventMessageFailedWaitingForRetry,
	EventMessageFailed,
	EventReactionReceived,
	EventReferralReceived,
	EventMessageTemplateStatusUpdated,
	EventMessageTemplateCategoryUpdated,
	EventPlacementNameStatusUpdated,
	EventPlacementQualityUpdated,
	EventWhatsAppAutomaticEvent,
}

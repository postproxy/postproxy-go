# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [1.14.0] - 2026-10-20

### Added

- **WhatsApp.** `PlatformWhatsApp` joins `Platform`. DMs: `ChatCreateOptions.PlacementID` (the phone number id, required for WhatsApp), and `MessageSendOptions` gains `Template` (`MessageTemplateInput`: `ID` or `Name` + `Language`, `Variables`, `ButtonParams`, `HeaderMedia`, `HeaderLocation`), `Interactive`, `Location` (`MessageLocation`), `Contacts`, plus `Category`, `LinkPreview`, `VoiceNote` and `Filename`. `Chat` gains `ExternalPlacementID`, `Group` and `WithinMessagingWindow`; `Message` gains `Template` (`MessageTemplateRef`), `Interactive`, `Location`, `Contacts` and `ErrorDetails`.
- **`client.WhatsAppTemplates`** — `List`, `Get`, `Create`, `Update`, `Delete`, `Library` for message templates (`WhatsAppTemplate`, `WhatsAppTemplateInput`, `WhatsAppTemplateUpdate`). Template ids and names are URL-escaped.
- **`client.WhatsApp`** — number administration: `NumberInfo`, `RegisterPhoneNumber`, `RequestVerificationCode`, `VerifyPhoneNumber`, `BusinessProfile`, `UpdateBusinessProfile`, `UpdateBusinessProfilePhoto`, `DisplayName`, `RequestDisplayNameChange`, `Username`, `SetUsername`, `DeleteUsername`, `UsernameSuggestions`, blocked users (`BlockedUsers`, `BlockedUserStatus`, `BlockUsers`, `UnblockUsers`), groups (`Groups`, `CreateGroup`, `Group`, `UpdateGroup`, `DeleteGroup`, `AddGroupParticipants`, `RemoveGroupParticipants`, `CreateGroupInviteLink`, `GroupJoinRequests`, `ApproveGroupJoinRequests`, `RejectGroupJoinRequests`), conversions (`Dataset`, `CreateDataset`, `SendConversionEvent`) and `Media`, which returns the raw bytes and content type (`WhatsAppMedia`).
- **`client.GoogleBusiness`** — manage a Google Business location's own data: `Location`, `UpdateLocation`, `AvailableLocationCategories`, `UpdateHours`, `Attributes`, `AvailableAttributes`, `UpdateAttributes`, `ServiceList`, `UpdateServiceList`, `FoodMenus`, `UpdateFoodMenus`, `PlaceActionLinks`, `AvailablePlaceActionTypes`, `CreatePlaceActionLink`, `UpdatePlaceActionLink`, `DeletePlaceActionLink`, `Media`, `CreateMedia`, `DeleteMedia`. Google's payloads pass through as `map[string]any`.
- **`client.Placements`** (`List`, `SetActiveProfile`, typed as `AccountPlacement`) and **`client.APIKeys`** (`List`, `Create`, `Delete`, typed as `APIKey`) — account-wide keys only.
- **Profiles:** `AssignToGroup`, `DailyStats` (`DailyStatsOptions.Attribution`, `DailyStatsResponse`), `FollowerHistory`, `ParticipantStatus` (Instagram), `ResolveMention` (LinkedIn), `InstagramAvailableAudio`, `FacebookPage` / `UpdateFacebookPage`, and chat backfill — `BackfillChats`, `ChatSyncs`, `ChatSync` (`ChatSync`, `ChatSyncListOptions`). `Profile` gains `Username`, `PlatformID`, `PlatformURL`, `AvatarURL`, `InactiveReason` and `PlacementsSyncStatus`; `Placement` gains `AvatarURL`; `PostSync` gains `PlacementID`.
- **Comments:** `Edit` (Facebook, YouTube), `Pin` / `Unpin` (TikTok), `PrivateReply` takes optional `PrivateReplyOptions` (`QuickReplies`, `Buttons`, `Card`, `Elements`), and `BulkCommentListOptions.PlacementID`. `Comment` and `BulkComment` gain `ExternalEditedAt`.
- **Chats:** `MarkRead`, `StartTyping`, `StopTyping`.
- **Carousels:** `MessageSendOptions.Elements` with the new `MessageElement` type (up to 10 cards, Facebook and Instagram). `Message.Elements` and `Message.StoryReply`.
- **Post params:** `AltText` (new `AltText` type — one entry per media item; a single entry is sent as a plain string) on Facebook, Instagram, LinkedIn, Bluesky, Threads, Pinterest and Twitter; `MadeWithAI` on Instagram and Facebook; Facebook `Link` and `Countries`; LinkedIn `LinkedInFormatPoll` with `PollQuestion` / `PollOptions` / `PollDurationDays`, plus `Countries` and `ReshareURN`; Threads `TopicTag`; TikTok `UploadAsDraft`, `CoverURL`, `ThumbOffset`; Instagram `AudioID`, `AudioVolume`, `VideoVolume`, `StripAudio`.
- **`CoverFile`** on Instagram, TikTok, YouTube and Pinterest params: a local path is uploaded with the post as `platforms[<platform>][cover_file]` (the request switches to multipart even without `MediaFiles`); a base64 `data:` URL is sent as-is. Works on `Posts.Create` and `Posts.Update`.
- `PostCreateOptions.SkipDuplicateCheck` sends `X-Skip-Duplicate-Check: true`.
- `PlatformResult` gains `ProfileID`, `ProfileName`, `Permalink` and `Attachments` (`PlatformAttachment`); `ThreadChild` gains `Platforms`, so you can see whether each item in a thread went out. `Post` gains `Draft` and `Source`.
- `ProfileGroups.InitializeConnection` takes optional `ConnectionOptions` — `Onboarding` (`ConnectionOnboardingInstagram` / `ConnectionOnboardingFacebook` for Instagram, `ConnectionOnboardingAPI` / `ConnectionOnboardingBusinessApp` for WhatsApp) and `RefreshConnection`.
- New constants: `InstagramTrialStrategy`, `InstagramAudioType`, `ConnectionOnboarding`, `DailyStatsAttribution`, `ProfileInactiveReason`.
- **Webhooks:** `platform_post.edited`, `platform_post.edit_failed`, `profile.placements_synced`, `placement.orphaned`, `placement.restored`, `comment.delivery_unconfirmed`, `comment.edited`, `comment.edit_failed`, `profile_comment.deleted`, `message_template.status_updated`, `message_template.category_updated`, `placement.name_status_updated`, `placement.quality_updated` and `whatsapp.automatic_event`, plus the previously missing `comment.failed`, `comment.failed_waiting_for_retry`, `profile_comment.failed`, `profile_comment.failed_waiting_for_retry` and `referral.received`. `ParseWebhookEvent` returned `ErrUnknownWebhookEvent` for these before. New decoders `AsPlacementsSynced`, `AsPlacement`, `AsReferral`, `AsMessageTemplate` and `AsWhatsAppAutomaticEvent`; `ProfileEventData` gains `PlatformID` and `Origin`, `CommentCreatedData` gains `ProfileID`, `Error`, `ErrorDetails` and `ExternalEditedAt`, `ProfileCommentCreatedData` gains `Error`, `ErrorDetails` and `Permalink`.
- `IsForbiddenError` (403 — account canceled or past due, or the platform refused) and `IsRateLimitError` (429), with `PostProxyError.RetryAfter` holding the `Retry-After` seconds.
- Examples: `examples/google-business`, `examples/whatsapp-templates`, `examples/whatsapp-groups`, `examples/dm-carousel`, `examples/edit-published-post`; a `CoverFile` post in `examples/create-post` and onboarding options in `examples/initialize-connection`.

### Changed

- `InstagramParams.TrialStrategy` is now `*InstagramTrialStrategy` (`"MANUAL"` / `"SS_PERFORMANCE"`); it was `*bool`, which the API never accepted.
- `Posts.Update` on a **published** post now edits it on the platforms: only `Body`, `Platforms.YouTube.Title`, `Platforms.Pinterest.Title` and `Platforms.Pinterest.DestinationLink` may change.
- `Comments.PrivateReply` and `ProfileGroups.InitializeConnection` take their new options as a trailing variadic argument, so existing calls compile unchanged.
- Multipart requests now also send non-pointer platform params such as `TelegramParams.ChatID`, which were dropped before, and `PlatformParams.GoogleBusiness`, which used to panic.

## [1.13.0] - 2026-08-17

### Added

- **Interactive DMs on Facebook and Instagram.** `MessageSendOptions` gained **`QuickReplies`** (up to 13 tappable chips above the participant's composer, gone once tapped — each `Title` + `Payload`, typed as the new `QuickReply`) and **`Buttons`** (up to 3 attached to the message, each a `web_url` or `postback`, typed as `MessageButton`). An optional **`Card`** (`MessageCard`: `Subtitle`, `ImageURL`, `DefaultAction` as `CardDefaultAction`) fills in the rest of the card carrying `Buttons`, and requires them. Meta's equivalent of what Telegram has had via `ReplyMarkup`.
- **`Message.TappedAction`** — set on inbound messages created by a tap, typed as the new `TappedAction` (`Kind`, `Payload`, `Title`) with the `TappedActionQuickReply` / `TappedActionPostback` / `TappedActionCallbackQuery` constants. Present on the `message.*` webhook payloads too, so you no longer dig through `PlatformData` for the payload you set. Derived rather than stored, so it also resolves for taps recorded earlier, including Instagram ice-breaker taps and Telegram callback queries (`TappedActionCallbackQuery` — the one part of this that isn't Meta-only).
- `Message.QuickReplies`, `Message.Buttons`, and `Message.Card`, echoing back what was sent.
- Quick-replies and buttons examples in `examples/manage-messages`, plus the previously undocumented `ReplyMarkup` / `ReplyToExternalID` in the README's Direct Messages section.

### Notes

- Buttons are delivered as a Meta generic template and your `Body` becomes its element title, so **`Body` is capped at 80 characters when buttons are present** — Meta's limit, not ours; longer text is rejected with a `422` naming the length. Buttons cannot be combined with media.
- **Instagram is stricter than Messenger**: it delivers quick replies only on a plain-text message, so `QuickReplies` with media or with `Buttons` returns `422` on Instagram while both are accepted on Facebook.
- Meta-only — `QuickReplies` / `Buttons` / `Card` return `422` on Telegram and Bluesky chats, where `ReplyMarkup` remains the Telegram equivalent.
- Validation is server-side and names the offending index (e.g. `buttons[1].url must be an https:// URL`), surfacing as the usual error for a `422`. The SDK does not duplicate the limits.
- The new options are sent on the JSON path only. To combine quick replies with an attachment, pass `Media` as a hosted URL rather than uploading via `MediaFiles`.

## [1.12.0] - 2026-08-06

### Added

- **Post syncs & backfill.** `Profiles.BackfillPosts(ctx, id, from, opts)` walks a profile's feed backwards from the newest post and imports the history behind it; `Profiles.PostSyncs(ctx, id, opts)` and `Profiles.PostSync(ctx, id, postSyncID, opts)` expose every post pull — the one fired on connect, the recurring poll, and backfills — as the new `PostSync` type, with `PostSyncTrigger` and `PostSyncStatus` constants and `PostSyncListOptions`.
- **`Comments.ListAll(ctx, opts)`** — comments across every post in the profile group in one request, via the new `BulkCommentListOptions`. Flat: replies are their own entries linked by `ParentExternalID`, typed as the new `BulkComment` (adds `PostID`, `ProfileID`, `Platform`).
- `From` and `To` on `CommentListOptions`, filtering on when PostProxy received the comment.
- **Idempotency.** `postproxy.WithIdempotencyKey(ctx, key)` returns a context that sends an `Idempotency-Key` header on every write made with it, so a dropped connection no longer forces a choice between a duplicate write and a lost one. `IdempotencyKeyFromContext` reads it back. Carried on the context rather than each options struct, so no method signature changed.
- `IsConflictError(err)` for 409 — a duplicate submission (`Response["duplicate_post_id"]`), a backfill already running (`Response["profile_sync_id"]`), or an in-flight idempotency key.
- **Instagram user tags.** `InstagramParams.UserTags` with the new `InstagramUserTag` type (`Username`, `X`, `Y`, `MediaIndex`) — tag accounts on feed posts, reels, and stories.
- `StatsRecord.RawStats` — every metric under its original platform name, alongside the normalized `Stats`.
- `examples/backfill-posts`, and cross-post comment listing in `examples/manage-comments`.

### Changed

- LinkedIn post stats now normalize `likes`, `comments`, `shares`, and `clicks` alongside `impressions` (server-side; `Stats` was already an open map).
- `HUMAN_AGENT` is now approved on **both** Facebook and Instagram and extends the reply window to 7 days. `MessageSendOptions.Tag` is unchanged — see the README for Meta's policy limits.

## [1.11.0] - 2026-07-14

### Added

- `Profiles.IceBreakers`, `Profiles.SetIceBreakers`, and `Profiles.DeleteIceBreakers` for managing Instagram DM ice breakers, with `IceBreaker` and `IceBreakersResponse` types.
- `Profiles.AssignPlacementToGroup(ctx, id, placementID, targetProfileGroupID, opts)` to move a placement (Facebook Page, Telegram channel, GBP location) to another profile group.
- `Placement.Metadata` and `Placement.ProfileGroupID` fields.
- Twitter polls: new `TwitterFormatPoll` constant, and `TwitterParams` gains `PollOptions` (2-4 choices, max 25 chars each) and `PollDurationMinutes` (5-10080).

## [1.10.0] - 2026-06-03

### Added

- Direct Messages API: `Chats` service (`List`, `Create`, `Get`, `Archive`, `Unarchive`) and `Messages` service (`List`, `Send`, `Get`, `Edit`, `React`, `Unreact`). `Messages.Send` uses multipart/form-data when `MediaFiles` are provided, else JSON.
- New models: `Attachment`, `Reaction`, `Chat`, `Message`, plus `MessageDirection` and `MessageStatus` enums.
- `Comments.PrivateReply` — send a private reply (direct message) to a comment author; returns a `*Message`.
- `Comment.Attachments` (`[]Attachment`) and `Comment.Metadata` (`map[string]any`).
- 10 new webhook event types: `profile_comment.created`, the eight `message.*` events, and `reaction.received`, all appended to `WebhookEventTypes`.
- Typed webhook payloads `MessageEventData`, `ReactionEventData`, and `ProfileCommentCreatedData` with `WebhookEvent.AsMessage()`, `AsReaction()`, and `AsProfileCommentCreated()` decoders.

## [1.9.0] - 2026-05-15

### Added

- `PlatformGoogleBusiness` platform value for posts and profiles.
- `ProfileComments` service: `List`, `Get`, `Create`, `Delete` for review replies via `/api/profiles/:profile_id/comments`.
- Per-media platform error reporting: `Media.Platforms` (`[]MediaPlatformError`) with `ErrorDetails`.
- `PlatformParams.GoogleBusiness` (`map[string]any`) for Google Business post parameters.

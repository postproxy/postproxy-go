# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

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

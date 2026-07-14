# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

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

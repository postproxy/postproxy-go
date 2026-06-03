# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

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

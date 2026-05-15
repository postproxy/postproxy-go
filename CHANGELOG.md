# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [1.9.0] - 2026-05-15

### Added

- `PlatformGoogleBusiness` platform value for posts and profiles.
- `ProfileComments` service: `List`, `Get`, `Create`, `Delete` for review replies via `/api/profiles/:profile_id/comments`.
- Per-media platform error reporting: `Media.Platforms` (`[]MediaPlatformError`) with `ErrorDetails`.
- `PlatformParams.GoogleBusiness` (`map[string]any`) for Google Business post parameters.

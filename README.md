# Postproxy Go SDK

The official Go SDK for the [Postproxy](https://postproxy.dev) API.

## Installation

```bash
go get github.com/postproxy/postproxy-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/postproxy/postproxy-go"
)

func main() {
	client := postproxy.NewClient("your-api-key", postproxy.WithProfileGroupID("your-profile-group-id"))
	ctx := context.Background()

	// Get profiles from the profile group
	profiles, err := client.Profiles.List(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	profileIDs := make([]string, len(profiles.Data))
	for i, p := range profiles.Data {
		fmt.Printf("%s (%s)\n", p.Name, p.Platform)
		profileIDs[i] = p.ID
	}

	// Create a post to all profiles
	post, err := client.Posts.Create(ctx, "Hello world!", profileIDs, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created post: %s\n", post.ID)
}
```

## Client Options

```go
// Custom base URL
client := postproxy.NewClient("key", postproxy.WithBaseURL("https://custom.api"))

// Custom HTTP client
client := postproxy.NewClient("key", postproxy.WithHTTPClient(&http.Client{
	Timeout: 30 * time.Second,
}))

// Default profile group ID
client := postproxy.NewClient("key", postproxy.WithProfileGroupID("pg-123"))
```

### Idempotency

`WithIdempotencyKey` returns a context that carries an `Idempotency-Key` header on every
write request (`POST`/`PUT`/`PATCH`/`DELETE`) made with it. If the connection drops before
you see the response, retry with the same key and you get the original response back
instead of a second post:

```go
ctx := postproxy.WithIdempotencyKey(context.Background(), uuid.NewString())

post, err := client.Posts.Create(ctx, "Hello", []string{"profile-id"}, nil)

// Retrying the same call with the same context replays the original response.
```

The key rides on the context rather than each method's options struct, so it works
uniformly across every write without changing any signature. Use a fresh key — and a fresh
context — per logical operation; a UUID is ideal. Keys are scoped to your account and may
be up to 255 characters. The SDK never generates keys or retries for you.

| Situation | Result |
|---|---|
| First request with the key | Runs normally |
| Retry after a success | Original status and body replayed |
| Retry while the first is still running | 409 — `IsConflictError(err)`; wait and retry |
| Same key, different request body | 422 — `IsValidationError(err)` |
| Retry after an error response | Runs normally — errors are not replayed |

Only successful (`2xx`) responses are stored, so a request that failed validation or hit a
quota leaves the key free — fix the payload and retry with the same key. Stored responses
are kept for **24 hours**. Requests made with a plain context are unaffected.

## Resources

### Posts

```go
// List posts with filtering
page := 1
status := postproxy.PostStatusDraft
posts, err := client.Posts.List(ctx, &postproxy.PostListOptions{
	Page:   &page,
	Status: &status,
})

// Get a post
post, err := client.Posts.Get(ctx, "post-id", nil)

// Create a post with media URLs
post, err := client.Posts.Create(ctx, "Caption", []string{"profile-id"}, &postproxy.PostCreateOptions{
	Media: []string{"https://example.com/image.jpg"},
})

// Create a post with local file upload
post, err := client.Posts.Create(ctx, "Caption", []string{"profile-id"}, &postproxy.PostCreateOptions{
	MediaFiles: []string{"./photo.jpg"},
})

// Create a post with platform-specific parameters
igFormat := postproxy.InstagramFormatReel
parseMode := postproxy.TelegramParseModeHTML
disablePreview := true
tagX, tagY := 0.5, 0.4
slide := 2
post, err := client.Posts.Create(ctx, "Caption", []string{"profile-id"}, &postproxy.PostCreateOptions{
	Platforms: &postproxy.PlatformParams{
		Instagram: &postproxy.InstagramParams{
			Format: &igFormat,
			// Tag public Instagram accounts. Images require X and Y (floats
			// 0.0–1.0 from the top-left corner); reels and video slides are
			// tagged by username only. MediaIndex picks the carousel slide
			// (0-based, defaults to 0).
			UserTags: []postproxy.InstagramUserTag{
				{Username: "natgeo", X: &tagX, Y: &tagY},
				{Username: "spacex", MediaIndex: &slide},
			},
		},
		Telegram: &postproxy.TelegramParams{
			ChatID:             "-1001234567890",
			ParseMode:          &parseMode,
			DisableLinkPreview: &disablePreview,
		},
		// Google Business: untyped map. Formats: "standard", "event", "offer".
		// CTA actions: LEARN_MORE, BOOK, ORDER, SHOP, SIGN_UP, CALL. Max 1 image (≤5 MB).
		GoogleBusiness: map[string]any{
			"format":          "standard",
			"location_id":     "accounts/123/locations/456",
			"cta_action_type": "LEARN_MORE",
			"cta_url":         "https://example.com",
		},
	},
})

// Create a Twitter poll (2-4 options, max 25 chars each; 5-10080 minutes)
pollFormat := postproxy.TwitterFormatPoll
pollDuration := 1440
post, err = client.Posts.Create(ctx, "Which framework?", []string{"twitter"}, &postproxy.PostCreateOptions{
	Platforms: &postproxy.PlatformParams{
		Twitter: &postproxy.TwitterParams{
			Format:              &pollFormat,
			PollOptions:         []string{"Rails", "Django", "Laravel"},
			PollDurationMinutes: &pollDuration,
		},
	},
})

// LinkedIn poll (2-4 options, max 30 chars each; 1, 3, 7 or 14 days)
liPoll := postproxy.LinkedInFormatPoll
question := "Tabs or spaces?"
days := 3
post, err = client.Posts.Create(ctx, "Settle this", []string{"linkedin"}, &postproxy.PostCreateOptions{
	Platforms: &postproxy.PlatformParams{
		LinkedIn: &postproxy.LinkedInParams{
			Format:           &liPoll,
			PollQuestion:     &question,
			PollOptions:      []string{"Tabs", "Spaces"},
			PollDurationDays: &days,
		},
	},
})

// Alt text: one entry per media item, in order. A single entry is sent as a
// plain string and applies to the first media item. Facebook, Instagram,
// LinkedIn, Bluesky, Threads, Pinterest and Twitter.
post, err = client.Posts.Create(ctx, "Foxes", []string{"bluesky"}, &postproxy.PostCreateOptions{
	Media: []string{"https://example.com/fox1.jpg", "https://example.com/fox2.jpg"},
	Platforms: &postproxy.PlatformParams{
		Bluesky: &postproxy.BlueskyParams{AltText: postproxy.AltText{"A red fox in the snow", "The same fox, mid-leap"}},
	},
})

// More platform params
link := "https://example.com/launch"
topic := "Coffee"
trial := postproxy.InstagramTrialStrategyManual // or InstagramTrialStrategySSPerformance
audioID := "1234567890"                         // from Profiles.InstagramAvailableAudio
madeWithAI, draftOnTikTok := true, true
platforms := &postproxy.PlatformParams{
	Facebook:  &postproxy.FacebookParams{Link: &link, Countries: []string{"US", "CA"}, MadeWithAI: &madeWithAI},
	Instagram: &postproxy.InstagramParams{TrialStrategy: &trial, AudioID: &audioID},
	Threads:   &postproxy.ThreadsParams{TopicTag: &topic},
	TikTok:    &postproxy.TikTokParams{UploadAsDraft: &draftOnTikTok},
	// LinkedIn: ReshareURN, Countries. Instagram: AudioVolume, VideoVolume, StripAudio, MadeWithAI.
	// TikTok: CoverURL, ThumbOffset.
}

// Instagram: turn comments off, tag a location (a Facebook Page id with an
// address; post and reel, not stories) and keep a reel off the feed/grid.
// LinkedIn: a title for an attached document and a first comment (not on
// polls). YouTube: add the video to a playlist from Profiles.YouTubePlaylists.
reel := postproxy.InstagramFormatReel
disableComments, shareToFeed := true, false
locationID := "110843418940484"
docTitle, firstComment := "Q3 report", "Full deck in the PDF"
ytTitle, playlistID := "Intro", "PL1234567890"
platforms = &postproxy.PlatformParams{
	Instagram: &postproxy.InstagramParams{
		Format:         &reel,
		DisableComment: &disableComments,
		LocationID:     &locationID,
		ShareToFeed:    &shareToFeed, // reels only
	},
	LinkedIn: &postproxy.LinkedInParams{DocumentTitle: &docTitle, FirstComment: &firstComment},
	YouTube:  &postproxy.YouTubeParams{Title: &ytTitle, PlaylistID: &playlistID},
}

// Custom cover from disk. CoverFile (Instagram, TikTok, YouTube, Pinterest) is
// a local path the SDK uploads with the post — the request switches to
// multipart even without MediaFiles. A base64 "data:image/jpeg;base64,..." URL
// is sent as-is instead. Works on Create and Update.
videoFormat := postproxy.TikTokFormatVideo
cover := "./cover.jpg"
post, err = client.Posts.Create(ctx, "Behind the scenes", []string{"tiktok"}, &postproxy.PostCreateOptions{
	Media: []string{"https://example.com/bts.mp4"},
	Platforms: &postproxy.PlatformParams{
		TikTok: &postproxy.TikTokParams{Format: &videoFormat, CoverFile: &cover},
	},
})

// Publish even if an identical post was just submitted (skips the 409)
post, err = client.Posts.Create(ctx, "Same again", []string{"profile-id"}, &postproxy.PostCreateOptions{
	SkipDuplicateCheck: true, // sends X-Skip-Duplicate-Check: true
})

// Create a thread post
post, err := client.Posts.Create(ctx, "Thread starts here", []string{"profile-id"}, &postproxy.PostCreateOptions{
	Thread: []postproxy.ThreadChildInput{
		{Body: "Second post in the thread"},
		{Body: "Third with media", Media: []string{"https://example.com/img.jpg"}},
	},
})
for _, child := range post.Thread {
	// Each thread item is delivered separately; Platforms shows whether it went out.
	fmt.Printf("%s: %s (%d platforms)\n", child.ID, child.Body, len(child.Platforms))
}

// Each PlatformResult carries ProfileID, ProfileName, Permalink and any
// Attachments PostProxy generated for that platform (e.g. a TikTok cover).
for _, p := range post.Platforms {
	fmt.Println(p.Platform, p.ProfileName, p.Status, p.Permalink)
}

// Update a draft or scheduled post
newBody := "Updated content!"
post, err := client.Posts.Update(ctx, "post-id", &postproxy.PostUpdateOptions{
	Body: &newBody,
})

// Edit a published post. Only Body, Platforms.YouTube.Title,
// Platforms.Pinterest.Title and Platforms.Pinterest.DestinationLink can change;
// the edit is pushed to the platforms in the background and reported through
// the platform_post.edited / platform_post.edit_failed webhooks.
ytTitle := "Launch day — full walkthrough"
post, err = client.Posts.Update(ctx, "published-post-id", &postproxy.PostUpdateOptions{
	Body:      &newBody,
	Platforms: &postproxy.PlatformParams{YouTube: &postproxy.YouTubeParams{Title: &ytTitle}},
})

// Update platform params only
ytPrivacy := "unlisted"
post, err := client.Posts.Update(ctx, "post-id", &postproxy.PostUpdateOptions{
	Platforms: &postproxy.PlatformParams{
		YouTube: &postproxy.YouTubeParams{
			PrivacyStatus: &ytPrivacy,
		},
	},
})

// Replace profiles and media
post, err := client.Posts.Update(ctx, "post-id", &postproxy.PostUpdateOptions{
	Profiles: []string{"twitter", "threads"},
	Media:    []string{"https://example.com/new-image.jpg"},
})

// Replace thread children
post, err := client.Posts.Update(ctx, "post-id", &postproxy.PostUpdateOptions{
	Thread: []postproxy.ThreadChildInput{
		{Body: "Updated first reply"},
		{Body: "Updated second reply", Media: []string{"https://example.com/img.jpg"}},
	},
})

// Remove all media
post, err := client.Posts.Update(ctx, "post-id", &postproxy.PostUpdateOptions{
	Media: []string{},
})

// Publish a draft
post, err := client.Posts.PublishDraft(ctx, "post-id", nil)

// Delete a post
result, err := client.Posts.Delete(ctx, "post-id", nil)

// Delete a post and also remove it from social platforms
truthy := true
result, err = client.Posts.Delete(ctx, "post-id", &postproxy.PostDeleteOptions{DeleteOnPlatform: &truthy})

// Delete from platforms only (keeps DB record). Defaults to all platforms.
r1, err := client.Posts.DeleteOnPlatform(ctx, "post-id", nil)
// Target a single network
network := "twitter"
r2, err := client.Posts.DeleteOnPlatform(ctx, "post-id", &postproxy.PostDeleteOnPlatformOptions{Network: &network})
// Target a specific profile
profID := "prof-abc"
r3, err := client.Posts.DeleteOnPlatform(ctx, "post-id", &postproxy.PostDeleteOnPlatformOptions{ProfileID: &profID})
// Target a specific post profile (covers entire thread for that profile)
ppID := "pp-abc"
r4, err := client.Posts.DeleteOnPlatform(ctx, "post-id", &postproxy.PostDeleteOnPlatformOptions{PostProfileID: &ppID})

// Get stats for posts
from := "2026-02-01T00:00:00Z"
to := "2026-02-24T00:00:00Z"
stats, err := client.Posts.Stats(ctx, []string{"post-id-1", "post-id-2"}, &postproxy.PostStatsOptions{
	Profiles: []string{"instagram", "twitter"},
	From:     &from,
	To:       &to,
})
for postID, postStats := range stats.Data {
	for _, plat := range postStats.Platforms {
		fmt.Printf("%s on %s: %d records\n", postID, plat.Platform, len(plat.Records))
	}
}
```

### Queues

```go
// List all queues
queues, err := client.Queues.List(ctx, nil)

// Get a queue
queue, err := client.Queues.Get(ctx, "queue-id")

// Get next available slot
nextSlot, err := client.Queues.NextSlot(ctx, "queue-id")
fmt.Println(nextSlot.NextSlot)

// Create a queue with timeslots
tz := "America/New_York"
jitter := 10
queue, err := client.Queues.Create(ctx, "Morning Posts", "profile-group-id", &postproxy.QueueCreateOptions{
	Timezone: &tz,
	Jitter:   &jitter,
	Timeslots: []postproxy.TimeslotInput{
		{Day: 1, Time: "09:00"},
		{Day: 2, Time: "09:00"},
		{Day: 3, Time: "09:00"},
	},
})

// Update a queue
newJitter := 15
queue, err := client.Queues.Update(ctx, "queue-id", &postproxy.QueueUpdateOptions{
	Jitter: &newJitter,
	Timeslots: []postproxy.TimeslotInput{
		{Day: 6, Time: "10:00"},            // add new timeslot
		{ID: 1, Destroy: true},             // remove existing timeslot
	},
})

// Pause/unpause a queue
enabled := false
queue, err := client.Queues.Update(ctx, "queue-id", &postproxy.QueueUpdateOptions{
	Enabled: &enabled,
})

// Delete a queue
result, err := client.Queues.Delete(ctx, "queue-id")

// Add a post to a queue
queueID := "queue-id"
priority := "high"
post, err := client.Posts.Create(ctx, "Queued post", []string{"profile-id"}, &postproxy.PostCreateOptions{
	QueueID:       &queueID,
	QueuePriority: &priority,
})
```

### Webhooks

```go
// List webhooks
webhooks, err := client.Webhooks.List(ctx)

// Get a webhook
webhook, err := client.Webhooks.Get(ctx, "wh-id")

// Create a webhook
webhook, err := client.Webhooks.Create(ctx, "https://example.com/webhook", []string{"post.published", "post.failed"}, &postproxy.WebhookCreateOptions{
	Description: strPtr("My webhook"),
})
fmt.Println(webhook.ID, webhook.Secret)

// Update a webhook
enabled := false
webhook, err := client.Webhooks.Update(ctx, "wh-id", &postproxy.WebhookUpdateOptions{
	Events:  []string{"post.published"},
	Enabled: &enabled,
})

// Delete a webhook
result, err := client.Webhooks.Delete(ctx, "wh-id")

// List deliveries
deliveries, err := client.Webhooks.Deliveries(ctx, "wh-id", nil)
for _, d := range deliveries.Data {
	fmt.Printf("%s: %v\n", d.EventType, d.Success)
}
```

#### Signature verification

Verify incoming webhook signatures using HMAC-SHA256:

```go
valid := postproxy.VerifyWebhookSignature(
	string(requestBody),
	r.Header.Get("X-PostProxy-Signature"),
	"whsec_...",
)
```

#### Event types and typed payloads

Subscribe to any of these events (or pass `[]string{"*"}` for all):

`post.processed`, `post.imported`, `platform_post.published`, `platform_post.failed`, `platform_post.failed_waiting_for_retry`, `platform_post.edited`, `platform_post.edit_failed`, `platform_post.insights`, `profile.connected`, `profile.disconnected`, `profile.stats`, `profile.placements_synced`, `placement.orphaned`, `placement.restored`, `media.failed`, `comment.created`, `comment.failed`, `comment.failed_waiting_for_retry`, `comment.delivery_unconfirmed`, `comment.edited`, `comment.edit_failed`, `profile_comment.created`, `profile_comment.failed`, `profile_comment.failed_waiting_for_retry`, `profile_comment.edited`, `profile_comment.deleted`, `message.received`, `message.sent`, `message.delivered`, `message.read`, `message.edited`, `message.deleted`, `message.failed_waiting_for_retry`, `message.failed`, `reaction.received`, `referral.received`, `chat.thread_control_changed`, `chat.thread_control_requested`, `message_template.status_updated`, `message_template.category_updated`, `placement.name_status_updated`, `placement.quality_updated`, `whatsapp.automatic_event`.

`ParseWebhookEvent` validates the envelope; the `As*` helpers decode `Data` into a typed payload:

```go
event, err := postproxy.ParseWebhookEvent(requestBody)
if err != nil {
	// errors.Is(err, postproxy.ErrUnknownWebhookEvent) for unknown types
	return
}
switch event.Type {
case postproxy.EventProfileStats:
	data, _ := event.AsProfileStats()
	fmt.Println(data.ProfileID, data.Stats)
case postproxy.EventPlatformPostPublished:
	data, _ := event.AsPlatformPost()
	fmt.Println("Published:", data.PlatformID)
case postproxy.EventCommentCreated:
	data, _ := event.AsCommentCreated()
	fmt.Println(*data.AuthorUsername, ":", data.Body)
case postproxy.EventMessageReceived,
	postproxy.EventMessageSent,
	postproxy.EventMessageDelivered,
	postproxy.EventMessageRead,
	postproxy.EventMessageEdited,
	postproxy.EventMessageDeleted,
	postproxy.EventMessageFailedWaitingForRetry,
	postproxy.EventMessageFailed:
	data, _ := event.AsMessage()
	fmt.Println(data.Message.Direction, data.Message.ID)
case postproxy.EventReactionReceived:
	data, _ := event.AsReaction()
	fmt.Println(data.Action, *data.Emoji)
case postproxy.EventProfileCommentCreated:
	data, _ := event.AsProfileCommentCreated()
	fmt.Println(data.ID, ":", data.Body)
case postproxy.EventPlatformPostEditFailed:
	data, _ := event.AsPlatformPost()
	fmt.Println("Edit failed on", data.Platform, data.Error)
case postproxy.EventCommentEditFailed:
	data, _ := event.AsCommentCreated()
	fmt.Println(data.ID, data.Error)
case postproxy.EventProfilePlacementsSynced:
	data, _ := event.AsPlacementsSynced()
	fmt.Println(data.ID, len(data.Placements))
case postproxy.EventPlacementQualityUpdated:
	data, _ := event.AsPlacement()
	fmt.Println(data.Name, data.Event, data.CurrentLimit)
case postproxy.EventMessageTemplateStatusUpdated:
	data, _ := event.AsMessageTemplate()
	fmt.Println(data.Name, data.Status)
case postproxy.EventReferralReceived:
	data, _ := event.AsReferral()
	fmt.Println(data.Chat.ID, data.CtwaClid)
case postproxy.EventChatThreadControlChanged:
	data, _ := event.AsThreadControlChanged()
	fmt.Println("Chat", data.Chat.ID, "now owned by", data.NewOwnerAppID)
case postproxy.EventProfileDisconnected:
	data, _ := event.AsProfileEvent()
	fmt.Println(data.Name, data.Origin, data.Error) // Error, ErrorDetails when PostProxy disconnected it
case postproxy.EventWhatsAppAutomaticEvent:
	data, _ := event.AsWhatsAppAutomaticEvent()
	fmt.Println(data.EventName, data.CustomData)
}
```

The typed payloads `MessageEventData` (the 8 `message.*` events), `ReactionEventData` (`reaction.received`), and `ProfileCommentCreatedData` (`profile_comment.created`) each embed the relevant resource. `MessageEventData` and `ReactionEventData` carry a full `Message` in their `Message` field.

| Events | Decoder | Payload |
|---|---|---|
| `platform_post.*` (incl. `edited`, `edit_failed`) | `AsPlatformPost` | `PlatformPostData` |
| `comment.*` | `AsCommentCreated` | `CommentCreatedData` (`Error`, `ErrorDetails` on failures) |
| `profile_comment.*` | `AsProfileCommentCreated` | `ProfileCommentCreatedData` |
| `profile.placements_synced` | `AsPlacementsSynced` | `PlacementsSyncedData` |
| `placement.*` | `AsPlacement` | `PlacementEventData` |
| `referral.received` | `AsReferral` | `ReferralEventData` |
| `chat.thread_control_changed` | `AsThreadControlChanged` | `ThreadControlChangedEventData` |
| `chat.thread_control_requested` | `AsThreadControlRequested` | `ThreadControlRequestedEventData` |
| `profile.connected`, `profile.disconnected` | `AsProfileEvent` | `ProfileEventData` (`Origin`, `Error`, `ErrorDetails` on disconnects) |
| `message_template.*` | `AsMessageTemplate` | `MessageTemplateEventData` |
| `whatsapp.automatic_event` | `AsWhatsAppAutomaticEvent` | `WhatsAppAutomaticEventData` |

### Comments

```go
// List comments on a post (paginated)
comments, err := client.Comments.List(ctx, "post-id", "profile-id", nil)
for _, comment := range comments.Data {
	fmt.Printf("%s: %s\n", comment.AuthorUsername, comment.Body)
	for _, reply := range comment.Replies {
		fmt.Printf("  %s: %s\n", reply.AuthorUsername, reply.Body)
	}
}

// List with pagination
page := 2
perPage := 10
comments, err := client.Comments.List(ctx, "post-id", "profile-id", &postproxy.CommentListOptions{
	Page:    &page,
	PerPage: &perPage,
})

// Filter by when PostProxy received the comment (created_at, not posted_at).
// A bare date means that date's start of day. Applies to top-level comments —
// one in range brings its full Replies slice with it.
from := "2026-03-25"
to := "2026-03-26T12:00:00Z"
recent, err := client.Comments.List(ctx, "post-id", "profile-id", &postproxy.CommentListOptions{
	From: &from,
	To:   &to,
})

// Get a single comment
comment, err := client.Comments.Get(ctx, "post-id", "comment-id", "profile-id")

// Create a comment
comment, err := client.Comments.Create(ctx, "post-id", "profile-id", &postproxy.CommentCreateOptions{
	Text: "Great post!",
})

// Reply to a comment
parentID := "comment-id"
reply, err := client.Comments.Create(ctx, "post-id", "profile-id", &postproxy.CommentCreateOptions{
	Text:     "Thanks!",
	ParentID: &parentID,
})

// Delete a comment
result, err := client.Comments.Delete(ctx, "post-id", "comment-id", "profile-id")
fmt.Println(result.Accepted) // true

// Hide / unhide a comment
client.Comments.Hide(ctx, "post-id", "comment-id", "profile-id")
client.Comments.Unhide(ctx, "post-id", "comment-id", "profile-id")

// Like / unlike a comment
client.Comments.Like(ctx, "post-id", "comment-id", "profile-id")
client.Comments.Unlike(ctx, "post-id", "comment-id", "profile-id")

// Edit a comment you published (Facebook, YouTube). Runs in the background —
// the outcome arrives as comment.edited / comment.edit_failed. The edited
// comment carries ExternalEditedAt.
client.Comments.Edit(ctx, "post-id", "comment-id", "profile-id", "Fixed typo")

// Pin / unpin a comment (TikTok)
client.Comments.Pin(ctx, "post-id", "comment-id", "profile-id")
client.Comments.Unpin(ctx, "post-id", "comment-id", "profile-id")
```

Comments may carry media `Attachments` (`[]Attachment`) and author signals in `Metadata` (`map[string]any`):

```go
for _, att := range comment.Attachments {
	fmt.Printf("%s -> %v (%s)\n", att.Type, att.URL, att.Status)
}
if verified, ok := comment.Metadata["is_verified_user"].(bool); ok && verified {
	fmt.Println("verified author")
}
```

Reply privately to a comment via direct message (Instagram/Facebook). This returns a `*Message`, not a comment:

```go
msg, err := client.Comments.PrivateReply(ctx, "post-id", "comment-id", "profile-id", "Thanks — DM-ing you the details.")
fmt.Println(msg.ID, msg.ChatID, msg.Status)

// One private reply per comment, within 7 days of it. QuickReplies, Buttons,
// Card and Elements work as in Messages.Send; text may be empty with Elements.
msg, err = client.Comments.PrivateReply(ctx, "post-id", "comment-id", "profile-id", "", &postproxy.PrivateReplyOptions{
	Elements: []postproxy.MessageElement{{
		Title:    "Red sneaker",
		ImageURL: "https://cdn.example.com/red.png",
		Buttons:  []postproxy.MessageButton{{Type: "web_url", Title: "Buy", URL: "https://shop.example.com/red"}},
	}},
})
```

#### Comments across posts

`Comments.ListAll` returns comments spanning every post in the profile group in one
request — the comments counterpart to `Posts.Stats`. Every filter is optional.

**This list is flat.** Unlike the per-post list, replies are not nested: every comment,
top-level or reply, is its own entry linked to its parent by `ParentExternalID`, so `Total`
counts every comment and paging is exact. Entries are `BulkComment`, which adds `PostID`,
`ProfileID`, and `Platform` and drops `Replies`.

```go
perPage := 50
all, err := client.Comments.ListAll(ctx, &postproxy.BulkCommentListOptions{
	Profiles: []string{"instagram", "prof-abc"},  // profile IDs or network names, mixed
	PostIDs:  []string{"post-1", "post-2"},       // omit for every post in scope
	// PlacementID: &pageID,                      // narrow to one page / organization
	From:     &from,
	PerPage:  &perPage,                           // max 100
})

for _, c := range all.Data {
	// Each entry says where it came from, so you can act on it with the
	// post-scoped methods above.
	fmt.Println(c.Platform, c.PostID, c.ProfileID, c.Body)

	if c.ParentExternalID != nil {
		fmt.Println("  ↳ reply to", *c.ParentExternalID)
	}
}
```

Unknown or out-of-scope IDs in `PostIDs` and `Profiles` are ignored rather than erroring.
Results are ordered newest first by receipt time.

### Direct Messages

Manage direct-message conversations (chats) and the messages within them across Facebook, Instagram, Telegram, Bluesky, and WhatsApp.

```go
// List chats for a profile (paginated)
perPage := 20
chats, err := client.Chats.List(ctx, "profile-id", &postproxy.ChatListOptions{PerPage: &perPage})
for _, chat := range chats.Data {
	fmt.Printf("%v (%s)\n", chat.ParticipantUsername, chat.Platform)
}

// Find or create a chat with a participant
username := "jane_doe"
chat, err := client.Chats.Create(ctx, "profile-id", "igsid_8675309", &postproxy.ChatCreateOptions{
	ParticipantUsername: &username,
})

// Get a single chat
chat, err = client.Chats.Get(ctx, chat.ID, nil)

// Archive / unarchive a chat (Bluesky only)
client.Chats.Archive(ctx, chat.ID, nil)
client.Chats.Unarchive(ctx, chat.ID, nil)

// Mark the latest inbound message read (shows the read receipt)
client.Chats.MarkRead(ctx, chat.ID, nil)

// Typing indicator (Facebook, Instagram, WhatsApp)
client.Chats.StartTyping(ctx, chat.ID, nil)
client.Chats.StopTyping(ctx, chat.ID, nil)

// Thread control (Facebook & Instagram): move the conversation between Meta
// apps sharing the inbox. Each returns the chat; Metadata["thread_control"]
// holds the owner.
note := "escalated to support"
client.Chats.PassThreadControl(ctx, chat.ID, "263902037430900", &postproxy.ThreadControlOptions{Metadata: &note})
client.Chats.TakeThreadControl(ctx, chat.ID, nil)
client.Chats.ReleaseThreadControl(ctx, chat.ID, nil)
client.Chats.RequestThreadControl(ctx, chat.ID, nil) // owner decides; watch chat.thread_control_changed
owner, err := client.Chats.ThreadOwner(ctx, chat.ID, nil) // AppID, ExpiresAt (both nil when idle)

// List messages in a chat
dir := postproxy.MessageDirectionInbound
messages, err := client.Messages.List(ctx, chat.ID, &postproxy.MessageListOptions{Direction: &dir})
for _, msg := range messages.Data {
	fmt.Printf("[%s] %v\n", msg.Direction, msg.Body)
}

// Send a text message
body := "Yes, we ship worldwide!"
sent, err := client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{Body: &body})

// Send outside the 24h window with a tag (Facebook/Instagram)
tag := "HUMAN_AGENT"
client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{Body: &body, Tag: &tag})

// Send media by hosted URL, or from a local file (multipart)
client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{Media: []string{"https://cdn.example.com/photo.png"}})
client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{MediaFiles: []string{"./photo.png"}})

// Get / edit a message (edit is Telegram only)
msg, err := client.Messages.Get(ctx, sent.ID, nil)
updated := "Updated answer."
client.Messages.Edit(ctx, sent.ID, &postproxy.MessageEditOptions{Body: &updated})

// React / unreact (Facebook & Instagram)
reaction := "love"
emoji := "❤️"
client.Messages.React(ctx, sent.ID, &postproxy.MessageReactOptions{Reaction: &reaction, Emoji: &emoji})
client.Messages.Unreact(ctx, sent.ID, nil)

// Telegram: thread under a message and attach an inline keyboard
// (ReplyToExternalID also works on WhatsApp and Instagram)
replyTo := "4821"
client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{
	Body:              &body,
	ReplyToExternalID: &replyTo,
	ReplyMarkup: map[string]any{
		"inline_keyboard": []any{[]any{map[string]any{"text": "Track order", "callback_data": "track:1"}}},
	},
})
```

#### Quick replies and buttons (Facebook & Instagram)

Meta's two interactive primitives. **Quick replies** are chips above the participant's
composer that disappear once tapped; **buttons** are attached to the message and stay in
the thread. Telegram's equivalent is `ReplyMarkup` above — passing `QuickReplies` or
`Buttons` on a Telegram or Bluesky chat returns `422`.

```go
// Quick replies — up to 13. Title ≤ 20 chars, Payload ≤ 1000.
prompt := "What can I help with?"
client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{
	Body: &prompt,
	QuickReplies: []postproxy.QuickReply{
		{Title: "Track order", Payload: "TRACK"},
		{Title: "Talk to support", Payload: "HELP"},
	},
})

// Buttons — up to 3, each either web_url or postback. Card is optional and
// requires Buttons.
shipped := "Your order shipped"
client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{
	Body: &shipped,
	Buttons: []postproxy.MessageButton{
		{Type: "web_url", Title: "Track", URL: "https://shop.example.com/o/123"},
		{Type: "postback", Title: "Cancel", Payload: "CANCEL:123"},
	},
	Card: &postproxy.MessageCard{
		Subtitle: "Arriving Friday",
		ImageURL: "https://cdn.example.com/shoe.png",
		DefaultAction: &postproxy.CardDefaultAction{
			Type: "web_url", URL: "https://shop.example.com/o/123",
		},
	},
})
```

Buttons are delivered as a Meta generic template and your `Body` becomes the template's
element title — so **`Body` is capped at 80 characters when buttons are present**. That is
Meta's limit, not Postproxy's, and a longer body is rejected with a `422` naming the
length. Buttons cannot be combined with media. Instagram is stricter than Messenger: it
delivers quick replies only on a plain-text message, so `QuickReplies` with media or with
`Buttons` returns `422` on Instagram while both are accepted on Facebook.

Validation happens server-side and names the offending index — `buttons[1].url must be an
https:// URL` — surfacing as the SDK's usual error for a `422`.

> The new options are sent on the JSON path only. To combine quick replies with an
> attachment, pass `Media` as a hosted URL rather than uploading via `MediaFiles`.

A tap comes back as an **inbound message** carrying `TappedAction`:

```go
inbound, _ := client.Messages.List(ctx, chat.ID, &postproxy.MessageListOptions{Direction: &dir})
for _, msg := range inbound.Data {
	if msg.TappedAction != nil {
		// postproxy.TappedActionQuickReply / TappedActionPostback / TappedActionCallbackQuery
		fmt.Println(msg.TappedAction.Kind, msg.TappedAction.Payload)
	}
}
```

Subscribe to `message.received` to react to taps as they happen — the same field is on the
webhook payload. `TappedAction` is derived rather than stored, so it also resolves for taps
recorded before Postproxy exposed it, including Instagram ice-breaker taps and Telegram
callback queries (`TappedActionCallbackQuery`). A tap also opens the 24h window.

#### Carousels (Facebook & Instagram)

`Elements` sends a horizontally scrolling carousel of up to 10 cards. It can't be combined
with `Body`, `Buttons` or `Card`. Sent messages echo it back in `Message.Elements`;
`Message.StoryReply` is set on inbound replies to your stories.

```go
client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{
	Elements: []postproxy.MessageElement{
		{
			Title:    "Red sneaker",
			Subtitle: "$89 — sizes 38-46",
			ImageURL: "https://cdn.example.com/red.png",
			Buttons:  []postproxy.MessageButton{{Type: "web_url", Title: "Buy", URL: "https://shop.example.com/red"}},
		},
		{
			Title:    "Blue sneaker",
			ImageURL: "https://cdn.example.com/blue.png",
			Buttons:  []postproxy.MessageButton{{Type: "postback", Title: "Notify me", Payload: "NOTIFY:blue"}},
		},
	},
})
```

#### WhatsApp

A WhatsApp chat runs through one of the profile's phone numbers, so `Chats.Create` needs
its `PlacementID`. Send exactly one of `Body`, `Media`, `Template`, `Interactive`,
`Location` or `Contacts`. Outside the 24-hour customer service window
(`Chat.WithinMessagingWindow`) only a `Template` gets through.

```go
numbers, _ := client.Profiles.Placements(ctx, "whatsapp-profile-id", nil)
chat, err := client.Chats.Create(ctx, "whatsapp-profile-id", "15551234567", &postproxy.ChatCreateOptions{
	PlacementID: &numbers.Data[0].ID,
})

// Template by name (+ language) or ID. Variables are positional: header, body,
// then dynamic URL buttons.
client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{
	Template: &postproxy.MessageTemplateInput{
		Name:      "order_shipped",
		Language:  "en_US",
		Variables: []string{"Jane", "#1234"},
		ButtonParams: []postproxy.TemplateButtonParam{
			{Index: 0, SubType: "url", Extra: map[string]any{"text": "1234"}},
		},
	},
})

// Location, contacts and interactive payloads
client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{
	Location: &postproxy.MessageLocation{Latitude: 52.52, Longitude: 13.405, Name: "Pickup point"},
})
client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{
	Contacts: []map[string]any{{"name": map[string]any{"formatted_name": "Jane"}}},
})

// Send options: Category ("utility"), LinkPreview, VoiceNote (audio as a voice
// note) and Filename (shown for a document)
voiceNote := true
client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{
	Media:     []string{"https://files.example.com/hello.ogg"},
	VoiceNote: &voiceNote,
})
```

`Chat` gains `ExternalPlacementID` and `Group` (a WhatsApp group chat); `Message` gains
`Template`, `Interactive`, `Location`, `Contacts` and `ErrorDetails`. Download an inbound
attachment's bytes with `WhatsApp.Media` (see [WhatsApp](#whatsapp-numbers)).

### Profile comments (Google Business reviews)

Profile-level comments expose Google Business reviews and replies. Reviews are user-generated — the SDK lets you list/get them and reply to or delete your own replies. Reviews sync twice daily.

```go
// List reviews for a profile (paginated)
reviews, err := client.ProfileComments.List(ctx, "profile-id", nil)

// Filter by placement (location)
placement := "accounts/123/locations/456"
reviews, err := client.ProfileComments.List(ctx, "profile-id", &postproxy.ProfileCommentListOptions{
    PlacementID: &placement,
})

// Get a single review
review, err := client.ProfileComments.Get(ctx, "profile-id", "review-id")

// Reply to a review (parentID is the review id)
reply, err := client.ProfileComments.Create(ctx, "profile-id", "review-id", "Thanks for visiting!")

// Delete your reply
_, err := client.ProfileComments.Delete(ctx, "profile-id", "reply-id")
```

### Profiles

```go
// List profiles
profiles, err := client.Profiles.List(ctx, nil)

// Get a profile
profile, err := client.Profiles.Get(ctx, "profile-id", nil)

// Get placements
placements, err := client.Profiles.Placements(ctx, "profile-id", nil)

// Move a placement (e.g. a Facebook Page or Telegram channel) to another group
placement, err := client.Profiles.AssignPlacementToGroup(ctx, "profile-id", "placement-external-id", "pg-other", nil)
fmt.Println(placement.ProfileGroupID) // "pg-other"

// Ice breakers (Instagram DMs): FAQ prompts shown when a user opens a chat
iceBreakers, err := client.Profiles.IceBreakers(ctx, "profile-id", nil)

_, err = client.Profiles.SetIceBreakers(ctx, "profile-id", []postproxy.IceBreaker{
	{Question: "What services do you offer?", Payload: "services"},
	{Question: "What are your hours?", Payload: "hours"},
}, nil) // 1-4 items

_, err = client.Profiles.DeleteIceBreakers(ctx, "profile-id", nil)

// Delete a profile
result, err := client.Profiles.Delete(ctx, "profile-id", nil)

// Profile stats timeseries — PlacementID required for facebook, linkedin, telegram
placementID := "108520199"
from := "2026-04-01T00:00:00Z"
stats, err := client.Profiles.GetProfileStats(ctx, "prof_li_001", &postproxy.ProfileStatsOptions{
	PlacementID: &placementID,
	From:        &from,
})
for _, r := range stats.Data.Records {
	fmt.Printf("%s: %v\n", r.RecordedAt, r.Stats["followerCount"])
}

// Bluesky — no placements
bsky, err := client.Profiles.GetProfileStats(ctx, "prof_bsky_001", nil)
fmt.Println(bsky.Data.Records[len(bsky.Data.Records)-1].Stats["followersCount"])

// Per-day series. Attribution "publish" counts post metrics toward the day
// the post was published instead of the day they were received.
attribution := postproxy.DailyStatsAttributionPublish
daily, err := client.Profiles.DailyStats(ctx, "profile-id", &postproxy.DailyStatsOptions{
	From:        &from,
	Attribution: &attribution,
})
for _, r := range daily.Data.Records {
	fmt.Println(r.Date, r.Metrics["impressions"])
}

// Daily follower gains, losses and totals
history, err := client.Profiles.FollowerHistory(ctx, "profile-id", &postproxy.FollowerHistoryOptions{From: &from})

// Move the profile (with its placements) to another profile group
moved, err := client.Profiles.AssignToGroup(ctx, "profile-id", "pg-other", nil)

// Instagram: whether a DM participant is verified and follows you
status, err := client.Profiles.ParticipantStatus(ctx, "profile-id", "igsid_8675309", nil)

// LinkedIn: company page URL -> URN to @-mention in a post
mention, err := client.Profiles.ResolveMention(ctx, "profile-id", "https://www.linkedin.com/company/acme", nil)
fmt.Println(mention.URN)

// Instagram: music and sounds a reel can attach via InstagramParams.AudioID
// (profiles connected via Facebook login only)
audioType := postproxy.InstagramAudioTypeMusic
query := "summer"
audio, err := client.Profiles.InstagramAvailableAudio(ctx, "profile-id", &postproxy.InstagramAvailableAudioOptions{
	AudioType:   &audioType,
	SearchQuery: &query,
})

// Facebook: read and update a Page's about / contact fields
fields := "about,website"
page, err := client.Profiles.FacebookPage(ctx, "profile-id", "page-id", &postproxy.FacebookPageOptions{Fields: &fields})
_, err = client.Profiles.UpdateFacebookPage(ctx, "profile-id", "page-id", map[string]any{"about": "Fresh coffee"}, nil)

// Instagram (Facebook login): incoming collaboration invites
invites, err := client.Profiles.InstagramCollaborationInvites(ctx, "ig-profile-id", nil) // NextCursor -> After
_, err = client.Profiles.AcceptInstagramCollaborationInvite(ctx, "ig-profile-id", invites.CollaborationInvites[0].MediaID, nil)
_, err = client.Profiles.DeclineInstagramCollaborationInvite(ctx, "ig-profile-id", "17890000000000001", nil)

// Pinterest boards: a board is a placement; its ID is the BoardID you publish to
boardDescription := "Light dinners for hot days"
privacy := postproxy.PinterestBoardPrivacyPublic // or PinterestBoardPrivacyProtected
board, err := client.Profiles.CreatePinterestBoard(ctx, "pin-profile-id", "Summer Recipes", &postproxy.PinterestBoardOptions{
	Description: &boardDescription,
	Privacy:     &privacy,
})
_, err = client.Profiles.DeletePinterestBoard(ctx, "pin-profile-id", board.ID, nil) // also deletes its Pins

// YouTube playlists: publish into one with YouTubeParams.PlaylistID
playlists, err := client.Profiles.YouTubePlaylists(ctx, "yt-profile-id", nil) // NextPageToken -> PageToken
public := postproxy.YouTubePrivacyPublic
playlist, err := client.Profiles.CreateYouTubePlaylist(ctx, "yt-profile-id", "Tutorials", &postproxy.YouTubePlaylistCreateOptions{
	PrivacyStatus: &public,
})
// Nil fields keep their value; a pointer to "" clears the description
empty := ""
_, err = client.Profiles.UpdateYouTubePlaylist(ctx, "yt-profile-id", playlist.ID, &postproxy.YouTubePlaylistUpdateOptions{Description: &empty})
_, err = client.Profiles.DeleteYouTubePlaylist(ctx, "yt-profile-id", playlist.ID, nil) // keeps the videos
```

`Profile` carries `Username`, `PlatformID`, `PlatformURL`, `AvatarURL`,
`PlacementsSyncStatus`, and `InactiveReason` (`ProfileInactiveReasonNetworkTaken` when the
same account was connected elsewhere). `Placement` carries `AvatarURL`.

Every stats record (post stats and profile stats alike) carries `RawStats` alongside the
normalized `Stats`, exposing each metric under its **original platform name**:

```go
stats, err := client.Posts.Stats(ctx, []string{"post-id"}, nil)
record := stats.Data["post-id"].Platforms[0].Records[0]

fmt.Println(record.Stats["impressions"])       // normalized
fmt.Println(record.RawStats["views"])          // Instagram's own name
fmt.Println(record.RawStats["impression_count"]) // Twitter/X's own name
```

LinkedIn post stats now normalize `likes`, `comments`, `shares`, and `clicks` alongside
`impressions` — previously only `impressions` was normalized.

#### Post syncs & backfill

Postproxy mirrors posts published natively on a platform into your account. Every one of
those pulls is recorded as a **post sync**: the one fired when the profile connects, the
recurring poll, and any backfill you start.

```go
// Start a backfill — walks the feed backwards from the newest post in batches
// of 25 until it reaches `from` or the platform stops returning posts.
sync, err := client.Profiles.BackfillPosts(ctx, "profile-id", "2025-01-01", nil)
fmt.Println(sync.ID, sync.Status) // "sync456def" "pending"

// Poll it to completion — finished when Status is PostSyncStatusCompleted or
// PostSyncStatusFailed.
run, err := client.Profiles.PostSync(ctx, "profile-id", sync.ID, nil)
fmt.Println(run.PostsImported, "of", run.PostsSeen, "back to", run.OldestPostedAt)

// List recent runs (kept for 30 days), newest first
trigger := postproxy.PostSyncTriggerBackfill
status := postproxy.PostSyncStatusCompleted
runs, err := client.Profiles.PostSyncs(ctx, "profile-id", &postproxy.PostSyncListOptions{
	Trigger: &trigger,  // connect | scheduled | backfill
	Status:  &status,   // pending | running | completed | failed
})
```

| `PostSync` field | Description |
|---|---|
| `ID` | Sync identifier |
| `ProfileID` | Profile this run belongs to |
| `PlacementID` | Placement the run pulled from, when there is one |
| `Kind` | Always `posts` today |
| `Trigger` | `connect`, `scheduled`, or `backfill` |
| `Status` | `pending`, `running`, `completed`, or `failed` |
| `StartedAt` / `CompletedAt` | ISO 8601 timestamps, `nil` until set |
| `PostsSeen` | Posts the platform returned across the run |
| `PostsImported` | Posts that were **new** and got created |
| `BackfillFrom` | The date floor requested; `nil` for `connect`/`scheduled` |
| `OldestPostedAt` | Publish date of the oldest post the run reached |
| `Error` | Platform error message when `Status` is failed |
| `CreatedAt` | ISO 8601 timestamp |

**How far back a backfill reaches depends on the platform's API**, not on Postproxy: where
history is pageable we follow it, otherwise the run ends early with whatever it got and
still reports `PostSyncStatusCompleted`.

Only one backfill runs per profile at a time — starting a second returns a 409 carrying the
running one's id:

```go
sync, err := client.Profiles.BackfillPosts(ctx, "profile-id", "2025-01-01", nil)
if postproxy.IsConflictError(err) {
	apiErr := err.(*postproxy.PostProxyError)
	runningID, _ := apiErr.Response["profile_sync_id"].(string)
	// Poll the run that's already going.
}
```

Posts you already have are skipped, so overlapping backfills are safe. Imported posts
behave exactly like ones the poll picks up (`source: "imported"`, `post.imported` webhook),
but a backfill's follow-up work is queued at a lower priority so a deep run can't slow down
publishing.

DM history works the same way: `BackfillChats` imports chats back to `from` and returns a
`ChatSync` (with `ChatsSeen` and `MessagesImported`), one run per profile at a time.

```go
chatSync, err := client.Profiles.BackfillChats(ctx, "profile-id", "2026-01-01T00:00:00Z", nil)
run, err := client.Profiles.ChatSync(ctx, "profile-id", chatSync.ID, nil)
runs, err := client.Profiles.ChatSyncs(ctx, "profile-id", &postproxy.ChatSyncListOptions{Status: &status})
```

### Profile Groups

```go
// List profile groups
groups, err := client.ProfileGroups.List(ctx)

// Get a profile group
group, err := client.ProfileGroups.Get(ctx, "group-id")

// Create a profile group
group, err := client.ProfileGroups.Create(ctx, "My Group")

// Delete a profile group
result, err := client.ProfileGroups.Delete(ctx, "group-id")

// Initialize an OAuth connection
conn, err := client.ProfileGroups.InitializeConnection(
	ctx, "group-id", postproxy.PlatformInstagram, "https://myapp.com/callback",
)
fmt.Println("Connect here:", conn.URL)

// Pick the onboarding flow: Instagram via Facebook login (needed for reel
// music), or WhatsApp via the Cloud API / the WhatsApp Business app.
// RefreshConnection re-authorizes an existing profile.
onboarding := postproxy.ConnectionOnboardingBusinessApp // ...Instagram, ...Facebook, ...API
conn, err = client.ProfileGroups.InitializeConnection(
	ctx, "group-id", postproxy.PlatformWhatsApp, "https://myapp.com/callback",
	&postproxy.ConnectionOptions{Onboarding: &onboarding},
)

// BlueSky — app password (synchronous, no OAuth)
bsky, err := client.ProfileGroups.ConnectBluesky(ctx, "group-id", postproxy.BlueskyConnectOptions{
	Identifier:  "yourname.bsky.social",
	AppPassword: "xxxx-xxxx-xxxx-xxxx",
})
fmt.Println(bsky.Profile.ID)

// Telegram — bring-your-own-bot. Channels populate asynchronously; poll
// placements until non-empty.
tg, err := client.ProfileGroups.ConnectTelegram(ctx, "group-id", postproxy.TelegramConnectOptions{
	BotToken: "123456789:ABCdef-GhIJklMnOpQrStUvWxYz",
})
if tg.NextStep != nil { fmt.Println(*tg.NextStep) }

var placements *postproxy.ListResponse[postproxy.Placement]
for {
	placements, err = client.Profiles.Placements(ctx, tg.Profile.ID, nil)
	if err != nil || len(placements.Data) > 0 { break }
	time.Sleep(3 * time.Second)
}
```

### Placements & API keys

Both require an account-wide API key.

```go
// Every placement (Facebook Page, LinkedIn organization, WhatsApp number, ...)
// across the account. One placement can be reachable through several
// connected profiles; ActiveProfile is the one PostProxy posts and reads through.
placements, err := client.Placements.List(ctx, nil)
for _, p := range placements.Data {
	fmt.Println(p.Name, p.Platform, p.ActiveProfile.Name, len(p.LinkedProfiles))
}
_, err = client.Placements.SetActiveProfile(ctx, "placement-id", "profile-id")

// Profile-group-scoped API keys
keys, err := client.APIKeys.List(ctx)
key, err := client.APIKeys.Create(ctx, "CI", "group-id")
fmt.Println(*key.Token) // shown only once — store it
_, err = client.APIKeys.Delete(ctx, key.ID)
```

### Google Business

Manage a Google Business Profile location's own data — hours, attributes, services, menus,
action links and photos — separate from posting to it. `locationID` is the placement id
from `Profiles.Placements`. Google's payloads are camelCased and passed through untouched
as `map[string]any`, so you can PATCH back exactly what you GET. Updates take a `fields`
list naming which top-level fields to change (Google's update mask).

```go
gbp := client.GoogleBusiness

readMask := "title,phoneNumbers,websiteUri"
location, err := gbp.Location(ctx, "profile-id", "location-id", &postproxy.GoogleBusinessLocationOptions{ReadMask: &readMask})
_, err = gbp.UpdateLocation(ctx, "profile-id", "location-id", []string{"websiteUri"}, map[string]any{
	"websiteUri": "https://example.com",
}, nil)
categories, err := gbp.AvailableLocationCategories(ctx, "profile-id", "location-id", "US", nil)

_, err = gbp.UpdateHours(ctx, "profile-id", "location-id", postproxy.GoogleBusinessHours{
	Fields: []string{"regularHours"},
	RegularHours: map[string]any{"periods": []any{map[string]any{
		"openDay": "MONDAY", "openTime": map[string]any{"hours": 9},
		"closeDay": "MONDAY", "closeTime": map[string]any{"hours": 17},
	}}},
}, nil)

// Attributes, services and food menus
attrs, err := gbp.Attributes(ctx, "profile-id", "location-id", nil)
available, err := gbp.AvailableAttributes(ctx, "profile-id", "location-id", nil)
_, err = gbp.UpdateAttributes(ctx, "profile-id", "location-id", []map[string]any{
	{"name": "attributes/has_wifi", "values": []any{true}},
}, &postproxy.GoogleBusinessAttributesOptions{AttributeMask: []string{"attributes/has_wifi"}})
services, err := gbp.ServiceList(ctx, "profile-id", "location-id", nil)
_, err = gbp.UpdateServiceList(ctx, "profile-id", "location-id", serviceItems, nil)
menus, err := gbp.FoodMenus(ctx, "profile-id", "location-id", nil)
_, err = gbp.UpdateFoodMenus(ctx, "profile-id", "location-id", menuList, nil)

// Place action links (booking, ordering, ...)
types, err := gbp.AvailablePlaceActionTypes(ctx, "profile-id", "location-id", nil)
links, err := gbp.PlaceActionLinks(ctx, "profile-id", "location-id", nil)
_, err = gbp.CreatePlaceActionLink(ctx, "profile-id", "location-id", postproxy.GoogleBusinessPlaceActionLink{
	URI: "https://book.example.com", PlaceActionType: "APPOINTMENT",
}, nil)
_, err = gbp.UpdatePlaceActionLink(ctx, "profile-id", "location-id", "locations/1/placeActionLinks/2",
	[]string{"uri"}, postproxy.GoogleBusinessPlaceActionLink{URI: "https://book2.example.com"}, nil)
_, err = gbp.DeletePlaceActionLink(ctx, "profile-id", "location-id", "locations/1/placeActionLinks/2", nil)

// Photos
category := "EXTERIOR"
_, err = gbp.CreateMedia(ctx, "profile-id", "location-id", "https://example.com/storefront.jpg",
	&postproxy.GoogleBusinessMediaOptions{Category: &category})
media, err := gbp.Media(ctx, "profile-id", "location-id", nil)
_, err = gbp.DeleteMedia(ctx, "profile-id", "location-id", "accounts/1/locations/1/media/9", nil)

// Search keywords that showed the location, month by month (max 12 months).
// Unlike the calls above, the response is typed (SearchKeywordsResponse).
from, to := "2026-07", "2026-09"
kw, err := gbp.SearchKeywords(ctx, "profile-id", "location-id", &postproxy.GoogleBusinessSearchKeywordsOptions{
	FromMonth: &from, ToMonth: &to,
})
for _, month := range kw.Months {
	for _, k := range month.Keywords {
		// Low-volume keywords have a nil Impressions and an ImpressionsThreshold bound
		fmt.Println(month.Month, k.Keyword, k.Impressions, k.ImpressionsThreshold)
	}
}
```

### WhatsApp templates

Templates live on the WhatsApp Business account (the profile) and are addressed by their
PostProxy id, Meta's id, or their name — pass `Language` with a name that exists in
several languages.

```go
template, err := client.WhatsAppTemplates.Create(ctx, "profile-id", postproxy.WhatsAppTemplateInput{
	Name:     "order_shipped",
	Category: "UTILITY", // MARKETING, UTILITY or AUTHENTICATION
	Language: "en_US",
	Components: []map[string]any{
		{"type": "BODY", "text": "Hi {{1}}, your order {{2}} has shipped."},
	},
}, nil)
fmt.Println(template.Status) // "PENDING" — follow via Get or message_template.status_updated

status, refresh := "APPROVED", true
approved, err := client.WhatsAppTemplates.List(ctx, "profile-id", &postproxy.WhatsAppTemplateListOptions{
	Status:  &status,
	Refresh: &refresh, // re-read from Meta first
})

lang := "en_US"
one, err := client.WhatsAppTemplates.Get(ctx, "profile-id", "order_shipped", &postproxy.WhatsAppTemplateOptions{Language: &lang})
_, err = client.WhatsAppTemplates.Update(ctx, "profile-id", one.ID, postproxy.WhatsAppTemplateUpdate{Components: components}, nil)
// Deleting by name without Language removes every language of the template.
_, err = client.WhatsAppTemplates.Delete(ctx, "profile-id", "order_shipped", nil)

// Meta's ready-made template library
library, err := client.WhatsAppTemplates.Library(ctx, "profile-id", nil)
```

### WhatsApp numbers

Number administration. `phoneNumberID` is the placement id from `Profiles.Placements`;
dataset and conversion calls are account-level and don't take one. Meta's responses come
back as `map[string]any`.

```go
wa := client.WhatsApp

// Registration
info, err := wa.NumberInfo(ctx, "profile-id", "phone-number-id", nil)
_, err = wa.RequestVerificationCode(ctx, "profile-id", "phone-number-id", nil)
_, err = wa.VerifyPhoneNumber(ctx, "profile-id", "phone-number-id", "123456", nil)
_, err = wa.RegisterPhoneNumber(ctx, "profile-id", "phone-number-id", "000000", nil) // two-step PIN

// Business profile, display name, username
about := "Fresh coffee, fast replies"
_, err = wa.UpdateBusinessProfile(ctx, "profile-id", "phone-number-id", postproxy.WhatsAppBusinessProfile{About: &about}, nil)
_, err = wa.UpdateBusinessProfilePhoto(ctx, "profile-id", "phone-number-id", postproxy.WhatsAppProfilePhoto{URL: "https://example.com/logo.png"}, nil)
_, err = wa.RequestDisplayNameChange(ctx, "profile-id", "phone-number-id", "Acme Coffee", nil) // placement.name_status_updated
_, err = wa.SetUsername(ctx, "profile-id", "phone-number-id", "acmecoffee", nil)
suggestions, err := wa.UsernameSuggestions(ctx, "profile-id", "phone-number-id", nil)

// Blocked users
_, err = wa.BlockUsers(ctx, "profile-id", "phone-number-id", []string{"15557654321"}, nil)
blocked, err := wa.BlockedUsers(ctx, "profile-id", "phone-number-id", nil)

// Groups
group, err := wa.CreateGroup(ctx, "profile-id", "phone-number-id", "VIP customers", nil)
groupID := fmt.Sprint(group["id"])
invite, err := wa.CreateGroupInviteLink(ctx, "profile-id", "phone-number-id", groupID, nil)
requests, err := wa.GroupJoinRequests(ctx, "profile-id", "phone-number-id", groupID, nil)
_, err = wa.ApproveGroupJoinRequests(ctx, "profile-id", "phone-number-id", groupID, []string{"15551234567"}, nil)

// Conversions for click-to-WhatsApp ads
value := 24.5
_, err = wa.SendConversionEvent(ctx, "profile-id", postproxy.WhatsAppConversionEvent{
	EventName: "Purchase", EventID: "order-1234", Phone: "15551234567", Currency: "USD", Value: &value,
}, nil)

// Download an inbound media file (e.g. an attachment's ExternalID) as raw bytes
file, err := wa.Media(ctx, "profile-id", "phone-number-id", "media-id", nil)
fmt.Println(file.ContentType, len(file.Data))
```

## Error Handling

```go
post, err := client.Posts.Get(ctx, "invalid-id", nil)
if err != nil {
	if postproxy.IsNotFoundError(err) {
		fmt.Println("Post not found")
	} else if postproxy.IsAuthenticationError(err) {
		fmt.Println("Invalid API key")
	} else if postproxy.IsRateLimitError(err) {
		apiErr := err.(*postproxy.PostProxyError)
		fmt.Printf("Rate limited — retry in %ds\n", apiErr.RetryAfter) // 0 when not given
	} else if postproxy.IsConflictError(err) {
		// Duplicate submission, a backfill already running, or an in-flight
		// Idempotency-Key. The details are in the parsed response body.
		apiErr := err.(*postproxy.PostProxyError)
		fmt.Println(apiErr.Response["duplicate_post_id"], apiErr.Response["profile_sync_id"])
	} else {
		fmt.Printf("Error: %v\n", err)
	}
}
```

| Status | Helper | Raised for |
|---|---|---|
| 400 | `IsBadRequestError` | Missing required parameters |
| 401 | `IsAuthenticationError` | Invalid, missing, or insufficient API key permissions |
| 403 | `IsForbiddenError` | Account canceled or past due, or the platform refused the action for lack of permission |
| 404 | `IsNotFoundError` | Resource does not exist or is not accessible |
| 409 | `IsConflictError` | Duplicate submission, a backfill already running, or an in-flight `Idempotency-Key` |
| 422 | `IsValidationError` | Validation failed |
| 429 | `IsRateLimitError` | Rate limited; `PostProxyError.RetryAfter` holds the `Retry-After` seconds |

## Examples

See the [examples](./examples) directory for complete working examples.

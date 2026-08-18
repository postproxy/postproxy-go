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

// Create a thread post
post, err := client.Posts.Create(ctx, "Thread starts here", []string{"profile-id"}, &postproxy.PostCreateOptions{
	Thread: []postproxy.ThreadChildInput{
		{Body: "Second post in the thread"},
		{Body: "Third with media", Media: []string{"https://example.com/img.jpg"}},
	},
})
for _, child := range post.Thread {
	fmt.Printf("%s: %s\n", child.ID, child.Body)
}

// Update a post (only drafts or scheduled posts)
newBody := "Updated content!"
post, err := client.Posts.Update(ctx, "post-id", &postproxy.PostUpdateOptions{
	Body: &newBody,
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

`post.processed`, `post.imported`, `platform_post.published`, `platform_post.failed`, `platform_post.failed_waiting_for_retry`, `platform_post.insights`, `profile.connected`, `profile.disconnected`, `profile.stats`, `media.failed`, `comment.created`, `profile_comment.created`, `message.received`, `message.sent`, `message.delivered`, `message.read`, `message.edited`, `message.deleted`, `message.failed_waiting_for_retry`, `message.failed`, `reaction.received`.

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
}
```

The typed payloads `MessageEventData` (the 8 `message.*` events), `ReactionEventData` (`reaction.received`), and `ProfileCommentCreatedData` (`profile_comment.created`) each embed the relevant resource. `MessageEventData` and `ReactionEventData` carry a full `Message` in their `Message` field.

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

Manage direct-message conversations (chats) and the messages within them across Facebook, Instagram, Telegram, and Bluesky.

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
```

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

## Error Handling

```go
post, err := client.Posts.Get(ctx, "invalid-id", nil)
if err != nil {
	if postproxy.IsNotFoundError(err) {
		fmt.Println("Post not found")
	} else if postproxy.IsAuthenticationError(err) {
		fmt.Println("Invalid API key")
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
| 404 | `IsNotFoundError` | Resource does not exist or is not accessible |
| 409 | `IsConflictError` | Duplicate submission, a backfill already running, or an in-flight `Idempotency-Key` |
| 422 | `IsValidationError` | Validation failed |
| 429 | — | Posting rate limit reached; check `PostProxyError.StatusCode` |

## Examples

See the [examples](./examples) directory for complete working examples.

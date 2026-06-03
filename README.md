# PostProxy Go SDK

The official Go SDK for the [PostProxy](https://postproxy.dev) API.

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
post, err := client.Posts.Create(ctx, "Caption", []string{"profile-id"}, &postproxy.PostCreateOptions{
	Platforms: &postproxy.PlatformParams{
		Instagram: &postproxy.InstagramParams{
			Format: &igFormat,
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
```

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
	} else {
		fmt.Printf("Error: %v\n", err)
	}
}
```

## Examples

See the [examples](./examples) directory for complete working examples.

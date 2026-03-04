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
post, err := client.Posts.Create(ctx, "Caption", []string{"profile-id"}, &postproxy.PostCreateOptions{
	Platforms: &postproxy.PlatformParams{
		Instagram: &postproxy.InstagramParams{
			Format: &igFormat,
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

// Publish a draft
post, err := client.Posts.PublishDraft(ctx, "post-id", nil)

// Delete a post
result, err := client.Posts.Delete(ctx, "post-id", nil)

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

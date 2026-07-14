package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/postproxy/postproxy-go"
)

func main() {
	apiKey := os.Getenv("POSTPROXY_API_KEY")
	if apiKey == "" {
		log.Fatal("POSTPROXY_API_KEY environment variable is required")
	}

	profileGroupID := os.Getenv("POSTPROXY_PROFILE_GROUP_ID")
	if profileGroupID == "" {
		log.Fatal("POSTPROXY_PROFILE_GROUP_ID environment variable is required")
	}

	client := postproxy.NewClient(apiKey, postproxy.WithProfileGroupID(profileGroupID))
	ctx := context.Background()

	// Fetch profiles from the profile group.
	profiles, err := client.Profiles.List(ctx, nil)
	if err != nil {
		log.Fatalf("failed to list profiles: %v", err)
	}
	if len(profiles.Data) == 0 {
		log.Fatal("no profiles found — connect a profile first")
	}

	fmt.Printf("Found %d profile(s):\n", len(profiles.Data))
	var igProfileID string
	for _, p := range profiles.Data {
		fmt.Printf("  - %s (%s, %s)\n", p.Name, p.Platform, p.ID)
		if p.Platform == postproxy.PlatformInstagram {
			igProfileID = p.ID
		}
	}
	if igProfileID == "" {
		log.Fatal("no Instagram profile found — connect one first")
	}

	// Create a post with platform-specific parameters.
	draft := true
	igFormat := postproxy.InstagramFormatPost
	post, err := client.Posts.Create(ctx, "Check out this reel!", []string{igProfileID}, &postproxy.PostCreateOptions{
		Media: []string{"https://example.com/image.jpg"},
		Platforms: &postproxy.PlatformParams{
			Instagram: &postproxy.InstagramParams{
				Format: &igFormat,
			},
		},
		Draft: &draft,
	})
	if err != nil {
		log.Fatalf("failed to create post: %v", err)
	}
	fmt.Printf("Created post with platform params: %s\n", post.ID)

	// Create a draft post with file upload.
	post, err = client.Posts.Create(ctx, "Draft with local file", []string{igProfileID}, &postproxy.PostCreateOptions{
		MediaFiles: []string{"./photo.jpg"},
		Draft:      &draft,
	})
	if err != nil {
		log.Fatalf("failed to create draft: %v", err)
	}
	fmt.Printf("Created draft: %s\n", post.ID)

	// Create a thread post.
	post, err = client.Posts.Create(ctx, "Here's a thread about PostProxy 🧵", []string{igProfileID}, &postproxy.PostCreateOptions{
		Thread: []postproxy.ThreadChildInput{
			{Body: "First, connect your social accounts."},
			{Body: "Then, create posts with media!", Media: []string{"https://example.com/demo.jpg"}},
			{Body: "Finally, schedule or publish instantly."},
		},
	})
	if err != nil {
		log.Fatalf("failed to create thread: %v", err)
	}
	fmt.Printf("Created thread: %s (%d children)\n", post.ID, len(post.Thread))

	// Twitter poll: 2-4 options (max 25 chars each), 5-10080 minutes
	pollFormat := postproxy.TwitterFormatPoll
	pollDuration := 1440
	post, err = client.Posts.Create(ctx, "Which framework?", []string{"twitter"}, &postproxy.PostCreateOptions{
		Platforms: &postproxy.PlatformParams{
			Twitter: &postproxy.TwitterParams{
				Format:              &pollFormat,
				PollOptions:         []string{"Rails", "Django", "Laravel", "Other"},
				PollDurationMinutes: &pollDuration,
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to create poll: %v", err)
	}
	fmt.Printf("Created poll post: %s\n", post.ID)
}

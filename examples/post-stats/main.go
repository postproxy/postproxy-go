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

	postIDs := os.Args[1:]
	if len(postIDs) == 0 {
		log.Fatal("usage: post-stats <post-id> [post-id...]")
	}

	client := postproxy.NewClient(apiKey)
	ctx := context.Background()

	// Fetch stats for all provided post IDs.
	from := "2026-02-01T00:00:00Z"
	stats, err := client.Posts.Stats(ctx, postIDs, &postproxy.PostStatsOptions{
		From: &from,
	})
	if err != nil {
		log.Fatalf("failed to get stats: %v", err)
	}

	for postID, postStats := range stats.Data {
		fmt.Printf("Post %s:\n", postID)
		for _, plat := range postStats.Platforms {
			fmt.Printf("  %s (profile %s):\n", plat.Platform, plat.ProfileID)
			for _, record := range plat.Records {
				fmt.Printf("    [%s]", record.RecordedAt)
				for key, val := range record.Stats {
					fmt.Printf(" %s=%v", key, val)
				}
				fmt.Println()
			}
		}
	}
}

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

	client := postproxy.NewClient(apiKey)
	ctx := context.Background()

	// List profile groups, use the first one or create a new one.
	groups, err := client.ProfileGroups.List(ctx)
	if err != nil {
		log.Fatalf("failed to list profile groups: %v", err)
	}

	var group postproxy.ProfileGroup
	if len(groups.Data) > 0 {
		group = groups.Data[0]
		fmt.Printf("Using existing profile group: %s (%s)\n", group.Name, group.ID)
	} else {
		created, err := client.ProfileGroups.Create(ctx, "My App")
		if err != nil {
			log.Fatalf("failed to create profile group: %v", err)
		}
		group = *created
		fmt.Printf("Created profile group: %s (%s)\n", group.Name, group.ID)
	}

	// Initialize a connection for Instagram.
	conn, err := client.ProfileGroups.InitializeConnection(
		ctx,
		group.ID,
		postproxy.PlatformInstagram,
		"https://myapp.com/callback",
	)
	if err != nil {
		log.Fatalf("failed to initialize connection: %v", err)
	}
	fmt.Printf("Connect your Instagram account: %s\n", conn.URL)
}

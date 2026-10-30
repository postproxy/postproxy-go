package main

import (
	"context"
	"fmt"
	"log"
	"os"

	postproxy "github.com/postproxy/postproxy-go"
)

func main() {
	apiKey := os.Getenv("POSTPROXY_API_KEY")
	profileGroupID := os.Getenv("POSTPROXY_PROFILE_GROUP_ID")

	client := postproxy.NewClient(apiKey, postproxy.WithProfileGroupID(profileGroupID))
	ctx := context.Background()

	profileID := "your-youtube-profile-id"

	// List the channel's playlists, page by page
	var pageToken *string
	for {
		page, err := client.Profiles.YouTubePlaylists(ctx, profileID, &postproxy.YouTubePlaylistsOptions{PageToken: pageToken})
		if err != nil {
			log.Fatalf("Error listing playlists: %v", err)
		}
		for _, p := range page.Data {
			fmt.Printf("%s: %s (%s, %d videos)\n", p.ID, p.Title, p.PrivacyStatus, p.ItemCount)
		}
		if page.NextPageToken == nil {
			break
		}
		pageToken = page.NextPageToken
	}

	// Create a playlist and publish a video into it
	description := "Step-by-step guides"
	public := postproxy.YouTubePrivacyPublic
	playlist, err := client.Profiles.CreateYouTubePlaylist(ctx, profileID, "Tutorials", &postproxy.YouTubePlaylistCreateOptions{
		Description:   &description,
		PrivacyStatus: &public,
	})
	if err != nil {
		log.Fatalf("Error creating playlist: %v", err)
	}

	title := "How to set up your account"
	post, err := client.Posts.Create(ctx, title, []string{profileID}, &postproxy.PostCreateOptions{
		Media: []string{"https://example.com/tutorial.mp4"},
		Platforms: &postproxy.PlatformParams{
			YouTube: &postproxy.YouTubeParams{Title: &title, PlaylistID: &playlist.ID},
		},
	})
	if err != nil {
		log.Fatalf("Error creating post: %v", err)
	}
	fmt.Printf("Post: %s\n", post.ID)

	// Rename it later; nil fields keep their value
	newTitle := "Getting started"
	if _, err := client.Profiles.UpdateYouTubePlaylist(ctx, profileID, playlist.ID, &postproxy.YouTubePlaylistUpdateOptions{
		Title: &newTitle,
	}); err != nil {
		log.Fatalf("Error updating playlist: %v", err)
	}

	// Deleting a playlist keeps its videos
	// client.Profiles.DeleteYouTubePlaylist(ctx, profileID, playlist.ID, nil)
}

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	postproxy "github.com/postproxy/postproxy-go"
)

func main() {
	apiKey := os.Getenv("POSTPROXY_API_KEY")
	profileGroupID := os.Getenv("POSTPROXY_PROFILE_GROUP_ID")

	client := postproxy.NewClient(apiKey, postproxy.WithProfileGroupID(profileGroupID))
	ctx := context.Background()

	// Only the body, the YouTube title and the Pinterest title / destination
	// link can change once a post is published.
	body := "Launch day! (now with the right link)"
	title := "Launch day — full walkthrough"
	link := "https://example.com/launch"
	post, err := client.Posts.Update(ctx, "published-post-id", &postproxy.PostUpdateOptions{
		Body: &body,
		Platforms: &postproxy.PlatformParams{
			YouTube:   &postproxy.YouTubeParams{Title: &title},
			Pinterest: &postproxy.PinterestParams{DestinationLink: &link},
		},
	})
	if err != nil {
		log.Fatalf("Error editing post: %v", err)
	}
	for _, p := range post.Platforms {
		permalink := ""
		if p.Permalink != nil {
			permalink = *p.Permalink
		}
		fmt.Printf("%s: %s\n", p.Platform, permalink)
	}
}

// handleWebhook learns how the edit went: it is pushed to each platform in
// the background.
func handleWebhook(w http.ResponseWriter, r *http.Request) {
	raw, _ := io.ReadAll(r.Body)
	event, err := postproxy.ParseWebhookEvent(raw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch event.Type {
	case postproxy.EventPlatformPostEdited:
		data, _ := event.AsPlatformPost()
		fmt.Printf("Edited on %s\n", data.Platform)
	case postproxy.EventPlatformPostEditFailed:
		data, _ := event.AsPlatformPost()
		fmt.Printf("Edit failed on %s: %v\n", data.Platform, data.Error)
	}
	w.WriteHeader(http.StatusOK)
}

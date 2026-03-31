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

	postID := "your-post-id"
	profileID := "your-profile-id"

	// List comments on a post
	comments, err := client.Comments.List(ctx, postID, profileID, nil)
	if err != nil {
		log.Fatalf("Error listing comments: %v", err)
	}
	fmt.Printf("Total comments: %d\n", comments.Total)
	for _, c := range comments.Data {
		fmt.Printf("  %v: %s\n", c.AuthorUsername, c.Body)
		for _, r := range c.Replies {
			fmt.Printf("    %v: %s\n", r.AuthorUsername, r.Body)
		}
	}

	// Create a comment
	newComment, err := client.Comments.Create(ctx, postID, profileID, "Thanks for the feedback!", nil)
	if err != nil {
		log.Fatalf("Error creating comment: %v", err)
	}
	fmt.Printf("Created comment: %s (status: %s)\n", newComment.ID, newComment.Status)

	// Reply to a comment
	parentID := newComment.ID
	reply, err := client.Comments.Create(ctx, postID, profileID, "Glad you liked it!", &postproxy.CommentCreateOptions{ParentID: &parentID})
	if err != nil {
		log.Fatalf("Error creating reply: %v", err)
	}
	fmt.Printf("Created reply: %s\n", reply.ID)

	// Hide / unhide
	_, _ = client.Comments.Hide(ctx, postID, newComment.ID, profileID)
	fmt.Println("Comment hidden")

	_, _ = client.Comments.Unhide(ctx, postID, newComment.ID, profileID)
	fmt.Println("Comment unhidden")

	// Like / unlike
	_, _ = client.Comments.Like(ctx, postID, newComment.ID, profileID)
	fmt.Println("Comment liked")

	_, _ = client.Comments.Unlike(ctx, postID, newComment.ID, profileID)
	fmt.Println("Comment unliked")

	// Delete
	_, _ = client.Comments.Delete(ctx, postID, newComment.ID, profileID)
	fmt.Println("Comment deleted")
}

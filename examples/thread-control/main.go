package main

import (
	"context"
	"fmt"
	"log"
	"os"

	postproxy "github.com/postproxy/postproxy-go"
)

// Meta's own inbox apps
const (
	pageInboxAppID      = "263902037430900"
	instagramInboxAppID = "1217981644879628"
)

func main() {
	apiKey := os.Getenv("POSTPROXY_API_KEY")
	profileGroupID := os.Getenv("POSTPROXY_PROFILE_GROUP_ID")

	client := postproxy.NewClient(apiKey, postproxy.WithProfileGroupID(profileGroupID))
	ctx := context.Background()

	// A Facebook or Instagram chat
	chatID := "your-chat-id"

	// Who controls the conversation right now? Both fields are nil when idle.
	owner, err := client.Chats.ThreadOwner(ctx, chatID, nil)
	if err != nil {
		log.Fatalf("Error reading thread owner: %v", err)
	}
	if owner.AppID != nil {
		fmt.Printf("Owner app: %s\n", *owner.AppID)
	} else {
		fmt.Println("Conversation is idle")
	}

	// Escalate to a human in Meta's inbox (use instagramInboxAppID for Instagram)
	metadata := "escalated to support"
	chat, err := client.Chats.PassThreadControl(ctx, chatID, pageInboxAppID, &postproxy.ThreadControlOptions{
		Metadata: &metadata,
	})
	if err != nil {
		log.Fatalf("Error passing thread control: %v", err)
	}
	fmt.Printf("Thread control: %v\n", chat.Metadata["thread_control"])
	_ = instagramInboxAppID

	// Later: take it back (needs Conversation Routing on the Page) ...
	if _, err := client.Chats.TakeThreadControl(ctx, chatID, nil); err != nil {
		log.Fatalf("Error taking thread control: %v", err)
	}
	// ... or, without Conversation Routing, ask the owner to pass it back and
	// wait for the chat.thread_control_changed webhook.
	// client.Chats.RequestThreadControl(ctx, chatID, nil)

	// Done: let the conversation go idle
	if _, err := client.Chats.ReleaseThreadControl(ctx, chatID, nil); err != nil {
		log.Fatalf("Error releasing thread control: %v", err)
	}
}

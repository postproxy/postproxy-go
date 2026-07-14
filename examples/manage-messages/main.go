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

	profileID := "your-profile-id" // a DM-capable profile (facebook/instagram/telegram/bluesky)

	// List existing chats for a profile
	perPage := 20
	chats, err := client.Chats.List(ctx, profileID, &postproxy.ChatListOptions{PerPage: &perPage})
	if err != nil {
		log.Fatalf("Error listing chats: %v", err)
	}
	fmt.Printf("Total chats: %d\n", chats.Total)
	for _, chat := range chats.Data {
		name := chat.ParticipantExternalID
		if chat.ParticipantUsername != nil {
			name = *chat.ParticipantUsername
		}
		fmt.Printf("  %s (platform: %s)\n", name, chat.Platform)
	}

	// Find or create a chat with a participant
	username := "jane_doe"
	chat, err := client.Chats.Create(ctx, profileID, "igsid_8675309", &postproxy.ChatCreateOptions{
		ParticipantUsername: &username,
	})
	if err != nil {
		log.Fatalf("Error creating chat: %v", err)
	}
	fmt.Printf("Chat: %s (platform: %s)\n", chat.ID, chat.Platform)

	// List inbound messages in the chat
	dir := postproxy.MessageDirectionInbound
	messages, err := client.Messages.List(ctx, chat.ID, &postproxy.MessageListOptions{Direction: &dir})
	if err != nil {
		log.Fatalf("Error listing messages: %v", err)
	}
	for _, msg := range messages.Data {
		body := ""
		if msg.Body != nil {
			body = *msg.Body
		}
		fmt.Printf("  [%s] %s\n", msg.Direction, body)
		for _, att := range msg.Attachments {
			url := ""
			if att.URL != nil {
				url = *att.URL
			}
			fmt.Printf("    attachment: %s -> %s\n", att.Type, url)
		}
	}

	// Send a text message (within the 24h window)
	text := "Yes, we ship worldwide!"
	sent, err := client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{Body: &text})
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}
	fmt.Printf("Sent message: %s (status: %s)\n", sent.ID, sent.Status)

	// Send outside the 24h window with a tag (Facebook/Instagram)
	followUp := "Following up on your order."
	tag := "HUMAN_AGENT"
	_, _ = client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{Body: &followUp, Tag: &tag})

	// Send an image by hosted URL
	_, _ = client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{
		Media: []string{"https://cdn.example.com/photo.png"},
	})

	// Send an image from a local file (multipart)
	// _, _ = client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{MediaFiles: []string{"./photo.png"}})

	// React / unreact (Facebook & Instagram)
	reaction := "love"
	emoji := "❤️"
	_, _ = client.Messages.React(ctx, sent.ID, &postproxy.MessageReactOptions{Reaction: &reaction, Emoji: &emoji})
	_, _ = client.Messages.Unreact(ctx, sent.ID, nil)

	// Edit an outbound message (Telegram only)
	// updated := "Updated answer."
	// _, _ = client.Messages.Edit(ctx, sent.ID, &postproxy.MessageEditOptions{Body: &updated})

	// Archive / unarchive a chat (Bluesky only)
	// _, _ = client.Chats.Archive(ctx, chat.ID, nil)
	// _, _ = client.Chats.Unarchive(ctx, chat.ID, nil)

	// Private reply to a comment (Instagram/Facebook) — returns a Message
	reply, err := client.Comments.PrivateReply(ctx, "your-post-id", "comment-id", profileID, "Thanks — DM-ing you the details.")
	if err != nil {
		log.Fatalf("Error sending private reply: %v", err)
	}
	fmt.Printf("Private reply queued: %s (chat: %s)\n", reply.ID, reply.ChatID)

	// Ice breakers (Instagram only): FAQ prompts shown when a user opens a chat
	_, err = client.Profiles.SetIceBreakers(ctx, profileID, []postproxy.IceBreaker{
		{Question: "What services do you offer?", Payload: "services"},
		{Question: "What are your hours?", Payload: "hours"},
	}, nil)
	if err != nil {
		log.Fatalf("Error setting ice breakers: %v", err)
	}
	iceBreakers, err := client.Profiles.IceBreakers(ctx, profileID, nil)
	if err != nil {
		log.Fatalf("Error listing ice breakers: %v", err)
	}
	for _, ib := range iceBreakers.IceBreakers {
		fmt.Printf("Ice breaker: %s\n", ib.Question)
	}
	// _, _ = client.Profiles.DeleteIceBreakers(ctx, profileID, nil)
}

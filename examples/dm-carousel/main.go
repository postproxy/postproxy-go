package main

import (
	"context"
	"log"
	"os"

	postproxy "github.com/postproxy/postproxy-go"
)

func main() {
	apiKey := os.Getenv("POSTPROXY_API_KEY")
	profileGroupID := os.Getenv("POSTPROXY_PROFILE_GROUP_ID")

	client := postproxy.NewClient(apiKey, postproxy.WithProfileGroupID(profileGroupID))
	ctx := context.Background()

	profileID := "your-instagram-profile-id"

	perPage := 1
	chats, err := client.Chats.List(ctx, profileID, &postproxy.ChatListOptions{PerPage: &perPage})
	if err != nil || len(chats.Data) == 0 {
		log.Fatalf("Error listing chats: %v", err)
	}
	chat := chats.Data[0]

	// Show the participant we're on it, then send a product carousel
	_, _ = client.Chats.MarkRead(ctx, chat.ID, nil)
	_, _ = client.Chats.StartTyping(ctx, chat.ID, nil)

	_, err = client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{
		Elements: []postproxy.MessageElement{
			{
				Title:    "Red sneaker",
				Subtitle: "$89 — sizes 38-46",
				ImageURL: "https://cdn.example.com/red.png",
				Buttons:  []postproxy.MessageButton{{Type: "web_url", Title: "Buy", URL: "https://shop.example.com/red"}},
			},
			{
				Title:    "Blue sneaker",
				Subtitle: "$92 — back in stock soon",
				ImageURL: "https://cdn.example.com/blue.png",
				Buttons:  []postproxy.MessageButton{{Type: "postback", Title: "Notify me", Payload: "NOTIFY:blue"}},
			},
		},
	})
	if err != nil {
		log.Fatalf("Error sending carousel: %v", err)
	}

	_, _ = client.Chats.StopTyping(ctx, chat.ID, nil)
}

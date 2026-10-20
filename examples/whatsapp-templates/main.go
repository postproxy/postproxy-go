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

	profileID := "your-whatsapp-profile-id"

	// Create a template — Meta reviews it before it can be sent
	template, err := client.WhatsAppTemplates.Create(ctx, profileID, postproxy.WhatsAppTemplateInput{
		Name:     "order_shipped",
		Category: "UTILITY",
		Language: "en_US",
		Components: []map[string]any{
			{"type": "BODY", "text": "Hi {{1}}, your order {{2}} has shipped."},
		},
	}, nil)
	if err != nil {
		log.Fatalf("Error creating template: %v", err)
	}
	fmt.Printf("Template %s: %s\n", template.Name, template.Status) // "order_shipped" "PENDING"

	// Check approval status (or subscribe to message_template.status_updated)
	status := "APPROVED"
	refresh := true
	approved, err := client.WhatsAppTemplates.List(ctx, profileID, &postproxy.WhatsAppTemplateListOptions{
		Status:  &status,
		Refresh: &refresh,
	})
	if err != nil {
		log.Fatalf("Error listing templates: %v", err)
	}
	for _, t := range approved.Data {
		fmt.Printf("Approved: %s\n", t.Name)
	}

	// Send it to a customer. The phone number id is a placement of the profile.
	numbers, err := client.Profiles.Placements(ctx, profileID, nil)
	if err != nil || len(numbers.Data) == 0 {
		log.Fatalf("Error listing phone numbers: %v", err)
	}
	chat, err := client.Chats.Create(ctx, profileID, "15551234567", &postproxy.ChatCreateOptions{
		PlacementID: &numbers.Data[0].ID,
	})
	if err != nil {
		log.Fatalf("Error creating chat: %v", err)
	}

	_, err = client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{
		Template: &postproxy.MessageTemplateInput{
			Name:      "order_shipped",
			Language:  "en_US",
			Variables: []string{"Jane", "#1234"},
		},
	})
	if err != nil {
		log.Fatalf("Error sending template: %v", err)
	}

	// Inside the 24h window, free-form messages work too
	fresh, err := client.Chats.Get(ctx, chat.ID, nil)
	if err != nil {
		log.Fatalf("Error fetching chat: %v", err)
	}
	if fresh.WithinMessagingWindow != nil && *fresh.WithinMessagingWindow {
		_, _ = client.Messages.Send(ctx, chat.ID, &postproxy.MessageSendOptions{
			Location: &postproxy.MessageLocation{Latitude: 52.52, Longitude: 13.405, Name: "Pickup point"},
		})
	}
}

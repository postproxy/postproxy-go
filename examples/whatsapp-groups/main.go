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
	wa := client.WhatsApp

	profileID := "your-whatsapp-profile-id"
	phoneNumberID := "your-phone-number-id"

	// Business profile
	about := "Fresh coffee, fast replies"
	_, err := wa.UpdateBusinessProfile(ctx, profileID, phoneNumberID, postproxy.WhatsAppBusinessProfile{
		About:    &about,
		Websites: []string{"https://example.com"},
	}, nil)
	if err != nil {
		log.Fatalf("Error updating business profile: %v", err)
	}

	// Create a group and share its invite link
	description := "Early access to new roasts"
	approval := "approval_required"
	group, err := wa.CreateGroup(ctx, profileID, phoneNumberID, "VIP customers", &postproxy.WhatsAppGroupCreateOptions{
		Description:      &description,
		JoinApprovalMode: &approval,
	})
	if err != nil {
		log.Fatalf("Error creating group: %v", err)
	}
	groupID := fmt.Sprint(group["id"])

	invite, err := wa.CreateGroupInviteLink(ctx, profileID, phoneNumberID, groupID, nil)
	if err != nil {
		log.Fatalf("Error creating invite link: %v", err)
	}
	fmt.Printf("Invite link: %v\n", invite)

	// Approve everyone waiting to join
	requests, err := wa.GroupJoinRequests(ctx, profileID, phoneNumberID, groupID, nil)
	if err != nil {
		log.Fatalf("Error listing join requests: %v", err)
	}
	fmt.Printf("Join requests: %v\n", requests)
	_, _ = wa.ApproveGroupJoinRequests(ctx, profileID, phoneNumberID, groupID, []string{"15551234567"}, nil)

	// Block a spammer
	_, _ = wa.BlockUsers(ctx, profileID, phoneNumberID, []string{"15557654321"}, nil)

	// Report a purchase back to Meta for a click-to-WhatsApp ad
	value := 24.5
	_, err = wa.SendConversionEvent(ctx, profileID, postproxy.WhatsAppConversionEvent{
		EventName: "Purchase",
		EventID:   "order-1234",
		Phone:     "15551234567",
		Currency:  "USD",
		Value:     &value,
	}, nil)
	if err != nil {
		log.Fatalf("Error sending conversion event: %v", err)
	}
}

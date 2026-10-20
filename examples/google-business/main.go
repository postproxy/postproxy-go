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
	gbp := client.GoogleBusiness

	profileID := "your-google-business-profile-id"

	// locationID is the placement id of the location
	placements, err := client.Profiles.Placements(ctx, profileID, nil)
	if err != nil {
		log.Fatalf("Error listing placements: %v", err)
	}
	if len(placements.Data) == 0 {
		log.Fatal("no locations found")
	}
	locationID := placements.Data[0].ID

	// Read and update the location
	readMask := "title,phoneNumbers,websiteUri,regularHours"
	location, err := gbp.Location(ctx, profileID, locationID, &postproxy.GoogleBusinessLocationOptions{ReadMask: &readMask})
	if err != nil {
		log.Fatalf("Error reading location: %v", err)
	}
	fmt.Printf("Location: %v\n", location["title"])

	_, err = gbp.UpdateLocation(ctx, profileID, locationID, []string{"websiteUri"}, map[string]any{
		"websiteUri": "https://example.com",
	}, nil)
	if err != nil {
		log.Fatalf("Error updating location: %v", err)
	}

	// Opening hours: Mon-Fri 9-17
	var periods []any
	for _, day := range []string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY"} {
		periods = append(periods, map[string]any{
			"openDay":   day,
			"openTime":  map[string]any{"hours": 9},
			"closeDay":  day,
			"closeTime": map[string]any{"hours": 17},
		})
	}
	_, err = gbp.UpdateHours(ctx, profileID, locationID, postproxy.GoogleBusinessHours{
		Fields:       []string{"regularHours"},
		RegularHours: map[string]any{"periods": periods},
	}, nil)
	if err != nil {
		log.Fatalf("Error updating hours: %v", err)
	}

	// Booking link
	types, err := gbp.AvailablePlaceActionTypes(ctx, profileID, locationID, nil)
	if err != nil {
		log.Fatalf("Error listing action types: %v", err)
	}
	fmt.Printf("Available action types: %v\n", types)

	preferred := true
	_, err = gbp.CreatePlaceActionLink(ctx, profileID, locationID, postproxy.GoogleBusinessPlaceActionLink{
		URI:             "https://book.example.com",
		PlaceActionType: "APPOINTMENT",
		IsPreferred:     &preferred,
	}, nil)
	if err != nil {
		log.Fatalf("Error creating action link: %v", err)
	}

	// Add a storefront photo
	category, format := "EXTERIOR", "PHOTO"
	_, err = gbp.CreateMedia(ctx, profileID, locationID, "https://example.com/storefront.jpg", &postproxy.GoogleBusinessMediaOptions{
		Category:    &category,
		MediaFormat: &format,
	})
	if err != nil {
		log.Fatalf("Error adding photo: %v", err)
	}

	media, err := gbp.Media(ctx, profileID, locationID, nil)
	if err != nil {
		log.Fatalf("Error listing media: %v", err)
	}
	fmt.Printf("Media: %v\n", media)
}

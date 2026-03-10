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
	profileGroupID := os.Getenv("POSTPROXY_PROFILE_GROUP_ID")

	client := postproxy.NewClient(apiKey, postproxy.WithProfileGroupID(profileGroupID))
	ctx := context.Background()

	// Create a queue with weekday morning timeslots
	tz := "America/New_York"
	desc := "Weekday morning content"
	jitter := 10
	queue, err := client.Queues.Create(ctx, "Morning Posts", profileGroupID, &postproxy.QueueCreateOptions{
		Description: &desc,
		Timezone:    &tz,
		Jitter:      &jitter,
		Timeslots: []postproxy.TimeslotInput{
			{Day: 1, Time: "09:00"},
			{Day: 2, Time: "09:00"},
			{Day: 3, Time: "09:00"},
			{Day: 4, Time: "09:00"},
			{Day: 5, Time: "09:00"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created queue: %s %s\n", queue.ID, queue.Name)
	fmt.Printf("Timeslots: %d\n", len(queue.Timeslots))

	// List all queues
	queues, err := client.Queues.List(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, q := range queues.Data {
		fmt.Printf("Queue: %s\n", q.Name)
	}

	// Get next available slot
	nextSlot, err := client.Queues.NextSlot(ctx, queue.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Next slot: %s\n", nextSlot.NextSlot)

	// Add a post to the queue
	profiles, err := client.Profiles.List(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	queueID := queue.ID
	priority := "high"
	post, err := client.Posts.Create(ctx, "This post will be scheduled by the queue", []string{profiles.Data[0].ID}, &postproxy.PostCreateOptions{
		QueueID:       &queueID,
		QueuePriority: &priority,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Queued post: %s scheduled at: %v\n", post.ID, post.ScheduledAt)

	// Update the queue
	newJitter := 15
	_, err = client.Queues.Update(ctx, queue.ID, &postproxy.QueueUpdateOptions{
		Jitter: &newJitter,
		Timeslots: []postproxy.TimeslotInput{
			{Day: 6, Time: "10:00"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Pause the queue
	enabled := false
	_, err = client.Queues.Update(ctx, queue.ID, &postproxy.QueueUpdateOptions{
		Enabled: &enabled,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Queue paused")

	// Delete the queue
	deleted, err := client.Queues.Delete(ctx, queue.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted: %v\n", deleted.Deleted)
}

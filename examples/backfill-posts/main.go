// Backfill a profile's older posts and follow the sync run to completion.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	postproxy "github.com/postproxy/postproxy-go"
)

func main() {
	apiKey := os.Getenv("POSTPROXY_API_KEY")
	profileGroupID := os.Getenv("POSTPROXY_PROFILE_GROUP_ID")

	client := postproxy.NewClient(apiKey, postproxy.WithProfileGroupID(profileGroupID))
	ctx := context.Background()

	profileID := "your-profile-id"

	// Start a backfill. It walks the profile's feed backwards from the newest
	// post in batches of 25 and stops at `from` — or earlier, if the platform
	// stops returning history. Runs in the background.
	sync, err := client.Profiles.BackfillPosts(ctx, profileID, "2025-01-01", nil)
	if err != nil {
		if !postproxy.IsConflictError(err) {
			log.Fatalf("Error starting backfill: %v", err)
		}

		// Only one backfill runs per profile at a time; the running one already
		// covers any window a second request could ask for.
		apiErr := err.(*postproxy.PostProxyError)
		runningID, _ := apiErr.Response["profile_sync_id"].(string)
		fmt.Printf("Backfill already running: %s\n", runningID)

		sync, err = client.Profiles.PostSync(ctx, profileID, runningID, nil)
		if err != nil {
			log.Fatalf("Error fetching the running sync: %v", err)
		}
	}

	fmt.Printf("Backfill %s — status: %s\n", sync.ID, sync.Status)

	// Poll until it finishes.
	for sync.Status == postproxy.PostSyncStatusPending || sync.Status == postproxy.PostSyncStatusRunning {
		time.Sleep(5 * time.Second)

		sync, err = client.Profiles.PostSync(ctx, profileID, sync.ID, nil)
		if err != nil {
			log.Fatalf("Error polling the sync: %v", err)
		}

		oldest := "—"
		if sync.OldestPostedAt != nil {
			oldest = *sync.OldestPostedAt
		}
		fmt.Printf("  %s: %d imported of %d seen, reached back to %s\n",
			sync.Status, sync.PostsImported, sync.PostsSeen, oldest)
	}

	if sync.Status == postproxy.PostSyncStatusFailed {
		fmt.Printf("Backfill failed: %v\n", sync.Error)
	} else {
		fmt.Printf("Done. Imported %d posts\n", sync.PostsImported)
	}

	// Every pull is recorded — the sync fired on connect, the recurring poll,
	// and each backfill. Runs are kept for 30 days.
	perPage := 10
	runs, err := client.Profiles.PostSyncs(ctx, profileID, &postproxy.PostSyncListOptions{PerPage: &perPage})
	if err != nil {
		log.Fatalf("Error listing post syncs: %v", err)
	}
	fmt.Printf("\nRecent post syncs (%d):\n", runs.Total)
	for _, run := range runs.Data {
		fmt.Printf("  %s %s → %s (%d/%d new)\n",
			run.CreatedAt, run.Trigger, run.Status, run.PostsImported, run.PostsSeen)
	}
}

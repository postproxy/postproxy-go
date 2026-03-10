package postproxy

import (
	"context"
	"encoding/json"
	"net/http"
)

// QueuesService handles communication with the queues related methods of the PostProxy API.
type QueuesService struct {
	client *Client
}

// List returns all queues.
func (s *QueuesService) List(ctx context.Context, opts *RequestOptions) (*ListResponse[Queue], error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/post_queues", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result ListResponse[Queue]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a single queue by ID.
func (s *QueuesService) Get(ctx context.Context, id string) (*Queue, error) {
	data, err := s.client.request(ctx, http.MethodGet, "/post_queues/"+id)
	if err != nil {
		return nil, err
	}

	var result Queue
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// NextSlot returns the next available timeslot for a queue.
func (s *QueuesService) NextSlot(ctx context.Context, id string) (*NextSlotResponse, error) {
	data, err := s.client.request(ctx, http.MethodGet, "/post_queues/"+id+"/next_slot")
	if err != nil {
		return nil, err
	}

	var result NextSlotResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new queue.
func (s *QueuesService) Create(ctx context.Context, name string, profileGroupID string, opts *QueueCreateOptions) (*Queue, error) {
	postQueue := map[string]any{
		"name": name,
	}
	if opts != nil {
		if opts.Description != nil {
			postQueue["description"] = *opts.Description
		}
		if opts.Timezone != nil {
			postQueue["timezone"] = *opts.Timezone
		}
		if opts.Jitter != nil {
			postQueue["jitter"] = *opts.Jitter
		}
		if opts.Timeslots != nil {
			postQueue["queue_timeslots_attributes"] = opts.Timeslots
		}
	}

	payload := map[string]any{
		"profile_group_id": profileGroupID,
		"post_queue":       postQueue,
	}

	data, err := s.client.request(ctx, http.MethodPost, "/post_queues", withJSON(payload))
	if err != nil {
		return nil, err
	}

	var result Queue
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates an existing queue.
func (s *QueuesService) Update(ctx context.Context, id string, opts *QueueUpdateOptions) (*Queue, error) {
	postQueue := map[string]any{}
	if opts != nil {
		if opts.Name != nil {
			postQueue["name"] = *opts.Name
		}
		if opts.Description != nil {
			postQueue["description"] = *opts.Description
		}
		if opts.Timezone != nil {
			postQueue["timezone"] = *opts.Timezone
		}
		if opts.Enabled != nil {
			postQueue["enabled"] = *opts.Enabled
		}
		if opts.Jitter != nil {
			postQueue["jitter"] = *opts.Jitter
		}
		if opts.Timeslots != nil {
			postQueue["queue_timeslots_attributes"] = opts.Timeslots
		}
	}

	payload := map[string]any{
		"post_queue": postQueue,
	}

	data, err := s.client.request(ctx, "PATCH", "/post_queues/"+id, withJSON(payload))
	if err != nil {
		return nil, err
	}

	var result Queue
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a queue by ID.
func (s *QueuesService) Delete(ctx context.Context, id string) (*DeleteResponse, error) {
	data, err := s.client.request(ctx, http.MethodDelete, "/post_queues/"+id)
	if err != nil {
		return nil, err
	}

	var result DeleteResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

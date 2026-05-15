package postproxy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ProfileCommentsService handles communication with the profile-level comments
// (Google Business reviews and replies) methods of the PostProxy API.
type ProfileCommentsService struct {
	client *Client
}

// List returns a paginated list of reviews and replies for a profile.
func (s *ProfileCommentsService) List(ctx context.Context, profileID string, opts *ProfileCommentListOptions) (*PaginatedResponse[ProfileComment], error) {
	params := url.Values{}
	reqOpts := []requestOption{}
	if opts != nil {
		if opts.PlacementID != nil {
			params.Set("placement_id", *opts.PlacementID)
		}
		if opts.Page != nil {
			params.Set("page", fmt.Sprintf("%d", *opts.Page))
		}
		if opts.PerPage != nil {
			params.Set("per_page", fmt.Sprintf("%d", *opts.PerPage))
		}
		if opts.ProfileGroupID != nil {
			reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
		}
	}
	if len(params) > 0 {
		reqOpts = append(reqOpts, withParams(params))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+profileID+"/comments", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result PaginatedResponse[ProfileComment]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a single profile comment by ID.
func (s *ProfileCommentsService) Get(ctx context.Context, profileID, commentID string) (*ProfileComment, error) {
	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+profileID+"/comments/"+commentID)
	if err != nil {
		return nil, err
	}

	var result ProfileComment
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a reply to a review.
func (s *ProfileCommentsService) Create(ctx context.Context, profileID, parentID, text string) (*ProfileComment, error) {
	body := map[string]string{"parent_id": parentID, "text": text}

	data, err := s.client.request(ctx, http.MethodPost, "/profiles/"+profileID+"/comments", withJSON(body))
	if err != nil {
		return nil, err
	}

	var result ProfileComment
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a reply.
func (s *ProfileCommentsService) Delete(ctx context.Context, profileID, commentID string) (*AcceptedResponse, error) {
	data, err := s.client.request(ctx, http.MethodDelete, "/profiles/"+profileID+"/comments/"+commentID)
	if err != nil {
		return nil, err
	}

	var result AcceptedResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

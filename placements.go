package postproxy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// PlacementsService handles account-wide placements. Requires an account-wide
// API key — placements span profile groups.
type PlacementsService struct {
	client *Client
}

// PlacementListOptions contains options for listing placements.
type PlacementListOptions struct {
	Page    *int
	PerPage *int
}

// List returns a paginated list of the account's placements.
func (s *PlacementsService) List(ctx context.Context, opts *PlacementListOptions) (*PaginatedResponse[AccountPlacement], error) {
	var reqOpts []requestOption
	if opts != nil {
		params := url.Values{}
		if opts.Page != nil {
			params.Set("page", strconv.Itoa(*opts.Page))
		}
		if opts.PerPage != nil {
			params.Set("per_page", strconv.Itoa(*opts.PerPage))
		}
		if len(params) > 0 {
			reqOpts = append(reqOpts, withParams(params))
		}
	}

	data, err := s.client.request(ctx, http.MethodGet, "/placements", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result PaginatedResponse[AccountPlacement]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetActiveProfile chooses which connected profile PostProxy posts and reads
// through, when the same placement is reachable through several.
func (s *PlacementsService) SetActiveProfile(ctx context.Context, id, profileID string) (*AccountPlacement, error) {
	data, err := s.client.request(ctx, http.MethodPatch, "/placements/"+id+"/set_active_profile", withJSON(map[string]string{"profile_id": profileID}))
	if err != nil {
		return nil, err
	}

	var result AccountPlacement
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

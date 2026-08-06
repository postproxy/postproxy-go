package postproxy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// ProfilesService handles communication with the profiles related methods of the PostProxy API.
type ProfilesService struct {
	client *Client
}

// List returns all profiles.
func (s *ProfilesService) List(ctx context.Context, opts *RequestOptions) (*ListResponse[Profile], error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/profiles", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result ListResponse[Profile]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a single profile by ID.
func (s *ProfilesService) Get(ctx context.Context, id string, opts *RequestOptions) (*Profile, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+id, reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Profile
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProfileStats returns the profile stats timeseries. `PlacementID` is required
// for facebook, linkedin, and telegram profiles.
func (s *ProfilesService) GetProfileStats(ctx context.Context, id string, opts *ProfileStatsOptions) (*ProfileStatsResponse, error) {
	var reqOpts []requestOption

	if opts != nil {
		params := url.Values{}
		if opts.PlacementID != nil {
			params.Set("placement_id", *opts.PlacementID)
		}
		if opts.From != nil {
			params.Set("from", *opts.From)
		}
		if opts.To != nil {
			params.Set("to", *opts.To)
		}
		if len(params) > 0 {
			reqOpts = append(reqOpts, withParams(params))
		}
		if opts.ProfileGroupID != nil {
			reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
		}
	}

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+id+"/stats", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result ProfileStatsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Placements returns the available placements for a profile.
func (s *ProfilesService) Placements(ctx context.Context, id string, opts *RequestOptions) (*ListResponse[Placement], error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+id+"/placements", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result ListResponse[Placement]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AssignPlacementToGroup moves a placement (e.g. a Facebook Page, Telegram
// channel, or Google Business location) to another profile group. placementID
// is the placement's external ID as returned by Placements.
func (s *ProfilesService) AssignPlacementToGroup(ctx context.Context, id, placementID, targetProfileGroupID string, opts *RequestOptions) (*Placement, error) {
	reqOpts := []requestOption{withJSON(map[string]any{
		"placement_id":            placementID,
		"target_profile_group_id": targetProfileGroupID,
	})}
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodPatch, "/profiles/"+id+"/assign_placement_to_group", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Placement
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BackfillPosts imports older posts from the platform. It walks the profile's
// feed backwards from the newest post until it reaches `from` or the platform
// stops returning posts. Runs in the background — poll PostSync with the
// returned ID for progress.
//
// Only one backfill runs per profile at a time; starting a second returns a
// 409 (see IsConflictError) whose Response carries the running one's
// "profile_sync_id". How far back a run reaches depends on the platform's API,
// not on PostProxy.
func (s *ProfilesService) BackfillPosts(ctx context.Context, id, from string, opts *RequestOptions) (*PostSync, error) {
	reqOpts := []requestOption{withJSON(map[string]string{"from": from})}
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodPost, "/profiles/"+id+"/backfill_posts", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result PostSync
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PostSyncs returns post sync runs for a profile, newest first. Runs are kept
// for 30 days.
func (s *ProfilesService) PostSyncs(ctx context.Context, id string, opts *PostSyncListOptions) (*PaginatedResponse[PostSync], error) {
	var reqOpts []requestOption

	if opts != nil {
		params := url.Values{}
		if opts.Trigger != nil {
			params.Set("trigger", string(*opts.Trigger))
		}
		if opts.Status != nil {
			params.Set("status", string(*opts.Status))
		}
		if opts.Page != nil {
			params.Set("page", strconv.Itoa(*opts.Page))
		}
		if opts.PerPage != nil {
			params.Set("per_page", strconv.Itoa(*opts.PerPage))
		}
		if len(params) > 0 {
			reqOpts = append(reqOpts, withParams(params))
		}
		if opts.ProfileGroupID != nil {
			reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
		}
	}

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+id+"/post_syncs", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result PaginatedResponse[PostSync]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PostSync returns a single run. Poll this to follow a backfill to completion —
// the run is finished when Status is PostSyncStatusCompleted or
// PostSyncStatusFailed.
func (s *ProfilesService) PostSync(ctx context.Context, id, postSyncID string, opts *RequestOptions) (*PostSync, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+id+"/post_syncs/"+postSyncID, reqOpts...)
	if err != nil {
		return nil, err
	}

	var result PostSync
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// IceBreakers returns the DM ice breakers for a profile. Supported for
// Instagram profiles only.
func (s *ProfilesService) IceBreakers(ctx context.Context, id string, opts *RequestOptions) (*IceBreakersResponse, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+id+"/ice_breakers", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result IceBreakersResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetIceBreakers replaces the DM ice breakers for a profile (1-4 items).
func (s *ProfilesService) SetIceBreakers(ctx context.Context, id string, iceBreakers []IceBreaker, opts *RequestOptions) (*SuccessResponse, error) {
	reqOpts := []requestOption{withJSON(map[string]any{"ice_breakers": iceBreakers})}
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodPost, "/profiles/"+id+"/ice_breakers", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result SuccessResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteIceBreakers removes all DM ice breakers from a profile.
func (s *ProfilesService) DeleteIceBreakers(ctx context.Context, id string, opts *RequestOptions) (*SuccessResponse, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodDelete, "/profiles/"+id+"/ice_breakers", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result SuccessResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a profile by ID.
func (s *ProfilesService) Delete(ctx context.Context, id string, opts *RequestOptions) (*SuccessResponse, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodDelete, "/profiles/"+id, reqOpts...)
	if err != nil {
		return nil, err
	}

	var result SuccessResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

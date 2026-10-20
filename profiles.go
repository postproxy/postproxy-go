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

// DailyStats returns the profile's per-day stats series. PlacementID is
// required for facebook, linkedin, and telegram profiles.
func (s *ProfilesService) DailyStats(ctx context.Context, id string, opts *DailyStatsOptions) (*DailyStatsResponse, error) {
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
		if opts.Attribution != nil {
			params.Set("attribution", string(*opts.Attribution))
		}
		if len(params) > 0 {
			reqOpts = append(reqOpts, withParams(params))
		}
		if opts.ProfileGroupID != nil {
			reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
		}
	}

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+id+"/daily_stats", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result DailyStatsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// FollowerHistory returns daily follower gains, losses and totals.
func (s *ProfilesService) FollowerHistory(ctx context.Context, id string, opts *FollowerHistoryOptions) (*FollowerHistoryResponse, error) {
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

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+id+"/follower_history", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result FollowerHistoryResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ParticipantStatus reports whether an Instagram DM participant is verified
// and follows you. Instagram only.
func (s *ProfilesService) ParticipantStatus(ctx context.Context, id, userID string, opts *RequestOptions) (*ParticipantStatus, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+id+"/participants/"+url.PathEscape(userID), reqOpts...)
	if err != nil {
		return nil, err
	}

	var result ParticipantStatus
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ResolveMention resolves a LinkedIn company page URL (or vanity name) to the
// URN used to @-mention it in a post. LinkedIn only.
func (s *ProfilesService) ResolveMention(ctx context.Context, id, pageURL string, opts *RequestOptions) (*ResolvedMention, error) {
	params := url.Values{}
	params.Set("url", pageURL)
	reqOpts := []requestOption{withParams(params)}
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+id+"/resolve_mention", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result ResolvedMention
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// InstagramAvailableAudio lists music and original sounds an Instagram reel
// can attach via InstagramParams.AudioID. Profiles connected via Facebook
// login only.
func (s *ProfilesService) InstagramAvailableAudio(ctx context.Context, id string, opts *InstagramAvailableAudioOptions) (*InstagramAvailableAudioResponse, error) {
	var reqOpts []requestOption

	if opts != nil {
		params := url.Values{}
		if opts.AudioType != nil {
			params.Set("audio_type", string(*opts.AudioType))
		}
		if opts.SearchQuery != nil {
			params.Set("search_query", *opts.SearchQuery)
		}
		if len(params) > 0 {
			reqOpts = append(reqOpts, withParams(params))
		}
		if opts.ProfileGroupID != nil {
			reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
		}
	}

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+id+"/instagram/available_audio", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result InstagramAvailableAudioResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// FacebookPage reads a Facebook Page's about / contact fields. Facebook only;
// pageID is the placement id from Placements.
func (s *ProfilesService) FacebookPage(ctx context.Context, id, pageID string, opts *FacebookPageOptions) (map[string]any, error) {
	params := url.Values{}
	params.Set("page_id", pageID)
	if opts != nil && opts.Fields != nil {
		params.Set("fields", *opts.Fields)
	}
	reqOpts := []requestOption{withParams(params)}
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+id+"/facebook/page", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateFacebookPage updates a Facebook Page's about / contact fields, e.g.
// "about" (255 characters or fewer), "description", "website", "phone",
// "emails".
func (s *ProfilesService) UpdateFacebookPage(ctx context.Context, id, pageID string, fields map[string]any, opts *RequestOptions) (map[string]any, error) {
	body := make(map[string]any, len(fields)+1)
	for k, v := range fields {
		body[k] = v
	}
	body["page_id"] = pageID

	reqOpts := []requestOption{withJSON(body)}
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodPatch, "/profiles/"+id+"/facebook/update_page", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// AssignToGroup moves the profile (with its placements) to another profile group.
func (s *ProfilesService) AssignToGroup(ctx context.Context, id, targetProfileGroupID string, opts *RequestOptions) (*Profile, error) {
	reqOpts := []requestOption{withJSON(map[string]string{"target_profile_group_id": targetProfileGroupID})}
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodPatch, "/profiles/"+id+"/assign_to_group", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Profile
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BackfillChats imports DM history back to `from`. Like BackfillPosts, it runs
// in the background and only one runs per profile — poll ChatSync with the
// returned ID for progress.
func (s *ProfilesService) BackfillChats(ctx context.Context, id, from string, opts *RequestOptions) (*ChatSync, error) {
	reqOpts := []requestOption{withJSON(map[string]string{"from": from})}
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodPost, "/profiles/"+id+"/backfill_chats", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result ChatSync
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ChatSyncs returns chat sync runs for a profile, newest first.
func (s *ProfilesService) ChatSyncs(ctx context.Context, id string, opts *ChatSyncListOptions) (*PaginatedResponse[ChatSync], error) {
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

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+id+"/chat_syncs", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result PaginatedResponse[ChatSync]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ChatSync returns a single chat sync run.
func (s *ProfilesService) ChatSync(ctx context.Context, id, chatSyncID string, opts *RequestOptions) (*ChatSync, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/profiles/"+id+"/chat_syncs/"+chatSyncID, reqOpts...)
	if err != nil {
		return nil, err
	}

	var result ChatSync
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

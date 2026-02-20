package postproxy

import (
	"context"
	"encoding/json"
	"net/http"
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

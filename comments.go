package postproxy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// CommentsService handles communication with the comments related methods of the PostProxy API.
type CommentsService struct {
	client *Client
}

// List returns a paginated list of comments for a post.
func (s *CommentsService) List(ctx context.Context, postID, profileID string, opts *CommentListOptions) (*PaginatedResponse[Comment], error) {
	params := url.Values{}
	params.Set("profile_id", profileID)
	if opts != nil {
		if opts.Page != nil {
			params.Set("page", fmt.Sprintf("%d", *opts.Page))
		}
		if opts.PerPage != nil {
			params.Set("per_page", fmt.Sprintf("%d", *opts.PerPage))
		}
		if opts.From != nil {
			params.Set("from", *opts.From)
		}
		if opts.To != nil {
			params.Set("to", *opts.To)
		}
	}

	data, err := s.client.request(ctx, http.MethodGet, "/posts/"+postID+"/comments", withParams(params))
	if err != nil {
		return nil, err
	}

	var result PaginatedResponse[Comment]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListAll returns comments spanning every post in the profile group. Flat:
// replies come back as their own entries linked by ParentExternalID, so Total
// counts every comment and paging is exact.
//
// Unknown or out-of-scope IDs in PostIDs and Profiles are ignored rather than
// erroring. Results are ordered newest first by receipt time.
func (s *CommentsService) ListAll(ctx context.Context, opts *BulkCommentListOptions) (*PaginatedResponse[BulkComment], error) {
	var reqOpts []requestOption

	if opts != nil {
		params := url.Values{}
		if len(opts.PostIDs) > 0 {
			params.Set("post_ids", strings.Join(opts.PostIDs, ","))
		}
		if len(opts.Profiles) > 0 {
			params.Set("profiles", strings.Join(opts.Profiles, ","))
		}
		if opts.From != nil {
			params.Set("from", *opts.From)
		}
		if opts.To != nil {
			params.Set("to", *opts.To)
		}
		if opts.Page != nil {
			params.Set("page", fmt.Sprintf("%d", *opts.Page))
		}
		if opts.PerPage != nil {
			params.Set("per_page", fmt.Sprintf("%d", *opts.PerPage))
		}
		if len(params) > 0 {
			reqOpts = append(reqOpts, withParams(params))
		}
		if opts.ProfileGroupID != nil {
			reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
		}
	}

	data, err := s.client.request(ctx, http.MethodGet, "/comments", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result PaginatedResponse[BulkComment]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a single comment by ID.
func (s *CommentsService) Get(ctx context.Context, postID, commentID, profileID string) (*Comment, error) {
	params := url.Values{}
	params.Set("profile_id", profileID)

	data, err := s.client.request(ctx, http.MethodGet, "/posts/"+postID+"/comments/"+commentID, withParams(params))
	if err != nil {
		return nil, err
	}

	var result Comment
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new comment or reply on a post.
func (s *CommentsService) Create(ctx context.Context, postID, profileID, text string, opts *CommentCreateOptions) (*Comment, error) {
	params := url.Values{}
	params.Set("profile_id", profileID)

	body := map[string]string{"text": text}
	if opts != nil && opts.ParentID != nil {
		body["parent_id"] = *opts.ParentID
	}

	data, err := s.client.request(ctx, http.MethodPost, "/posts/"+postID+"/comments", withParams(params), withJSON(body))
	if err != nil {
		return nil, err
	}

	var result Comment
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a comment.
func (s *CommentsService) Delete(ctx context.Context, postID, commentID, profileID string) (*AcceptedResponse, error) {
	params := url.Values{}
	params.Set("profile_id", profileID)

	data, err := s.client.request(ctx, http.MethodDelete, "/posts/"+postID+"/comments/"+commentID, withParams(params))
	if err != nil {
		return nil, err
	}

	var result AcceptedResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Hide hides a comment.
func (s *CommentsService) Hide(ctx context.Context, postID, commentID, profileID string) (*AcceptedResponse, error) {
	return s.commentAction(ctx, postID, commentID, profileID, "hide")
}

// Unhide unhides a comment.
func (s *CommentsService) Unhide(ctx context.Context, postID, commentID, profileID string) (*AcceptedResponse, error) {
	return s.commentAction(ctx, postID, commentID, profileID, "unhide")
}

// Like likes a comment.
func (s *CommentsService) Like(ctx context.Context, postID, commentID, profileID string) (*AcceptedResponse, error) {
	return s.commentAction(ctx, postID, commentID, profileID, "like")
}

// Unlike unlikes a comment.
func (s *CommentsService) Unlike(ctx context.Context, postID, commentID, profileID string) (*AcceptedResponse, error) {
	return s.commentAction(ctx, postID, commentID, profileID, "unlike")
}

// PrivateReply sends a private reply (direct message) to the author of a comment.
// Returns the created Message (Instagram/Facebook).
func (s *CommentsService) PrivateReply(ctx context.Context, postID, commentID, profileID, text string) (*Message, error) {
	params := url.Values{}
	params.Set("profile_id", profileID)

	path := fmt.Sprintf("/posts/%s/comments/%s/private_reply", postID, commentID)
	data, err := s.client.request(ctx, http.MethodPost, path, withParams(params), withJSON(map[string]string{"text": text}))
	if err != nil {
		return nil, err
	}

	var result Message
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *CommentsService) commentAction(ctx context.Context, postID, commentID, profileID, action string) (*AcceptedResponse, error) {
	params := url.Values{}
	params.Set("profile_id", profileID)

	path := fmt.Sprintf("/posts/%s/comments/%s/%s", postID, commentID, action)
	data, err := s.client.request(ctx, http.MethodPost, path, withParams(params))
	if err != nil {
		return nil, err
	}

	var result AcceptedResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

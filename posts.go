package postproxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

// PostsService handles communication with the posts related methods of the PostProxy API.
type PostsService struct {
	client *Client
}

// List returns a paginated list of posts.
func (s *PostsService) List(ctx context.Context, opts *PostListOptions) (*PaginatedResponse[Post], error) {
	var reqOpts []requestOption
	if opts != nil {
		params := url.Values{}
		if opts.Page != nil {
			params.Set("page", strconv.Itoa(*opts.Page))
		}
		if opts.PerPage != nil {
			params.Set("per_page", strconv.Itoa(*opts.PerPage))
		}
		if opts.Status != nil {
			params.Set("status", string(*opts.Status))
		}
		if opts.Platforms != nil {
			for _, p := range opts.Platforms {
				params.Add("platforms[]", string(p))
			}
		}
		if opts.ScheduledAfter != nil {
			params.Set("scheduled_after", *opts.ScheduledAfter)
		}
		reqOpts = append(reqOpts, withParams(params))
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/posts", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result PaginatedResponse[Post]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a single post by ID.
func (s *PostsService) Get(ctx context.Context, id string, opts *RequestOptions) (*Post, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, "/posts/"+id, reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Post
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new post. If MediaFiles is set in opts, the request uses multipart/form-data.
func (s *PostsService) Create(ctx context.Context, body string, profiles []string, opts *PostCreateOptions) (*Post, error) {
	var reqOpts []requestOption

	if opts != nil && needsFormData(opts) {
		formOpts, err := s.buildFormData(body, profiles, opts)
		if err != nil {
			return nil, err
		}
		reqOpts = append(reqOpts, formOpts...)
	} else {
		jsonOpts := s.buildJSON(body, profiles, opts)
		reqOpts = append(reqOpts, jsonOpts...)
	}

	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodPost, "/posts", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Post
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates an existing post.
func (s *PostsService) Update(ctx context.Context, id string, opts *PostUpdateOptions) (*Post, error) {
	var reqOpts []requestOption

	if opts != nil && needsUpdateFormData(opts) {
		formOpts, err := s.buildUpdateFormData(opts)
		if err != nil {
			return nil, err
		}
		reqOpts = append(reqOpts, formOpts...)
	} else if opts != nil {
		jsonOpts := s.buildUpdateJSON(opts)
		reqOpts = append(reqOpts, jsonOpts...)
	}

	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, "PATCH", "/posts/"+id, reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Post
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func needsUpdateFormData(opts *PostUpdateOptions) bool {
	if len(opts.MediaFiles) > 0 {
		return true
	}
	for _, t := range opts.Thread {
		if len(t.MediaFiles) > 0 {
			return true
		}
	}
	return false
}

func (s *PostsService) buildUpdateJSON(opts *PostUpdateOptions) []requestOption {
	payload := map[string]any{}

	post := map[string]any{}
	if opts.Body != nil {
		post["body"] = *opts.Body
	}
	if opts.ScheduledAt != nil {
		post["scheduled_at"] = *opts.ScheduledAt
	}
	if opts.Draft != nil {
		post["draft"] = *opts.Draft
	}
	if len(post) > 0 {
		payload["post"] = post
	}

	if opts.Profiles != nil {
		payload["profiles"] = opts.Profiles
	}
	if opts.Media != nil {
		payload["media"] = opts.Media
	}
	if opts.Platforms != nil {
		payload["platforms"] = opts.Platforms
	}
	if opts.Thread != nil {
		payload["thread"] = opts.Thread
	}
	if opts.QueueID != nil {
		payload["queue_id"] = *opts.QueueID
	}
	if opts.QueuePriority != nil {
		payload["queue_priority"] = *opts.QueuePriority
	}

	return []requestOption{withJSON(payload)}
}

func (s *PostsService) buildUpdateFormData(opts *PostUpdateOptions) ([]requestOption, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	if opts.Body != nil {
		_ = w.WriteField("post[body]", *opts.Body)
	}

	if opts.Profiles != nil {
		for _, p := range opts.Profiles {
			_ = w.WriteField("profiles[]", p)
		}
	}

	if opts.ScheduledAt != nil {
		_ = w.WriteField("post[scheduled_at]", *opts.ScheduledAt)
	}
	if opts.Draft != nil {
		_ = w.WriteField("post[draft]", strconv.FormatBool(*opts.Draft))
	}
	if opts.QueueID != nil {
		_ = w.WriteField("queue_id", *opts.QueueID)
	}
	if opts.QueuePriority != nil {
		_ = w.WriteField("queue_priority", *opts.QueuePriority)
	}

	if opts.Media != nil {
		for _, u := range opts.Media {
			_ = w.WriteField("media[]", u)
		}
	}

	if opts.Platforms != nil {
		writePlatformParams(w, opts.Platforms)
	}

	for i, child := range opts.Thread {
		prefix := fmt.Sprintf("thread[%d]", i)
		_ = w.WriteField(prefix+"[body]", child.Body)
		for _, u := range child.Media {
			_ = w.WriteField(prefix+"[media][]", u)
		}
		for _, filePath := range child.MediaFiles {
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read thread media file %s: %w", filePath, err)
			}

			contentType := mimeTypeFromPath(filePath)
			filename := filepath.Base(filePath)

			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", fmt.Sprintf(`form-data; name=%q; filename=%q`, prefix+"[media][]", filename))
			h.Set("Content-Type", contentType)
			part, err := w.CreatePart(h)
			if err != nil {
				return nil, fmt.Errorf("failed to create form file: %w", err)
			}
			if _, err := part.Write(fileData); err != nil {
				return nil, fmt.Errorf("failed to write file data: %w", err)
			}
		}
	}

	for _, filePath := range opts.MediaFiles {
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read media file %s: %w", filePath, err)
		}

		contentType := mimeTypeFromPath(filePath)
		filename := filepath.Base(filePath)

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="media[]"; filename=%q`, filename))
		h.Set("Content-Type", contentType)
		part, err := w.CreatePart(h)
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}
		if _, err := part.Write(fileData); err != nil {
			return nil, fmt.Errorf("failed to write file data: %w", err)
		}
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return []requestOption{withFormData(&buf, w.FormDataContentType())}, nil
}

// PublishDraft publishes a draft post.
func (s *PostsService) PublishDraft(ctx context.Context, id string, opts *RequestOptions) (*Post, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodPost, "/posts/"+id+"/publish", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result Post
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a post by ID.
func (s *PostsService) Delete(ctx context.Context, id string, opts *RequestOptions) (*DeleteResponse, error) {
	var reqOpts []requestOption
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodDelete, "/posts/"+id, reqOpts...)
	if err != nil {
		return nil, err
	}

	var result DeleteResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Stats retrieves stats snapshots for one or more posts.
func (s *PostsService) Stats(ctx context.Context, postIDs []string, opts *PostStatsOptions) (*StatsResponse, error) {
	params := url.Values{}
	params.Set("post_ids", strings.Join(postIDs, ","))

	if opts != nil {
		if len(opts.Profiles) > 0 {
			params.Set("profiles", strings.Join(opts.Profiles, ","))
		}
		if opts.From != nil {
			params.Set("from", *opts.From)
		}
		if opts.To != nil {
			params.Set("to", *opts.To)
		}
	}

	data, err := s.client.request(ctx, http.MethodGet, "/posts/stats", withParams(params))
	if err != nil {
		return nil, err
	}

	var result StatsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func needsFormData(opts *PostCreateOptions) bool {
	if len(opts.MediaFiles) > 0 {
		return true
	}
	for _, t := range opts.Thread {
		if len(t.MediaFiles) > 0 {
			return true
		}
	}
	return false
}

func (s *PostsService) buildJSON(body string, profiles []string, opts *PostCreateOptions) []requestOption {
	payload := map[string]any{
		"post": map[string]any{
			"body": body,
		},
		"profiles": profiles,
	}

	post := payload["post"].(map[string]any)
	if opts != nil {
		if opts.Media != nil {
			payload["media"] = opts.Media
		}
		if opts.ScheduledAt != nil {
			post["scheduled_at"] = *opts.ScheduledAt
		}
		if opts.Draft != nil {
			post["draft"] = *opts.Draft
		}
		if opts.Platforms != nil {
			payload["platforms"] = opts.Platforms
		}
		if opts.Thread != nil {
			payload["thread"] = opts.Thread
		}
		if opts.QueueID != nil {
			payload["queue_id"] = *opts.QueueID
		}
		if opts.QueuePriority != nil {
			payload["queue_priority"] = *opts.QueuePriority
		}
	}

	return []requestOption{withJSON(payload)}
}

func (s *PostsService) buildFormData(body string, profiles []string, opts *PostCreateOptions) ([]requestOption, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	_ = w.WriteField("post[body]", body)

	for _, p := range profiles {
		_ = w.WriteField("profiles[]", p)
	}

	if opts.ScheduledAt != nil {
		_ = w.WriteField("post[scheduled_at]", *opts.ScheduledAt)
	}
	if opts.Draft != nil {
		_ = w.WriteField("post[draft]", strconv.FormatBool(*opts.Draft))
	}
	if opts.QueueID != nil {
		_ = w.WriteField("queue_id", *opts.QueueID)
	}
	if opts.QueuePriority != nil {
		_ = w.WriteField("queue_priority", *opts.QueuePriority)
	}

	if opts.Media != nil {
		for _, u := range opts.Media {
			_ = w.WriteField("media[]", u)
		}
	}

	if opts.Platforms != nil {
		writePlatformParams(w, opts.Platforms)
	}

	for i, child := range opts.Thread {
		prefix := fmt.Sprintf("thread[%d]", i)
		_ = w.WriteField(prefix+"[body]", child.Body)
		for _, u := range child.Media {
			_ = w.WriteField(prefix+"[media][]", u)
		}
		for _, filePath := range child.MediaFiles {
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read thread media file %s: %w", filePath, err)
			}

			contentType := mimeTypeFromPath(filePath)
			filename := filepath.Base(filePath)

			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", fmt.Sprintf(`form-data; name=%q; filename=%q`, prefix+"[media][]", filename))
			h.Set("Content-Type", contentType)
			part, err := w.CreatePart(h)
			if err != nil {
				return nil, fmt.Errorf("failed to create form file: %w", err)
			}
			if _, err := part.Write(fileData); err != nil {
				return nil, fmt.Errorf("failed to write file data: %w", err)
			}
		}
	}

	for _, filePath := range opts.MediaFiles {
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read media file %s: %w", filePath, err)
		}

		contentType := mimeTypeFromPath(filePath)
		filename := filepath.Base(filePath)

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="media[]"; filename=%q`, filename))
		h.Set("Content-Type", contentType)
		part, err := w.CreatePart(h)
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}
		if _, err := part.Write(fileData); err != nil {
			return nil, fmt.Errorf("failed to write file data: %w", err)
		}
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return []requestOption{withFormData(&buf, w.FormDataContentType())}, nil
}

func writePlatformParams(w *multipart.Writer, pp *PlatformParams) {
	if pp == nil {
		return
	}

	v := reflect.ValueOf(*pp)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.IsNil() {
			continue
		}

		tag := t.Field(i).Tag.Get("json")
		platform := strings.Split(tag, ",")[0]
		params := field.Elem()
		paramsType := params.Type()

		for j := 0; j < params.NumField(); j++ {
			pField := params.Field(j)
			pTag := paramsType.Field(j).Tag.Get("json")
			key := strings.Split(pTag, ",")[0]
			prefix := fmt.Sprintf("platforms[%s][%s]", platform, key)

			if pField.Kind() == reflect.Slice {
				for k := 0; k < pField.Len(); k++ {
					_ = w.WriteField(prefix+"[]", fmt.Sprintf("%v", pField.Index(k).Interface()))
				}
			} else if pField.Kind() == reflect.Ptr {
				if !pField.IsNil() {
					_ = w.WriteField(prefix, fmt.Sprintf("%v", pField.Elem().Interface()))
				}
			}
		}
	}
}

var mimeTypes = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
	".mp4":  "video/mp4",
	".mov":  "video/quicktime",
	".avi":  "video/x-msvideo",
	".webm": "video/webm",
}

func mimeTypeFromPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if ct, ok := mimeTypes[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}

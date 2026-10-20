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

// Create creates a new post. If MediaFiles (or a local CoverFile in the
// platform params) is set in opts, the request uses multipart/form-data.
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
		if opts.SkipDuplicateCheck {
			reqOpts = append(reqOpts, withHeader("X-Skip-Duplicate-Check", "true"))
		}
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

// Update updates an existing post. On a published post only the body, the
// YouTube title and the Pinterest title / destination link can change; the
// edit is pushed to the platforms in the background and reported through the
// platform_post.edited / platform_post.edit_failed webhooks.
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
	if len(opts.MediaFiles) > 0 || hasLocalCoverFile(opts.Platforms) {
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

	if err := writePlatformParams(w, opts.Platforms); err != nil {
		return nil, err
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

// Delete deletes a post by ID. Pass DeleteOnPlatform=true in opts to also delete the post from all published platforms.
func (s *PostsService) Delete(ctx context.Context, id string, opts *PostDeleteOptions) (*DeleteResponse, error) {
	var reqOpts []requestOption
	if opts != nil {
		if opts.DeleteOnPlatform != nil {
			params := url.Values{}
			params.Set("delete_on_platform", strconv.FormatBool(*opts.DeleteOnPlatform))
			reqOpts = append(reqOpts, withParams(params))
		}
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

// DeleteOnPlatform deletes a published post from one or more social media platforms without removing it from the database.
// If opts is nil (or all targeting fields are nil), the post is deleted from all published platforms.
func (s *PostsService) DeleteOnPlatform(ctx context.Context, id string, opts *PostDeleteOnPlatformOptions) (*DeleteOnPlatformResponse, error) {
	var reqOpts []requestOption
	payload := map[string]any{}
	if opts != nil {
		if opts.PostProfileID != nil {
			payload["post_profile_id"] = *opts.PostProfileID
		}
		if opts.ProfileID != nil {
			payload["profile_id"] = *opts.ProfileID
		}
		if opts.Network != nil {
			payload["network"] = *opts.Network
		}
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}
	if len(payload) > 0 {
		reqOpts = append(reqOpts, withJSON(payload))
	}

	data, err := s.client.request(ctx, http.MethodPost, "/posts/"+id+"/delete_on_platform", reqOpts...)
	if err != nil {
		return nil, err
	}

	var result DeleteOnPlatformResponse
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
	if len(opts.MediaFiles) > 0 || hasLocalCoverFile(opts.Platforms) {
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

	if err := writePlatformParams(w, opts.Platforms); err != nil {
		return nil, err
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

// isLocalCoverFile reports whether a cover_file value is a local path. A
// base64 data URL is sent to the API as a plain string instead.
func isLocalCoverFile(value *string) bool {
	return value != nil && *value != "" && !strings.HasPrefix(*value, "data:")
}

func hasLocalCoverFile(pp *PlatformParams) bool {
	if pp == nil {
		return false
	}
	return (pp.Instagram != nil && isLocalCoverFile(pp.Instagram.CoverFile)) ||
		(pp.TikTok != nil && isLocalCoverFile(pp.TikTok.CoverFile)) ||
		(pp.YouTube != nil && isLocalCoverFile(pp.YouTube.CoverFile)) ||
		(pp.Pinterest != nil && isLocalCoverFile(pp.Pinterest.CoverFile))
}

// writeFilePart attaches a local file to the form under name.
func writeFilePart(w *multipart.Writer, name, filePath string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name=%q; filename=%q`, name, filepath.Base(filePath)))
	h.Set("Content-Type", mimeTypeFromPath(filePath))
	part, err := w.CreatePart(h)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(fileData); err != nil {
		return fmt.Errorf("failed to write file data: %w", err)
	}
	return nil
}

// writePlatformParams writes platform params as flat form fields. A local
// cover_file is attached as a file.
func writePlatformParams(w *multipart.Writer, pp *PlatformParams) error {
	if pp == nil {
		return nil
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

		if field.Kind() == reflect.Map {
			for _, key := range field.MapKeys() {
				prefix := fmt.Sprintf("platforms[%s][%v]", platform, key.Interface())
				writeFormValue(w, prefix, reflect.ValueOf(field.MapIndex(key).Interface()))
			}
			continue
		}

		params := field.Elem()
		paramsType := params.Type()

		for j := 0; j < params.NumField(); j++ {
			pField := params.Field(j)
			pTag := paramsType.Field(j).Tag.Get("json")
			key := strings.Split(pTag, ",")[0]
			prefix := fmt.Sprintf("platforms[%s][%s]", platform, key)

			if key == "cover_file" {
				if cover, ok := pField.Interface().(*string); ok && isLocalCoverFile(cover) {
					if err := writeFilePart(w, prefix, *cover); err != nil {
						return err
					}
					continue
				}
			}
			// A single alt text is sent as a plain string, as in JSON.
			if alt, ok := pField.Interface().(AltText); ok && len(alt) == 1 {
				_ = w.WriteField(prefix, alt[0])
				continue
			}
			writeFormValue(w, prefix, pField)
		}
	}
	return nil
}

func writeFormValue(w *multipart.Writer, prefix string, value reflect.Value) {
	switch value.Kind() {
	case reflect.Slice:
		for k := 0; k < value.Len(); k++ {
			_ = w.WriteField(prefix+"[]", fmt.Sprintf("%v", value.Index(k).Interface()))
		}
	case reflect.Ptr:
		if !value.IsNil() {
			_ = w.WriteField(prefix, fmt.Sprintf("%v", value.Elem().Interface()))
		}
	case reflect.Invalid:
	default:
		_ = w.WriteField(prefix, fmt.Sprintf("%v", value.Interface()))
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

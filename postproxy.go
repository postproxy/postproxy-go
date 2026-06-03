package postproxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
)

var userAgent = fmt.Sprintf("postproxy-go/%s (go/%s; %s/%s)", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)

const defaultBaseURL = "https://api.postproxy.dev"

// Client is the PostProxy API client.
type Client struct {
	apiKey         string
	baseURL        string
	httpClient     *http.Client
	profileGroupID string

	Posts           *PostsService
	Profiles        *ProfilesService
	ProfileGroups   *ProfileGroupsService
	Webhooks        *WebhooksService
	Queues          *QueuesService
	Comments        *CommentsService
	ProfileComments *ProfileCommentsService
	Chats           *ChatsService
	Messages        *MessagesService
}

// Option configures the Client.
type Option func(*Client)

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithProfileGroupID sets the default profile group ID for all requests.
func WithProfileGroupID(id string) Option {
	return func(c *Client) {
		c.profileGroupID = id
	}
}

// NewClient creates a new PostProxy API client.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	c.Posts = &PostsService{client: c}
	c.Profiles = &ProfilesService{client: c}
	c.ProfileGroups = &ProfileGroupsService{client: c}
	c.Webhooks = &WebhooksService{client: c}
	c.Queues = &QueuesService{client: c}
	c.Comments = &CommentsService{client: c}
	c.ProfileComments = &ProfileCommentsService{client: c}
	c.Chats = &ChatsService{client: c}
	c.Messages = &MessagesService{client: c}
	return c
}

type requestOptions struct {
	params         url.Values
	body           any
	formBody       io.Reader
	contentType    string
	profileGroupID *string
}

type requestOption func(*requestOptions)

func withParams(params url.Values) requestOption {
	return func(o *requestOptions) {
		o.params = params
	}
}

func withJSON(body any) requestOption {
	return func(o *requestOptions) {
		o.body = body
	}
}

func withFormData(body io.Reader, contentType string) requestOption {
	return func(o *requestOptions) {
		o.formBody = body
		o.contentType = contentType
	}
}

func withProfileGroupID(id *string) requestOption {
	return func(o *requestOptions) {
		o.profileGroupID = id
	}
}

func (c *Client) request(ctx context.Context, method, path string, opts ...requestOption) ([]byte, error) {
	ro := &requestOptions{}
	for _, opt := range opts {
		opt(ro)
	}

	u, err := url.Parse(c.baseURL + "/api" + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Build query params.
	q := u.Query()
	if ro.params != nil {
		for k, vs := range ro.params {
			for _, v := range vs {
				q.Add(k, v)
			}
		}
	}

	// Add profile group ID.
	pgID := c.profileGroupID
	if ro.profileGroupID != nil {
		pgID = *ro.profileGroupID
	}
	if pgID != "" {
		q.Set("profile_group_id", pgID)
	}
	u.RawQuery = q.Encode()

	// Build request body.
	var bodyReader io.Reader
	contentType := ""

	if ro.formBody != nil {
		bodyReader = ro.formBody
		contentType = ro.contentType
	} else if ro.body != nil {
		data, err := json.Marshal(ro.body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
		contentType = "application/json"
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, parseError(resp.StatusCode, respBody)
	}

	return respBody, nil
}

func parseError(statusCode int, body []byte) error {
	apiErr := &PostProxyError{
		StatusCode: statusCode,
		Message:    http.StatusText(statusCode),
	}

	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err == nil {
		apiErr.Response = parsed
		if msg, ok := parsed["error"].(string); ok {
			apiErr.Message = msg
		} else if msg, ok := parsed["message"].(string); ok {
			apiErr.Message = msg
		}
	}

	return apiErr
}

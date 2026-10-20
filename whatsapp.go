package postproxy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// WhatsAppService handles WhatsApp number administration: registration,
// business profile, display name, username, blocked users, groups and
// conversion events.
//
// phoneNumberID is the placement id from ProfilesService.Placements. Dataset
// and conversion calls are account-level and don't take one. Meta's responses
// are passed through as map[string]any.
type WhatsAppService struct {
	client *Client
}

// WhatsAppVerificationCodeOptions contains options for
// WhatsAppService.RequestVerificationCode.
type WhatsAppVerificationCodeOptions struct {
	// Method is "SMS" or "VOICE".
	Method         *string
	Language       *string
	ProfileGroupID *string
}

// WhatsAppBusinessProfile holds the business profile fields to change.
type WhatsAppBusinessProfile struct {
	About       *string
	Address     *string
	Description *string
	Email       *string
	Websites    []string
	Vertical    *string
}

// WhatsAppProfilePhoto is a business profile photo: a public URL, or base64
// Data with its ContentType (image/jpeg or image/png).
type WhatsAppProfilePhoto struct {
	URL         string
	Data        string
	ContentType string
}

// WhatsAppSetUsernameOptions contains options for WhatsAppService.SetUsername.
type WhatsAppSetUsernameOptions struct {
	// TransferAction is "none" or "force_transfer".
	TransferAction *string
	ProfileGroupID *string
}

// WhatsAppPageOptions contains cursor paging options for list calls.
type WhatsAppPageOptions struct {
	Limit          *int
	After          *string
	ProfileGroupID *string
}

// WhatsAppGroupCreateOptions contains options for WhatsAppService.CreateGroup.
type WhatsAppGroupCreateOptions struct {
	Description *string
	// JoinApprovalMode is "approval_required" or "auto_approve".
	JoinApprovalMode *string
	ProfileGroupID   *string
}

// WhatsAppGroupUpdate holds the group fields to change.
type WhatsAppGroupUpdate struct {
	Subject     *string
	Description *string
	// JoinApprovalMode is "approval_required" or "auto_approve".
	JoinApprovalMode *string
}

// WhatsAppConversionEvent reports a conversion (e.g. Purchase, LeadSubmitted)
// back to Meta for a click-to-WhatsApp ad. Identify the person by ChatID or
// Phone.
type WhatsAppConversionEvent struct {
	EventName string `json:"event_name"`
	EventID   string `json:"event_id"`
	ChatID    string `json:"chat_id,omitempty"`
	Phone     string `json:"phone,omitempty"`
	DatasetID string `json:"dataset_id,omitempty"`
	// EventTime is a Unix timestamp or an ISO 8601 string.
	EventTime  any      `json:"event_time,omitempty"`
	Email      string   `json:"email,omitempty"`
	ExternalID string   `json:"external_id,omitempty"`
	Currency   string   `json:"currency,omitempty"`
	Value      *float64 `json:"value,omitempty"`
	ContentIDs []string `json:"content_ids,omitempty"`
	TestCode   string   `json:"test_code,omitempty"`
}

// WhatsAppMedia is a downloaded media file.
type WhatsAppMedia struct {
	Data        []byte
	ContentType string
}

// call sends GET and DELETE data as query params and everything else as JSON.
func (s *WhatsAppService) call(ctx context.Context, method, profileID, action string, data map[string]any, profileGroupID *string) (map[string]any, error) {
	path := "/profiles/" + profileID + "/whatsapp/" + action
	reqOpts := []requestOption{withProfileGroupID(profileGroupID)}

	if method == http.MethodGet || method == http.MethodDelete {
		params := url.Values{}
		for k, v := range data {
			params.Set(k, fmt.Sprintf("%v", v))
		}
		reqOpts = append(reqOpts, withParams(params))
	} else {
		reqOpts = append(reqOpts, withJSON(data))
	}

	body, err := s.client.request(ctx, method, path, reqOpts...)
	if err != nil {
		return nil, err
	}
	return decodeMap(body)
}

func phoneNumber(phoneNumberID string) map[string]any {
	return map[string]any{"phone_number_id": phoneNumberID}
}

func phoneNumberGroup(phoneNumberID, groupID string) map[string]any {
	return map[string]any{"phone_number_id": phoneNumberID, "group_id": groupID}
}

func setOptional(data map[string]any, key string, value *string) {
	if value != nil {
		data[key] = *value
	}
}

// --- Phone number ---

// NumberInfo returns the phone number's details.
func (s *WhatsAppService) NumberInfo(ctx context.Context, profileID, phoneNumberID string, opts *RequestOptions) (map[string]any, error) {
	return s.call(ctx, http.MethodGet, profileID, "number_info", phoneNumber(phoneNumberID), pgID(opts))
}

// RegisterPhoneNumber registers the number. pin is the 6-digit two-step
// verification PIN.
func (s *WhatsAppService) RegisterPhoneNumber(ctx context.Context, profileID, phoneNumberID, pin string, opts *RequestOptions) (map[string]any, error) {
	data := phoneNumber(phoneNumberID)
	data["pin"] = pin
	return s.call(ctx, http.MethodPost, profileID, "register_phone_number", data, pgID(opts))
}

// RequestVerificationCode sends a verification code by SMS or voice call.
func (s *WhatsAppService) RequestVerificationCode(ctx context.Context, profileID, phoneNumberID string, opts *WhatsAppVerificationCodeOptions) (map[string]any, error) {
	data := phoneNumber(phoneNumberID)
	var pg *string
	if opts != nil {
		setOptional(data, "method", opts.Method)
		setOptional(data, "language", opts.Language)
		pg = opts.ProfileGroupID
	}
	return s.call(ctx, http.MethodPost, profileID, "request_verification_code", data, pg)
}

// VerifyPhoneNumber verifies the number with the code received.
func (s *WhatsAppService) VerifyPhoneNumber(ctx context.Context, profileID, phoneNumberID, code string, opts *RequestOptions) (map[string]any, error) {
	data := phoneNumber(phoneNumberID)
	data["code"] = code
	return s.call(ctx, http.MethodPost, profileID, "verify_phone_number", data, pgID(opts))
}

// --- Business profile ---

// BusinessProfile returns the number's business profile.
func (s *WhatsAppService) BusinessProfile(ctx context.Context, profileID, phoneNumberID string, opts *RequestOptions) (map[string]any, error) {
	return s.call(ctx, http.MethodGet, profileID, "business_profile", phoneNumber(phoneNumberID), pgID(opts))
}

// UpdateBusinessProfile changes the number's business profile.
func (s *WhatsAppService) UpdateBusinessProfile(ctx context.Context, profileID, phoneNumberID string, fields WhatsAppBusinessProfile, opts *RequestOptions) (map[string]any, error) {
	data := phoneNumber(phoneNumberID)
	setOptional(data, "about", fields.About)
	setOptional(data, "address", fields.Address)
	setOptional(data, "description", fields.Description)
	setOptional(data, "email", fields.Email)
	if fields.Websites != nil {
		data["websites"] = fields.Websites
	}
	setOptional(data, "vertical", fields.Vertical)
	return s.call(ctx, http.MethodPatch, profileID, "update_business_profile", data, pgID(opts))
}

// UpdateBusinessProfilePhoto sets the business profile photo.
func (s *WhatsAppService) UpdateBusinessProfilePhoto(ctx context.Context, profileID, phoneNumberID string, photo WhatsAppProfilePhoto, opts *RequestOptions) (map[string]any, error) {
	data := phoneNumber(phoneNumberID)
	if photo.URL != "" {
		data["url"] = photo.URL
	}
	if photo.Data != "" {
		data["data"] = photo.Data
	}
	if photo.ContentType != "" {
		data["content_type"] = photo.ContentType
	}
	return s.call(ctx, http.MethodPost, profileID, "update_business_profile_photo", data, pgID(opts))
}

// --- Display name ---

// DisplayName returns the number's display name and its review status.
func (s *WhatsAppService) DisplayName(ctx context.Context, profileID, phoneNumberID string, opts *RequestOptions) (map[string]any, error) {
	return s.call(ctx, http.MethodGet, profileID, "display_name", phoneNumber(phoneNumberID), pgID(opts))
}

// RequestDisplayNameChange asks Meta to review a new display name; the
// outcome arrives as the placement.name_status_updated webhook.
func (s *WhatsAppService) RequestDisplayNameChange(ctx context.Context, profileID, phoneNumberID, displayName string, opts *RequestOptions) (map[string]any, error) {
	data := phoneNumber(phoneNumberID)
	data["display_name"] = displayName
	return s.call(ctx, http.MethodPost, profileID, "request_display_name_change", data, pgID(opts))
}

// --- Username ---

// Username returns the number's username.
func (s *WhatsAppService) Username(ctx context.Context, profileID, phoneNumberID string, opts *RequestOptions) (map[string]any, error) {
	return s.call(ctx, http.MethodGet, profileID, "username", phoneNumber(phoneNumberID), pgID(opts))
}

// SetUsername sets the number's username.
func (s *WhatsAppService) SetUsername(ctx context.Context, profileID, phoneNumberID, username string, opts *WhatsAppSetUsernameOptions) (map[string]any, error) {
	data := phoneNumber(phoneNumberID)
	data["username"] = username
	var pg *string
	if opts != nil {
		setOptional(data, "transfer_action", opts.TransferAction)
		pg = opts.ProfileGroupID
	}
	return s.call(ctx, http.MethodPut, profileID, "set_username", data, pg)
}

// DeleteUsername removes the number's username.
func (s *WhatsAppService) DeleteUsername(ctx context.Context, profileID, phoneNumberID string, opts *RequestOptions) (map[string]any, error) {
	return s.call(ctx, http.MethodDelete, profileID, "delete_username", phoneNumber(phoneNumberID), pgID(opts))
}

// UsernameSuggestions returns available usernames for the number.
func (s *WhatsAppService) UsernameSuggestions(ctx context.Context, profileID, phoneNumberID string, opts *RequestOptions) (map[string]any, error) {
	return s.call(ctx, http.MethodGet, profileID, "username_suggestions", phoneNumber(phoneNumberID), pgID(opts))
}

// --- Blocked users ---

func pageData(data map[string]any, opts *WhatsAppPageOptions) *string {
	if opts == nil {
		return nil
	}
	if opts.Limit != nil {
		data["limit"] = strconv.Itoa(*opts.Limit)
	}
	setOptional(data, "after", opts.After)
	return opts.ProfileGroupID
}

// BlockedUsers lists the users the number has blocked.
func (s *WhatsAppService) BlockedUsers(ctx context.Context, profileID, phoneNumberID string, opts *WhatsAppPageOptions) (map[string]any, error) {
	data := phoneNumber(phoneNumberID)
	pg := pageData(data, opts)
	return s.call(ctx, http.MethodGet, profileID, "blocked_users", data, pg)
}

// BlockedUserStatus reports whether a user is blocked.
func (s *WhatsAppService) BlockedUserStatus(ctx context.Context, profileID, phoneNumberID, user string, opts *RequestOptions) (map[string]any, error) {
	data := phoneNumber(phoneNumberID)
	data["user"] = user
	return s.call(ctx, http.MethodGet, profileID, "blocked_user_status", data, pgID(opts))
}

// BlockUsers blocks users, given as phone numbers or WhatsApp ids.
func (s *WhatsAppService) BlockUsers(ctx context.Context, profileID, phoneNumberID string, users []string, opts *RequestOptions) (map[string]any, error) {
	data := phoneNumber(phoneNumberID)
	data["users"] = users
	return s.call(ctx, http.MethodPost, profileID, "block_users", data, pgID(opts))
}

// UnblockUsers unblocks users, given as phone numbers or WhatsApp ids.
func (s *WhatsAppService) UnblockUsers(ctx context.Context, profileID, phoneNumberID string, users []string, opts *RequestOptions) (map[string]any, error) {
	data := phoneNumber(phoneNumberID)
	data["users"] = users
	return s.call(ctx, http.MethodPost, profileID, "unblock_users", data, pgID(opts))
}

// --- Groups ---

// Groups lists the number's groups.
func (s *WhatsAppService) Groups(ctx context.Context, profileID, phoneNumberID string, opts *WhatsAppPageOptions) (map[string]any, error) {
	data := phoneNumber(phoneNumberID)
	pg := pageData(data, opts)
	return s.call(ctx, http.MethodGet, profileID, "groups", data, pg)
}

// CreateGroup creates a group with the given subject.
func (s *WhatsAppService) CreateGroup(ctx context.Context, profileID, phoneNumberID, subject string, opts *WhatsAppGroupCreateOptions) (map[string]any, error) {
	data := phoneNumber(phoneNumberID)
	data["subject"] = subject
	var pg *string
	if opts != nil {
		setOptional(data, "description", opts.Description)
		setOptional(data, "join_approval_mode", opts.JoinApprovalMode)
		pg = opts.ProfileGroupID
	}
	return s.call(ctx, http.MethodPost, profileID, "create_group", data, pg)
}

// Group returns a single group.
func (s *WhatsAppService) Group(ctx context.Context, profileID, phoneNumberID, groupID string, opts *RequestOptions) (map[string]any, error) {
	return s.call(ctx, http.MethodGet, profileID, "group", phoneNumberGroup(phoneNumberID, groupID), pgID(opts))
}

// UpdateGroup changes a group's subject, description or join approval mode.
func (s *WhatsAppService) UpdateGroup(ctx context.Context, profileID, phoneNumberID, groupID string, changes WhatsAppGroupUpdate, opts *RequestOptions) (map[string]any, error) {
	data := phoneNumberGroup(phoneNumberID, groupID)
	setOptional(data, "subject", changes.Subject)
	setOptional(data, "description", changes.Description)
	setOptional(data, "join_approval_mode", changes.JoinApprovalMode)
	return s.call(ctx, http.MethodPatch, profileID, "update_group", data, pgID(opts))
}

// DeleteGroup deletes a group.
func (s *WhatsAppService) DeleteGroup(ctx context.Context, profileID, phoneNumberID, groupID string, opts *RequestOptions) (map[string]any, error) {
	return s.call(ctx, http.MethodDelete, profileID, "delete_group", phoneNumberGroup(phoneNumberID, groupID), pgID(opts))
}

// AddGroupParticipants adds phone numbers to a group.
func (s *WhatsAppService) AddGroupParticipants(ctx context.Context, profileID, phoneNumberID, groupID string, phoneNumbers []string, opts *RequestOptions) (map[string]any, error) {
	data := phoneNumberGroup(phoneNumberID, groupID)
	data["phone_numbers"] = phoneNumbers
	return s.call(ctx, http.MethodPost, profileID, "add_group_participants", data, pgID(opts))
}

// RemoveGroupParticipants removes phone numbers from a group.
func (s *WhatsAppService) RemoveGroupParticipants(ctx context.Context, profileID, phoneNumberID, groupID string, phoneNumbers []string, opts *RequestOptions) (map[string]any, error) {
	data := phoneNumberGroup(phoneNumberID, groupID)
	data["phone_numbers"] = phoneNumbers
	return s.call(ctx, http.MethodPost, profileID, "remove_group_participants", data, pgID(opts))
}

// CreateGroupInviteLink creates an invite link for a group.
func (s *WhatsAppService) CreateGroupInviteLink(ctx context.Context, profileID, phoneNumberID, groupID string, opts *RequestOptions) (map[string]any, error) {
	return s.call(ctx, http.MethodPost, profileID, "create_group_invite_link", phoneNumberGroup(phoneNumberID, groupID), pgID(opts))
}

// GroupJoinRequests lists pending requests to join a group.
func (s *WhatsAppService) GroupJoinRequests(ctx context.Context, profileID, phoneNumberID, groupID string, opts *RequestOptions) (map[string]any, error) {
	return s.call(ctx, http.MethodGet, profileID, "group_join_requests", phoneNumberGroup(phoneNumberID, groupID), pgID(opts))
}

// ApproveGroupJoinRequests approves join requests from phone numbers.
func (s *WhatsAppService) ApproveGroupJoinRequests(ctx context.Context, profileID, phoneNumberID, groupID string, phoneNumbers []string, opts *RequestOptions) (map[string]any, error) {
	data := phoneNumberGroup(phoneNumberID, groupID)
	data["phone_numbers"] = phoneNumbers
	return s.call(ctx, http.MethodPost, profileID, "approve_group_join_requests", data, pgID(opts))
}

// RejectGroupJoinRequests rejects join requests from phone numbers.
func (s *WhatsAppService) RejectGroupJoinRequests(ctx context.Context, profileID, phoneNumberID, groupID string, phoneNumbers []string, opts *RequestOptions) (map[string]any, error) {
	data := phoneNumberGroup(phoneNumberID, groupID)
	data["phone_numbers"] = phoneNumbers
	return s.call(ctx, http.MethodPost, profileID, "reject_group_join_requests", data, pgID(opts))
}

// --- Conversions (account-level) ---

// Dataset returns the account's conversions dataset.
func (s *WhatsAppService) Dataset(ctx context.Context, profileID string, opts *RequestOptions) (map[string]any, error) {
	return s.call(ctx, http.MethodGet, profileID, "dataset", map[string]any{}, pgID(opts))
}

// CreateDataset creates the account's conversions dataset.
func (s *WhatsAppService) CreateDataset(ctx context.Context, profileID string, opts *RequestOptions) (map[string]any, error) {
	return s.call(ctx, http.MethodPost, profileID, "create_dataset", map[string]any{}, pgID(opts))
}

// SendConversionEvent reports a conversion back to Meta for a
// click-to-WhatsApp ad.
func (s *WhatsAppService) SendConversionEvent(ctx context.Context, profileID string, event WhatsAppConversionEvent, opts *RequestOptions) (map[string]any, error) {
	reqOpts := []requestOption{withJSON(event), withProfileGroupID(pgID(opts))}

	data, err := s.client.request(ctx, http.MethodPost, "/profiles/"+profileID+"/whatsapp/send_conversion_event", reqOpts...)
	if err != nil {
		return nil, err
	}
	return decodeMap(data)
}

// --- Media ---

// Media downloads an inbound media file (e.g. by a message attachment's
// ExternalID) and returns the raw bytes.
func (s *WhatsAppService) Media(ctx context.Context, profileID, phoneNumberID, mediaID string, opts *RequestOptions) (*WhatsAppMedia, error) {
	params := url.Values{}
	params.Set("phone_number_id", phoneNumberID)
	params.Set("media_id", mediaID)

	data, contentType, err := s.client.requestRaw(ctx, http.MethodGet, "/profiles/"+profileID+"/whatsapp/media", withParams(params), withProfileGroupID(pgID(opts)))
	if err != nil {
		return nil, err
	}
	return &WhatsAppMedia{Data: data, ContentType: contentType}, nil
}

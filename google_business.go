package postproxy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GoogleBusinessService manages a Google Business Profile location's own data
// — hours, attributes, services, menus, action links and photos — separate
// from posting to it.
//
// locationID is the placement id from ProfilesService.Placements. Google's
// payloads are camelCased and passed through untouched as map[string]any, so
// you can PATCH back exactly what you GET. Updates take a fields list naming
// which top-level fields to change (Google's update mask).
type GoogleBusinessService struct {
	client *Client
}

// GoogleBusinessLocationOptions contains options for GoogleBusinessService.Location.
type GoogleBusinessLocationOptions struct {
	// ReadMask is a comma-separated list of fields to read, e.g. "title,phoneNumbers".
	ReadMask       *string
	ProfileGroupID *string
}

// GoogleBusinessCategoriesOptions contains options for
// GoogleBusinessService.AvailableLocationCategories.
type GoogleBusinessCategoriesOptions struct {
	LanguageCode   *string
	Filter         *string
	View           *string
	PageSize       *int
	PageToken      *string
	ProfileGroupID *string
}

// GoogleBusinessHours is the body of GoogleBusinessService.UpdateHours, in
// Google's camelCase shape.
type GoogleBusinessHours struct {
	Fields       []string
	RegularHours map[string]any
	SpecialHours map[string]any
	MoreHours    []map[string]any
}

// GoogleBusinessAvailableAttributesOptions contains options for
// GoogleBusinessService.AvailableAttributes.
type GoogleBusinessAvailableAttributesOptions struct {
	CategoryName   *string
	RegionCode     *string
	LanguageCode   *string
	ShowAll        *bool
	PageSize       *int
	PageToken      *string
	ProfileGroupID *string
}

// GoogleBusinessAttributesOptions contains options for
// GoogleBusinessService.UpdateAttributes.
type GoogleBusinessAttributesOptions struct {
	AttributeMask  []string
	ProfileGroupID *string
}

// GoogleBusinessFieldsOptions contains options for the service list and food
// menu calls.
type GoogleBusinessFieldsOptions struct {
	Fields         []string
	ProfileGroupID *string
}

// GoogleBusinessPageOptions contains paging options for list calls.
type GoogleBusinessPageOptions struct {
	PageSize       *int
	PageToken      *string
	ProfileGroupID *string
}

// GoogleBusinessActionTypesOptions contains options for
// GoogleBusinessService.AvailablePlaceActionTypes.
type GoogleBusinessActionTypesOptions struct {
	LanguageCode   *string
	PageSize       *int
	PageToken      *string
	ProfileGroupID *string
}

// GoogleBusinessPlaceActionLink is a booking / ordering link. PlaceActionType
// is e.g. APPOINTMENT, ONLINE_APPOINTMENT, DINING_RESERVATION, FOOD_ORDERING,
// FOOD_DELIVERY, FOOD_TAKEOUT, SHOP_ONLINE.
type GoogleBusinessPlaceActionLink struct {
	URI             string
	PlaceActionType string
	IsPreferred     *bool
}

// GoogleBusinessMediaOptions contains options for GoogleBusinessService.CreateMedia.
type GoogleBusinessMediaOptions struct {
	// Category is e.g. COVER, PROFILE, LOGO, EXTERIOR, INTERIOR, PRODUCT, ADDITIONAL.
	Category *string
	// MediaFormat is "PHOTO" or "VIDEO".
	MediaFormat    *string
	Description    *string
	ProfileGroupID *string
}

func (s *GoogleBusinessService) path(profileID, action string) string {
	return "/profiles/" + profileID + "/google_business/" + action
}

func (s *GoogleBusinessService) read(ctx context.Context, profileID, action, locationID string, params url.Values, profileGroupID *string) (map[string]any, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("location_id", locationID)

	data, err := s.client.request(ctx, http.MethodGet, s.path(profileID, action), withParams(params), withProfileGroupID(profileGroupID))
	if err != nil {
		return nil, err
	}
	return decodeMap(data)
}

func (s *GoogleBusinessService) write(ctx context.Context, method, profileID, action, locationID string, body map[string]any, profileGroupID *string) (map[string]any, error) {
	payload := map[string]any{"location_id": locationID}
	for k, v := range body {
		payload[k] = v
	}

	data, err := s.client.request(ctx, method, s.path(profileID, action), withJSON(payload), withProfileGroupID(profileGroupID))
	if err != nil {
		return nil, err
	}
	return decodeMap(data)
}

func decodeMap(data []byte) (map[string]any, error) {
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func pgID(opts *RequestOptions) *string {
	if opts == nil {
		return nil
	}
	return opts.ProfileGroupID
}

func setString(params url.Values, key string, value *string) {
	if value != nil {
		params.Set(key, *value)
	}
}

func setInt(params url.Values, key string, value *int) {
	if value != nil {
		params.Set(key, strconv.Itoa(*value))
	}
}

// --- Location ---

// Location reads the location.
func (s *GoogleBusinessService) Location(ctx context.Context, profileID, locationID string, opts *GoogleBusinessLocationOptions) (map[string]any, error) {
	params := url.Values{}
	var pg *string
	if opts != nil {
		setString(params, "read_mask", opts.ReadMask)
		pg = opts.ProfileGroupID
	}
	return s.read(ctx, profileID, "location", locationID, params, pg)
}

// UpdateLocation changes the location. location holds the Google fields to
// change (e.g. "title", "phoneNumbers", "websiteUri"); fields names them.
func (s *GoogleBusinessService) UpdateLocation(ctx context.Context, profileID, locationID string, fields []string, location map[string]any, opts *RequestOptions) (map[string]any, error) {
	body := make(map[string]any, len(location)+1)
	for k, v := range location {
		body[k] = v
	}
	body["fields"] = fields
	return s.write(ctx, http.MethodPatch, profileID, "update_location", locationID, body, pgID(opts))
}

// AvailableLocationCategories lists the business categories available in a region.
func (s *GoogleBusinessService) AvailableLocationCategories(ctx context.Context, profileID, locationID, regionCode string, opts *GoogleBusinessCategoriesOptions) (map[string]any, error) {
	params := url.Values{}
	params.Set("region_code", regionCode)
	var pg *string
	if opts != nil {
		setString(params, "language_code", opts.LanguageCode)
		setString(params, "filter", opts.Filter)
		setString(params, "view", opts.View)
		setInt(params, "page_size", opts.PageSize)
		setString(params, "page_token", opts.PageToken)
		pg = opts.ProfileGroupID
	}
	return s.read(ctx, profileID, "available_location_categories", locationID, params, pg)
}

// UpdateHours changes the location's regular, special or more hours.
func (s *GoogleBusinessService) UpdateHours(ctx context.Context, profileID, locationID string, hours GoogleBusinessHours, opts *RequestOptions) (map[string]any, error) {
	body := map[string]any{}
	if hours.Fields != nil {
		body["fields"] = hours.Fields
	}
	if hours.RegularHours != nil {
		body["regularHours"] = hours.RegularHours
	}
	if hours.SpecialHours != nil {
		body["specialHours"] = hours.SpecialHours
	}
	if hours.MoreHours != nil {
		body["moreHours"] = hours.MoreHours
	}
	return s.write(ctx, http.MethodPatch, profileID, "update_hours", locationID, body, pgID(opts))
}

// --- Attributes ---

// Attributes reads the location's attributes.
func (s *GoogleBusinessService) Attributes(ctx context.Context, profileID, locationID string, opts *RequestOptions) (map[string]any, error) {
	return s.read(ctx, profileID, "attributes", locationID, nil, pgID(opts))
}

// AvailableAttributes lists the attributes the location can set.
func (s *GoogleBusinessService) AvailableAttributes(ctx context.Context, profileID, locationID string, opts *GoogleBusinessAvailableAttributesOptions) (map[string]any, error) {
	params := url.Values{}
	var pg *string
	if opts != nil {
		setString(params, "category_name", opts.CategoryName)
		setString(params, "region_code", opts.RegionCode)
		setString(params, "language_code", opts.LanguageCode)
		if opts.ShowAll != nil {
			params.Set("show_all", strconv.FormatBool(*opts.ShowAll))
		}
		setInt(params, "page_size", opts.PageSize)
		setString(params, "page_token", opts.PageToken)
		pg = opts.ProfileGroupID
	}
	return s.read(ctx, profileID, "available_attributes", locationID, params, pg)
}

// UpdateAttributes changes the location's attributes.
func (s *GoogleBusinessService) UpdateAttributes(ctx context.Context, profileID, locationID string, attributes []map[string]any, opts *GoogleBusinessAttributesOptions) (map[string]any, error) {
	body := map[string]any{"attributes": attributes}
	var pg *string
	if opts != nil {
		if opts.AttributeMask != nil {
			body["attribute_mask"] = opts.AttributeMask
		}
		pg = opts.ProfileGroupID
	}
	return s.write(ctx, http.MethodPatch, profileID, "update_attributes", locationID, body, pg)
}

// --- Services & food menus ---

// ServiceList reads the location's services.
func (s *GoogleBusinessService) ServiceList(ctx context.Context, profileID, locationID string, opts *RequestOptions) (map[string]any, error) {
	return s.read(ctx, profileID, "service_list", locationID, nil, pgID(opts))
}

// UpdateServiceList replaces the location's services.
func (s *GoogleBusinessService) UpdateServiceList(ctx context.Context, profileID, locationID string, serviceItems []map[string]any, opts *GoogleBusinessFieldsOptions) (map[string]any, error) {
	body := map[string]any{"serviceItems": serviceItems}
	var pg *string
	if opts != nil {
		if opts.Fields != nil {
			body["fields"] = opts.Fields
		}
		pg = opts.ProfileGroupID
	}
	return s.write(ctx, http.MethodPatch, profileID, "update_service_list", locationID, body, pg)
}

// FoodMenus reads the location's food menus.
func (s *GoogleBusinessService) FoodMenus(ctx context.Context, profileID, locationID string, opts *GoogleBusinessFieldsOptions) (map[string]any, error) {
	params := url.Values{}
	var pg *string
	if opts != nil {
		if opts.Fields != nil {
			params.Set("fields", strings.Join(opts.Fields, ","))
		}
		pg = opts.ProfileGroupID
	}
	return s.read(ctx, profileID, "food_menus", locationID, params, pg)
}

// UpdateFoodMenus replaces the location's food menus.
func (s *GoogleBusinessService) UpdateFoodMenus(ctx context.Context, profileID, locationID string, menus []map[string]any, opts *GoogleBusinessFieldsOptions) (map[string]any, error) {
	body := map[string]any{"menus": menus}
	var pg *string
	if opts != nil {
		if opts.Fields != nil {
			body["fields"] = opts.Fields
		}
		pg = opts.ProfileGroupID
	}
	return s.write(ctx, http.MethodPatch, profileID, "update_food_menus", locationID, body, pg)
}

// --- Place action links (booking, ordering, ...) ---

// PlaceActionLinks lists the location's place action links.
func (s *GoogleBusinessService) PlaceActionLinks(ctx context.Context, profileID, locationID string, opts *GoogleBusinessPageOptions) (map[string]any, error) {
	params := url.Values{}
	var pg *string
	if opts != nil {
		setInt(params, "page_size", opts.PageSize)
		setString(params, "page_token", opts.PageToken)
		pg = opts.ProfileGroupID
	}
	return s.read(ctx, profileID, "place_action_links", locationID, params, pg)
}

// AvailablePlaceActionTypes lists the place action types the location supports.
func (s *GoogleBusinessService) AvailablePlaceActionTypes(ctx context.Context, profileID, locationID string, opts *GoogleBusinessActionTypesOptions) (map[string]any, error) {
	params := url.Values{}
	var pg *string
	if opts != nil {
		setString(params, "language_code", opts.LanguageCode)
		setInt(params, "page_size", opts.PageSize)
		setString(params, "page_token", opts.PageToken)
		pg = opts.ProfileGroupID
	}
	return s.read(ctx, profileID, "available_place_action_types", locationID, params, pg)
}

func placeActionLinkBody(link GoogleBusinessPlaceActionLink) map[string]any {
	body := map[string]any{}
	if link.URI != "" {
		body["uri"] = link.URI
	}
	if link.PlaceActionType != "" {
		body["place_action_type"] = link.PlaceActionType
	}
	if link.IsPreferred != nil {
		body["is_preferred"] = *link.IsPreferred
	}
	return body
}

// CreatePlaceActionLink adds a booking / ordering link.
func (s *GoogleBusinessService) CreatePlaceActionLink(ctx context.Context, profileID, locationID string, link GoogleBusinessPlaceActionLink, opts *RequestOptions) (map[string]any, error) {
	return s.write(ctx, http.MethodPost, profileID, "create_place_action_link", locationID, placeActionLinkBody(link), pgID(opts))
}

// UpdatePlaceActionLink changes a link. name is the link's resource name from
// PlaceActionLinks; fields names what to change.
func (s *GoogleBusinessService) UpdatePlaceActionLink(ctx context.Context, profileID, locationID, name string, fields []string, link GoogleBusinessPlaceActionLink, opts *RequestOptions) (map[string]any, error) {
	body := placeActionLinkBody(link)
	body["name"] = name
	body["fields"] = fields
	return s.write(ctx, http.MethodPatch, profileID, "update_place_action_link", locationID, body, pgID(opts))
}

// DeletePlaceActionLink removes a link by its resource name.
func (s *GoogleBusinessService) DeletePlaceActionLink(ctx context.Context, profileID, locationID, name string, opts *RequestOptions) (map[string]any, error) {
	params := url.Values{}
	params.Set("location_id", locationID)
	params.Set("name", name)
	return s.remove(ctx, profileID, "delete_place_action_link", params, pgID(opts))
}

// --- Media ---

// Media lists the location's photos and videos.
func (s *GoogleBusinessService) Media(ctx context.Context, profileID, locationID string, opts *GoogleBusinessPageOptions) (map[string]any, error) {
	params := url.Values{}
	var pg *string
	if opts != nil {
		setInt(params, "page_size", opts.PageSize)
		setString(params, "page_token", opts.PageToken)
		pg = opts.ProfileGroupID
	}
	return s.read(ctx, profileID, "media", locationID, params, pg)
}

// CreateMedia adds a photo or video from a public URL.
func (s *GoogleBusinessService) CreateMedia(ctx context.Context, profileID, locationID, mediaURL string, opts *GoogleBusinessMediaOptions) (map[string]any, error) {
	body := map[string]any{"media_url": mediaURL}
	var pg *string
	if opts != nil {
		if opts.Category != nil {
			body["category"] = *opts.Category
		}
		if opts.MediaFormat != nil {
			body["media_format"] = *opts.MediaFormat
		}
		if opts.Description != nil {
			body["description"] = *opts.Description
		}
		pg = opts.ProfileGroupID
	}
	return s.write(ctx, http.MethodPost, profileID, "create_media", locationID, body, pg)
}

// DeleteMedia removes a media item. mediaName is the item's resource name from Media.
func (s *GoogleBusinessService) DeleteMedia(ctx context.Context, profileID, locationID, mediaName string, opts *RequestOptions) (map[string]any, error) {
	params := url.Values{}
	params.Set("location_id", locationID)
	params.Set("media_name", mediaName)
	return s.remove(ctx, profileID, "delete_media", params, pgID(opts))
}

func (s *GoogleBusinessService) remove(ctx context.Context, profileID, action string, params url.Values, profileGroupID *string) (map[string]any, error) {
	data, err := s.client.request(ctx, http.MethodDelete, s.path(profileID, action), withParams(params), withProfileGroupID(profileGroupID))
	if err != nil {
		return nil, err
	}
	return decodeMap(data)
}

// --- Search keywords ---

// GoogleBusinessSearchKeywordsOptions contains options for
// GoogleBusinessService.SearchKeywords. Months are "YYYY-MM"; the range is at
// most 12 months.
type GoogleBusinessSearchKeywordsOptions struct {
	FromMonth      *string
	ToMonth        *string
	ProfileGroupID *string
}

// SearchKeyword is one search term that showed the location.
type SearchKeyword struct {
	Keyword string `json:"keyword"`
	// Impressions is nil for low-volume keywords; see ImpressionsThreshold.
	Impressions *int `json:"impressions"`
	// ImpressionsThreshold is the upper bound when Google reports no exact count.
	ImpressionsThreshold *int `json:"impressions_threshold"`
}

// SearchKeywordsMonth holds the keywords for one month ("YYYY-MM").
type SearchKeywordsMonth struct {
	Month    string          `json:"month"`
	Keywords []SearchKeyword `json:"keywords"`
}

// SearchKeywordsResponse is returned by GoogleBusinessService.SearchKeywords.
type SearchKeywordsResponse struct {
	Months []SearchKeywordsMonth `json:"months"`
}

// SearchKeywords returns the searches that showed the location in Google
// Search and Maps, month by month. Unlike the rest of this service, the
// response is PostProxy-shaped (snake_case) and typed.
func (s *GoogleBusinessService) SearchKeywords(ctx context.Context, profileID, locationID string, opts *GoogleBusinessSearchKeywordsOptions) (*SearchKeywordsResponse, error) {
	params := url.Values{}
	params.Set("location_id", locationID)
	var pg *string
	if opts != nil {
		setString(params, "from_month", opts.FromMonth)
		setString(params, "to_month", opts.ToMonth)
		pg = opts.ProfileGroupID
	}

	data, err := s.client.request(ctx, http.MethodGet, s.path(profileID, "search_keywords"), withParams(params), withProfileGroupID(pg))
	if err != nil {
		return nil, err
	}

	var result SearchKeywordsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

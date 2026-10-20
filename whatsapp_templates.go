package postproxy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// WhatsAppTemplatesService manages WhatsApp message templates, which live on
// the WhatsApp Business account (the profile). A template is addressed by its
// PostProxy id, Meta's id, or its name — pass Language with a name that exists
// in several languages.
type WhatsAppTemplatesService struct {
	client *Client
}

// WhatsAppTemplateInput is the body of WhatsAppTemplatesService.Create.
type WhatsAppTemplateInput struct {
	Name string `json:"name"`
	// Category is MARKETING, UTILITY or AUTHENTICATION.
	Category string `json:"category"`
	Language string `json:"language"`
	// ParameterFormat is "POSITIONAL" or "NAMED".
	ParameterFormat string `json:"parameter_format,omitempty"`
	// Components are Meta's component objects (HEADER, BODY, FOOTER, BUTTONS).
	Components []map[string]any `json:"components,omitempty"`
	// LibraryTemplateName creates the template from Meta's template library
	// instead of Components.
	LibraryTemplateName         string           `json:"library_template_name,omitempty"`
	LibraryTemplateBodyInputs   map[string]any   `json:"library_template_body_inputs,omitempty"`
	LibraryTemplateButtonInputs []map[string]any `json:"library_template_button_inputs,omitempty"`
	MessageSendTTLSeconds       *int             `json:"message_send_ttl_seconds,omitempty"`
}

// WhatsAppTemplateUpdate is the body of WhatsAppTemplatesService.Update.
type WhatsAppTemplateUpdate struct {
	Components            []map[string]any `json:"components,omitempty"`
	MessageSendTTLSeconds *int             `json:"message_send_ttl_seconds,omitempty"`
}

// WhatsAppTemplateListOptions contains options for WhatsAppTemplatesService.List.
type WhatsAppTemplateListOptions struct {
	Name     *string
	Language *string
	Status   *string
	// Refresh re-reads the templates from Meta first.
	Refresh        *bool
	ProfileGroupID *string
}

// WhatsAppTemplateOptions picks a template's language and profile group.
type WhatsAppTemplateOptions struct {
	Language       *string
	ProfileGroupID *string
}

// WhatsAppTemplateLibraryOptions contains options for WhatsAppTemplatesService.Library.
type WhatsAppTemplateLibraryOptions struct {
	Name           *string
	Language       *string
	ProfileGroupID *string
}

func (s *WhatsAppTemplatesService) path(profileID, id string) string {
	base := "/profiles/" + profileID + "/templates"
	if id == "" {
		return base
	}
	return base + "/" + url.PathEscape(id)
}

func templateOptions(opts *WhatsAppTemplateOptions) []requestOption {
	if opts == nil {
		return nil
	}
	reqOpts := []requestOption{withProfileGroupID(opts.ProfileGroupID)}
	if opts.Language != nil {
		params := url.Values{}
		params.Set("language", *opts.Language)
		reqOpts = append(reqOpts, withParams(params))
	}
	return reqOpts
}

// List returns the profile's templates.
func (s *WhatsAppTemplatesService) List(ctx context.Context, profileID string, opts *WhatsAppTemplateListOptions) (*ListResponse[WhatsAppTemplate], error) {
	var reqOpts []requestOption
	if opts != nil {
		params := url.Values{}
		setString(params, "name", opts.Name)
		setString(params, "language", opts.Language)
		setString(params, "status", opts.Status)
		if opts.Refresh != nil {
			params.Set("refresh", strconv.FormatBool(*opts.Refresh))
		}
		if len(params) > 0 {
			reqOpts = append(reqOpts, withParams(params))
		}
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, s.path(profileID, ""), reqOpts...)
	if err != nil {
		return nil, err
	}

	var result ListResponse[WhatsAppTemplate]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a single template by PostProxy id, Meta id or name.
func (s *WhatsAppTemplatesService) Get(ctx context.Context, profileID, id string, opts *WhatsAppTemplateOptions) (*WhatsAppTemplate, error) {
	data, err := s.client.request(ctx, http.MethodGet, s.path(profileID, id), templateOptions(opts)...)
	if err != nil {
		return nil, err
	}

	var result WhatsAppTemplate
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create submits a template to Meta for review; follow Status via Get or the
// message_template.status_updated webhook.
func (s *WhatsAppTemplatesService) Create(ctx context.Context, profileID string, template WhatsAppTemplateInput, opts *RequestOptions) (*WhatsAppTemplate, error) {
	reqOpts := []requestOption{withJSON(template)}
	if opts != nil {
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodPost, s.path(profileID, ""), reqOpts...)
	if err != nil {
		return nil, err
	}

	var result WhatsAppTemplate
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update changes a template's components or message TTL.
func (s *WhatsAppTemplatesService) Update(ctx context.Context, profileID, id string, changes WhatsAppTemplateUpdate, opts *WhatsAppTemplateOptions) (*WhatsAppTemplate, error) {
	reqOpts := append(templateOptions(opts), withJSON(changes))

	data, err := s.client.request(ctx, http.MethodPatch, s.path(profileID, id), reqOpts...)
	if err != nil {
		return nil, err
	}

	var result WhatsAppTemplate
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a template. Deleting by name without Language removes every
// language of the template.
func (s *WhatsAppTemplatesService) Delete(ctx context.Context, profileID, id string, opts *WhatsAppTemplateOptions) (map[string]any, error) {
	data, err := s.client.request(ctx, http.MethodDelete, s.path(profileID, id), templateOptions(opts)...)
	if err != nil {
		return nil, err
	}
	return decodeMap(data)
}

// Library browses Meta's ready-made template library.
func (s *WhatsAppTemplatesService) Library(ctx context.Context, profileID string, opts *WhatsAppTemplateLibraryOptions) (map[string]any, error) {
	var reqOpts []requestOption
	if opts != nil {
		params := url.Values{}
		setString(params, "name", opts.Name)
		setString(params, "language", opts.Language)
		if len(params) > 0 {
			reqOpts = append(reqOpts, withParams(params))
		}
		reqOpts = append(reqOpts, withProfileGroupID(opts.ProfileGroupID))
	}

	data, err := s.client.request(ctx, http.MethodGet, s.path(profileID, "")+"/library", reqOpts...)
	if err != nil {
		return nil, err
	}
	return decodeMap(data)
}

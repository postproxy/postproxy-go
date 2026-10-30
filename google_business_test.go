package postproxy

import (
	"context"
	"testing"
)

const gbpBase = "/api/profiles/prof_gbp/google_business"

func TestGoogleBusinessLocation(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"title": "Cafe"})

	mask := "title,phoneNumbers"
	result, err := c.GoogleBusiness.Location(context.Background(), "prof_gbp", "loc_1", &GoogleBusinessLocationOptions{ReadMask: &mask})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["title"] != "Cafe" {
		t.Errorf("unexpected result: %v", result)
	}

	req := requests()[0]
	if req.Path != gbpBase+"/location" {
		t.Errorf("unexpected path: %s", req.Path)
	}
	if req.Query.Get("location_id") != "loc_1" || req.Query.Get("read_mask") != "title,phoneNumbers" {
		t.Errorf("unexpected query: %v", req.Query)
	}
}

func TestGoogleBusinessUpdateLocation(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{})

	_, err := c.GoogleBusiness.UpdateLocation(context.Background(), "prof_gbp", "loc_1", []string{"websiteUri"}, map[string]any{"websiteUri": "https://acme.test"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := requests()[0]
	if req.Method != "PATCH" || req.Path != gbpBase+"/update_location" {
		t.Errorf("unexpected request: %s %s", req.Method, req.Path)
	}
	assertJSON(t, req.JSON(t), `{"location_id":"loc_1","websiteUri":"https://acme.test","fields":["websiteUri"]}`)
}

func TestGoogleBusinessAvailableLocationCategories(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{"categories": []any{}})

	pageSize := 50
	if _, err := c.GoogleBusiness.AvailableLocationCategories(context.Background(), "prof_gbp", "loc_1", "US", &GoogleBusinessCategoriesOptions{PageSize: &pageSize}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	q := requests()[0].Query
	if q.Get("region_code") != "US" || q.Get("page_size") != "50" {
		t.Errorf("unexpected query: %v", q)
	}
	if q.Has("filter") {
		t.Error("expected no filter param")
	}
}

func TestGoogleBusinessUpdateHours(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{})

	regularHours := map[string]any{"periods": []any{map[string]any{
		"openDay": "MONDAY", "openTime": map[string]any{"hours": 9},
		"closeDay": "MONDAY", "closeTime": map[string]any{"hours": 17},
	}}}
	_, err := c.GoogleBusiness.UpdateHours(context.Background(), "prof_gbp", "loc_1", GoogleBusinessHours{
		Fields:       []string{"regularHours"},
		RegularHours: regularHours,
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertJSON(t, requests()[0].JSON(t), `{"location_id":"loc_1","fields":["regularHours"],"regularHours":{"periods":[{"openDay":"MONDAY","openTime":{"hours":9},"closeDay":"MONDAY","closeTime":{"hours":17}}]}}`)
}

func TestGoogleBusinessUpdateAttributesServicesMenus(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{})
	ctx := context.Background()

	_, _ = c.GoogleBusiness.UpdateAttributes(ctx, "prof_gbp", "loc_1", []map[string]any{{"name": "attributes/has_wifi", "values": []any{true}}}, &GoogleBusinessAttributesOptions{AttributeMask: []string{"attributes/has_wifi"}})
	_, _ = c.GoogleBusiness.UpdateServiceList(ctx, "prof_gbp", "loc_1", []map[string]any{{"freeFormServiceItem": map[string]any{}}}, nil)
	_, _ = c.GoogleBusiness.UpdateFoodMenus(ctx, "prof_gbp", "loc_1", []map[string]any{{"labels": []any{}}}, &GoogleBusinessFieldsOptions{Fields: []string{"menus"}})

	reqs := requests()
	if reqs[0].Path != gbpBase+"/update_attributes" {
		t.Errorf("unexpected path: %s", reqs[0].Path)
	}
	assertJSON(t, reqs[0].JSON(t), `{"location_id":"loc_1","attributes":[{"name":"attributes/has_wifi","values":[true]}],"attribute_mask":["attributes/has_wifi"]}`)
	assertJSON(t, reqs[1].JSON(t), `{"location_id":"loc_1","serviceItems":[{"freeFormServiceItem":{}}]}`)
	assertJSON(t, reqs[2].JSON(t), `{"location_id":"loc_1","menus":[{"labels":[]}],"fields":["menus"]}`)
}

func TestGoogleBusinessReads(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{})
	ctx := context.Background()

	showAll := true
	_, _ = c.GoogleBusiness.Attributes(ctx, "prof_gbp", "loc_1", nil)
	_, _ = c.GoogleBusiness.AvailableAttributes(ctx, "prof_gbp", "loc_1", &GoogleBusinessAvailableAttributesOptions{ShowAll: &showAll})
	_, _ = c.GoogleBusiness.ServiceList(ctx, "prof_gbp", "loc_1", nil)
	_, _ = c.GoogleBusiness.FoodMenus(ctx, "prof_gbp", "loc_1", &GoogleBusinessFieldsOptions{Fields: []string{"menus", "name"}})

	reqs := requests()
	for i, want := range []string{"/attributes", "/available_attributes", "/service_list", "/food_menus"} {
		if reqs[i].Method != "GET" || reqs[i].Path != gbpBase+want {
			t.Errorf("request %d: unexpected %s %s", i, reqs[i].Method, reqs[i].Path)
		}
	}
	if reqs[1].Query.Get("show_all") != "true" {
		t.Errorf("expected show_all=true, got %v", reqs[1].Query)
	}
	if reqs[3].Query.Get("fields") != "menus,name" {
		t.Errorf("expected fields=menus,name, got %v", reqs[3].Query)
	}
}

func TestGoogleBusinessPlaceActionLinks(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{})
	ctx := context.Background()

	lang := "en"
	preferred := true
	name := "locations/1/placeActionLinks/2"
	_, _ = c.GoogleBusiness.PlaceActionLinks(ctx, "prof_gbp", "loc_1", nil)
	_, _ = c.GoogleBusiness.AvailablePlaceActionTypes(ctx, "prof_gbp", "loc_1", &GoogleBusinessActionTypesOptions{LanguageCode: &lang})
	_, _ = c.GoogleBusiness.CreatePlaceActionLink(ctx, "prof_gbp", "loc_1", GoogleBusinessPlaceActionLink{URI: "https://book.test", PlaceActionType: "APPOINTMENT", IsPreferred: &preferred}, nil)
	_, _ = c.GoogleBusiness.UpdatePlaceActionLink(ctx, "prof_gbp", "loc_1", name, []string{"uri"}, GoogleBusinessPlaceActionLink{URI: "https://book2.test"}, nil)
	_, _ = c.GoogleBusiness.DeletePlaceActionLink(ctx, "prof_gbp", "loc_1", name, nil)

	reqs := requests()
	if reqs[0].Method != "GET" {
		t.Errorf("expected GET, got %s", reqs[0].Method)
	}
	if reqs[1].Query.Get("language_code") != "en" {
		t.Errorf("unexpected query: %v", reqs[1].Query)
	}
	if reqs[2].Method != "POST" {
		t.Errorf("expected POST, got %s", reqs[2].Method)
	}
	assertJSON(t, reqs[2].JSON(t), `{"location_id":"loc_1","uri":"https://book.test","place_action_type":"APPOINTMENT","is_preferred":true}`)
	if reqs[3].Method != "PATCH" {
		t.Errorf("expected PATCH, got %s", reqs[3].Method)
	}
	assertJSON(t, reqs[3].JSON(t), `{"location_id":"loc_1","name":"locations/1/placeActionLinks/2","fields":["uri"],"uri":"https://book2.test"}`)
	if reqs[4].Method != "DELETE" || reqs[4].Path != gbpBase+"/delete_place_action_link" {
		t.Errorf("unexpected request: %s %s", reqs[4].Method, reqs[4].Path)
	}
	if reqs[4].Query.Get("name") != name || reqs[4].Query.Get("location_id") != "loc_1" {
		t.Errorf("unexpected query: %v", reqs[4].Query)
	}
}

func TestGoogleBusinessMedia(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{})
	ctx := context.Background()

	pageSize := 10
	category, format := "COVER", "PHOTO"
	_, _ = c.GoogleBusiness.Media(ctx, "prof_gbp", "loc_1", &GoogleBusinessPageOptions{PageSize: &pageSize})
	_, _ = c.GoogleBusiness.CreateMedia(ctx, "prof_gbp", "loc_1", "https://img.test/a.jpg", &GoogleBusinessMediaOptions{Category: &category, MediaFormat: &format})
	_, _ = c.GoogleBusiness.DeleteMedia(ctx, "prof_gbp", "loc_1", "accounts/1/locations/1/media/9", nil)

	reqs := requests()
	if reqs[0].Query.Get("page_size") != "10" {
		t.Errorf("unexpected query: %v", reqs[0].Query)
	}
	assertJSON(t, reqs[1].JSON(t), `{"location_id":"loc_1","media_url":"https://img.test/a.jpg","category":"COVER","media_format":"PHOTO"}`)
	if reqs[2].Method != "DELETE" || reqs[2].Query.Get("media_name") != "accounts/1/locations/1/media/9" {
		t.Errorf("unexpected request: %s %v", reqs[2].Method, reqs[2].Query)
	}
}

func TestGoogleBusinessSearchKeywords(t *testing.T) {
	c, requests := newRecordingClient(t, map[string]any{
		"months": []any{map[string]any{
			"month": "2026-07",
			"keywords": []any{
				map[string]any{"keyword": "coffee near me", "impressions": 245, "impressions_threshold": nil},
				map[string]any{"keyword": "acme coffee", "impressions": nil, "impressions_threshold": 15},
			},
		}},
	})

	result, err := c.GoogleBusiness.SearchKeywords(context.Background(), "prof_gbp", "loc_1", &GoogleBusinessSearchKeywordsOptions{
		FromMonth: ptr("2026-07"), ToMonth: ptr("2026-09"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	kw := result.Months[0].Keywords
	if *kw[0].Impressions != 245 || kw[1].Impressions != nil || *kw[1].ImpressionsThreshold != 15 {
		t.Errorf("unexpected keywords: %+v", kw)
	}

	req := requests()[0]
	if req.Method != "GET" || req.Path != gbpBase+"/search_keywords" {
		t.Errorf("unexpected request: %s %s", req.Method, req.Path)
	}
	q := req.Query
	if q.Get("location_id") != "loc_1" || q.Get("from_month") != "2026-07" || q.Get("to_month") != "2026-09" {
		t.Errorf("unexpected query: %v", q)
	}
}

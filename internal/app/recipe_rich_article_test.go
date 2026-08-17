package app

import (
	"encoding/json"
	webassets "geliumui/site/web"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRichArticleRouteStructureAndMetadata(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/rich-article", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d", res.Code)
	}
	body := res.Body.String()
	if strings.Count(body, "<h1") != 1 {
		t.Fatalf("want exactly one h1, got %d", strings.Count(body, "<h1"))
	}
	for _, marker := range []string{"recipe-rich-article-eyebrow", "recipe-rich-article-lead", "8 min read", "<h2", "<h3", "<ol>", "<ul>", "<blockquote", "<code>", "<pre>", "ui-alert", "<picture>", "alt=\"", "poster=\"/static/rich-article-image.svg\"", "kind=\"captions\"", "<audio controls", "id=\"transcript\"", "Embed unavailable", "ui-data-table-table", "recipe-rich-article-feed", "Loading the next related story", "No related stories", "role=\"alert\"", `"@type":"BlogPosting"`, "recipe-rich-article-layout"} {
		if !strings.Contains(body, marker) {
			t.Errorf("missing %q", marker)
		}
	}
	if !strings.Contains(body, `<meta name="robots" content="noindex, nofollow">`) {
		t.Error("missing noindex")
	}
}

func TestRichArticleJSONLDIsValidArticle(t *testing.T) {
	ld := richArticleJSONLD()
	var value map[string]any
	if err := json.Unmarshal([]byte(ld), &value); err != nil {
		t.Fatal(err)
	}
	if value["@type"] != "BlogPosting" {
		t.Fatalf("type = %v", value["@type"])
	}
	if value["url"] != siteBaseURL+"/recipes/rich-article" {
		t.Fatalf("url = %v", value["url"])
	}
}

func TestRichArticleRouteIsGETOnlyAndNavigationIsStable(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut} {
		res := httptest.NewRecorder()
		New().ServeHTTP(res, httptest.NewRequest(method, "/recipes/rich-article", nil))
		if res.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s status = %d", method, res.Code)
		}
	}
	groups := docsNavFor("", "", "")
	found := false
	for _, group := range groups.Groups {
		for _, link := range group.Links {
			if link.Href == "/recipes/rich-article" {
				found = true
			}
		}
	}
	if !found {
		t.Error("rich article missing from docs navigation")
	}
	packBytes, err := webassets.Assets.ReadFile("static/llms-ux.txt")
	if err != nil {
		t.Fatal(err)
	}
	pack := string(packBytes)
	for _, id := range []string{"RICH-ARTICLE", "RICH-STRUCTURE", "RICH-MEDIA", "RICH-STATES", "RICH-DATA", "RICH-SEO"} {
		if !strings.Contains(pack, id) {
			t.Errorf("missing LLM id %s", id)
		}
	}
}

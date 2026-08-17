package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSEOIndexabilityPolicyDrivesMetaAndSitemap(t *testing.T) {
	for _, route := range []string{"/recipes/admin-resource", "/demo/whatsapp", "/examples/newsletter", "/components/dialog/confirm"} {
		if isIndexablePath(route) || indexabilityPolicy(route) != "noindex, nofollow" {
			t.Fatalf("%s must be excluded by the shared indexability policy", route)
		}
		body := getOKBody(t, route)
		if !strings.Contains(body, `<meta name="robots" content="noindex, nofollow">`) {
			t.Errorf("%s must render noindex", route)
		}
		for _, p := range sitemapPaths() {
			if p == route {
				t.Errorf("%s must not be in sitemap paths", route)
			}
		}
	}
}

func TestBlogSitemapAndBlogPostingRegistryConsistency(t *testing.T) {
	body := getOKBody(t, "/sitemap.xml")
	for _, post := range blogPosts {
		path := "/blog/" + post.Slug
		if !strings.Contains(body, "<loc>"+siteBaseURL+path+"</loc>") {
			t.Errorf("sitemap missing %s", path)
		}
		page := getOKBody(t, path)
		var ld map[string]any
		if err := json.Unmarshal([]byte(extractJSONLD(t, page)), &ld); err != nil {
			t.Fatalf("%s JSON-LD: %v", path, err)
		}
		if ld["@type"] != "BlogPosting" || ld["headline"] != post.Title || ld["datePublished"] != post.Date {
			t.Errorf("%s BlogPosting does not match registry: %#v", path, ld)
		}
		if _, ok := ld["dateModified"]; ok {
			t.Errorf("%s must not invent dateModified", path)
		}
		if ld["mainEntityOfPage"] != siteBaseURL+path {
			t.Errorf("%s mainEntityOfPage mismatch", path)
		}
	}
	if !strings.Contains(body, "<loc>"+siteBaseURL+"/blog</loc>") {
		t.Error("sitemap missing /blog")
	}
}

func TestDocsJSONLDAndSocialMetadataCoherence(t *testing.T) {
	for _, route := range []string{"/docs", "/docs/seo", "/docs/aeo"} {
		res := httptest.NewRecorder()
		New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, route, nil))
		if res.Code != http.StatusOK {
			t.Fatalf("%s status = %d", route, res.Code)
		}
		body := res.Body.String()
		var ld map[string]any
		if err := json.Unmarshal([]byte(extractJSONLD(t, body)), &ld); err != nil {
			t.Fatalf("%s JSON-LD: %v", route, err)
		}
		wantType := "WebPage"
		if route == "/docs" {
			wantType = "CollectionPage"
		}
		if ld["@type"] != wantType {
			t.Errorf("%s @type = %v, want %s", route, ld["@type"], wantType)
		}
		canonical := siteBaseURL + route
		for _, tag := range []string{
			`<link rel="canonical" href="` + canonical + `">`,
			`<meta property="og:url" content="` + canonical + `">`,
			`<meta name="twitter:url" content="` + canonical + `">`,
		} {
			if !strings.Contains(body, tag) {
				t.Errorf("%s missing coherent tag %s", route, tag)
			}
		}
		if ld["url"] != canonical {
			t.Errorf("%s JSON-LD url mismatch", route)
		}
	}
}

func TestTwitterMetadataMatchesOpenGraph(t *testing.T) {
	body := getOKBody(t, "/blog/genesis")
	for _, pair := range []string{
		`<meta name="twitter:title" content="How Gelium UI was born · Gelium UI">`,
		`<meta name="twitter:description" content="Gelium UI started as a side project and grew into a contract-tested component library for server-rendered Go apps.">`,
		`<meta name="twitter:url" content="https://gelium-ui.example/blog/genesis">`,
		`<meta name="twitter:card" content="summary_large_image">`,
	} {
		if !strings.Contains(body, pair) {
			t.Errorf("blog social metadata missing %s", pair)
		}
	}
}

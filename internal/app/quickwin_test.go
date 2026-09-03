package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestDocsSearchQHonoredServerSide proves the 0-JS search fallback actually
// searches: GET /docs?q=<term> renders a server-built results section from the
// same nav model as the client-side index (docs_search_test.go documents the
// fallback as a live GET form; the hub now reads q instead of ignoring it).
func TestDocsSearchQHonoredServerSide(t *testing.T) {
	body := getOKBody(t, "/docs?q=table")
	for _, contract := range []string{
		`Search results for “table”`,
		`/components/data-table`,
		`Data`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("/docs?q=table missing search result contract %q", contract)
		}
	}
}

// TestDocsSearchQNoMatch proves an empty result set is honest and guided
// instead of pretending the search ran.
func TestDocsSearchQNoMatch(t *testing.T) {
	body := getOKBody(t, "/docs?q=zzzznonexistentterm")
	if !strings.Contains(body, "No matches for") {
		t.Error("/docs?q=zzzznonexistentterm must render a no-matches section")
	}
}

// TestErrorPageDocsShellNoindex proves doc-route failures render through the
// docs shell (sidebar + topbar search) and never carry an indexable robots
// meta: a 404 with "index, follow" is a soft-404 defect (SEO §4, §16).
func TestErrorPageDocsShellNoindex(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/nope", nil))
	if res.Code != http.StatusNotFound {
		t.Fatalf("GET /docs/nope status = %d, want %d", res.Code, http.StatusNotFound)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<form class="docs-search" method="get" action="/docs" role="search">`,
		`name="robots" content="noindex, nofollow"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("/docs/nope missing shell/noindex contract %q", contract)
		}
	}
}

// TestTemplatePagesInSitemap proves /docs/templates/* are indexable AND listed
// in the sitemap: they were routed and canonical before but invisible to
// crawler discovery (SEO §4).
func TestTemplatePagesInSitemap(t *testing.T) {
	body := getOKBody(t, "/sitemap.xml")
	for _, path := range []string{"/docs/templates/product", "/docs/templates/design"} {
		if !strings.Contains(body, path) {
			t.Errorf("sitemap must include %s", path)
		}
	}
}

// TestMediaFixturesServe proves the /docs/media demo fixtures exist and serve
// with correct content types: the page referenced two assets that were absent
// from the embedded FS (live 404s on the demo).
func TestMediaFixturesServe(t *testing.T) {
	svg := getOKBody(t, "/static/media/editorial-placeholder.svg")
	if !strings.Contains(svg, "<svg") {
		t.Error("editorial-placeholder.svg must be a real SVG document")
	}
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/static/media/empty-audio.ogg", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("GET /static/media/empty-audio.ogg status = %d, want %d", res.Code, http.StatusOK)
	}
	if ct := res.Header().Get("Content-Type"); !strings.HasPrefix(ct, "audio/") {
		t.Errorf("empty-audio.ogg Content-Type = %q, want audio/*", ct)
	}
}

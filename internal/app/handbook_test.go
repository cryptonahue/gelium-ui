package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// handbookRoute is one destination in the Gelium Handbook IA. Contract is a
// content marker unique to the page — the thing that makes the page real
// rather than a stub; H1 is the rendered heading level 1 of that page.
type handbookRoute struct {
	path     string
	label    string
	h1       string
	contract string
}

// handbookRoutes is the single source of truth for the five Handbook pages
// (four Handbook sections + Design Principles), mirroring the sidebar group.
var handbookRoutes = []handbookRoute{
	{path: "/docs/themes", label: "Themes", h1: "Themes", contract: "?theme=basecoat"},
	{path: "/docs/tokens", label: "Tokens", h1: "Tokens", contract: "--ui-color-primary"},
	{path: "/docs/server-contracts", label: "Server contracts", h1: "Server contracts", contract: "HX-Trigger"},
	{path: "/docs/accessibility", label: "Accessibility", h1: "Accessibility", contract: "forced-colors"},
	{path: "/docs/principles", label: "Design principles", h1: "Design principles", contract: "Native semantics first"},
}

// TestHandbookPagesRender proves every Handbook destination returns 200 under
// the docs shell, renders its own content (h1 + a page-specific contract
// marker), and is reachable from the Handbook sidebar group.
func TestHandbookPagesRender(t *testing.T) {
	for _, hb := range handbookRoutes {
		t.Run(hb.path, func(t *testing.T) {
			res := httptest.NewRecorder()
			New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, hb.path, nil))
			if res.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
			}
			body := res.Body.String()
			if !strings.Contains(body, "<h1>"+hb.h1+"</h1>") {
				t.Errorf("page must render <h1>%s</h1>", hb.h1)
			}
			if !strings.Contains(body, hb.contract) {
				t.Errorf("page is missing its content contract %q", hb.contract)
			}
			if !strings.Contains(body, `class="docs-nav-group-label">Handbook`) {
				t.Error("sidebar must include a Handbook nav group")
			}
			if !strings.Contains(body, `href="`+hb.path+`"`) {
				t.Errorf("sidebar must link to %s", hb.path)
			}
		})
	}
}

// TestDocsNavForHandbookGroup proves the nav model carries exactly the five
// Handbook destinations with the right labels and correct current marking.
func TestDocsNavForHandbookGroup(t *testing.T) {
	nav := docsNavFor("/docs/tokens", "", "")
	var handbook *docsNavGroup
	for i := range nav.Groups {
		if nav.Groups[i].Title == "Handbook" {
			handbook = &nav.Groups[i]
			break
		}
	}
	if handbook == nil {
		t.Fatal("docsNavFor must include a Handbook group")
	}
	if len(handbook.Links) != len(handbookRoutes) {
		t.Fatalf("Handbook links = %d, want %d", len(handbook.Links), len(handbookRoutes))
	}
	for _, hb := range handbookRoutes {
		found := false
		for _, link := range handbook.Links {
			if link.Path != hb.path {
				continue
			}
			found = true
			if link.Label != hb.label {
				t.Errorf("Handbook link %s label = %q, want %q", hb.path, link.Label, hb.label)
			}
			if link.Path == "/docs/tokens" && !link.Current {
				t.Error("Tokens must be the current Handbook link on /docs/tokens")
			}
			if link.Path != "/docs/tokens" && link.Current {
				t.Errorf("peer %s must not be current on /docs/tokens", link.Path)
			}
		}
		if !found {
			t.Errorf("Handbook group is missing %s", hb.path)
		}
	}
}

// TestHandbookPagesInSitemap proves every Handbook page is an indexable route
// in the generated sitemap (SEO contract: one <url> per public docs page).
func TestHandbookPagesInSitemap(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, hb := range handbookRoutes {
		want := "<loc>https://gelium-ui.example" + hb.path + "</loc>"
		if !strings.Contains(body, want) {
			t.Errorf("sitemap is missing %q", want)
		}
	}
}

// TestDocsIndexListsHandbook proves the /docs hub leads readers to the
// Handbook: a Handbook section with every page linked from the hub body.
func TestDocsIndexListsHandbook(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, "<h2>Handbook</h2>") {
		t.Error("docs hub must include a Handbook section")
	}
	for _, hb := range handbookRoutes {
		if !strings.Contains(body, `href="`+hb.path+`"`) {
			t.Errorf("docs hub must link to %s", hb.path)
		}
	}
}

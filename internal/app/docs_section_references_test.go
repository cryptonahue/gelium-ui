package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func TestDocsSectionReferencesIndex(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/section-references", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		">Section references</h1>",
		"not a moodboard",
		"How to use this catalog",
		"Official source",
		"Gelium remake",
		"Keep / adapt",
		"Browse by type",
		`method="get" action="/docs/section-references"`,
		`class="ui-select ui-select-outlined"`,
		`<label class="ui-select-label" for="section-ref-type">Type</label>`,
		`class="ui-button ui-button-outline"`,
		">Showing 6 results.</p>",
		`<nav aria-label="Section reference results">`,
		`class="ui-list"`,
		`class="ui-list-item ui-list-item--two-line"`,
		`class="ui-list-item-link"`,
		`class="ui-list-item-text"`,
		`class="ui-list-item-headline"`,
		`class="ui-list-item-supporting"`,
		`class="ui-list-item-icon ui-list-item-icon--end"`,
		`href="/docs/section-references/article"`,
		`href="/docs/section-references/404-vercel"`,
		"Vercel",
		">Apply</span></button>",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("index missing %q", contract)
		}
	}
	if !strings.Contains(body, `href="/docs/section-references"`) {
		t.Error("sidebar must link to /docs/section-references")
	}
	if count := strings.Count(body, `<h1`); count != 1 {
		t.Errorf("index must render exactly one h1, got %d", count)
	}
	resultsStart := strings.Index(body, `<nav aria-label="Section reference results">`)
	if resultsStart < 0 {
		t.Fatal("index results navigation is required")
	}
	results := body[resultsStart:]
	if count := strings.Count(results, `class="ui-list-item ui-list-item--two-line"`); count != len(sectionRefCatalog) {
		t.Errorf("index result rows = %d, want %d", count, len(sectionRefCatalog))
	}
	if strings.Contains(body, "\n- [") {
		t.Error("index must not render bare Markdown bullet results")
	}
}

func TestDocsSectionReferencesClosedFilters(t *testing.T) {
	tests := []struct {
		name      string
		typeValue string
		path      string
	}{
		{name: "auth", typeValue: "auth", path: "/docs/section-references/auth-register"},
		{name: "faq", typeValue: "faq", path: "/docs/section-references/faq-vercel"},
		{name: "hero", typeValue: "hero", path: "/docs/section-references/hero-linear"},
		{name: "pricing", typeValue: "pricing", path: "/docs/section-references/pricing-linear"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/section-references?type="+tt.typeValue, nil))
			if res.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
			}
			body := res.Body.String()
			if !strings.Contains(body, `href="`+tt.path+`"`) {
				t.Errorf("%s filter must list %s", tt.typeValue, tt.path)
			}
			if !strings.Contains(body, `option value="`+tt.typeValue+`" selected`) {
				t.Errorf("%s option must be selected", tt.typeValue)
			}
			if strings.Contains(body, `href="/docs/section-references/article"`) {
				t.Errorf("%s filter must not list the article ficha", tt.typeValue)
			}
			resultsStart := strings.Index(body, `<nav aria-label="Section reference results">`)
			if resultsStart < 0 {
				t.Fatal("filtered results navigation is required")
			}
			results := body[resultsStart:]
			if count := strings.Count(results, `class="ui-list-item ui-list-item--two-line"`); count != 1 {
				t.Errorf("%s filter rows = %d, want 1", tt.typeValue, count)
			}
			if !strings.Contains(body, ">Showing 1 result for "+tt.typeValue+".</p>") {
				t.Errorf("%s filter must report its result count", tt.typeValue)
			}
		})
	}
}

func TestDocsSectionReferencesFilter404(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/section-references?type=404", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `href="/docs/section-references/404-vercel"`) {
		t.Error("404 filter must list the Vercel 404 ficha")
	}
	if strings.Contains(body, `href="/docs/section-references/article"`) {
		t.Error("404 filter must not list the article ficha")
	}
	if !strings.Contains(body, `option value="404" selected`) {
		t.Error("404 option must be selected")
	}
}

func TestDocsSectionReferencesUnknownTypeFallsBackToAll(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/section-references?type=gradient", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `href="/docs/section-references/article"`) {
		t.Error("unknown type must fall back to all and still list article")
	}
}

func TestDocsSectionReferencesEmptyResultsUseCompactState(t *testing.T) {
	selected := closedSectionRefType("article")
	body := sectionRefResultsMarkdown(selected, []sectionRefEntry{})
	for _, contract := range []string{
		"Showing 0 results for article.",
		`class="ui-empty-state ui-empty-state--compact"`,
		`role="status"`,
		">No references found</p>",
		">There are no published references for this type yet.</p>",
		`href="/docs/section-references">Clear filter</a>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("empty branch missing %q", contract)
		}
	}
	if strings.Contains(body, `class="ui-list"`) {
		t.Error("empty branch must not render a results list")
	}
}
func TestDocsSectionReferencesEmbeddedContentFailureKeepsRecoveryState(t *testing.T) {
	tests := []struct {
		name string
		path string
		call func(*server, http.ResponseWriter, *http.Request)
	}{
		{
			name: "index",
			path: "/docs/section-references",
			call: func(s *server, w http.ResponseWriter, r *http.Request) {
				s.docsSectionReferences(w, r)
			},
		},
		{
			name: "detail",
			path: "/docs/section-references/article",
			call: func(s *server, w http.ResponseWriter, r *http.Request) {
				r.SetPathValue("id", "article")
				s.docsSectionReferenceDetail(w, r)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newDocsServer(t)
			s.assets = fstest.MapFS{}
			res := httptest.NewRecorder()
			tt.call(s, res, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if res.Code != http.StatusInternalServerError {
				t.Fatalf("status = %d, want %d", res.Code, http.StatusInternalServerError)
			}
			body := res.Body.String()
			for _, contract := range []string{
				`class="ui-error-state" role="alert"`,
				"Something went wrong",
				"Please try again later.",
				`href="/"`,
			} {
				if !strings.Contains(body, contract) {
					t.Errorf("recovery response missing %q", contract)
				}
			}
		})
	}
}

func TestDocsSectionReferencesArticleDetail(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/section-references/article", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		">Article / rich post</h1>",
		`id="gelium-remake"`,
		"Access as a request",
		"ui-data-table-table",
		"https://vercel.com/blog/the-end-of-credential-sprawl-for-agents",
		"Keep / adapt",
		"Ready to deploy?",
		"Ask before copying",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("detail missing %q", contract)
		}
	}
}

func TestDocsSectionReferencesVercel404Detail(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/section-references/404-vercel", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		">404 / missing route</h1>",
		`id="gelium-remake"`,
		"https://vercel.com/docs/errors/does-not-exist",
		`class="ui-error-state" role="alert"`,
		`class="ui-error-state-code" aria-hidden="true">404`,
		`<h2 class="ui-error-state-title">Page not found</h2>`,
		`<a class="ui-button" href="/docs/section-references">Back to section references</a>`,
		"Keep / adapt",
		"Ask before copying",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("404 detail missing %q", contract)
		}
	}
	if count := strings.Count(body, "<h1"); count != 1 {
		t.Errorf("404 detail must render exactly one h1, got %d", count)
	}
}

func TestDocsSectionReferencesPublishedDetails(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		title    string
		official string
		markers  []string
	}{
		{
			name:     "auth register",
			path:     "/docs/section-references/auth-register",
			title:    "Auth / register",
			official: "https://vercel.com/signup",
			markers: []string{
				`id="gelium-remake"`,
				`class="ui-card ui-card-outlined"`,
				`<form class="auth-register-form" method="post"`,
				`class="ui-text-field ui-text-field-filled"`,
				`type="email"`,
				`type="password"`,
				`class="ui-checkbox"`,
				`type="checkbox"`,
				`class="ui-validation-summary" role="alert"`,
				"Forgot your password?",
				"Keep / adapt",
				"Ask before copying",
			},
		},
		{
			name:     "vercel faq",
			path:     "/docs/section-references/faq-vercel",
			title:    "FAQ / pricing questions",
			official: "https://vercel.com/pricing",
			markers: []string{
				`id="gelium-remake"`,
				`class="ui-accordion ui-accordion--behavior-native`,
				`<details class="ui-accordion-item"`,
				`name="pricing-faq"`,
				`<summary class="ui-accordion-trigger"`,
				"Frequently asked questions",
				`class="ui-button ui-button-primary"`,
				"Keep / adapt",
				"Ask before copying",
			},
		},
		{
			name:     "linear hero",
			path:     "/docs/section-references/hero-linear",
			title:    "Hero / product direction",
			official: "https://linear.app",
			markers: []string{
				`id="gelium-remake"`,
				`class="ui-hero"`,
				`class="ui-hero-content"`,
				`class="ui-hero-title"`,
				`class="ui-hero-actions"`,
				`class="ui-button ui-button-primary"`,
				`class="ui-list"`,
				"Keep / adapt",
				"Ask before copying",
			},
		},
		{
			name:     "linear pricing",
			path:     "/docs/section-references/pricing-linear",
			title:    "Pricing / plan comparison",
			official: "https://linear.app/pricing",
			markers: []string{
				`id="gelium-remake"`,
				`class="ui-card ui-card-outlined"`,
				`class="ui-data-table-scroll"`,
				`class="ui-data-table-table"`,
				"Plan comparison",
				`class="ui-button ui-button-primary"`,
				"Keep / adapt",
				"Ask before copying",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if res.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
			}
			body := res.Body.String()
			if !strings.Contains(body, ">"+tt.title+"</h1>") {
				t.Errorf("detail missing H1 title %q", tt.title)
			}
			if !strings.Contains(body, tt.official) {
				t.Errorf("detail missing official URL %q", tt.official)
			}
			if strings.Count(body, "<h1") != 1 {
				t.Errorf("detail must render exactly one h1, got %d", strings.Count(body, "<h1"))
			}
			remakeStart := strings.Index(body, `id="gelium-remake"`)
			if remakeStart < 0 {
				t.Fatal("detail remake landmark is required")
			}
			keepStart := strings.Index(body[remakeStart:], "Keep / adapt")
			if keepStart < 0 {
				t.Fatal("detail Keep / adapt landmark is required")
			}
			remake := body[remakeStart : remakeStart+keepStart]
			if strings.Contains(remake, "/recipes/rich-article") || strings.Contains(remake, "rich-article") {
				t.Error("detail remake must not use the unrelated rich-article fixture")
			}
			for _, marker := range tt.markers {
				if !strings.Contains(body, marker) {
					t.Errorf("detail missing %q", marker)
				}
			}
		})
	}
}

func TestDocsSectionReferencesUnknownIDIs404(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/section-references/missing", nil))
	if res.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusNotFound)
	}
}

func TestDocsSectionReferencesInSitemap(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, want := range []string{
		"<loc>https://gelium-ui.example/docs/section-references/article</loc>",
		"<loc>https://gelium-ui.example/docs/section-references/404-vercel</loc>",
		"<loc>https://gelium-ui.example/docs/section-references/auth-register</loc>",
		"<loc>https://gelium-ui.example/docs/section-references/faq-vercel</loc>",
		"<loc>https://gelium-ui.example/docs/section-references/hero-linear</loc>",
		"<loc>https://gelium-ui.example/docs/section-references/pricing-linear</loc>",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("sitemap missing %q", want)
		}
	}
}

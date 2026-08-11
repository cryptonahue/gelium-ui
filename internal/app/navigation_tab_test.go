package app

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestNavTabDocsRouteRendersStandaloneAndComposedDestinations(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-tab", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-tab docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	// Demo scaffolding lives in the component's own grid vocabulary, without the
	// ui- prefix (reserved for component primitives).
	for _, contract := range []string{
		`navigation-tab-demo-grid`,
		`navigation-tab-demo-group`,
		`navigation-tab-demo-heading`,
		`navigation-tab-demo-row`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("navigation-tab demo is missing %q", contract)
		}
	}
	// The destination contract is the same one the navigation bar uses: the
	// composed demo reuses the delivered .ui-nav-bar root.
	if !strings.Contains(body, `class="ui-nav-bar"`) {
		t.Error("navigation-tab demo must compose the existing .ui-nav-bar for the in-bar variant")
	}
}

func TestNavTabDocsRouteServerDerivesActiveTab(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-tab", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-tab docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	// The active tab is the current page, marked server-side: it must carry
	// aria-current="page" and the --active class, and its href is the route
	// being served.
	if !strings.Contains(body, `class="ui-nav-tab ui-nav-tab--active" href="/components/navigation-tab" aria-current="page"`) {
		t.Error("the current page tab must be marked active with aria-current=\"page\"")
	}
	// Every active tab must be the current page: the --active class and
	// aria-current must always co-occur and never appear on another tab. The
	// regexp anchors on the class-name boundary so the hide-inactive-label row
	// (where --active is followed by another modifier) is still counted, while
	// the escaped code samples in the Markdown anatomy (quote-escaped to
	// &#34;) are excluded. Count only inside the nav-tab block (the page
	// breadcrumb also uses aria-current="page" for its current crumb and must
	// not skew this check).
	navStart := strings.Index(body, `class="ui-nav-tab"`)
	navBody := body[navStart:]
	navEnd := strings.Index(navBody, "</nav>")
	if navStart < 0 || navEnd < 0 {
		t.Fatal("nav-tab block not found in rendered body")
	}
	navBlock := navBody[:navEnd]
	activeClass := regexp.MustCompile(`ui-nav-tab--active(?:\s|")`).FindAllString(navBlock, -1)
	if got, want := len(activeClass), strings.Count(navBlock, `aria-current="page"`); got != want {
		t.Errorf("--active class count %d != aria-current count %d", got, want)
	}
}

func TestNavTabSingleGlyphRenderedWhenSlotsMatch(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-tab", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-tab docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	// The "Home" tab uses the same trusted glyph for inactive and active slots;
	// the template must render exactly ONE copy per demo row. Rendering both
	// would show the icon twice because there is nothing to swap. There are
	// four demo rows, so expect four copies total — never eight.
	homeSVG := string(navTabHomeSVG)
	if got := strings.Count(body, homeSVG); got != 4 {
		t.Errorf("home tab must render its glyph once per demo row (inactive==active), got %d copies, want 4", got)
	}

	// The "Navigation tab" tab uses DISTINCT inactive/active glyphs
	// (navTabMenuSVG / navTabAppsSVG); both copies must render so the CSS swap
	// can hide one based on the active state.
	if !strings.Contains(body, string(navTabMenuSVG)) {
		t.Error("navigation-tab tab must render its inactive glyph")
	}
	if !strings.Contains(body, string(navTabAppsSVG)) {
		t.Error("navigation-tab tab must render its distinct active glyph for the CSS swap")
	}
}

func TestNavTabDocsRouteReusesBadgeAndComposesNavBar(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-tab", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-tab docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	// Badge destinations reuse the existing .ui-badge primitive: the dot form
	// (aria-hidden) and the large count form.
	for _, contract := range []string{
		`class="ui-badge" aria-hidden="true"`,
		`class="ui-badge ui-badge-large"`,
		`>3</span>`,
		`>12</span>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("navigation-tab badges must reuse .ui-badge, missing %q", contract)
		}
	}
	// Every destination is a real link and the in-bar variant is a real <nav>.
	if got := strings.Count(body, `class="ui-nav-tab`); got < 3 {
		t.Errorf("expected at least three navigation tabs, got %d", got)
	}
}

func TestNavTabDocsRouteHidesInactiveLabels(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-tab", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-tab docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	// The hide-inactive-label variant keeps the modifier on each inactive tab
	// (upstream hideInactiveLabel is per-tab).
	if !strings.Contains(body, `class="ui-nav-tab ui-nav-tab--hide-inactive-label"`) {
		t.Error("hide-inactive-label demo must set the ui-nav-tab--hide-inactive-label modifier on inactive tabs")
	}
	if !strings.Contains(body, `class="ui-nav-tab ui-nav-tab--active ui-nav-tab--hide-inactive-label"`) {
		t.Error("the active tab keeps its label and must carry both modifiers")
	}
}

func TestNavTabDocsRouteDogfoodsLinkSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-tab", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-tab docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Navigation tab</h1>`,
		`href="/components/navigation-tab"`,
		`aria-label="Navigation tab examples"`,
		`class="ui-nav-tab"`,
		`<nav`,
		`<a class="ui-nav-tab`,
		`<span class="ui-nav-tab-icon"`,
		`<span class="ui-nav-tab-indicator"`,
		`<span class="ui-nav-tab-label"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("navigation-tab docs are missing %q", contract)
		}
	}
	// No fake-tab ARIA and no roving focus inside the rendered demo: the
	// roadmap contract is a real semantic link, never role="tab" /
	// role="tablist" / tabindex. The check is scoped to the demo preview
	// section because the Markdown anatomy legitimately documents the upstream
	// button/roving-focus behavior.
	const previewMarker = `<section class="component-preview" aria-label="Navigation tab examples">`
	start := strings.Index(body, previewMarker)
	if start < 0 {
		t.Fatal("navigation-tab demo preview section not found")
	}
	preview := body[start:]
	for _, forbidden := range []string{
		`role="tab"`,
		`role="tablist"`,
		`tabindex`,
	} {
		if strings.Contains(preview, forbidden) {
			t.Errorf("navigation-tab demo must not contain %q (link semantics, not a fake tab)", forbidden)
		}
	}
	// The decorative glyph slots (inactive + active) are trusted internal SVG
	// constants; the classes must land in the rendered page.
	for _, contract := range []string{
		`class="ui-nav-tab-glyph"`,
		`class="ui-nav-tab-glyph ui-nav-tab-glyph--active"`,
		`aria-hidden="true"`,
		`focusable="false"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("navigation-tab docs are missing %q", contract)
		}
	}
}

func TestNavTabDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/navigation-tab", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST navigation-tab status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

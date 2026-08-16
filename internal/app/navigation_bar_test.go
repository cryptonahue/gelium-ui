package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNavBarDocsRouteRendersVariantsAndBadgeReuse(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-bar", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-bar docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	// Standard, hide-inactive-labels, and badges demo groups.
	for _, contract := range []string{
		`navigation-bar-demo-grid`,
		`navigation-bar-demo-group`,
		`navigation-bar-demo-heading`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("navigation-bar demo is missing %q", contract)
		}
	}
	// The hide-inactive-labels variant keeps the modifier on the nav root.
	if !strings.Contains(body, `class="ui-nav-bar ui-nav-bar--hide-inactive-labels"`) {
		t.Error("hide-inactive-labels demo must set the ui-nav-bar--hide-inactive-labels modifier")
	}
	// Badge destinations reuse the existing .ui-badge primitive: the dot form
	// (aria-hidden) and the large count form.
	for _, contract := range []string{
		`class="ui-badge" aria-hidden="true"`,
		`class="ui-badge ui-badge-large"`,
		`>3</span>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("navigation-bar badges must reuse .ui-badge, missing %q", contract)
		}
	}
	// Every destination is a real link.
	if got := strings.Count(body, `class="ui-nav-bar-destination"`); got < 3 {
		t.Errorf("expected at least three inert destinations, got %d", got)
	}
}

func TestNavBarDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/navigation-bar", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST navigation-bar status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestNavBarDocsRouteServerDerivesActiveDestination(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-bar", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-bar docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	// The active destination is the current page, marked server-side: it must
	// carry aria-current="page" and the --active class, and its href is the
	// route being served.
	if !strings.Contains(body, `class="ui-nav-bar-destination ui-nav-bar-destination--active" href="/components/navigation-bar" aria-current="page"`) {
		t.Error("the current page destination must be marked active with aria-current=\"page\"")
	}
	// Every active destination must be the current page: the --active class and
	// aria-current must always co-occur and never appear on another destination.
	// Count only inside the nav bar block (the page breadcrumb also uses
	// aria-current="page" for its current crumb and must not skew this check).
	navStart := strings.Index(body, `class="ui-nav-bar"`)
	navBody := body[navStart:]
	navEnd := strings.Index(navBody, "</nav>")
	if navStart < 0 || navEnd < 0 {
		t.Fatal("nav bar block not found in rendered body")
	}
	navBlock := navBody[:navEnd]
	if got, want := strings.Count(navBlock, `ui-nav-bar-destination--active"`), strings.Count(navBlock, `aria-current="page"`); got != want {
		t.Errorf("--active class count %d != aria-current count %d", got, want)
	}
}

func TestNavBarSingleGlyphRenderedWhenSlotsMatch(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-bar", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-bar docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	// The "Home" destination uses the same trusted glyph for inactive and
	// active slots; the template must render exactly ONE copy per demo bar.
	// Rendering both would show the icon twice because there is nothing to
	// swap. There are three demo bars, so expect three copies total — never six.
	homeSVG := string(navBarHomeSVG)
	if got := strings.Count(body, homeSVG); got != 3 {
		t.Errorf("home destination must render its glyph once per demo bar (inactive==active), got %d copies, want 3", got)
	}

	// The "Navigation bar" destination uses DISTINCT inactive/active glyphs
	// (navBarMenuSVG / navBarAppsSVG); both copies must render so the CSS swap
	// can hide one based on the active state.
	if !strings.Contains(body, string(navBarMenuSVG)) {
		t.Error("navigation-bar destination must render its inactive glyph")
	}
	if !strings.Contains(body, string(navBarAppsSVG)) {
		t.Error("navigation-bar destination must render its distinct active glyph for the CSS swap")
	}
}

func TestNavBarDocsRouteDogfoodsNavSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-bar", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-bar docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Navigation bar</h1>`,
		`href="/components/navigation-bar"`,
		`aria-label="Navigation bar examples"`,
		`class="ui-nav-bar"`,
		`<nav`,
		`<ul`,
		`class="ui-nav-bar-destination"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("navigation-bar docs are missing %q", contract)
		}
	}
	// The decorative glyph slots (inactive + active) are trusted internal SVG
	// constants; the classes must land in the rendered page.
	for _, contract := range []string{
		`class="ui-nav-bar-glyph"`,
		`class="ui-nav-bar-glyph ui-nav-bar-glyph--active"`,
		`aria-hidden="true"`,
		`focusable="false"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("navigation-bar docs are missing %q", contract)
		}
	}
}

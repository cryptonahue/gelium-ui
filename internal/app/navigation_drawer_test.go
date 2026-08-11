package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNavigationDrawerDocsRouteRendersBothVariants(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-drawer", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-drawer docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	// Demo scaffolding groups plus both roadmap variants: the permanent
	// standard variant is a real <nav> embedded in the layout; the modal
	// variant is a native <dialog> with an invoker-command trigger.
	for _, contract := range []string{
		`navigation-drawer-demo-grid`,
		`navigation-drawer-demo-group`,
		`navigation-drawer-demo-heading`,
		`class="ui-navigation-drawer ui-navigation-drawer--standard"`,
		`class="ui-navigation-drawer ui-navigation-drawer--modal"`,
		`<dialog`,
		`<nav`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("navigation-drawer demo is missing %q", contract)
		}
	}
}

func TestNavigationDrawerDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/navigation-drawer", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST navigation-drawer status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestNavigationDrawerDocsRouteServerDerivesActiveDestination(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-drawer", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-drawer docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	// The active destination is the current page, marked server-side: it must
	// carry aria-current="page" and the --active class, and its href is the
	// route being served.
	if !strings.Contains(body, `class="ui-navigation-drawer-destination ui-navigation-drawer-destination--active" href="/components/navigation-drawer" aria-current="page"`) {
		t.Error("the current page destination must be marked active with aria-current=\"page\"")
	}
	// Every active destination must be the current page: the --active class and
	// aria-current must always co-occur and never appear on another destination.
	// Count only inside the standard drawer <nav> block (the page breadcrumb
	// also uses aria-current="page" for its current crumb and must not skew
	// this check).
	navStart := strings.Index(body, `<nav class="ui-navigation-drawer`)
	navBody := body[navStart:]
	navEnd := strings.Index(navBody, "</nav>")
	if navStart < 0 || navEnd < 0 {
		t.Fatal("standard drawer nav block not found in rendered body")
	}
	navBlock := navBody[:navEnd]
	if got, want := strings.Count(navBlock, `ui-navigation-drawer-destination--active"`), strings.Count(navBlock, `aria-current="page"`); got != want {
		t.Errorf("--active class count %d != aria-current count %d", got, want)
	}
}

func TestNavigationDrawerDocsRouteModalIsNativeDialogWithInvokerTrigger(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-drawer", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-drawer docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	// The modal variant is a native <dialog> with a stable id and the
	// closedby="any" light-dismiss hint used by the Dialog component.
	if !strings.Contains(body, `<dialog class="ui-navigation-drawer ui-navigation-drawer--modal" id="navigation-drawer-modal" closedby="any"`) {
		t.Error("modal variant must be a native <dialog> with id and closedby=\"any\"")
	}
	// The trigger is the existing button primitive carrying the declarative
	// invoker command: no component JavaScript.
	if !strings.Contains(body, `command="show-modal" commandfor="navigation-drawer-modal"`) {
		t.Error("modal trigger must use the native invoker command show-modal")
	}
	if !strings.Contains(body, `>Open navigation drawer</span>`) {
		t.Error("modal trigger label must be rendered by the button primitive")
	}
}

func TestNavigationDrawerDocsRouteDestinationsAreRealLinksAndReuseBadge(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-drawer", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-drawer docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	// Every destination is a real link.
	if got := strings.Count(body, `class="ui-navigation-drawer-destination"`); got < 3 {
		t.Errorf("expected at least three inert destinations, got %d", got)
	}
	// Badge destinations reuse the existing .ui-badge primitive: the dot form
	// (aria-hidden) and the large count form.
	for _, contract := range []string{
		`class="ui-badge" aria-hidden="true"`,
		`class="ui-badge ui-badge-large"`,
		`>3</span>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("navigation-drawer badges must reuse .ui-badge, missing %q", contract)
		}
	}
}

func TestNavigationDrawerDocsRouteSingleGlyphPerDestination(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-drawer", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-drawer docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()

	// The drawer keeps ONE glyph per destination (unlike the navigation bar it
	// never swaps an active/inactive icon pair), so a shared trusted glyph must
	// render exactly as many copies as drawers that contain it. drawerHomeSVG
	// appears in all three demo drawers.
	homeSVG := string(drawerHomeSVG)
	if got := strings.Count(body, homeSVG); got != 3 {
		t.Errorf("home destination must render its glyph once per demo drawer, got %d copies, want 3", got)
	}
}

func TestNavigationDrawerDocsRouteDogfoodsNavAndDialogSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/navigation-drawer", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("navigation-drawer docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Navigation drawer</h1>`,
		`href="/components/navigation-drawer"`,
		`aria-label="Navigation drawer examples"`,
		`class="ui-navigation-drawer ui-navigation-drawer--`,
		`<nav`,
		`<dialog`,
		`<ul`,
		`class="ui-navigation-drawer-destination"`,
		`class="ui-navigation-drawer-list"`,
		`class="ui-navigation-drawer-item"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("navigation-drawer docs are missing %q", contract)
		}
	}
	// The decorative glyph slots are trusted internal SVG constants; the
	// classes must land in the rendered page.
	for _, contract := range []string{
		`class="ui-navigation-drawer-glyph"`,
		`class="ui-navigation-drawer-glyph-svg"`,
		`aria-hidden="true"`,
		`focusable="false"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("navigation-drawer docs are missing %q", contract)
		}
	}
}

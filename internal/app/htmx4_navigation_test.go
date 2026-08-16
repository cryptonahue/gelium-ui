package app

import (
	"strings"
	"testing"

	webassets "geliumui/web"
)

func TestHTMX4NavigationContract(t *testing.T) {
	layout := readAsset(t, "templates/layout.html")
	for _, contract := range []string{
		`hx-boost:inherited="true"`,
		`hx-swap:inherited="innerMorph"`,
		`hx-push-url:inherited="true"`,
		`hx-history-elt`,
		`name="htmx-config"`,
		`&quot;history&quot;:true`,
		`&quot;morphSkip&quot;`,
		`&quot;morphIgnore&quot;`,
		`/static/morph-afterswap.js?v=0.5.2`,
	} {
		if !strings.Contains(layout, contract) {
			t.Errorf("layout is missing HTMX 4 contract %q", contract)
		}
	}
}

func TestHTMX4RuntimeAndEnhancementsArePresent(t *testing.T) {
	runtime := readAsset(t, "static/htmx.min.js")
	if !strings.Contains(runtime, "4.0.0-beta6") || !strings.Contains(runtime, "innerMorph") {
		t.Fatal("embedded runtime is not the official HTMX 4.0.0-beta6 build")
	}
	app := readAsset(t, "static/app.js")
	for _, legacy := range []string{"htmx:beforeSwap", "htmx:beforeRequest", "htmx:afterSwap", "htmx:responseError", "htmx:sendError", "event.detail.xhr", "responseText", "getResponseHeader", "this.submit();"} {
		if strings.Contains(app, legacy) {
			t.Errorf("app.js must not use the HTMX 2 API %q", legacy)
		}
	}
	for _, contract := range []string{"htmx:before:swap", "htmx:before:request", "htmx:after:swap", "htmx:response:error", "ctx.response", "ctx.text", "X-Gelium-Validation", "shouldSwap = true", "isError = false"} {
		if !strings.Contains(app, contract) {
			t.Errorf("app.js is missing the HTMX 4 response contract %q", contract)
		}
	}
	search := readAsset(t, "static/search.js")
	for name, content := range map[string]string{"app.js": app, "search.js": search} {
		for _, contract := range []string{"htmx:after:swap", "htmx:before:history:restore", "data-gelium-"} {
			if !strings.Contains(content, contract) {
				t.Errorf("%s is missing idempotent post-swap contract %q", name, contract)
			}
		}
	}
	for _, contract := range []string{
		"htmx:before:swap", // server-authority reconciliation runs pre-swap
		"applyOptimisticChrome",
		"requestSubmit",      // fires submit so htmx intercepts (form.submit() = native reload)
		"keepPreservedState", // keeps the other form's hidden input in sync
		"refreshChromeHrefs", // rewrites chrome hrefs so sidebar clicks keep ?scheme= after an optimistic toggle
		"initOnThisPage",     // On-this-page scrollspy (progressive enhancement)
		"IntersectionObserver",
		"data-class",
		"theme-dark",
		"data-theme",
	} {
		if !strings.Contains(app, contract) {
			t.Errorf("app.js is missing the optimistic chrome contract %q", contract)
		}
	}
}

// TestLayer3RuntimeProgressiveEnhancementContract pins the Layer 3
// (mobile/runtime) contracts that live in embedded assets: same-document
// view transitions are an OPT-IN progressive enhancement — activated only
// when the browser supports document.startViewTransition AND the user has
// not requested reduced motion (WCAG 2.3.3). Unlike the cross-document
// at-rule, same-document VT is not auto-gated by prefers-reduced-motion, so
// the JS guard is the enforcement point. Layout carries the safe-area
// viewport (GOV.UK pattern) and a named mobile-nav disclosure.
func TestLayer3RuntimeProgressiveEnhancementContract(t *testing.T) {
	app := readAsset(t, "static/app.js")
	for _, contract := range []string{
		"document.startViewTransition",
		"htmx.config.transitions",
		"matchMedia",
		"prefers-reduced-motion: reduce",
	} {
		if !strings.Contains(app, contract) {
			t.Errorf("app.js is missing the view-transition activation contract %q", contract)
		}
	}
	if !strings.Contains(app, "!window.matchMedia") {
		t.Error("app.js must gate same-document view transitions on reduced motion (same-document VT is not auto-disabled)")
	}

	layout := readAsset(t, "templates/layout.html")
	for _, contract := range []string{
		`viewport-fit=cover`,
		`aria-label="Open navigation menu"`,
	} {
		if !strings.Contains(layout, contract) {
			t.Errorf("layout.html is missing the mobile/safe-area contract %q", contract)
		}
	}
}

func TestHTMX4NavigationPreservesNativeLinkBoundaries(t *testing.T) {
	layout := readAsset(t, "templates/layout.html")
	if !strings.Contains(layout, `href="#main-content"`) {
		t.Error("layout fixture lost the native skip-link anchor")
	}
	external := readAsset(t, "templates/docs-topbar.html")
	if !strings.Contains(external, "https://github.com") {
		t.Error("external navigation fixture is missing")
	}
	for _, path := range []string{"templates/recipe-public-feed.html", "templates/recipe-admin-resource.html", "templates/recipe-ops-queue.html"} {
		body := readAsset(t, path)
		if !strings.Contains(body, `hx-swap="outerHTML"`) || !strings.Contains(body, `hx-target=`) {
			t.Errorf("%s no longer declares explicit partial recipe swaps", path)
		}
	}
}

func readAsset(t *testing.T, path string) string {
	t.Helper()
	b, err := webassets.Assets.ReadFile(path)
	if err != nil {
		t.Fatalf("read embedded asset %s: %v", path, err)
	}
	return string(b)
}

package app

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	webassets "geliumui/web"
)

func openingTagWithID(t *testing.T, body, element, id string) string {
	t.Helper()
	pattern := `<` + element + `\b[^>]*\bid="` + regexp.QuoteMeta(id) + `"[^>]*>`
	tag := regexp.MustCompile(pattern).FindString(body)
	if tag == "" {
		t.Fatalf("body is missing <%s> with id %q", element, id)
	}
	return tag
}

func renderButton(t *testing.T, view buttonView) string {
	t.Helper()
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/button.html"))
	var rendered bytes.Buffer
	if err := tmpl.ExecuteTemplate(&rendered, "button", view); err != nil {
		t.Fatalf("execute button template: %v", err)
	}
	return rendered.String()
}

func renderTextField(t *testing.T, view textFieldView) string {
	t.Helper()
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/text-field.html"))
	var rendered bytes.Buffer
	if err := tmpl.ExecuteTemplate(&rendered, "text-field", view); err != nil {
		t.Fatalf("execute text field template: %v", err)
	}
	return rendered.String()
}

func TestHealthzReturnsPlainTextOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	res := httptest.NewRecorder()

	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if got := res.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Errorf("Content-Type = %q, want plain UTF-8 text", got)
	}
	if got := res.Body.String(); got != "ok\n" {
		t.Errorf("body = %q, want %q", got, "ok\\n")
	}
}

func TestHomeRendersMarketingLanding(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`class="ui-landing"`,
		`class="ui-hero`,
		`Server-rendered components`,
		`class="ui-button ui-button-primary"`,
		`href="/docs"`,
		`href="/components/button"`,
		`class="ui-feature-card`,
		`class="ui-split`,
		`Admin Resource`,
		`class="site-header"`,
		`aria-label="Appearance"`,
		`src="/static/htmx.min.js?v=0.4.0"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("home does not contain contract %q", contract)
		}
	}
	// Must not keep the old Markdown-only home article as the primary surface.
	if strings.Contains(body, `<article class="prose">`) {
		t.Error("marketing landing must not render prose article shell")
	}
	// Must not use the docs two-pane chrome on home.
	if strings.Contains(body, `class="docs-topbar"`) {
		t.Error("home must not render docs shell topbar")
	}
}

func TestLayoutCacheBustsEmbeddedAssetsAcrossExeUpgrades(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	body := res.Body.String()

	for _, asset := range []string{
		`href="/static/app.css?v=0.4.0"`,
		`src="/static/htmx.min.js?v=0.4.0"`,
		`src="/static/app.js?v=0.4.0"`,
	} {
		if !strings.Contains(body, asset) {
			t.Errorf("layout must cache-bust upgraded embedded asset %s", asset)
		}
	}
}

func TestStaticBuildArtifactsAreServedFromEmbeddedFilesystem(t *testing.T) {
	tests := []struct {
		path        string
		contentType string
		contract    string
	}{
		{path: "/static/app.css", contentType: "text/css; charset=utf-8", contract: ".ui-button"},
		{path: "/static/htmx.min.js", contentType: "text/javascript; charset=utf-8", contract: "htmx"},
		{path: "/static/app.js", contentType: "text/javascript; charset=utf-8", contract: "X-Loom-Validation"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			res := httptest.NewRecorder()
			New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if res.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
			}
			if got := res.Header().Get("Content-Type"); got != tt.contentType {
				t.Errorf("Content-Type = %q, want %q", got, tt.contentType)
			}
			if got := res.Header().Get("Cache-Control"); got != "no-cache" {
				t.Errorf("Cache-Control = %q, want revalidation with no-cache", got)
			}
			if !strings.Contains(res.Body.String(), tt.contract) {
				t.Errorf("asset does not contain build contract %q", tt.contract)
			}
		})
	}
}

func TestMaterialDarkThemeKeepsFilledFieldDistinctFromSurface(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/static/app.css", nil))
	css := res.Body.String()

	// Single dark mechanism: each dark value is declared exactly once in the
	// compiled bundle (the explicit .theme-material.theme-dark class route),
	// never duplicated by a dark media block.
	if got := strings.Count(css, "--ui-color-surface:#211f26"); got != 1 {
		t.Fatalf("compiled dark theme surface declarations = %d, want 1 (single class route)", got)
	}
	if got := strings.Count(css, "--ui-field-container:#36343b"); got != 1 {
		t.Errorf("compiled dark theme filled container declarations = %d, want 1 (single class route)", got)
	}
	if strings.Contains(css, "--ui-field-container:#211f26") {
		t.Error("dark filled field container must differ from the #211f26 surface")
	}
}

func TestMaterialThemeDefinesTextFieldTypescaleTokens(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/static/app.css", nil))
	css := res.Body.String()

	// Phase B decomposition (R1): the compiled bundle carries the decomposed
	// per-step tokens (size/weight/line-height) in the theme, and the core
	// alias composes them into the font: shorthand consumers use.
	for _, token := range []string{
		`--ui-type-body-lg-size:1rem`,
		`--ui-type-body-lg-weight:400`,
		`--ui-type-body-lg:var(--ui-type-body-lg-weight) var(--ui-type-body-lg-size)/var(--ui-type-body-lg-line-height) var(--ui-type-body-lg-family)`,
		`--ui-type-body-sm-size:.75rem`,
		`--ui-type-body-sm-weight:400`,
		`--ui-type-label-sm-size:.75rem`,
		`--ui-type-label-sm-weight:500`,
	} {
		if !strings.Contains(css, token) {
			t.Errorf("compiled Material theme is missing text-field typescale definition %q", token)
		}
	}
}

func TestMaterialThemeExposesSemanticFoundationContracts(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/static/app.css", nil))
	css := res.Body.String()

	for _, token := range []string{
		"--ui-color-primary:",
		"--ui-type-display-sm:",
		"--ui-radius-full:",
		"--ui-shadow-5:",
		"--ui-focus-thickness:3px",
		"--ui-focus-offset:2px",
		"--ui-field-container:",
		"--ui-field-border:",
		"--ui-field-error:",
		`.theme-material.dark`,
	} {
		if !strings.Contains(css, token) {
			t.Errorf("compiled Material theme is missing %q", token)
		}
	}
}
func TestThemeClassIsServerDrivenAndAllowlisted(t *testing.T) {
	tests := []struct {
		name  string
		theme string
		want  string
	}{
		{name: "empty falls back to default", theme: "", want: "theme-material"},
		{name: "default passes the allowlist", theme: "theme-material", want: "theme-material"},
		{name: "basecoat passes the allowlist", theme: "theme-basecoat", want: "theme-basecoat"},
		{name: "unknown theme falls back to default", theme: "theme-unknown", want: "theme-material"},
		{name: "arbitrary injection falls back to default", theme: "onload=alert(1)", want: "theme-material"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := themeClass(tt.theme); got != tt.want {
				t.Errorf("themeClass(%q) = %q, want %q", tt.theme, got, tt.want)
			}
		})
	}
}

func TestLayoutRendersThemeClassOnHTMLElement(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	body := res.Body.String()
	for _, contract := range []string{
		`<html lang="en" class="theme-material">`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("home does not render the server-driven theme class %q", contract)
		}
	}
}

func renderToast(t *testing.T, view toastView) string {
	t.Helper()
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/toast.html"))
	var rendered bytes.Buffer
	if err := tmpl.ExecuteTemplate(&rendered, "toast", view); err != nil {
		t.Fatalf("execute toast template: %v", err)
	}
	return rendered.String()
}

func renderBanner(t *testing.T, view bannerView) string {
	t.Helper()
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/banner.html"))
	var rendered bytes.Buffer
	if err := tmpl.ExecuteTemplate(&rendered, "banner", view); err != nil {
		t.Fatalf("execute banner template: %v", err)
	}
	return rendered.String()
}

// renderInlineAlert drives the inline-alert partial directly with the
// production view model (Phase D pattern 3), pinning the tone→role
// derivation (error → alert, rest → status).
func renderInlineAlert(t *testing.T, view inlineAlertView) string {
	t.Helper()
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/inline-alert.html"))
	var rendered bytes.Buffer
	if err := tmpl.ExecuteTemplate(&rendered, "inline-alert", view); err != nil {
		t.Fatalf("execute inline-alert template: %v", err)
	}
	return rendered.String()
}

// TestLayoutOmitsBannerWhenNil proves the pageView.Banner slot is optional: a
// nil banner renders nothing and the layout still builds cleanly.
func TestLayoutOmitsBannerWhenNil(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if strings.Contains(res.Body.String(), "ui-banner") {
		t.Error("layout must not render a banner when pageView.Banner is nil")
	}
}

// TestLayoutRendersBannerSlotBetweenHeaderAndMain proves the global banner slot
// sits between the </header> and <main> landmarks and renders the partial when
// pageView.Banner is set.
func TestLayoutRendersBannerSlotBetweenHeaderAndMain(t *testing.T) {
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/*.html"))
	var page bytes.Buffer
	data := pageView{
		Title:      "banner slot",
		ThemeClass: "theme-material",
		Nav:        []navLink{{Path: "/", Label: "Home"}},
		Banner: &bannerView{
			Tone:  "info",
			Title: "Scheduled maintenance",
			Body:  "Tonight 22:00-23:00 UTC.",
		},
	}
	if err := tmpl.ExecuteTemplate(&page, "layout", data); err != nil {
		t.Fatalf("execute layout with banner: %v", err)
	}
	body := page.String()
	for _, contract := range []string{
		`class="ui-banner ui-banner--info"`,
		`role="status"`,
		`class="ui-banner-title">Scheduled maintenance`,
		`class="ui-banner-body">Tonight 22:00-23:00 UTC.`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("layout does not render banner contract %q", contract)
		}
	}

	headerEnd := strings.Index(body, "</header>")
	mainStart := strings.Index(body, "<main")
	bannerPos := strings.Index(body, `class="ui-banner ui-banner--info"`)
	if headerEnd < 0 || mainStart < 0 {
		t.Fatal("layout is missing the header or main landmarks")
	}
	if bannerPos < headerEnd || bannerPos > mainStart {
		t.Error("banner slot must render between </header> and <main>")
	}
}

// TestBannerRoleIsDerivedFromTone proves role="alert" is reserved for the
// error tone and every other tone announces politely as status.
func TestBannerRoleIsDerivedFromTone(t *testing.T) {
	for _, tone := range []string{"error", "success", "info", "warning"} {
		got := renderBanner(t, bannerView{Tone: tone, Body: "x"})
		want := "role=\"status\""
		if tone == "error" {
			want = "role=\"alert\""
		}
		if !strings.Contains(got, want) {
			t.Errorf("tone %q must render %s", tone, want)
		}
	}
}

// TestInlineAlertRoleIsDerivedFromTone proves the inline-alert partial mirrors
// the banner contract: role="alert" is reserved for the error tone and every
// other tone announces politely as status — success included. This is the
// render-level confirmation for the success persistent reuse (the tone CSS
// alone was already pinned by TestInlineAlertTonesUseCoreTokens).
func TestInlineAlertRoleIsDerivedFromTone(t *testing.T) {
	for _, tone := range []string{"error", "success", "info", "warning"} {
		got := renderInlineAlert(t, inlineAlertView{Tone: tone, Body: "x"})
		want := "role=\"status\""
		if tone == "error" {
			want = "role=\"alert\""
		}
		if !strings.Contains(got, want) {
			t.Errorf("tone %q must render %s", tone, want)
		}
	}
}

// TestBannerMarkupContracts proves the decorative icon is aria-hidden, the CTA
// is a real link, the dismiss is a POST form with a submit button, and an empty
// DismissHref omits the dismiss form entirely.
func TestBannerMarkupContracts(t *testing.T) {
	icon := template.HTML(`<svg aria-hidden="true"><title>info</title></svg>`)
	got := renderBanner(t, bannerView{
		Tone:        "error",
		Icon:        icon,
		Title:       "Session expired",
		Body:        "Re-authenticate to continue.",
		CTA:         true,
		CTAHref:     "/login",
		CTALabel:    "Re-authenticate",
		DismissHref: "/dismiss",
		DismissIcon: template.HTML("×"),
	})
	for _, contract := range []string{
		`<span class="ui-banner-icon" aria-hidden="true">`,
		`class="ui-banner-title">Session expired`,
		`class="ui-banner-body">Re-authenticate to continue.`,
		`<a class="ui-button" href="/login">Re-authenticate</a>`,
		`<form class="ui-banner-dismiss" method="post" action="/dismiss">`,
		`<button class="ui-icon-button" type="submit" aria-label="Dismiss">×</button>`,
	} {
		if !strings.Contains(got, contract) {
			t.Errorf("banner markup is missing contract %q", contract)
		}
	}

	noDismiss := renderBanner(t, bannerView{Tone: "success", Body: "Saved."})
	if strings.Contains(noDismiss, "ui-banner-dismiss") {
		t.Error("banner must omit the dismiss form when DismissHref is empty")
	}
}

func renderErrorState(t *testing.T, view errorStateView) string {
	t.Helper()
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/error-state.html"))
	var rendered bytes.Buffer
	if err := tmpl.ExecuteTemplate(&rendered, "error-state", view); err != nil {
		t.Fatalf("execute error-state template: %v", err)
	}
	return rendered.String()
}

// TestErrorStateMarkupContracts proves the status code is the decorative
// anchor (aria-hidden), the title is the single h1, the body is the muted
// message, and the retry is a real GET link rendered only when Retry is set.
func TestErrorStateMarkupContracts(t *testing.T) {
	got := renderErrorState(t, errorStateView{
		StatusCode: 404,
		Title:      "Page not found",
		Body:       "The page you are looking for does not exist or has moved.",
		Retry:      true,
		Href:       "/",
		Label:      "Back to home",
	})
	for _, contract := range []string{
		`class="ui-error-state" role="alert"`,
		`class="ui-error-state-code" aria-hidden="true">404`,
		`<h1 class="ui-error-state-title">Page not found</h1>`,
		`class="ui-error-state-body">The page you are looking for does not exist or has moved.`,
		`<a class="ui-button" href="/">Back to home</a>`,
	} {
		if !strings.Contains(got, contract) {
			t.Errorf("error-state markup is missing contract %q", contract)
		}
	}
	if count := strings.Count(got, "<h1"); count != 1 {
		t.Errorf("error-state must render exactly one h1, got %d", count)
	}

	noRetry := renderErrorState(t, errorStateView{
		StatusCode: 500,
		Title:      "Something went wrong",
		Body:       "This page could not be loaded. Please try again later.",
	})
	if strings.Contains(noRetry, "ui-button") {
		t.Error("error-state must omit the retry link when Retry is false")
	}
}

// TestUnknownRouteRendersErrorStatePage proves the catch-all serves the ERROR
// STATE slot with the real 404 status and the full layout, instead of the
// plain-text net/http default.
func TestUnknownRouteRendersErrorStatePage(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/does-not-exist", nil))
	if res.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusNotFound)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`class="ui-error-state" role="alert"`,
		`class="ui-error-state-code" aria-hidden="true">404`,
		`<h1 class="ui-error-state-title">Page not found</h1>`,
		`class="ui-error-state-body">The page you are looking for does not exist or has moved.`,
		`<a class="ui-button" href="/">Back to home</a>`,
		`<main`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("404 page is missing contract %q", contract)
		}
	}
	if strings.Contains(body, "404 page not found") {
		t.Error("404 must not fall back to the plain net/http body")
	}
}

// TestErrorSlotOmittedWhenNil proves the Error slot is optional: a nil
// pageView.Error renders the normal content and no error state at all.
func TestErrorSlotOmittedWhenNil(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if strings.Contains(res.Body.String(), "ui-error-state") {
		t.Error("layout must not render an error state when pageView.Error is nil")
	}
}

// TestHomeRendersServerDrivenMetadata proves the home page emits the full
// server-driven metadata contract: description, clean canonical, index robots,
// a minimal OG set and the WebSite JSON-LD (SEO contract §1, §3, §6, §12).
func TestHomeRendersServerDrivenMetadata(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<meta name="description" content="Gelium UI — themeable,`,
		`<link rel="canonical" href="https://gelium-ui.example/">`,
		`<meta name="robots" content="index, follow">`,
		`<meta property="og:type" content="website">`,
		`<meta property="og:url" content="https://gelium-ui.example/">`,
		`<script type="application/ld+json">{"@context":"https://schema.org","@type":"WebSite"`,
		`class="theme-material"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("home is missing metadata contract %q", contract)
		}
	}
	// OG title follows page title for the marketing landing.
	if !strings.Contains(body, `<meta property="og:title"`) {
		t.Error("home must emit og:title")
	}
}

// TestComponentPageRendersMetadata proves component pages resolve per-route
// metadata (article OG type, clean canonical, index robots) from the render
// choke point without any handler change.
func TestComponentPageRendersMetadata(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/button", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<link rel="canonical" href="https://gelium-ui.example/components/button">`,
		`<meta name="robots" content="index, follow">`,
		`<meta property="og:type" content="article">`,
		`<meta property="og:title" content="Button · Gelium UI">`,
		`<meta property="og:url" content="https://gelium-ui.example/components/button">`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("button docs is missing metadata contract %q", contract)
		}
	}
}

// TestCanonicalIsCleanWithoutQuery proves the canonical never carries query
// state: a GET with ?foo=bar still resolves to the clean route path (contract
// §16), because the canonical derives from the route, not the request query.
// renderBreadcrumb drives the breadcrumb partial directly using the
// production breadcrumbItem view model (server.go).
func renderBreadcrumb(t *testing.T, items []breadcrumbItem) string {
	t.Helper()
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/breadcrumb.html"))
	var rendered bytes.Buffer
	if err := tmpl.ExecuteTemplate(&rendered, "breadcrumb", struct{ Items []breadcrumbItem }{Items: items}); err != nil {
		t.Fatalf("execute breadcrumb template: %v", err)
	}
	return rendered.String()
}

func renderFooter(t *testing.T, view footerView) string {
	t.Helper()
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/footer.html"))
	var rendered bytes.Buffer
	if err := tmpl.ExecuteTemplate(&rendered, "footer", view); err != nil {
		t.Fatalf("execute footer template: %v", err)
	}
	return rendered.String()
}

// TestBreadcrumbMarkupContracts proves the P1 markup contract
// (seo-patterns.md:50-64): nav → ol → li, every crumb but the last is a real
// link, and the current crumb is a span with aria-current="page" — never an <a>.
func TestBreadcrumbMarkupContracts(t *testing.T) {
	got := renderBreadcrumb(t, []breadcrumbItem{
		{Href: "/", Label: "Home"},
		{Href: "/docs", Label: "Docs"},
		{Label: "Button", Current: true},
	})
	for _, contract := range []string{
		`<nav aria-label="Breadcrumb">`,
		`<ol class="ui-breadcrumb">`,
		`<li class="ui-breadcrumb-item"><a href="/">Home</a></li>`,
		`<li class="ui-breadcrumb-item"><a href="/docs">Docs</a></li>`,
		`<li class="ui-breadcrumb-item"><span aria-current="page">Button</span></li>`,
	} {
		if !strings.Contains(got, contract) {
			t.Errorf("breadcrumb markup is missing contract %q", contract)
		}
	}

	// A trail with no current crumb renders only links.
	noCurrent := renderBreadcrumb(t, []breadcrumbItem{{Href: "/docs", Label: "Docs"}})
	if !strings.Contains(noCurrent, `<a href="/docs">Docs</a>`) {
		t.Error("breadcrumb must render non-current crumbs as links")
	}
	if strings.Contains(noCurrent, "aria-current") {
		t.Error("breadcrumb must not render aria-current when no crumb is current")
	}
}

// TestFooterMarkupContracts proves the footer partial renders the brand, the
// secondary nav with native details/summary sections (collapsed by default),
// and the legal line — zero JS.
func TestFooterMarkupContracts(t *testing.T) {
	got := renderFooter(t, footerView{
		Brand: "Gelium UI",
		Sections: []footerSection{
			{Title: "Documentation", Links: []navLink{{Path: "/docs", Label: "Docs"}}},
		},
		Legal: "© 2026 Gelium UI · MIT",
	})
	for _, contract := range []string{
		`<footer class="ui-footer">`,
		`<p class="ui-footer-brand">Gelium UI</p>`,
		`<nav class="ui-footer-nav" aria-label="Footer">`,
		`<section class="ui-footer-section">`,
		`<details class="ui-footer-details">`,
		`<summary class="ui-footer-heading">Documentation</summary>`,
		`<ul class="ui-footer-list">`,
		`<a href="/docs">Docs</a>`,
		`<p class="ui-footer-legal">© 2026 Gelium UI · MIT</p>`,
	} {
		if !strings.Contains(got, contract) {
			t.Errorf("footer markup is missing contract %q", contract)
		}
	}
	if strings.Contains(got, "<details open") {
		t.Error("footer must render <details> collapsed by default (zero-JS accordion)")
	}

	// Nil-safe: brand and legal are optional and omitted when empty.
	minimal := renderFooter(t, footerView{Sections: []footerSection{{Title: "S", Links: []navLink{{Path: "/", Label: "Home"}}}}})
	if strings.Contains(minimal, "ui-footer-brand") || strings.Contains(minimal, "ui-footer-legal") {
		t.Error("footer must omit brand and legal when empty")
	}
}

// TestLayoutOmitsFooterWhenNil proves the pageView.Footer slot is optional: a
// nil footer renders no <footer> and the layout still builds cleanly.
func TestLayoutOmitsFooterWhenNil(t *testing.T) {
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/*.html"))
	var page bytes.Buffer
	data := pageView{Title: "no footer", ThemeClass: "theme-material", Nav: []navLink{{Path: "/", Label: "Home"}}}
	if err := tmpl.ExecuteTemplate(&page, "layout", data); err != nil {
		t.Fatalf("execute layout without footer: %v", err)
	}
	if strings.Contains(page.String(), "<footer") {
		t.Error("layout must not render a footer when pageView.Footer is nil")
	}
}

// TestLayoutRendersFooterAfterMain proves the footer slot renders after the
// </main> landmark and before the toast region when pageView.Footer is set.
func TestLayoutRendersFooterAfterMain(t *testing.T) {
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/*.html"))
	var page bytes.Buffer
	data := pageView{
		Title:      "footer slot",
		ThemeClass: "theme-material",
		Nav:        []navLink{{Path: "/", Label: "Home"}},
		Footer: &footerView{
			Brand: "Gelium UI",
			Sections: []footerSection{
				{Title: "Components", Links: []navLink{{Path: "/components/button", Label: "Button"}}},
			},
			Legal: "© 2026 Gelium UI · MIT",
		},
	}
	if err := tmpl.ExecuteTemplate(&page, "layout", data); err != nil {
		t.Fatalf("execute layout with footer: %v", err)
	}
	body := page.String()
	for _, contract := range []string{
		`<footer class="ui-footer">`,
		`<p class="ui-footer-brand">Gelium UI</p>`,
		`class="ui-footer-heading">Components</summary>`,
		`class="ui-footer-legal">© 2026 Gelium UI · MIT</p>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("layout does not render footer contract %q", contract)
		}
	}

	mainEnd := strings.Index(body, "</main>")
	toastPos := strings.Index(body, "loom-toast-region")
	footerPos := strings.Index(body, `<footer class="ui-footer">`)
	if mainEnd < 0 || toastPos < 0 {
		t.Fatal("layout is missing the main landmark or toast region")
	}
	if footerPos < mainEnd || footerPos > toastPos {
		t.Error("footer slot must render after </main> and before the toast region")
	}
}

// TestHomeRendersDefaultFooter proves the real home page ships the footer
// chrome with the default site data: brand, docsNavFor IA groups (Getting
// started + docsSections + Patterns/Recipes/Themes), and the legal line.
func TestHomeRendersDefaultFooter(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<footer class="ui-footer">`,
		`<p class="ui-footer-brand">Gelium UI</p>`,
		`<nav class="ui-footer-nav" aria-label="Footer">`,
		`<summary class="ui-footer-heading">Getting started</summary>`,
		`<summary class="ui-footer-heading">Foundation</summary>`,
		`<summary class="ui-footer-heading">Actions</summary>`,
		`<summary class="ui-footer-heading">Patterns</summary>`,
		`<summary class="ui-footer-heading">Recipes</summary>`,
		`<summary class="ui-footer-heading">Themes</summary>`,
		`<a href="/components/button">Button</a>`,
		`<a href="/docs">Documentation</a>`,
		`<a href="/docs/patterns">Patterns</a>`,
		`<a href="/recipes/admin-resource">Admin Resource</a>`,
		`<p class="ui-footer-legal">© 2026 Gelium UI · MIT</p>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("home is missing footer contract %q", contract)
		}
	}
	if strings.Contains(body, "<details open") {
		t.Error("live footer must not render expanded <details>")
	}
}

func TestCanonicalIsCleanWithoutQuery(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/button?foo=bar", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<link rel="canonical" href="https://gelium-ui.example/components/button">`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("button docs is missing clean canonical %q", contract)
		}
	}
	if strings.Contains(body, "?foo=bar") {
		t.Error("canonical must not carry query state")
	}
}

// TestRobotsTxtPolicy proves GET /robots.txt serves the crawl policy over plain
// text: the public docs site is allowed, the demo/example/recipe surfaces are
// disallowed, and the generated sitemap is advertised (SEO contract §4).
func TestRobotsTxtPolicy(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/robots.txt", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if got := res.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Errorf("Content-Type = %q, want plain UTF-8 text", got)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"User-agent: *",
		"Allow: /",
		"Disallow: /demo/",
		"Disallow: /examples/",
		"Disallow: /recipes/",
		"Sitemap: https://gelium-ui.example/sitemap.xml",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("robots.txt is missing %q", contract)
		}
	}
}

// TestSitemapXMLDerivedFromRegistry proves GET /sitemap.xml lists exactly the
// indexable pages from the route registry — home, /docs and every component —
// once each, with absolute URLs, and never lists the noindex/recipe surfaces.
func TestSitemapXMLDerivedFromRegistry(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if got := res.Header().Get("Content-Type"); got != "application/xml; charset=utf-8" {
		t.Errorf("Content-Type = %q, want application/xml", got)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<?xml version="1.0" encoding="UTF-8"?>`,
		`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`,
		`<loc>https://gelium-ui.example/</loc>`,
		`<loc>https://gelium-ui.example/docs</loc>`,
		`<loc>https://gelium-ui.example/docs/patterns</loc>`,
		`<loc>https://gelium-ui.example/docs/themes</loc>`,
		`<loc>https://gelium-ui.example/components/button</loc>`,
		`<loc>https://gelium-ui.example/components/data-table</loc>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("sitemap is missing %q", contract)
		}
	}
	// home + /docs + patterns + themes + all components
	if got := strings.Count(body, "<url>"); got != len(componentRoutes())+4 {
		t.Errorf("sitemap <url> entries = %d, want %d (home + /docs + stubs + all components)", got, len(componentRoutes())+4)
	}
	for _, excluded := range []string{"/demo/", "/examples/", "/recipes/", "/components/dialog/confirm"} {
		if strings.Contains(body, excluded) {
			t.Errorf("sitemap must not list noindex/form surface %q", excluded)
		}
	}
}

// TestComponentPageRendersStructuredData proves every /components/* page emits
// a single valid JSON-LD block with the BreadcrumbList trail
// (Home > Components > page) and a TechArticle carrying the page headline and
// canonical URL (SEO contract §12).
func TestComponentPageRendersStructuredData(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/button", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<script type="application/ld+json">`,
		`"@context":"https://schema.org"`,
		`"@type":"BreadcrumbList"`,
		`"name":"Home"`,
		`"name":"Components"`,
		`"item":"https://gelium-ui.example/components/button"`,
		`"@type":"TechArticle"`,
		`"headline":"Button · Gelium UI"`,
		`"url":"https://gelium-ui.example/components/button"`,
		`"name":"Gelium UI"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("button docs JSON-LD is missing %q", contract)
		}
	}
	if !json.Valid([]byte(extractJSONLD(t, body))) {
		t.Error("component JSON-LD must parse with encoding/json")
	}
}

// extractJSONLD returns the text inside the first application/ld+json script.
func extractJSONLD(t *testing.T, body string) string {
	t.Helper()
	const open = `<script type="application/ld+json">`
	start := strings.Index(body, open)
	if start < 0 {
		t.Fatal("body has no ld+json script")
	}
	start += len(open)
	end := strings.Index(body[start:], `</script>`)
	if end < 0 {
		t.Fatal("ld+json script is not closed")
	}
	return body[start : start+end]
}

// TestHomeRendersOGImageAndTwitterCard proves the home page ships the default
// og:image placeholder and the matching large-image twitter:card (SEO contract
// §6).
func TestHomeRendersOGImageAndTwitterCard(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<meta property="og:image" content="https://gelium-ui.example/og.png">`,
		`<meta name="twitter:card" content="summary_large_image">`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("home is missing %q", contract)
		}
	}
}

// TestComponentPageRendersOGImage proves component pages also carry the default
// og:image placeholder so every indexable page ships a social image.
func TestComponentPageRendersOGImage(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/button", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<meta property="og:image" content="https://gelium-ui.example/og.png">`,
		`<meta name="twitter:card" content="summary_large_image">`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("button docs is missing %q", contract)
		}
	}
}

// TestLayoutRendersSkipLinkToMain proves every layout page ships the skip link
// as the first focusable element targeting the main landmark (G7). Home keeps
// the legacy centered column; docs shell routes use main.docs-shell-content
// inside the two-pane frame (task 3.2).
func TestLayoutRendersSkipLinkToMain(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		mainClass string
	}{
		{
			name:      "home marketing landing main",
			path:      "/",
			mainClass: `<main id="main-content" class="ui-landing-main">`,
		},
		{
			name:      "docs shell content column",
			path:      "/docs",
			mainClass: `<main id="main-content" class="docs-shell-content">`,
		},
		{
			name:      "component shell content column",
			path:      "/components/button",
			mainClass: `<main id="main-content" class="docs-shell-content">`,
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
			for _, contract := range []string{
				`<a class="ui-skip-link" href="#main-content">Skip to main content</a>`,
				tt.mainClass,
			} {
				if !strings.Contains(body, contract) {
					t.Errorf("%s is missing %q", tt.path, contract)
				}
			}
			// Shell pages must not keep the home-only centered main utility.
			if tt.path != "/" {
				if strings.Contains(body, `class="docs-shell docs-content"`) {
					t.Errorf("%s must not use home main.docs-shell.docs-content", tt.path)
				}
				if strings.Contains(body, `class="site-header"`) {
					t.Errorf("%s must not render legacy site-header", tt.path)
				}
			}
			skipPos := strings.Index(body, `class="ui-skip-link"`)
			mainPos := strings.Index(body, `id="main-content"`)
			if skipPos < 0 || mainPos < 0 || skipPos > mainPos {
				t.Error("skip link must render before the main landmark")
			}
		})
	}
}

// TestComponentBreadcrumbTrailMatchesStructuredData proves the visible
// breadcrumb on a component page is Home > Components > <label>, matching the
// BreadcrumbList JSON-LD so the rendered navigation and the structured data
// never disagree (SEO contract §11).
func TestComponentBreadcrumbTrailMatchesStructuredData(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/button", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<nav aria-label="Breadcrumb">`,
		`<a href="/">Home</a>`,
		`<a href="/docs">Components</a>`,
		`<span aria-current="page">Button</span>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("component breadcrumb is missing %q", contract)
		}
	}
}

// TestThemeQueryParamSelectsTheme proves the document-root theme selection
// (Phase H): a valid ?theme= query renders that theme's class on the document
// root without JS, an unknown theme falls back to the default, and the
// default stays theme-material when no query is present.
func TestThemeQueryParamSelectsTheme(t *testing.T) {
	cases := []struct {
		name  string
		query string
		want  string
	}{
		{name: "basecoat via query", query: "?theme=basecoat", want: `class="theme-basecoat"`},
		{name: "material via query", query: "?theme=material", want: `class="theme-material"`},
		{name: "unknown falls back to default", query: "?theme=unknown", want: `class="theme-material"`},
		{name: "no query keeps default", query: "", want: `class="theme-material"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/button"+tc.query, nil))
			if res.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
			}
			if !strings.Contains(res.Body.String(), tc.want) {
				t.Errorf("rendered page missing %q", tc.want)
			}
		})
	}
}

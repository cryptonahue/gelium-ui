package app

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	webassets "loomui/web"
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

func TestHomeRendersMarkdownInsideDogfoodedLayout(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Gelium UI</h1>`,
		`<main`,
		`class="ui-button ui-button-primary"`,
		`href="/components/button"`,
		`src="/static/htmx.min.js?v=0.4.0"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("home does not contain contract %q", contract)
		}
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

	if got := strings.Count(css, "--ui-color-surface:#211f26"); got != 2 {
		t.Fatalf("compiled dark theme surface declarations = %d, want 2", got)
	}
	if got := strings.Count(css, "--ui-field-container:#36343b"); got != 2 {
		t.Errorf("compiled dark theme filled container declarations = %d, want 2", got)
	}
	if strings.Contains(css, "--ui-field-container:#211f26") {
		t.Error("dark filled field container must differ from the #211f26 surface")
	}
}

func TestMaterialThemeDefinesTextFieldTypescaleTokens(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/static/app.css", nil))
	css := res.Body.String()

	for _, token := range []string{
		`--ui-type-body-lg:400 1rem/1.5rem var(--ui-font-sans)`,
		`--ui-type-body-sm:400 .75rem/1rem var(--ui-font-sans)`,
		`--ui-type-label-sm:500 .75rem/1rem var(--ui-font-sans)`,
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
		{name: "unknown theme falls back to default", theme: "theme-basecoat", want: "theme-material"},
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
		`<meta property="og:title" content="Gelium UI">`,
		`<meta property="og:url" content="https://gelium-ui.example/">`,
		`<script type="application/ld+json">{"@context":"https://schema.org","@type":"WebSite"`,
		`<html lang="en" class="theme-material">`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("home is missing metadata contract %q", contract)
		}
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
// chrome with the default site data: brand, Documentation/Components groups
// derived from the nav, and the legal line.
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
		`<summary class="ui-footer-heading">Documentation</summary>`,
		`<summary class="ui-footer-heading">Components</summary>`,
		`<a href="/components/button">Button</a>`,
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

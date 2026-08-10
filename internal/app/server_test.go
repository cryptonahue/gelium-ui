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
		`<h1>Gelidium UI</h1>`,
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

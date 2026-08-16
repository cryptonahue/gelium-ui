package app

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	webassets "geliumui/site/web"
)

// TestThemeClassRenderedAndServedCSSCarriesRootSelector proves the Phase H
// minimum end-to-end: the page is rendered with the server-driven theme class
// on the document root and the single served bundle carries the matching root
// selector. That combination is what makes selection class-driven: swap the
// class and the cascade picks the other theme — no JS, no rebuild.
func TestThemeClassRenderedAndServedCSSCarriesRootSelector(t *testing.T) {
	htmlRes := httptest.NewRecorder()
	New().ServeHTTP(htmlRes, httptest.NewRequest(http.MethodGet, "/", nil))
	body := htmlRes.Body.String()
	for _, contract := range []string{
		`<html lang="en" class="theme-material">`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("home does not render the server-driven theme class %q", contract)
		}
	}

	cssRes := httptest.NewRecorder()
	New().ServeHTTP(cssRes, httptest.NewRequest(http.MethodGet, "/static/app.css", nil))
	if cssRes.Code != http.StatusOK {
		t.Fatalf("app.css status = %d, want %d", cssRes.Code, http.StatusOK)
	}
	css := cssRes.Body.String()
	if !strings.Contains(css, ".theme-material{") {
		t.Error("served app.css must carry the .theme-material root selector for class-driven selection")
	}
}

// TestBasecoatThemeClassRenderedAndServedCSSCarriesRootSelector proves the
// second theme entered the Phase H mechanism end-to-end: a page rendered with
// ThemeClass "theme-basecoat" puts the class on the document root and the
// single served bundle carries the .theme-basecoat root selector, so swapping
// the server-driven class switches the visual direction without a rebuild.
func TestBasecoatThemeClassRenderedAndServedCSSCarriesRootSelector(t *testing.T) {
	tmpl := template.Must(template.ParseFS(webassets.Assets, "templates/*.html"))
	var page bytes.Buffer
	data := pageView{
		Title:      "basecoat render",
		ThemeClass: "theme-basecoat",
		Nav:        []navLink{{Path: "/", Label: "Home"}},
	}
	if err := tmpl.ExecuteTemplate(&page, "layout", data); err != nil {
		t.Fatalf("execute layout with theme-basecoat: %v", err)
	}
	if !strings.Contains(page.String(), `<html lang="en" class="theme-basecoat">`) {
		t.Error("layout must render the requested class on the document root")
	}

	cssRes := httptest.NewRecorder()
	New().ServeHTTP(cssRes, httptest.NewRequest(http.MethodGet, "/static/app.css", nil))
	if cssRes.Code != http.StatusOK {
		t.Fatalf("app.css status = %d, want %d", cssRes.Code, http.StatusOK)
	}
	if !strings.Contains(cssRes.Body.String(), ".theme-basecoat{") {
		t.Error("served app.css must carry the .theme-basecoat root selector for class-driven selection")
	}
}

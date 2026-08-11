package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestDocsThemesPageRendersOwnership proves GET /docs/themes documents
// family ownership, flat lib/npm theme paths, and dark as a single class
// route only (never media-only).
func TestDocsThemesPageRendersOwnership(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/themes", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"Themes",
		"Token-family ownership",
		"lib/themes/",
		"theme-material.css",
		"gelium-ui/themes/",
		"single dark class route",
		"never a media-only",
		"prefers-color-scheme",
		"?theme=basecoat",
		`href="/docs/tokens"`,
		`href="/docs/themes"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("themes page missing %q", contract)
		}
	}
	if strings.Contains(body, "themes/theme-material/theme.css") {
		t.Error("themes page must not document the retired nested path themes/theme-material/theme.css")
	}
	if strings.Contains(body, "themes/<name>/theme.css") {
		t.Error("themes page must not recommend themes/<name>/theme.css (flat lib/themes/<name>.css only)")
	}
}

package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestDocsTokensPageRendersOwnership proves GET /docs/tokens carries the
// family table with an Owner column and the must-define / core-never rules.
func TestDocsTokensPageRendersOwnership(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/tokens", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"Tokens",
		"--ui-color-primary",
		"Owner",
		"What a valid theme must define",
		"What core never overrides",
		"Semantic color",
		"single dark class route",
		`href="/docs/themes"`,
		`href="/docs/tokens"`,
		"lib/themes/",
		// Owner cells render as <strong>theme</strong> / core / component.
		"<strong>theme</strong>",
		"<strong>core</strong>",
		"<strong>component</strong>",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("tokens page missing %q", contract)
		}
	}
}

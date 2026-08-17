package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestDocsResponsivePageRendersContract proves GET /docs/responsive serves the
// English responsive-design handbook: viewports not device names, containment
// without overflow-x:hidden masking, Gelium touch/container tokens, and links
// into related handbook pages.
func TestDocsResponsivePageRendersContract(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/responsive", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"Design for screen sizes, not devices",
		"overflow-x: hidden must not mask",
		"--ui-touch-target",
		"--ui-container-max",
		"65ch",
		"min-width: 0",
		"viewports",
		`href="/docs/forms"`,
		`href="/docs/performance"`,
		`href="/docs/compare"`,
		`href="/docs/choose-the-right-control"`,
		`href="/docs/responsive"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("responsive page missing %q", contract)
		}
	}
}

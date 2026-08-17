package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestDocsIndexIsOrientationHub proves /docs is a start hub, not a second
// sidebar: Start here links, sidebar pointer, curated recipes/demos — and it
// must NOT dump every docsSections category heading (Foundation, Actions…).
func TestDocsIndexIsOrientationHub(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("docs index status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Documentation</h1>`,
		`>Start here</h2>`,
		`>Install and agents</h3>`,
		`>Core (screen criteria)</h3>`,
		`>System (design system)</h3>`,
		`>Use the sidebar</h2>`,
		`>Try a full screen</h2>`,
		`>Deep dive</h2>`,
		`href="/docs/screens"`,
		`href="/docs/journeys"`,
		`href="/docs/data-display"`,
		`href="/docs/feedback"`,
		`href="/docs/patterns"`,
		`href="/docs/ui-definition-of-done"`,
		`href="/docs/agent-workflow"`,
		`href="/docs/compare"`,
		`href="/docs/forms"`,
		`href="/docs/themes"`,
		`href="/docs/tokens"`,
		`href="/llms.txt"`,
		`href="/llms-ux.txt"`,
		`href="/recipes/admin-resource"`,
		`href="/demo/whatsapp"`,
		`npmjs.com/package/gelium-ui`,
		"Core (screen criteria)",
		"System (design system)",
		"Open",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("docs hub missing %q", contract)
		}
	}
	// Sidebar tier names appear in the orientation copy (not as H2 catalog dumps).
	if !strings.Contains(body, "Core") || !strings.Contains(body, "System") || !strings.Contains(body, "Meta") {
		t.Error("docs hub must name Core / System / Meta tiers")
	}
	// Redundancy guard: category catalog headings belong in the sidebar IA, not the hub body.
	for _, banned := range []string{
		`>Foundation</h2>`,
		`>Actions</h2>`,
		`>Input</h2>`,
		`>Navigation</h2>`,
		`>Data</h2>`,
		`>Handbook</h2>`,
		`>Meta</h2>`, // Meta is sidebar-only; hub must not dump the tier as H2 catalog
	} {
		if strings.Contains(body, banned) {
			t.Errorf("docs hub must not re-list sidebar catalog heading %q", banned)
		}
	}
}

func TestDocsIndexInNav(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/button", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("button docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `href="/docs"`) || !strings.Contains(body, `Documentation`) {
		t.Error("docs shell nav must include a Documentation link to /docs")
	}
}

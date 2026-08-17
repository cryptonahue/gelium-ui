package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestDocsPerformancePageRendersStance proves GET /docs/performance serves
// the performance product stance: JS as progressive enhancement, measured
// ~50 KB docs JS / large CSS by design, and how to measure — not a zero-KB contest.
func TestDocsPerformancePageRendersStance(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/performance", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"Performance stance",
		"CSS is the biggest asset by design",
		"progressive enhancement",
		"~50 KB",
		"~210 KB",
		"npm pack",
		"~87 KB",
		"What we do not chase",
		"unused utilities",
		`href="/docs/compare"`,
		`href="/docs/performance"`,
		"/llms.txt",
		"gelium-ui",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("performance page missing %q", contract)
		}
	}
}

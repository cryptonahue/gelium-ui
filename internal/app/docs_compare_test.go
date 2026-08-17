package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestDocsComparePageRendersHonestEvaluation proves GET /docs/compare serves
// the evaluator-facing comparison: payload orders of magnitude, when to use,
// and explicit no-gos (not a competitor dunk).
func TestDocsComparePageRendersHonestEvaluation(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/compare", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"Why Gelium",
		"When NOT to choose Gelium",
		"Radix",
		"shadcn",
		"Base UI",
		"~50 KB",
		"no-JS",
		"npm install gelium-ui",
		`href="/docs/compare"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("compare page missing %q", contract)
		}
	}
}

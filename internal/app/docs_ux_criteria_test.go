package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestDocsScreensPageRendersSourcedCriteria proves GET /docs/screens serves
// screen types, hierarchy, nav patterns, and cites external sources.
func TestDocsScreensPageRendersSourcedCriteria(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/screens", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		">Screens</h1>",
		"Screen types",
		"Build checklist",
		"GOV.UK",
		"USWDS",
		"design-system.service.gov.uk",
		"designsystem.digital.gov",
		"nngroup.com",
		`href="/docs/feedback"`,
		`href="/llms-ux.txt"`,
		`href="/docs/screens"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("screens page missing %q", contract)
		}
	}
}

// TestDocsFeedbackPageRendersDecisionMatrix proves GET /docs/feedback serves
// the sourced feedback matrix (GOV.UK error summary / banner rules).
func TestDocsFeedbackPageRendersDecisionMatrix(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/feedback", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		">Feedback</h1>",
		"Decision matrix",
		"error summary",
		"validation",
		"toast",
		"FEED-VAL",
		"FEED-OK-PAGE",
		"FEED-ROW",
		"FEED-PARTIAL",
		"FEED-LOAD-LIST",
		"Toast rules",
		"Loading rules",
		"design-system.service.gov.uk",
		"designsystem.digital.gov",
		"nngroup.com",
		`href="/docs/screens"`,
		`href="/llms-ux.txt"`,
		`href="/docs/feedback"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("feedback page missing %q", contract)
		}
	}
}

// TestLlmsUXTxtServesAgentDecisionPack proves GET /llms-ux.txt serves dense
// agent tables linked from the human UX handbook.
func TestLlmsUXTxtServesAgentDecisionPack(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/llms-ux.txt", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if got := res.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Errorf("Content-Type = %q, want plain UTF-8 text", got)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"UX decision pack",
		"SCREEN TYPE",
		"FEEDBACK matrix",
		"validation-summary",
		"FEED-VAL",
		"FEED-OK-PAGE",
		"FEED-ROW",
		"FEED-PARTIAL",
		"FEED-LOAD-LIST",
		"Toast rules",
		"/docs/screens",
		"/docs/feedback",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("llms-ux.txt missing %q", contract)
		}
	}
}

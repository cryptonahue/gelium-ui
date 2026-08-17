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
		"Primary action",
		"Nav decision cues",
		"GOV.UK",
		"USWDS",
		"design-system.service.gov.uk",
		"designsystem.digital.gov",
		"nngroup.com",
		`href="/docs/feedback"`,
		`href="/docs/journeys"`,
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
		"FEED-VAL",
		"FEED-PARTIAL",
		"Toast rules",
		"design-system.service.gov.uk",
		`href="/llms-ux.txt"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("feedback page missing %q", contract)
		}
	}
}

// TestDocsJourneysPageRendersShapes proves GET /docs/journeys serves journey
// shapes and post-submit landing rules from GOV.UK-adapted criteria.
func TestDocsJourneysPageRendersShapes(t *testing.T) {
	body := getOKBody(t, "/docs/journeys")
	for _, contract := range []string{
		">Journeys</h1>",
		"JOURNEY-LINEAR",
		"JOURNEY-TASKLIST",
		"After submit",
		"design-system.service.gov.uk",
		`href="/docs/feedback"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("journeys page missing %q", contract)
		}
	}
}

// TestDocsDataDisplayPageRendersPatterns proves GET /docs/data-display serves
// DATA-* when/when-not rules.
func TestDocsDataDisplayPageRendersPatterns(t *testing.T) {
	body := getOKBody(t, "/docs/data-display")
	for _, contract := range []string{
		">Data display</h1>",
		"DATA-TABLE",
		"DATA-CARDS",
		"FEED-EMPTY",
		"design-system.service.gov.uk",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("data-display page missing %q", contract)
		}
	}
}

// TestDocsPatternsDomainSkeletons proves patterns is no longer a stub and
// includes forum/catalog skeletons plus recipe links.
func TestDocsPatternsDomainSkeletons(t *testing.T) {
	body := getOKBody(t, "/docs/patterns")
	for _, contract := range []string{
		">Patterns</h1>",
		"SKEL-FORUM",
		"SKEL-CATALOG",
		"SKEL-ADMIN-RESOURCE",
		`href="/recipes/admin-resource"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("patterns page missing %q", contract)
		}
	}
}

// TestDocsDensityMotionAndDoD proves supporting UX criteria pages render.
func TestDocsDensityMotionAndDoD(t *testing.T) {
	cases := map[string][]string{
		"/docs/density": {
			"Density and shell",
			"Density modes",
			"--ui-touch-target",
		},
		"/docs/motion": {
			">Motion</h1>",
			"prefers-reduced-motion",
			"MOTION-NONE",
		},
		"/docs/ui-definition-of-done": {
			"UI definition of done",
			"DoD checklist",
			"FEED-VAL",
		},
	}
	for path, contracts := range cases {
		body := getOKBody(t, path)
		for _, c := range contracts {
			if !strings.Contains(body, c) {
				t.Errorf("%s missing %q", path, c)
			}
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
		"JOURNEY-LINEAR",
		"DATA-TABLE",
		"FEED-VAL",
		"SKEL-FORUM",
		"DoD",
		"/docs/journeys",
		"/docs/data-display",
		"/docs/ui-definition-of-done",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("llms-ux.txt missing %q", contract)
		}
	}
}

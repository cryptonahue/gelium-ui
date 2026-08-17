package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSkeletonDocsRendersShellTitleAndLiveVariants(t *testing.T) {
	s := newDocsServer(t)
	res := httptest.NewRecorder()
	s.skeletonDocs(res, httptest.NewRequest(http.MethodGet, "/components/skeleton", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("skeleton docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Skeleton</h1>`,
		`<title>Skeleton · Gelium UI</title>`,
		`href="https://gelium-ui.example/components/skeleton"`,
		`class="docs-shell"`,
		`class="ui-skeleton"`,
		`role="status"`,
		`class="sr-only">Loading conversations`,
		`class="ui-skeleton ui-skeleton--avatar"`,
		`class="ui-skeleton-block ui-skeleton-block--title"`,
		`class="ui-skeleton-block ui-skeleton-block--line"`,
		`ui-skeleton-block--short`,
		`ui-skeleton-block--circle`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("skeleton docs are missing %q", contract)
		}
	}
}

func TestSkeletonDocsLeadNamesComponentAndWhen(t *testing.T) {
	s := newDocsServer(t)
	res := httptest.NewRecorder()
	s.skeletonDocs(res, httptest.NewRequest(http.MethodGet, "/components/skeleton", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("skeleton docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"Use a skeleton while data loads",
		"No component JavaScript",
		"instead of a silent freeze",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("skeleton answer-first lead is missing %q", contract)
		}
	}
}
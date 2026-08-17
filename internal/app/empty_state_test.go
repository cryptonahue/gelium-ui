package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEmptyStateDocsRendersShellTitleAndLiveVariants(t *testing.T) {
	s := newDocsServer(t)
	res := httptest.NewRecorder()
	s.emptyStateDocs(res, httptest.NewRequest(http.MethodGet, "/components/empty-state", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("empty-state docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Empty state</h1>`,
		`<title>Empty state · Gelium UI</title>`,
		`href="https://gelium-ui.example/components/empty-state"`,
		`class="docs-shell"`,
		`class="ui-empty-state"`,
		`role="status"`,
		`class="ui-empty-state ui-empty-state--compact"`,
		`class="ui-empty-state-title"`,
		`class="ui-empty-state-body"`,
		`class="ui-button"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("empty-state docs are missing %q", contract)
		}
	}
}

func TestEmptyStateDocsLeadNamesComponentAndWhen(t *testing.T) {
	s := newDocsServer(t)
	res := httptest.NewRecorder()
	s.emptyStateDocs(res, httptest.NewRequest(http.MethodGet, "/components/empty-state", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("empty-state docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"Use an empty state when",
		"instead of a silent blank area",
		"no component JavaScript",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("empty-state answer-first lead is missing %q", contract)
		}
	}
}
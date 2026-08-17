package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestErrorStateDocsRendersShellTitleAndLiveSpecimens(t *testing.T) {
	s := newDocsServer(t)
	res := httptest.NewRecorder()
	s.errorStateDocs(res, httptest.NewRequest(http.MethodGet, "/components/error-state", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("error-state docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Error state</h1>`,
		`<title>Error state · Gelium UI</title>`,
		`href="https://gelium-ui.example/components/error-state"`,
		`class="docs-shell"`,
		`class="ui-error-state"`,
		`role="alert"`,
		`class="ui-error-state-code"`,
		`aria-hidden="true">404`,
		`class="ui-error-state-title"`,
		`class="ui-error-state-body"`,
		`class="ui-button"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("error-state docs are missing %q", contract)
		}
	}
}

func TestErrorStateDocsLeadNamesComponentAndWhen(t *testing.T) {
	s := newDocsServer(t)
	res := httptest.NewRecorder()
	s.errorStateDocs(res, httptest.NewRequest(http.MethodGet, "/components/error-state", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("error-state docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"Use an error state when",
		"no component JavaScript",
		"knows what happened and how to recover",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("error-state answer-first lead is missing %q", contract)
		}
	}
}
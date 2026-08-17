package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInlineAlertDocsRendersShellTitleAndLiveToneSpecimens(t *testing.T) {
	s := newDocsServer(t)
	res := httptest.NewRecorder()
	s.inlineAlertDocs(res, httptest.NewRequest(http.MethodGet, "/components/inline-alert", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("inline-alert docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Inline alert</h1>`,
		`<title>Inline alert · Gelium UI</title>`,
		`href="https://gelium-ui.example/components/inline-alert"`,
		`class="docs-shell"`,
		`class="ui-inline-alert ui-inline-alert--error"`,
		`role="alert"`,
		`class="ui-inline-alert ui-inline-alert--success"`,
		`role="status"`,
		`class="ui-inline-alert-title"`,
		`class="ui-inline-alert-body"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("inline-alert docs are missing %q", contract)
		}
	}
}

func TestInlineAlertDocsLeadNamesComponentAndWhen(t *testing.T) {
	s := newDocsServer(t)
	res := httptest.NewRecorder()
	s.inlineAlertDocs(res, httptest.NewRequest(http.MethodGet, "/components/inline-alert", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("inline-alert docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"section- or form-level",
		"Use an inline alert when",
		"no component JavaScript",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("inline-alert answer-first lead is missing %q", contract)
		}
	}
}
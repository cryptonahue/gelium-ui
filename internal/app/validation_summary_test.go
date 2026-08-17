package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidationSummaryDocsRendersShellTitleAndLiveSpecimen(t *testing.T) {
	s := newDocsServer(t)
	res := httptest.NewRecorder()
	s.validationSummaryDocs(res, httptest.NewRequest(http.MethodGet, "/components/validation-summary", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("validation-summary docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Validation summary</h1>`,
		`<title>Validation summary · Gelium UI</title>`,
		`href="https://gelium-ui.example/components/validation-summary"`,
		`class="docs-shell"`,
		`class="ui-validation-summary"`,
		`role="alert"`,
		`class="ui-validation-summary-title"`,
		`class="ui-validation-summary-list"`,
		`class="ui-validation-summary-item"`,
		`href="#email"`,
		`href="#password"`,
		`href="#region"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("validation-summary docs are missing %q", contract)
		}
	}
}

func TestValidationSummaryDocsLeadNamesComponentAndWhen(t *testing.T) {
	s := newDocsServer(t)
	res := httptest.NewRecorder()
	s.validationSummaryDocs(res, httptest.NewRequest(http.MethodGet, "/components/validation-summary", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("validation-summary docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"Use a validation summary when",
		"native anchors, no JavaScript",
		"navigation landmark",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("validation-summary answer-first lead is missing %q", contract)
		}
	}
}
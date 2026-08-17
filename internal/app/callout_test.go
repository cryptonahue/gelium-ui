package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalloutDocsRendersShellTitleAndLiveVariants(t *testing.T) {
	s := newDocsServer(t)
	res := httptest.NewRecorder()
	s.calloutDocs(res, httptest.NewRequest(http.MethodGet, "/components/callout", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("callout docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Callout</h1>`,
		`<title>Callout · Gelium UI</title>`,
		`href="https://gelium-ui.example/components/callout"`,
		`class="docs-shell"`,
		`class="ui-callout ui-callout--tip"`,
		`class="ui-callout ui-callout--info"`,
		`class="ui-callout-heading"`,
		`class="ui-callout-body"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("callout docs are missing %q", contract)
		}
	}
}

func TestCalloutDocsLeadNamesComponentAndWhen(t *testing.T) {
	s := newDocsServer(t)
	res := httptest.NewRecorder()
	s.calloutDocs(res, httptest.NewRequest(http.MethodGet, "/components/callout", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("callout docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"Use a callout when",
		"not a state signal",
		"no component JavaScript",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("callout answer-first lead is missing %q", contract)
		}
	}
}
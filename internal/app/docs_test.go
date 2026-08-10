package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDocsIndexRendersSectionsAndComponentLinks(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("docs index status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Documentation</h1>`,
		`<h2>Foundation</h2>`,
		`<h2>Actions</h2>`,
		`<h2>Input</h2>`,
		`<h2>Navigation</h2>`,
		`<h2>Data</h2>`,
		`href="/components/button"`,
		`href="/components/tooltip"`,
		`href="/components/data-table"`,
		`href="/demo/whatsapp"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("docs index is missing %q", contract)
		}
	}
}

func TestDocsIndexInNav(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/button", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("button docs status = %d, want %d", res.Code, http.StatusOK)
	}
	if !strings.Contains(res.Body.String(), `href="/docs">Docs`) {
		t.Error("nav must include a Docs link")
	}
}

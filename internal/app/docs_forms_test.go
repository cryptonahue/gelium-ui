package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestDocsFormsPageRendersContract proves GET /docs/forms serves the forms
// contract handbook: labels above fields, type/inputmode, autocomplete,
// validate-after-interaction, and links into related handbook pages.
func TestDocsFormsPageRendersContract(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/docs/forms", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		"Forms contract",
		"label above",
		"inputmode",
		"autocomplete",
		"after the user interacts",
		"Do not block paste",
		"Disabled wins over error",
		`href="/docs/content-style"`,
		`href="/docs/server-contracts"`,
		`href="/docs/forms"`,
		`href="/components/text-field"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("forms page missing %q", contract)
		}
	}
}

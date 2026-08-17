package app

import (
	"net/http"
	"strings"
	"testing"
)

func TestFooterDocsRouteDogfoodsLiveSpecimenAndShell(t *testing.T) {
	res := renderDocsPage(t, "/components/footer", (*server).footerDocs)

	if res.Code != http.StatusOK {
		t.Fatalf("footer docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Footer</h1>`,
		`class="docs-shell"`,
		`class="ui-footer"`,
		`class="ui-footer-brand"`,
		`class="ui-footer-nav"`,
		`aria-label="Footer"`,
		`class="ui-footer-section"`,
		`class="ui-footer-details"`,
		`class="ui-footer-heading"`,
		`class="ui-footer-list"`,
		`class="ui-footer-legal"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("footer docs are missing %q", contract)
		}
	}
}

func TestFooterDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := renderDocsMethod(t, "/components/footer", http.MethodPost, (*server).footerDocs)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST footer status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}
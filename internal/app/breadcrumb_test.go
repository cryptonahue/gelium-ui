package app

import (
	"net/http"
	"strings"
	"testing"
)

func TestBreadcrumbDocsRouteDogfoodsLiveSpecimenAndShell(t *testing.T) {
	res := renderDocsPage(t, "/components/breadcrumb", (*server).breadcrumbDocs)

	if res.Code != http.StatusOK {
		t.Fatalf("breadcrumb docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Breadcrumb</h1>`,
		`class="docs-shell"`,
		`aria-label="Breadcrumb"`,
		`class="ui-breadcrumb"`,
		`class="ui-breadcrumb-item"`,
		`aria-current="page"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("breadcrumb docs are missing %q", contract)
		}
	}
}

func TestBreadcrumbDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := renderDocsMethod(t, "/components/breadcrumb", http.MethodPost, (*server).breadcrumbDocs)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST breadcrumb status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}
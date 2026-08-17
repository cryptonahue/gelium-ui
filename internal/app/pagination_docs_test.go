package app

import (
	"net/http"
	"strings"
	"testing"
)

// Docs-page contract tests for Pagination. The companion recipe-primitive
// tests (renderPartial + paginationView) live in pagination_test.go; this
// file only covers the /components/pagination documentation page.

func TestPaginationDocsRouteDogfoodsLiveSpecimenAndShell(t *testing.T) {
	res := renderDocsPage(t, "/components/pagination", (*server).paginationDocs)

	if res.Code != http.StatusOK {
		t.Fatalf("pagination docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Pagination</h1>`,
		`class="docs-shell"`,
		`class="ui-pagination"`,
		`class="ui-pagination-page"`,
		`class="ui-pagination-page ui-pagination-page--current"`,
		`aria-current="page"`,
		`class="ui-pagination-page ui-pagination-page--disabled"`,
		`aria-disabled="true"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("pagination docs are missing %q", contract)
		}
	}
}

func TestPaginationDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := renderDocsMethod(t, "/components/pagination", http.MethodPost, (*server).paginationDocs)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST pagination status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}
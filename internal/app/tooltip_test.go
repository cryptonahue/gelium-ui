package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTooltipDocsRouteDogfoodsVariants(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/tooltip", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("tooltip docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Tooltip</h1>`,
		`href="/components/tooltip"`,
		`aria-label="Tooltip examples"`,
		`class="ui-tooltip"`,
		`role="tooltip"`,
		`aria-describedby`,
		`ui-tooltip--rich`,
		`ui-tooltip--top`,
		`tooltip-demo-grid`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("tooltip docs are missing %q", contract)
		}
	}
}

func TestTooltipDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/tooltip", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST tooltip status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

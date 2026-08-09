package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIconDocsRouteDogfoodsTrustedSVGContracts(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/icon", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("icon docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Icon</h1>`,
		`href="/components/icon"`,
		`aria-label="Icon examples"`,
		`class="ui-icon"`,
		`aria-hidden="true"`,
		`focusable="false"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("icon docs are missing %q", contract)
		}
	}
}

func TestIconDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/icon", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST icon status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSwitchDocsRouteDogfoodsStates(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/switch", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("switch docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Switch</h1>`,
		`href="/components/switch"`,
		`aria-label="Switch examples"`,
		`type="checkbox"`,
		`class="ui-switch"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("switch docs are missing %q", contract)
		}
	}
}

func TestSwitchDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/switch", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST switch status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

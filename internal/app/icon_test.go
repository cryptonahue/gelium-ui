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
		`>Icon</h1>`,
		`href="/components/icon"`,
		`aria-label="Icon examples"`,
		`class="ui-icon"`,
		`aria-hidden="true"`,
		`focusable="false"`,
		`viewBox="0 -960 960 960"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("icon docs are missing %q", contract)
		}
	}
	if got := strings.Count(body, `viewBox="0 -960 960 960"`); got < 36 {
		t.Errorf("icon docs must render the curated Material Symbols set, got %d glyphs, want >= 36", got)
	}
}

func TestIconDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/icon", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST icon status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

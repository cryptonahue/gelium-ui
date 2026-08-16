package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBadgeDocsRouteDogfoodsDotAndCountStates(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/badge", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("badge docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Badge</h1>`,
		`href="/components/badge"`,
		`aria-label="Badge examples"`,
		`class="ui-badge"`,
		`class="ui-badge ui-badge-large"`,
		`>3</span>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("badge docs are missing %q", contract)
		}
	}
}

func TestBadgeDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/badge", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST badge status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

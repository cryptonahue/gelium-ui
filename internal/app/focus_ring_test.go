package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFocusRingDocsRouteDogfoodsSharedFocusVisibleContract(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/focus-ring", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("focus ring docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Focus ring</h1>`,
		`href="/components/focus-ring"`,
		`aria-label="Focus ring example"`,
		`class="focus-demo-grid"`,
		`class="focus-demo-link"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("focus ring docs are missing %q", contract)
		}
	}
}

func TestFocusRingDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/focus-ring", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST focus ring status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

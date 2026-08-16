package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheckboxDocsRouteDogfoodsStates(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/checkbox", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("checkbox docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Checkbox</h1>`,
		`href="/components/checkbox"`,
		`aria-label="Checkbox examples"`,
		`type="checkbox"`,
		`class="ui-checkbox"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("checkbox docs are missing %q", contract)
		}
	}
}

func TestCheckboxDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/checkbox", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST checkbox status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

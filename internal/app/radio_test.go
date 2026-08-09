package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRadioDocsRouteDogfoodsStates(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/radio", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("radio docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Radio</h1>`,
		`href="/components/radio"`,
		`aria-label="Radio examples"`,
		`type="radio"`,
		`class="ui-radio"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("radio docs are missing %q", contract)
		}
	}
}

func TestRadioDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/radio", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST radio status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

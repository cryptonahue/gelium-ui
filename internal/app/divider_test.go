package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDividerDocsRouteDogfoodsNativeHRSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/divider", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("divider docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Divider</h1>`,
		`href="/components/divider"`,
		`aria-label="Divider examples"`,
		`class="ui-divider"`,
		`class="ui-divider ui-divider-inset"`,
		`class="ui-divider ui-divider-inset-start"`,
		`class="ui-divider ui-divider-inset-end"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("divider docs are missing %q", contract)
		}
	}
}

func TestDividerDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/divider", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST divider status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

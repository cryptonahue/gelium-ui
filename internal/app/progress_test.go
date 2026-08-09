package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProgressDocsRouteDogfoodsNativeProgressStates(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/progress", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("progress docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Progress</h1>`,
		`href="/components/progress"`,
		`aria-label="Progress examples"`,
		`<progress`,
		`class="ui-progress"`,
		`value="65"`,
		`aria-label="Indeterminate`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("progress docs are missing %q", contract)
		}
	}
}

func TestProgressDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/progress", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST progress status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

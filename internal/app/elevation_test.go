package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestElevationDocsRouteDogfoodsTokenMappedUtilityLevels(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/elevation", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("elevation docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Elevation</h1>`,
		`href="/components/elevation"`,
		`aria-label="Elevation example"`,
		`class="ui-elevation-0"`,
		`class="ui-elevation-1"`,
		`class="ui-elevation-2"`,
		`class="ui-elevation-3"`,
		`class="ui-elevation-4"`,
		`class="ui-elevation-5"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("elevation docs are missing %q", contract)
		}
	}
}

func TestElevationDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/elevation", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST elevation status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

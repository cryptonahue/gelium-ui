package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSliderDocsRouteDogfoodsNativeRangeAndStates(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/slider", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("slider docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Slider</h1>`,
		`href="/components/slider"`,
		`aria-label="Slider examples"`,
		`type="range"`,
		`class="ui-slider"`,
		`--ui-slider-fill`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("slider docs are missing %q", contract)
		}
	}
}

func TestSliderDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/slider", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST slider status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestAppJSEnhancesNativeRangeSliderFill(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/static/app.js", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("app.js status = %d, want %d", res.Code, http.StatusOK)
	}
	js := res.Body.String()
	for _, contract := range []string{
		`data-ui-slider`,
		`--ui-slider-fill`,
		`addEventListener("input"`,
		`closest`,
	} {
		if !strings.Contains(js, contract) {
			t.Errorf("app.js is missing slider enhancement contract %q", contract)
		}
	}
}

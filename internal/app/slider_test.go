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

// TestSliderDemoLabelsRangesFromVisibleCaptions closes gap G6: each range input
// must be named by its visible caption label, never overridden by an aria-label
// that shadows the visible text.
func TestSliderDemoLabelsRangesFromVisibleCaptions(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/slider", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("slider docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if strings.Contains(body, `aria-label="`) && strings.Contains(body, `slider-demo-caption`) {
		for _, forbidden := range []string{"Unselected slider", "Populated slider", "Disabled slider"} {
			if strings.Contains(body, forbidden) {
				t.Errorf("slider demo must not override the visible caption with aria-label %q (G6)", forbidden)
			}
		}
	}
	for _, pair := range []struct{ id, label string }{
		{id: "slider-unselected", label: "Unselected"},
		{id: "slider-disabled", label: "Disabled"},
	} {
		label := `<label class="slider-demo-caption" for="` + pair.id + `">` + pair.label + `</label>`
		if !strings.Contains(body, label) {
			t.Errorf("slider demo must label %s with its visible caption %q", pair.id, label)
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

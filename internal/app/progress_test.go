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
		`>Progress</h1>`,
		`href="/components/progress"`,
		`aria-label="Progress examples"`,
		`<progress`,
		`class="ui-progress"`,
		`value="65"`,
		`<label class="progress-demo-caption" for="progress-indeterminate">Indeterminate</label>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("progress docs are missing %q", contract)
		}
	}
}

// TestProgressDemoLabelsProgressFromVisibleCaptions closes gap G6: the demo
// must name each <progress> from its visible caption label, never overriding it
// with an aria-label that shadows the visible text.
func TestProgressDemoLabelsProgressFromVisibleCaptions(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/progress", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("progress docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, forbidden := range []string{
		`aria-label="Progress 30 percent"`,
		`aria-label="Progress 65 percent"`,
		`aria-label="Progress 100 percent"`,
		`aria-label="Indeterminate progress"`,
	} {
		if strings.Contains(body, forbidden) {
			t.Errorf("progress demo must not override the visible caption with aria-label %q (G6)", forbidden)
		}
	}
	for _, pair := range []struct{ id, label string }{
		{id: "progress-30", label: "30%"},
		{id: "progress-100", label: "100%"},
	} {
		label := `<label class="progress-demo-caption" for="` + pair.id + `">` + pair.label + `</label>`
		if !strings.Contains(body, label) {
			t.Errorf("progress demo must label %s with its visible caption %q", pair.id, label)
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

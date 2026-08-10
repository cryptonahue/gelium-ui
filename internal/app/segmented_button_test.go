package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSegmentedButtonDocsRouteDogfoodsStates(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/segmented-button", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("segmented button docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1>Segmented buttons</h1>`,
		`href="/components/segmented-button"`,
		`aria-label="Segmented button examples"`,
		`class="ui-segmented-button"`,
		`class="ui-segmented-button-set"`,
		`<fieldset`,
		`role="group"`,
		`<form method="get"`,
		`type="radio"`,
		`type="checkbox"`,
		`type="button"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("segmented button docs are missing %q", contract)
		}
	}
}

func TestSegmentedButtonDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/segmented-button", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST segmented button status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

// TestSegmentedButtonDocsSingleSelectUsesNativeRadioNoJS is the TDD guard for
// the roadmap contract "Single/multi select debe preferir radios/checkboxes sin
// JS": a single-select set is a native radio group sharing one name, with the
// selected segment derived from the native checked attribute — never from an
// aria-pressed toggle (the upstream Lit pattern, which needs component JS).
func TestSegmentedButtonDocsSingleSelectUsesNativeRadioNoJS(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/segmented-button", nil))

	body := res.Body.String()
	for _, contract := range []string{
		`<input type="radio" name="transport" value="walk"`,
		`<input type="radio" name="transport" value="transit"`,
		`<input type="radio" name="transport" value="drive" checked`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("single-select radio set is missing %q", contract)
		}
	}
	for _, forbidden := range []string{
		`aria-pressed="`,
		`onclick=`,
		`onchange=`,
	} {
		if strings.Contains(body, forbidden) {
			t.Errorf("segmented button docs must not rely on %q", forbidden)
		}
	}
}

// TestSegmentedButtonDocsMultiSelectUsesNativeCheckboxNoJS guards the
// multi-select contract: a native checkbox group sharing one name, with checked
// and disabled states expressed natively so the form submits the values without
// JavaScript.
func TestSegmentedButtonDocsMultiSelectUsesNativeCheckboxNoJS(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/segmented-button", nil))

	body := res.Body.String()
	for _, contract := range []string{
		`<input type="checkbox" name="formatting" value="bold" checked>`,
		`<input type="checkbox" name="formatting" value="italic" checked>`,
		`<input type="checkbox" name="formatting" value="underline" disabled>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("multi-select checkbox set is missing %q", contract)
		}
	}
}

// TestSegmentedButtonDocsActionSetUsesNativeButton guards the non-selection
// variant: the action set is a group of real <button type="button"> elements
// (no radio/checkbox, no checkmark) inside an accessible role="group" with an
// aria-label — the roadmap's "fieldset o grupo accesible" alternative.
func TestSegmentedButtonDocsActionSetUsesNativeButton(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/segmented-button", nil))

	body := res.Body.String()
	for _, contract := range []string{
		`role="group" aria-label="View actions"`,
		`<button type="button" class="ui-segmented-button ui-segmented-button--action">`,
		`<button type="button" class="ui-segmented-button ui-segmented-button--action" disabled>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("action button set is missing %q", contract)
		}
	}
}

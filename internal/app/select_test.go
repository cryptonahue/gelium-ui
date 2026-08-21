package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSelectDocsRouteDogfoodsNativeVariantsAndStates(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/select", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("select docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Select</h1>`,
		`href="/components/select"`,
		`>Examples</h2>`,
		`class="example-block"`,
		`<select`,
		`ui-select-filled`,
		`ui-select-outlined`,
		`<option`,
		`disabled`,
		`aria-invalid="true"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("select docs are missing %q", contract)
		}
	}
}

func TestSelectDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/select", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST select status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestSelectMenuDocsRouteDogfoodsNativeSelectField(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/select", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("select docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Examples</h2>`,
		`class="ui-select ui-select-filled"`,
		`<select id="select-menu" name="value">`,
		`<option value="standard">Standard</option>`,
		`<option value="priority" selected>Priority</option>`,
		`<option value="enterprise">Enterprise</option>`,
		`hx-post="/examples/select/menu"`,
		`hx-target="this"`,
		`hx-swap="outerHTML"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("select menu docs are missing %q", contract)
		}
	}
	for _, gone := range []string{`command=`, `commandfor=`, `closedby=`, `class="ui-select-menu"`, `aria-selected="true"`, `<dialog`} {
		if strings.Contains(body, gone) {
			t.Errorf("select menu docs must not contain dead M3 markup %q", gone)
		}
	}
}

func TestSelectMenuChangeNoHReflectsSelectionInNativeSelect(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/examples/select/menu", strings.NewReader("value=standard"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `<!doctype html>`) || !strings.Contains(body, `<title>Select · Gelium UI</title>`) {
		t.Error("no-HX response must be a complete documentation page")
	}
	menu := body[strings.Index(body, `id="select-menu-field"`):]
	if !strings.Contains(menu, `<option value="standard" selected>Standard</option>`) {
		t.Error("the chosen option must be marked selected in the native select")
	}
	if strings.Contains(menu, `<option value="priority" selected>Priority</option>`) {
		t.Error("the previously selected option must not stay selected")
	}
}

func TestSelectMenuChangeHXReturnsFragmentReflectingSelection(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/examples/select/menu", strings.NewReader("value=priority"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if strings.Contains(body, `<!doctype html>`) || strings.Contains(body, `<title>`) {
		t.Error("HX response must be a fragment, not a complete document")
	}
	if !strings.Contains(body, `class="ui-select ui-select-filled"`) {
		t.Errorf("fragment must return the native select field wrapper, got %q", body)
	}
	if !strings.Contains(body, `<select id="select-menu" name="value">`) {
		t.Error("fragment must return the native select control")
	}
	if !strings.Contains(body, `<option value="priority" selected>Priority</option>`) {
		t.Error("fragment must persist the newly selected option on the native select")
	}
}

func TestSelectMenuChangeRejectsUnknownValueWith422(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/examples/select/menu", strings.NewReader("value=not-a-plan"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	New().ServeHTTP(res, req)

	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnprocessableEntity)
	}
	if got := res.Header().Get("X-Gelium-Validation"); got != "true" {
		t.Errorf("X-Gelium-Validation = %q, want true", got)
	}
	body := res.Body.String()
	if !strings.Contains(body, "Select a valid option") {
		t.Error("422 fragment must carry a visible validation error")
	}
	if !strings.Contains(body, `aria-invalid="true"`) {
		t.Error("422 fragment must mark the native select as aria-invalid")
	}
}

// TestSelectPopupExperimentHTMXDisclosureContract protects the separately
// rendered HTMX-only experiment: the existing native Select remains the
// baseline while this disclosure uses only native details/summary and buttons.
func TestSelectPopupExperimentHTMXDisclosureContract(t *testing.T) {
	h := New()

	native := httptest.NewRecorder()
	h.ServeHTTP(native, httptest.NewRequest(http.MethodGet, "/components/select?execution=native", nil))
	if strings.Contains(native.Body.String(), `id="select-popup-experiment"`) {
		t.Error("native execution must not render the HTMX popup experiment")
	}
	invalidExecution := httptest.NewRecorder()
	h.ServeHTTP(invalidExecution, httptest.NewRequest(http.MethodGet, "/components/select?execution=not-real", nil))
	if strings.Contains(invalidExecution.Body.String(), `id="select-popup-experiment"`) {
		t.Error("invalid execution must resolve to native and omit the HTMX popup experiment")
	}

	page := httptest.NewRecorder()
	h.ServeHTTP(page, httptest.NewRequest(http.MethodGet, "/components/select?execution=htmx", nil))
	if page.Code != http.StatusOK {
		t.Fatalf("HTMX select docs status = %d, want %d", page.Code, http.StatusOK)
	}
	body := page.Body.String()
	popupStart := strings.Index(body, `id="select-popup-experiment"`)
	if popupStart < 0 {
		t.Fatal("HTMX select docs must render the popup experiment")
	}
	popup := body[popupStart:]
	for _, contract := range []string{
		`<details class="ui-select-popup-disclosure">`,
		`<summary>Choose a plan</summary>`,
		`action="/examples/select/popup#select-popup-experiment"`,
		`hx-post="/examples/select/popup"`,
		`hx-target="#select-popup-experiment"`,
		`hx-swap="outerHTML"`,
		`hx-sync="this:replace"`,
		`name="popup_value" value="standard"`,
		`role="status">Selected plan: Priority`,
	} {
		if !strings.Contains(popup, contract) {
			t.Errorf("HTMX experiment is missing %q", contract)
		}
	}
	if got := strings.Count(popup, `<form`); got != 1 {
		t.Errorf("popup experiment forms = %d, want one non-nested fallback form", got)
	}
	for _, forbidden := range []string{`role="listbox"`, `role="option"`, `<dialog`, `hx-on`, `onclick=`} {
		if strings.Contains(popup, forbidden) {
			t.Errorf("HTMX experiment must not fake popup behavior with %q", forbidden)
		}
	}

	fallback := httptest.NewRecorder()
	fallbackReq := httptest.NewRequest(http.MethodPost, "/examples/select/popup", strings.NewReader("popup_value=standard"))
	fallbackReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(fallback, fallbackReq)
	if fallback.Code != http.StatusSeeOther {
		t.Fatalf("non-HX popup POST status = %d, want %d", fallback.Code, http.StatusSeeOther)
	}
	if got := fallback.Header().Get("Location"); got != "/components/select?execution=htmx&select_popup_value=standard#select-popup-experiment" {
		t.Errorf("non-HX popup redirect = %q, want canonical experiment URL", got)
	}
	if got := fallback.Header().Get("Vary"); got != "HX-Request" {
		t.Errorf("non-HX popup Vary = %q, want HX-Request", got)
	}

	missingValue := httptest.NewRecorder()
	missingValueReq := httptest.NewRequest(http.MethodPost, "/examples/select/popup", nil)
	h.ServeHTTP(missingValue, missingValueReq)
	if missingValue.Code != http.StatusUnprocessableEntity {
		t.Errorf("non-HX popup POST without a submitted option status = %d, want %d", missingValue.Code, http.StatusUnprocessableEntity)
	}
	if got := missingValue.Header().Get("Vary"); got != "HX-Request" {
		t.Errorf("missing-value non-HX popup Vary = %q, want HX-Request", got)
	}
	if got := missingValue.Body.String(); !strings.Contains(got, `<!doctype html>`) || !strings.Contains(got, "Select a valid plan") {
		t.Error("missing-value non-HX popup POST must return the full 422 validation page")
	}

	hx := httptest.NewRecorder()
	hxReq := httptest.NewRequest(http.MethodPost, "/examples/select/popup", strings.NewReader("popup_value=enterprise"))
	hxReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	hxReq.Header.Set("HX-Request", "true")
	h.ServeHTTP(hx, hxReq)
	if hx.Code != http.StatusOK {
		t.Fatalf("HX popup POST status = %d, want %d", hx.Code, http.StatusOK)
	}
	if got := hx.Header().Get("Vary"); got != "HX-Request" {
		t.Errorf("HX popup Vary = %q, want HX-Request", got)
	}
	if strings.Contains(hx.Body.String(), `<!doctype html>`) || !strings.Contains(hx.Body.String(), `role="status">Selected plan: Enterprise`) {
		t.Error("HX popup POST must return the closed updated experiment fragment")
	}

	nonHTMXHeader := httptest.NewRecorder()
	nonHTMXHeaderReq := httptest.NewRequest(http.MethodPost, "/examples/select/popup", strings.NewReader("popup_value=standard"))
	nonHTMXHeaderReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	nonHTMXHeaderReq.Header.Set("HX-Request", "TRUE")
	h.ServeHTTP(nonHTMXHeader, nonHTMXHeaderReq)
	if nonHTMXHeader.Code != http.StatusSeeOther {
		t.Errorf("HX-Request TRUE status = %d, want non-HX %d", nonHTMXHeader.Code, http.StatusSeeOther)
	}
	if got := nonHTMXHeader.Header().Get("Vary"); got != "HX-Request" {
		t.Errorf("HX-Request TRUE Vary = %q, want HX-Request", got)
	}

	invalidHX := httptest.NewRecorder()
	invalidHXReq := httptest.NewRequest(http.MethodPost, "/examples/select/popup", strings.NewReader("popup_value=not-a-plan"))
	invalidHXReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	invalidHXReq.Header.Set("HX-Request", "true")
	h.ServeHTTP(invalidHX, invalidHXReq)
	if invalidHX.Code != http.StatusUnprocessableEntity {
		t.Fatalf("invalid HX popup POST status = %d, want %d", invalidHX.Code, http.StatusUnprocessableEntity)
	}
	if got := invalidHX.Header().Get("Vary"); got != "HX-Request" {
		t.Errorf("invalid HX popup Vary = %q, want HX-Request", got)
	}
	if got := invalidHX.Body.String(); !strings.Contains(got, `Select a valid plan`) || !strings.Contains(got, `Submitted value: not-a-plan`) || strings.Contains(got, `<!doctype html>`) {
		t.Error("invalid HX popup POST must return a swappable fragment with visible preserved value")
	}

	invalidFallback := httptest.NewRecorder()
	invalidFallbackReq := httptest.NewRequest(http.MethodPost, "/examples/select/popup", strings.NewReader("popup_value=not-a-plan"))
	invalidFallbackReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(invalidFallback, invalidFallbackReq)
	if invalidFallback.Code != http.StatusUnprocessableEntity {
		t.Fatalf("invalid non-HX popup POST status = %d, want %d", invalidFallback.Code, http.StatusUnprocessableEntity)
	}
	if got := invalidFallback.Body.String(); !strings.Contains(got, `<!doctype html>`) || !strings.Contains(got, `Select a valid plan`) || !strings.Contains(got, `Submitted value: not-a-plan`) {
		t.Error("invalid non-HX popup POST must return a full page with visible preserved value")
	}
}

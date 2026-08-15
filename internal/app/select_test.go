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
		`<h1>Select</h1>`,
		`href="/components/select"`,
		`aria-label="Select examples"`,
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
		`aria-label="Select menu example"`,
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

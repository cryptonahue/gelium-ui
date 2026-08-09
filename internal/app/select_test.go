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

func TestSelectMenuDocsRouteDogfoodsServerDrivenMenu(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/select", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("select docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`class="ui-select-menu"`,
		`aria-label="Select menu example"`,
		`command="show-modal"`,
		`commandfor="select-menu"`,
		`hx-post="/examples/select/menu"`,
		`hx-target="this"`,
		`hx-swap="outerHTML"`,
		`command="request-close"`,
		`class="ui-select-menu-item"`,
		`aria-selected="true"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("select menu docs are missing %q", contract)
		}
	}
}

func TestSelectMenuChangeNoHSetsValueAndRendersSessionState(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/examples/select/menu", strings.NewReader("value=priority&id=select-menu"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, `<!doctype html>`) || !strings.Contains(body, `<title>Select · Loom UI</title>`) {
		t.Error("no-HX response must be a complete documentation page")
	}
	if !strings.Contains(body, `value="priority"`) {
		t.Error("selected value must be reflected in the session state")
	}
	if !strings.Contains(body, `aria-selected="true"`) {
		t.Error("selected menu item must carry aria-selected")
	}
}

func TestSelectMenuChangeHXReturnsFragmentReflectingSelection(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/examples/select/menu", strings.NewReader("value=priority&id=select-menu"))
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
	if !strings.Contains(body, `class="ui-select m3-select"`) {
		t.Errorf("fragment must return the M3 select wrapper, got %q", body)
	}
	if !strings.Contains(body, `value="priority"`) {
		t.Error("fragment must persist the newly selected value on the hidden input")
	}
	if strings.Contains(body, `aria-selected="true"`) {
		t.Error("closed-menu fragment must not contain the open menu list")
	}
}

func TestSelectMenuChangeRejectsUnknownValueWith422(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/examples/select/menu", strings.NewReader("value=not-a-plan&id=select-menu"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	New().ServeHTTP(res, req)

	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnprocessableEntity)
	}
	if got := res.Header().Get("X-Loom-Validation"); got != "true" {
		t.Errorf("X-Loom-Validation = %q, want true", got)
	}
	body := res.Body.String()
	if !strings.Contains(body, "Select a valid option") {
		t.Error("422 fragment must carry a visible validation error")
	}
}

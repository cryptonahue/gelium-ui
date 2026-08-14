package app

// Wire compatibility guard. The Gelium rename deliberately keeps the legacy
// loom:* / X-Loom-* wire contracts frozen so existing consumers (the served
// HTMX hook in web/static/app.js and any server integration) keep working.
// See docs/gelium-ui-wire-compatibility.md. These assertions re-pin the frozen
// contracts so the rename is never silently reverted.

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWireCompatValidationHeaderStaysXLoomValidation(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/examples/text-field/validate", strings.NewReader("name=+++++"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnprocessableEntity)
	}
	if got := res.Header().Get("X-Loom-Validation"); got != "true" {
		t.Errorf("frozen X-Loom-Validation header = %q, want \"true\"", got)
	}
}

func TestWireCompatToastTriggerKeyStaysLoomToast(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/examples/toast/demo", strings.NewReader("message=Record+updated&type=success"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if got := res.Header().Get("HX-Trigger"); !strings.Contains(got, `"loom:toast"`) {
		t.Errorf("frozen HX-Trigger payload must keep the loom:toast key, got %q", got)
	}
}

func TestWireCompatToastRegionKeepsLoomId(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/toast", nil))
	body := res.Body.String()
	if !strings.Contains(body, `id="loom-toast-region"`) {
		t.Error("frozen toast live region must keep id=\"loom-toast-region\"")
	}
	if !strings.Contains(body, `data-loom-toast-dismiss`) {
		t.Error("frozen toast dismiss button must keep data-loom-toast-dismiss")
	}
}

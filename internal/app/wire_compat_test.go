package app

// Wire contract guard. The wire contracts use the gelium:* / X-Gelium-*
// names (the product rename migrated them in v0.4.x; see
// docs/gelium-ui-wire-compatibility.md). These assertions re-pin the
// canonical contracts so they are never silently reverted.

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

func TestWireCompatToastTriggerKeyStaysGeliumToast(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/examples/toast/demo", strings.NewReader("message=Record+updated&type=success"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if got := res.Header().Get("HX-Trigger"); !strings.Contains(got, `"gelium:toast"`) {
		t.Errorf("HX-Trigger payload must keep the gelium:toast key, got %q", got)
	}
}

func TestWireCompatToastRegionKeepsGeliumId(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/toast", nil))
	body := res.Body.String()
	if !strings.Contains(body, `id="gelium-toast-region"`) {
		t.Error("toast live region must keep id=\"gelium-toast-region\"")
	}
	if !strings.Contains(body, `data-gelium-toast-dismiss`) {
		t.Error("toast dismiss button must keep data-gelium-toast-dismiss")
	}
}

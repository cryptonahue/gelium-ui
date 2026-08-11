package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDialogDocsRouteDogfoodsPageVariantTriggerLink(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/dialog", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<title>Dialog · Gelium UI</title>`, `<h1>Dialog</h1>`, `href="/components/dialog"`,
		`<a class="ui-button ui-button-primary" href="/components/dialog/confirm">`,
		`<span>Open confirmation dialog</span></a>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("dialog docs are missing %q", contract)
		}
	}
	preview := body[strings.Index(body, `aria-label="Dialog example"`):]
	for _, forbidden := range []string{`command=`, `commandfor=`, `closedby=`, `<dialog`, `ui-dialog-`} {
		if strings.Contains(preview, forbidden) {
			t.Errorf("dialog preview must not ship inert command-based controls: %q", forbidden)
		}
	}
}

func TestDialogConfirmRouteRendersInlineAction(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/dialog/confirm", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<title>Dialog · Gelium UI</title>`, `<h1>Dialog</h1>`,
		`class="ui-dialog-page"`,
		`id="confirm-dialog-title"`, `id="confirm-dialog-description"`,
		`<form method="post" action="/components/dialog/confirm"`,
		`<input type="hidden" name="action" value="confirm">`,
		`<a class="ui-button ui-button-text" href="/components/dialog">`,
		`<span>Cancel</span></a>`,
		`class="ui-button ui-button-text" type="submit"`,
		`<span>Confirm</span></button>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("dialog confirm page is missing %q", contract)
		}
	}
}

func TestDialogConfirmPostRedirectsBackWith303(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/components/dialog/confirm", strings.NewReader("action=confirm"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	New().ServeHTTP(res, req)

	if res.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusSeeOther)
	}
	if got := res.Header().Get("Location"); got != "/components/dialog?confirmed=1" {
		t.Errorf("Location = %q, want /components/dialog?confirmed=1", got)
	}
}

func TestDialogDocsShowPersistentAlertAfterConfirmed(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/dialog?confirmed=1", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`class="ui-inline-alert ui-inline-alert--success"`,
		`role="status"`,
		`class="ui-inline-alert-body">Action confirmed.`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("dialog docs after confirm are missing %q", contract)
		}
	}
}

func TestDialogDocsExplainPageVariantAndModalCompatibility(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/dialog", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}

	body := res.Body.String()
	for _, contract := range []string{
		"page variant",
		"supporting browsers",
		"Baseline 2025",
		"no component JavaScript",
		"opt-in",
		"request-close",
		"not Baseline",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("dialog docs are missing compatibility contract %q", contract)
		}
	}
}

func TestDialogDocsRouteKeepsMethodAndUnknownRouteSemantics(t *testing.T) {
	for _, tt := range []struct {
		method, path string
		want         int
	}{
		{http.MethodPost, "/components/dialog", http.StatusMethodNotAllowed},
		{http.MethodGet, "/components/dialog/confirm/missing", http.StatusNotFound},
		{http.MethodGet, "/components/dialog/missing", http.StatusNotFound},
	} {
		res := httptest.NewRecorder()
		New().ServeHTTP(res, httptest.NewRequest(tt.method, tt.path, nil))
		if res.Code != tt.want {
			t.Errorf("%s %s status = %d, want %d", tt.method, tt.path, res.Code, tt.want)
		}
	}
}

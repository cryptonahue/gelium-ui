package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDialogDocsRouteDogfoodsNativeDeclarativeDialog(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/dialog", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<title>Dialog · Gelidium UI</title>`, `<h1>Dialog</h1>`, `href="/components/dialog"`,
		`id="confirm-dialog-title"`, `id="confirm-dialog-description"`,
		`command="show-modal" commandfor="confirm-dialog"`,
		`command="request-close" commandfor="confirm-dialog"`,
		`command="close" commandfor="confirm-dialog" value="confirm"`,
		`class="ui-button ui-button-text"`, `autofocus`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("dialog docs are missing %q", contract)
		}
	}
	dialog := openingTagWithID(t, body, "dialog", "confirm-dialog")
	for _, attribute := range []string{`closedby="any"`, `aria-labelledby="confirm-dialog-title"`, `aria-describedby="confirm-dialog-description"`} {
		if !strings.Contains(dialog, attribute) {
			t.Errorf("dialog = %q, want %s", dialog, attribute)
		}
	}
	for _, forbidden := range []string{" open", " role=", " aria-modal=", " tabindex="} {
		if strings.Contains(dialog, forbidden) {
			t.Errorf("dialog = %q, must omit redundant attribute %s", dialog, forbidden)
		}
	}
	for _, id := range []string{"confirm-dialog", "confirm-dialog-title", "confirm-dialog-description"} {
		if got := strings.Count(body, `id="`+id+`"`); got != 1 {
			t.Errorf("id %q occurs %d times, want exactly once", id, got)
		}
	}
}

func TestDialogDocsExplainProgressiveBrowserCompatibility(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/dialog", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}

	body := res.Body.String()
	for _, contract := range []string{
		"supporting browsers",
		"Baseline Low",
		"no component JavaScript fallback",
		"server-rendered fallback or adapter",
		"request-close",
		"newer than the invoker commands",
		"not Baseline",
		"Chromium-only",
		"instant or asymmetric",
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
		{http.MethodGet, "/components/dialog/missing", http.StatusNotFound},
	} {
		res := httptest.NewRecorder()
		New().ServeHTTP(res, httptest.NewRequest(tt.method, tt.path, nil))
		if res.Code != tt.want {
			t.Errorf("%s %s status = %d, want %d", tt.method, tt.path, res.Code, tt.want)
		}
	}
}

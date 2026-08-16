package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTabsDocsRouteServesBothVariantsAndSelectedState(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/tabs", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("tabs docs status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`>Tabs</h1>`,
		`href="/components/tabs"`,
		`aria-label="Tabs examples"`,
		`class="ui-tab ui-tab-primary"`,
		`class="ui-tab ui-tab-secondary"`,
		`ui-tab-stacked`,
		`ui-tabs-list`,
		`ui-tab-indicator`,
		`href="/components/tabs?tab=photos" aria-current="page"`,
		`href="/components/tabs?sub=travel" aria-current="page"`,
		`tabs-demo-grid`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("tabs docs are missing %q", contract)
		}
	}
}

func TestTabsDocsRouteSelectsPrimaryTabFromQuery(t *testing.T) {
	for _, tab := range []string{"videos", "music"} {
		res := httptest.NewRecorder()
		New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/tabs?tab="+tab, nil))

		body := res.Body.String()
		selected := `href="/components/tabs?tab=` + tab + `" aria-current="page"`
		if !strings.Contains(body, selected) {
			t.Errorf("?tab=%s must mark that tab selected, missing %q", tab, selected)
		}
		if strings.Contains(body, `href="/components/tabs?tab=photos" aria-current="page"`) {
			t.Errorf("?tab=%s must not keep photos selected", tab)
		}
	}
}

func TestTabsDocsRouteSelectsSecondaryTabFromQuery(t *testing.T) {
	for _, sub := range []string{"hotel", "activities"} {
		res := httptest.NewRecorder()
		New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/components/tabs?sub="+sub, nil))

		body := res.Body.String()
		selected := `href="/components/tabs?sub=` + sub + `" aria-current="page"`
		if !strings.Contains(body, selected) {
			t.Errorf("?sub=%s must mark that tab selected, missing %q", sub, selected)
		}
		if strings.Contains(body, `href="/components/tabs?sub=travel" aria-current="page"`) {
			t.Errorf("?sub=%s must not keep travel selected", sub)
		}
	}
}

func TestTabsDocsRouteFallsBackToDefaultsForUnknownSelection(t *testing.T) {
	for _, target := range []string{
		"/components/tabs?tab=bogus",
		"/components/tabs?sub=bogus",
		"/components/tabs?tab=%22%3E%3Cscript%3Ealert(1)%3C%2Fscript%3E",
	} {
		res := httptest.NewRecorder()
		New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, target, nil))

		if res.Code != http.StatusOK {
			t.Fatalf("status for %s = %d, want %d", target, res.Code, http.StatusOK)
		}
		body := res.Body.String()
		for _, contract := range []string{
			`href="/components/tabs?tab=photos" aria-current="page"`,
			`href="/components/tabs?sub=travel" aria-current="page"`,
		} {
			if !strings.Contains(body, contract) {
				t.Errorf("unknown %s must fall back to defaults; missing %q", target, contract)
			}
		}
		if strings.Contains(body, `href="/components/tabs?tab=bogus"`) ||
			strings.Contains(body, `href="/components/tabs?sub=bogus"`) ||
			strings.Contains(body, "alert(1)") {
			t.Errorf("query input must never be reflected for %s", target)
		}
	}
}

func TestTabsDocsRouteKeepsOnlyGETSemantics(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/components/tabs", nil))

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST tabs status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

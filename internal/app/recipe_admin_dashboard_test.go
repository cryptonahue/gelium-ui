package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecipeAdminDashboardRendersServerMetrics(t *testing.T) {
	resetRecipeResourceStore()
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-dashboard", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1 class="recipe-ad-title">Projects dashboard</h1>`,
		`<section class="recipe-ad-metrics" aria-labelledby="recipe-ad-metrics-title">`,
		`<article class="ui-card ui-card-outlined recipe-ad-metric">`,
		`Total projects`, `12`, `Active projects`, `Pending projects`, `Done projects`,
 `<section class="recipe-ad-recent" aria-labelledby="recipe-ad-recent-title">`,
 `<ul class="ui-list">`,
 `<span class="ui-list-item-headline">Mu navigation</span>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("dashboard does not contain contract %q", contract)
		}
	}
}

func TestRecipeAdminDashboardEmptyStateAndGETOnly(t *testing.T) {
	resetRecipeResourceStore()
	for _, item := range resourceDemoStore.snapshot() {
		resourceDemoStore.delete(item.ID)
	}
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-dashboard", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `No project metrics available`) {
		t.Fatalf("empty dashboard status/body = %d/%q", res.Code, res.Body.String())
	}
	res = httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/recipes/admin-dashboard", nil))
	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST dashboard status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
}

func TestRecipeAdminDashboardServerStatesAndTheme(t *testing.T) {
	resetRecipeResourceStore()
	for _, tc := range []struct {
		state  string
		marker string
	}{
		{"loading", `aria-busy="true"`},
		{"error", `role="alert"`},
	} {
		res := httptest.NewRecorder()
		New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-dashboard?state="+tc.state, nil))
		if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), tc.marker) {
			t.Errorf("state %s status/body missing %q: %d", tc.state, tc.marker, res.Code)
		}
	}
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-dashboard?theme=dark&scheme=dark", nil))
	if !strings.Contains(res.Body.String(), `class="theme-material theme-dark"`) {
		t.Error("dashboard must honor the server-driven dark theme query")
	}
}

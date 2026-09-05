package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecipeAdminResourceActivityIsReadOnlyAndProjectScoped(t *testing.T) {
	resetRecipeResourceStore()
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource/alpha/activity", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{"Alpha release activity", "Project moved to Active", "Release checklist reviewed", "All activity", "Read-only activity"} {
		if !strings.Contains(body, contract) {
			t.Errorf("activity page missing %q", contract)
		}
	}
	if strings.Contains(body, "Rollout window scheduled") {
		t.Error("activity page must not leak another project's activity")
	}
}

func TestRecipeAdminResourceActivityFiltersClosedType(t *testing.T) {
	resetRecipeResourceStore()
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource/alpha/activity?type=status", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if !strings.Contains(body, "Project moved to Active") || strings.Contains(body, "Release checklist reviewed") {
		t.Error("status filter must include only status activity")
	}
}

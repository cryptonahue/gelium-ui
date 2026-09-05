package app

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRecipeAdminResourceTasksListIsProjectScoped(t *testing.T) {
	resetRecipeResourceStore()
	resetRecipeTaskStore()
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/admin-resource/alpha/tasks", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{"Alpha release tasks", "Confirm release checklist", "Publish release notes", "Add task", "name=\"status\""} {
		if !strings.Contains(body, contract) {
			t.Errorf("tasks page missing %q", contract)
		}
	}
	if strings.Contains(body, "Review rollout metrics") {
		t.Error("tasks page must not leak another project's tasks")
	}
}

func TestRecipeAdminResourceTaskCreateUsesPost303And422(t *testing.T) {
	resetRecipeResourceStore()
	resetRecipeTaskStore()
	form := url.Values{"title": {"Document rollout"}, "status": {"Pending"}, "due_date": {"2026-02-01"}, "assignee": {"Alicia R."}}
	req := httptest.NewRequest(http.MethodPost, "/recipes/admin-resource/alpha/tasks", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)
	if res.Code != http.StatusSeeOther || res.Header().Get("Location") != "/recipes/admin-resource/alpha/tasks" {
		t.Fatalf("create response = %d %q, want 303 to tasks", res.Code, res.Header().Get("Location"))
	}

	bad := httptest.NewRequest(http.MethodPost, "/recipes/admin-resource/alpha/tasks", strings.NewReader("title=&status=Unknown&due_date="))
	bad.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	badRes := httptest.NewRecorder()
	New().ServeHTTP(badRes, bad)
	if badRes.Code != http.StatusUnprocessableEntity {
		t.Fatalf("invalid create status = %d, want 422", badRes.Code)
	}
}

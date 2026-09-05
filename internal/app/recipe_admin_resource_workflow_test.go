package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecipeAdminResourceWorkflowTransitionsWithPost303(t *testing.T) {
	resetRecipeResourceStore()
	form := strings.NewReader("status=Active")
	req := httptest.NewRequest(http.MethodPost, "/recipes/admin-resource/beta/transition", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)
	if res.Code != http.StatusSeeOther || res.Header().Get("Location") != "/recipes/admin-resource/beta" {
		t.Fatalf("transition response = %d %q, want 303 to detail", res.Code, res.Header().Get("Location"))
	}
	item, _ := resourceDemoStore.get("beta")
	if item.Status != "Active" {
		t.Fatalf("status = %q, want Active", item.Status)
	}
}

func TestRecipeAdminResourceWorkflowRejectsInvalidTransition(t *testing.T) {
	resetRecipeResourceStore()
	req := httptest.NewRequest(http.MethodPost, "/recipes/admin-resource/gamma/transition", strings.NewReader("status=Active"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)
	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", res.Code)
	}
	if !strings.Contains(res.Body.String(), "That status transition is not allowed") {
		t.Error("invalid transition should explain the blocked state")
	}
}

package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecipeAdminResourceRelatedActionsHonorConsumerAuthorization(t *testing.T) {
	resetRecipeResourceStore()
	resetRecipeTaskStore()
	s := newWithRecipeAdminAuthorizer(func(_ *http.Request, action string, _ *recipeResource) bool {
		return action == recipeAdminDeleteAction
	})

	for _, path := range []string{
		"/recipes/admin-resource/alpha/tasks",
		"/recipes/admin-resource/alpha/activity",
		"/recipes/admin-resource/alpha/transition",
	} {
		method := http.MethodGet
		if strings.HasSuffix(path, "transition") {
			method = http.MethodPost
		}
		req := httptest.NewRequest(method, path, strings.NewReader("status=Active"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)
		if res.Code != http.StatusForbidden {
			t.Errorf("%s %s status = %d, want %d", method, path, res.Code, http.StatusForbidden)
		}
	}
}

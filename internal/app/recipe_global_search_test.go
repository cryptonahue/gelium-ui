package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecipeGlobalSearchGroupsCrossResourceResults(t *testing.T) {
	resetRecipeResourceStore()
	resetRecipeQueueStore()
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/search?q=Bob", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{"Global search", "Beta rollout", "Ops Queue", "Credit note for cancelled plan", "/recipes/admin-resource/beta", "/recipes/ops-queue/billing-credit-12"} {
		if !strings.Contains(body, contract) {
			t.Errorf("search missing %q", contract)
		}
	}
}

func TestRecipeGlobalSearchNoResultsIsNavigable(t *testing.T) {
	res := httptest.NewRecorder()
	New().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/recipes/search?q=does-not-exist", nil))
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "No results found") {
		t.Fatalf("no-result search response = %d, body missing recovery state", res.Code)
	}
}

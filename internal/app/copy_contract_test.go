package app

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// Copy contract (content style guide, web/content/handbook-content-style.md):
// every recipe error message follows the action pattern — the fix is the
// message ("Enter the project name."), never a consequence-only "Name is
// required." — and every recipe empty state pairs the reason with a real next
// step. The style guide is enforced here and by the Handbook route tests
// (handbook_test.go).
func TestRecipeErrorCopyUsesActionPattern(t *testing.T) {
	resetRecipeResourceStore()

	res := postRecipeAdminResource(t, "/recipes/admin-resource", url.Values{
		"name":   {"  "},
		"status": {"On hold"},
		"date":   {""},
		"owner":  {"Alicia R."},
	})
	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnprocessableEntity)
	}
	body := res.Body.String()
	for _, want := range []string{
		`href="#recipe-ar-name-error">Enter the project name.</a>`,
		`href="#recipe-ar-status-error">Choose a status from the list.</a>`,
		`href="#recipe-ar-date-error">Enter the project date.</a>`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("422 form is missing the action-pattern error %q", want)
		}
	}
	for _, banned := range []string{
		"Name is required", "Choose a valid status", "Date is required",
		"is required.", "Invalid input", "is invalid",
	} {
		if strings.Contains(body, banned) {
			t.Errorf("422 form must not contain consequence-only copy %q", banned)
		}
	}
}

// TestRecipeEmptyStatesCarryActionLanguage proves the recipe empty states
// answer "what is here / why it is empty / what to do next" with real action
// language, never a bare "No data" (content style guide §Empty states).
func TestRecipeEmptyStatesCarryActionLanguage(t *testing.T) {
	t.Run("admin search no match", func(t *testing.T) {
		resetRecipeResourceStore()
		body := getOKBody(t, "/recipes/admin-resource?q=zzzznothing")
		for _, want := range []string{`No results`, "Try adjusting the filters."} {
			if !strings.Contains(body, want) {
				t.Errorf("admin search empty state is missing %q", want)
			}
		}
	})

	t.Run("admin empty dataset", func(t *testing.T) {
		resetRecipeResourceStore()
		for _, it := range resourceDemoStore.snapshot() {
			resourceDemoStore.delete(it.ID)
		}
		body := getOKBody(t, "/recipes/admin-resource")
		for _, want := range []string{`No projects yet`, "Create your first project to get started."} {
			if !strings.Contains(body, want) {
				t.Errorf("admin empty dataset is missing %q", want)
			}
		}
	})

	t.Run("feed following view", func(t *testing.T) {
		resetRecipeFeedStore()
		for _, it := range feedDemoStore.snapshot() {
			feedDemoStore.delete(it.ID)
		}
		body := getOKBody(t, "/recipes/public-feed?view=following")
		for _, want := range []string{`No posts from people you follow`, "Follow more people to fill this feed."} {
			if !strings.Contains(body, want) {
				t.Errorf("following empty state is missing %q", want)
			}
		}
	})
}

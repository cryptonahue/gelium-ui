package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRecipePublicFeedListRendersCardsWithAvatarAndStates proves the feed screen
// composes the primitives: Card items with a decorative Avatar, the "New" tone
// badge, the server-side Tabs view selector, the documented Skeleton loading
// placeholder, real react forms and server-side pagination.
func TestRecipePublicFeedListRendersCardsWithAvatarAndStates(t *testing.T) {
	resetRecipeFeedStore()

	res := getRecipe(t, "/recipes/public-feed")
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<h1 class="recipe-pf-title">Latest activity</h1>`,
		`<nav class="ui-tabs" aria-label="Feed views">`,
		`aria-current="page"`,
		`<article class="ui-card ui-card-outlined recipe-pf-card" aria-label="Post by Alicia R.">`,
		`class="ui-avatar ui-avatar--sm"`,
		`aria-hidden="true"`,
		`class="ui-avatar-initials">AR</span>`,
		`class="ui-badge ui-badge-large ui-badge--info">New</span>`,
		`method="post" action="/recipes/public-feed/post-01/react"`,
		`class="ui-skeleton ui-skeleton--avatar"`,
		`role="status"`,
		`id="feed-panel"`,
		`id="gelium-toast-region"`,
		`<span class="ui-pagination-page ui-pagination-page--current" aria-current="page">1</span>`,
		`6 posts · page 1 of 2`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("feed list is missing %q", contract)
		}
	}

	// Reverse chronological: the newest post leads the page.
	if idx01, idx02, idx03 := strings.Index(body, "The new design system guide"), strings.Index(body, "Rolling the ops queue recipe"), strings.Index(body, "Maintenance window"); idx01 < 0 || idx02 < 0 || idx03 < 0 || !(idx01 < idx02 && idx02 < idx03) {
		t.Errorf("feed must be reverse-chronological (post-01=%d post-02=%d post-03=%d)", idx01, idx02, idx03)
	}
}

// TestRecipePublicFeedSectionOrderContract proves purpose-bound feed regions
// render in the reader's order, independently of their visual skin classes.
func TestRecipePublicFeedSectionOrderContract(t *testing.T) {
	resetRecipeFeedStore()

	res := getRecipe(t, "/recipes/public-feed")
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	sections := []struct {
		name   string
		marker string
	}{
		{"page header", `>Latest activity</h1>`},
		{"feed view navigation", `aria-label="Feed views"`},
		{"feed list", `aria-label="Feed"`},
		{"post action zone", `aria-label="Post actions"`},
		{"loading and recovery", `aria-labelledby="recipe-pf-loading-heading"`},
		{"refresh action", `aria-label="Remote refresh"`},
	}

	previous := -1
	for _, section := range sections {
		index := strings.Index(body, section.marker)
		if index < 0 {
			t.Errorf("section contract is missing %s marker %q", section.name, section.marker)
			continue
		}
		if index < previous {
			t.Errorf("section %s rendered before its prerequisite", section.name)
		}
		previous = index
	}
	if list, cards := strings.Count(body, "<ol"), strings.Count(body, "<article"); list != 1 || cards == 0 {
		t.Errorf("feed contract requires one list containing repeated cards, lists=%d cards=%d", list, cards)
	}
}

// TestRecipePublicFeedViewsFilterServerSide proves the closed view vocabulary
// filters the feed server-side and marks the active Tab with aria-current.
func TestRecipePublicFeedViewsFilterServerSide(t *testing.T) {
	resetRecipeFeedStore()

	res := getRecipe(t, "/recipes/public-feed?view=following")
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{`href="?view=following"`, `3 posts · page 1 of 1`} {
		if !strings.Contains(body, contract) {
			t.Errorf("following view is missing %q", contract)
		}
	}
	if strings.Contains(body, "Dev Ops") {
		t.Error("following view must not include non-followed authors")
	}
	if strings.Contains(body, "ui-pagination") {
		t.Error("a single page of results must not render pagination")
	}

	res = getRecipe(t, "/recipes/public-feed?view=new")
	body = res.Body.String()
	if !strings.Contains(body, `3 posts · page 1 of 1`) {
		t.Errorf("new view must filter to the marked-new posts, caption %q", body)
	}

	// Unknown views sanitize to the default for-you.
	res = getRecipe(t, "/recipes/public-feed?view=bogus")
	body = res.Body.String()
	if !strings.Contains(body, `6 posts · page 1 of 2`) {
		t.Errorf("unknown view must sanitize to for-you, caption %q", body)
	}
}

// TestRecipePublicFeedPaginates proves page=2 shows the next slice of the
// reverse-chronological feed.
func TestRecipePublicFeedPaginates(t *testing.T) {
	resetRecipeFeedStore()

	res := getRecipe(t, "/recipes/public-feed?page=2")
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`<span class="ui-pagination-page ui-pagination-page--current" aria-current="page">2</span>`,
		`Skeleton states are the first thing I check`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("feed page 2 is missing %q", contract)
		}
	}
	if strings.Contains(body, "The new design system guide") {
		t.Error("feed page 2 must not contain a row from page 1")
	}
}

// TestRecipePublicFeedListHXFragment proves the HX-Request bifurcation: HTMX
// swaps only the #feed-panel fragment, not the whole page.
func TestRecipePublicFeedListHXFragment(t *testing.T) {
	resetRecipeFeedStore()

	req := httptest.NewRequest(http.MethodGet, "/recipes/public-feed?view=new", nil)
	req.Header.Set("HX-Request", "true")
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	if strings.Contains(body, "<html") || strings.Contains(body, "<!doctype html>") {
		t.Error("HX request must return the feed panel fragment, not a full document")
	}
	for _, contract := range []string{`id="feed-panel"`, `The new design system guide`} {
		if !strings.Contains(body, contract) {
			t.Errorf("feed panel fragment is missing %q", contract)
		}
	}
}

// TestRecipePublicFeedReact303 proves the like mutation follows POST+303 and the
// following render flashes the transient toast and increments the like count.
func TestRecipePublicFeedReact303(t *testing.T) {
	resetRecipeFeedStore()

	res := postRecipe(t, "/recipes/public-feed/post-01/react")
	if res.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusSeeOther)
	}
	if loc := res.Header().Get("Location"); loc != "/recipes/public-feed" {
		t.Errorf("Location = %q, want /recipes/public-feed", loc)
	}

	body := getRecipe(t, "/recipes/public-feed").Body.String()
	for _, contract := range []string{
		`class="ui-toast ui-toast-success"`,
		`You liked Alicia R.&#39;s post.`,
		`Like · 25</button>`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("react result is missing %q", contract)
		}
	}

	// The flash toast is consumed after one render.
	if body := getRecipe(t, "/recipes/public-feed").Body.String(); strings.Contains(body, "You liked") {
		t.Error("the flash toast must be consumed after one render")
	}
}

// TestRecipePublicFeedNotFound proves reacting to an unknown post renders the
// shared error-state page with the real 404 status.
func TestRecipePublicFeedNotFound(t *testing.T) {
	resetRecipeFeedStore()

	res := postRecipe(t, "/recipes/public-feed/nope/react")
	if res.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusNotFound)
	}
	if body := res.Body.String(); !strings.Contains(body, `class="ui-error-state"`) {
		t.Error("unknown post must render the error-state page")
	}
}

// TestRecipePublicFeedEmptyState proves the empty feed reuses the shared
// empty-state primitive with a real CTA.
func TestRecipePublicFeedEmptyState(t *testing.T) {
	resetRecipeFeedStore()
	for _, it := range feedDemoStore.snapshot() {
		feedDemoStore.delete(it.ID)
	}

	res := getRecipe(t, "/recipes/public-feed")
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{`No posts yet`, `Be the first to share something with the community.`, `href="/recipes/public-feed"`} {
		if !strings.Contains(body, contract) {
			t.Errorf("empty feed state is missing %q", contract)
		}
	}
}

// TestRecipePublicFeedRefreshPostOnly proves the refresh action is POST-only
// (GET answers 405 with Allow: POST), the no-JS refresh re-renders the page
// with a persistent inline toast + progress, and the HTMX refresh returns the
// fragment plus an HX-Trigger gelium:toast.
func TestRecipePublicFeedRefreshPostOnly(t *testing.T) {
	resetRecipeFeedStore()

	res := getRecipe(t, "/recipes/public-feed/refresh")
	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET refresh status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
	if got := res.Header().Get("Allow"); got != "POST" {
		t.Errorf("Allow = %q, want POST", got)
	}

	res = postRecipe(t, "/recipes/public-feed/refresh")
	if res.Code != http.StatusOK {
		t.Fatalf("POST refresh status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`Feed refreshed.`,
		`<progress value="100" max="100"></progress>`,
		`class="ui-toast ui-toast-success"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("no-JS refresh is missing %q", contract)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/recipes/public-feed/refresh", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	res = httptest.NewRecorder()
	New().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("HX refresh status = %d, want %d", res.Code, http.StatusOK)
	}
	if trigger := res.Header().Get("HX-Trigger"); !strings.Contains(trigger, "gelium:toast") || !strings.Contains(trigger, "Feed refreshed.") {
		t.Errorf("HX-Trigger = %q, want gelium:toast payload", trigger)
	}
	body = res.Body.String()
	if strings.Contains(body, "<html") {
		t.Error("HX refresh must return the refresh form fragment, not a full document")
	}
	if !strings.Contains(body, `class="recipe-pf-refresh"`) {
		t.Error("HX refresh must return the refresh form fragment")
	}
}

// TestRecipePublicFeedNoindexProvesRecipeSurfacesAreNeverIndexed proves the feed
// route emits robots noindex, nofollow plus a canonical and a per-route
// description in the recipe's own document head.
func TestRecipePublicFeedNoindexProvesRecipeSurfacesAreNeverIndexed(t *testing.T) {
	resetRecipeFeedStore()

	res := getRecipe(t, "/recipes/public-feed")
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	body := res.Body.String()
	for _, contract := range []string{
		`lang="en"`,
		`<meta name="robots" content="noindex, nofollow">`,
		`<link rel="canonical" href="https://gelium-ui.example/recipes/public-feed">`,
		`<title>Latest activity · Public/Social Feed recipe</title>`,
		`<meta name="description"`,
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("feed page is missing %q", contract)
		}
	}
}

package app

import (
	"strings"
	"testing"
)

// TestRecipesHonorChromeQuery proves the standalone recipe pages apply the
// document-root theme AND scheme selection exactly like the docs layout.
// Without this, sidebar navigation to a recipe silently flipped dark back to
// light (the recipe views only read ?theme=, never ?scheme=).
func TestRecipesHonorChromeQuery(t *testing.T) {
	for _, path := range []string{
		"/recipes/admin-resource",
		"/recipes/ops-queue",
		"/recipes/public-feed",
	} {
		t.Run(path, func(t *testing.T) {
			dark := getOKBody(t, path+"?theme=basecoat&scheme=dark")
			for _, contract := range []string{
				`class="theme-basecoat theme-dark" data-theme="dark"`,
			} {
				if !strings.Contains(dark, contract) {
					t.Errorf("recipe must render dark chrome %q", contract)
				}
			}
			plain := getOKBody(t, path)
			if strings.Contains(plain, "theme-dark") || strings.Contains(plain, `data-theme=`) {
				t.Error("recipe must render default chrome without theme/scheme")
			}
		})
	}
}

// TestRecipeChromeKeepsLightScheme proves the recipe pages emit the explicit
// light scheme (data-theme="light") like the docs layout when requested.
func TestRecipeChromeKeepsLightScheme(t *testing.T) {
	body := getOKBody(t, "/recipes/ops-queue?theme=material&scheme=light")
	if !strings.Contains(body, `class="theme-material" data-theme="light"`) {
		t.Error("recipe must render data-theme=light under ?scheme=light")
	}
}

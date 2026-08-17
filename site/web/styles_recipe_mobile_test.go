package web

import (
	"regexp"
	"strings"
	"testing"
)

// Recipes mobile layout contracts (roadmap: admin-resource ~646px / ops-queue
// ~780px expanded the viewport). Containment only — stack headers, wrap
// appbars, internal table scroll. Never overflow-x: hidden on the document root.

func TestRecipeAdminResourceMobileContainment(t *testing.T) {
	css := sourceComponentCSS(t, "recipe-admin-resource.css")

	if !strings.Contains(css, "flex-wrap: wrap") {
		t.Error("recipe-admin-resource.css must declare flex-wrap on appbar/header rows")
	}
	if !strings.Contains(css, "min-width: 0") {
		t.Error("recipe-admin-resource.css must declare min-width: 0 so flex children can shrink")
	}
	if !strings.Contains(css, "@media (max-width: 40rem)") {
		t.Error("recipe-admin-resource.css must stack the page header under max-width: 40rem")
	}
	start := strings.Index(css, "@media (max-width: 40rem)")
	if start < 0 {
		t.Fatal("missing 40rem media query")
	}
	media := css[start:]
	if !strings.Contains(media, ".recipe-ar-header") || !strings.Contains(media, "flex-direction: column") {
		t.Error("40rem media block must set .recipe-ar-header { flex-direction: column }")
	}
	if strings.Contains(css, "overflow-x: hidden") {
		t.Error("recipe-admin-resource.css must never mask overflow with overflow-x: hidden")
	}

	tmpl := repositoryFile(t, "site", "web", "templates", "recipe-admin-resource.html")
	if !regexp.MustCompile(`<div class="ui-data-table-scroll">\s*<table class="ui-data-table-table">`).MatchString(tmpl) {
		t.Error("recipe-admin-resource.html must wrap the table in .ui-data-table-scroll for internal horizontal scroll")
	}
}

func TestRecipeOpsQueueMobileContainment(t *testing.T) {
	css := sourceComponentCSS(t, "recipe-ops-queue.css")

	if !strings.Contains(css, "flex-wrap: wrap") {
		t.Error("recipe-ops-queue.css must declare flex-wrap on appbar/queue rows")
	}
	if !strings.Contains(css, "min-width: 0") {
		t.Error("recipe-ops-queue.css must declare min-width: 0 so flex children can shrink")
	}
	if !strings.Contains(css, ".recipe-oq-item") {
		t.Error("recipe-ops-queue.css must style .recipe-oq-item")
	}
	// Base rule lets the list row wrap trailing badge/actions.
	itemIdx := strings.Index(css, ".recipe-oq-item {")
	if itemIdx < 0 {
		t.Fatal("missing .recipe-oq-item rule")
	}
	itemBlock := css[itemIdx:]
	if end := strings.Index(itemBlock, "}"); end >= 0 {
		itemBlock = itemBlock[:end]
	}
	if !strings.Contains(itemBlock, "flex-wrap: wrap") {
		t.Error(".recipe-oq-item must flex-wrap so trailing slots do not expand the page")
	}
	if !strings.Contains(css, "@media (max-width: 48rem)") {
		t.Error("recipe-ops-queue.css must reflow queue rows under max-width: 48rem")
	}
	mediaStart := strings.Index(css, "@media (max-width: 48rem)")
	media := css[mediaStart:]
	if !strings.Contains(media, ".recipe-oq-item") {
		t.Error("48rem media block must target .recipe-oq-item")
	}
	if !strings.Contains(media, "flex: 1 1 100%") {
		t.Error("48rem media block must full-bleed trailing badge/actions (flex: 1 1 100%)")
	}
	if strings.Contains(css, "overflow-x: hidden") {
		t.Error("recipe-ops-queue.css must never mask overflow with overflow-x: hidden")
	}
}

func TestRecipePublicFeedMobileContainment(t *testing.T) {
	css := sourceComponentCSS(t, "recipe-public-feed.css")

	if !strings.Contains(css, ".recipe-pf-appbar") || !strings.Contains(css, "flex-wrap: wrap") {
		t.Error("recipe-public-feed.css appbar must flex-wrap")
	}
	if !strings.Contains(css, "min-width: 0") {
		t.Error("recipe-public-feed.css must declare min-width: 0")
	}
	if !strings.Contains(css, ".recipe-pf-card-actions") {
		t.Error("recipe-public-feed.css missing .recipe-pf-card-actions")
	} else {
		idx := strings.Index(css, ".recipe-pf-card-actions")
		block := css[idx:]
		if end := strings.Index(block, "}"); end >= 0 {
			block = block[:end]
		}
		if !strings.Contains(block, "flex-wrap: wrap") {
			t.Error(".recipe-pf-card-actions must flex-wrap")
		}
	}
	if strings.Contains(css, "overflow-x: hidden") {
		t.Error("recipe-public-feed.css must never mask overflow with overflow-x: hidden")
	}
}

func TestRecipeMobileContainmentReachesCompiledBundle(t *testing.T) {
	compiled := compiledAppCSS(t)
	for _, needle := range []string{
		"recipe-ar-appbar",
		"recipe-oq-item",
		"recipe-pf-card-actions",
		"ui-data-table-scroll",
	} {
		if !strings.Contains(compiled, needle) {
			t.Errorf("compiled app.css missing recipe mobile marker %q (run npm run build)", needle)
		}
	}
}

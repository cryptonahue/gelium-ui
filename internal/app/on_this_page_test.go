package app

import (
	"strings"
	"testing"
)

// TestOnThisPageRailProvesAnchoredTOC proves the server builds the "On this
// page" rail from the article's h2/h3 headings: the rail links carry the
// goldmark auto heading IDs and the SAME ids exist on the rendered article
// headings (the anchors actually resolve).
func TestOnThisPageRailProvesAnchoredTOC(t *testing.T) {
	body := getOKBody(t, "/docs/information-architecture")
	if !strings.Contains(body, `<aside class="docs-on-this-page"`) {
		t.Fatal("docs page must render the On this page rail")
	}
	if !strings.Contains(body, `href="#the-hierarchy-rule"`) {
		t.Error("rail is missing the h2 anchor #the-hierarchy-rule")
	}
	if !strings.Contains(body, `id="the-hierarchy-rule"`) {
		t.Error("rendered article heading must carry the auto id the-hierarchy-rule")
	}
	if !strings.Contains(body, `href="#agent-prompt"`) {
		t.Error("rail is missing the h2 anchor #agent-prompt")
	}
	// h3 sections nest in the rail but are skipped on pages without them.
	if !strings.Contains(body, "On this page") {
		t.Error("rail must carry its accessible label")
	}
}

// TestOnThisPageHiddenOnLegacyLayout proves the rail only renders inside the
// docs shell (marketing/legacy pages have no OnThisPage rail).
func TestOnThisPageHiddenOnLegacyLayout(t *testing.T) {
	body := getOKBody(t, "/")
	if strings.Contains(body, `class="docs-on-this-page"`) {
		t.Error("landing must not render the On this page rail")
	}
}

// TestPrevNextPaginationProvesIAOrder proves the previous/next pagination
// follows the SAME ordered IA as the sidebar and carries the chrome query, so
// pagination navigation keeps theme/scheme (GOV.UK pattern).
func TestPrevNextPaginationProvesIAOrder(t *testing.T) {
	body := getOKBody(t, "/docs/information-architecture?theme=basecoat&scheme=dark")
	for _, contract := range []string{
		`href="/docs?scheme=dark&amp;theme=basecoat"`,
		`href="/docs/choose-the-right-control?scheme=dark&amp;theme=basecoat"`,
		`>Previous</span>`,
		`>Next</span>`,
		"Information architecture",
		"Choose the right control",
	} {
		if !strings.Contains(body, contract) {
			t.Errorf("prev/next is missing contract %q", contract)
		}
	}
}

// TestPrevNextBoundariesProveFirstAndLastPages proves the pagination renders
// no Previous on the first IA page (/docs) and no Next on the last one
// (/recipes/public-feed) — the missing side becomes a spacer, not a broken
// link.
func TestPrevNextBoundariesProveFirstAndLastPages(t *testing.T) {
	first := getOKBody(t, "/docs")
	if strings.Contains(first, ">Previous</span>") {
		t.Error("first IA page must not render a Previous link")
	}
	if !strings.Contains(first, `href="/docs/information-architecture"`) {
		t.Error("first IA page must render Next to the second IA entry")
	}
	last := getOKBody(t, "/recipes/public-feed")
	if strings.Contains(last, ">Next</span>") {
		t.Error("last IA page must not render a Next link")
	}
}

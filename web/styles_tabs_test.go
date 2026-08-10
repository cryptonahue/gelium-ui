package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestTabsPrimitiveCSSMapsMaterialAnatomy(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "tabs.css"), " ")

	for _, contract := range []string{
		`.ui-tabs {`,
		`display: flex;`,
		`flex-direction: column;`,
		`.ui-tabs-list {`,
		`align-items: end;`,
		`.ui-tabs-item {`,
		`flex: 1;`,
		`.ui-tab {`,
		`height: var(--ui-tabs-height);`,
		`padding: 0 var(--ui-space-4);`,
		`gap: var(--ui-space-2);`,
		`font-weight: 500;`,
		`font-size: .875rem;`,
		`line-height: 1.25rem;`,
		`letter-spacing: .00625rem;`,
		`var(--ui-color-fg-muted)`,
		`.ui-tab-stacked {`,
		`flex-direction: column;`,
		`gap: 2px;`,
		`height: var(--ui-tabs-stacked-height);`,
		`.ui-tab-icon {`,
		`width: var(--ui-size-icon);`,
		`height: var(--ui-size-icon);`,
		`.ui-tab-indicator {`,
		`height: var(--ui-tabs-indicator-height);`,
		`border-radius: 3px 3px 0 0;`,
		`var(--ui-color-primary)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source tabs CSS is missing contract %q", contract)
		}
	}
}

func TestTabsVariantAndSelectionSelectors(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "tabs.css"), " ")

	for _, contract := range []string{
		`.ui-tab-secondary .ui-tab-indicator {`,
		`height: var(--ui-tabs-indicator-height-secondary);`,
		`border-radius: 0;`,
		`.ui-tab[aria-current="page"] .ui-tab-indicator { opacity: 1;`,
		`.ui-tab[aria-current="page"] {`,
		`color: var(--ui-color-primary);`,
		`.ui-tab-secondary[aria-current="page"] {`,
		`color: var(--ui-color-fg);`,
		`.ui-tab:hover::before { opacity: var(--ui-state-hover-opacity);`,
		`.ui-tab:focus-visible::before { opacity: var(--ui-state-focus-opacity);`,
		`.ui-tab:active::before { opacity: var(--ui-state-pressed-opacity);`,
		`.ui-tab:focus-visible {`,
		`outline: var(--ui-focus-thickness) solid var(--ui-color-focus-ring);`,
		`.ui-tab[aria-current="page"]::before { background: var(--ui-color-primary);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source tabs CSS is missing selection/state contract %q", contract)
		}
	}
}

// TestTabsDemoGridClassVocabularyIsClosed is the TDD regression for the closed
// vocabulary rule: the demo/preview scaffolding classes in tabs.html MUST match
// the CSS selectors exactly, and the ui- prefix is reserved for the component
// primitives (.ui-tab*), never for the demo grid.
func TestTabsDemoGridClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "web", "templates", "tabs.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "tabs.css"), " ")

	for _, cls := range []string{
		"tabs-demo-grid",
		"tabs-demo-group",
		"tabs-demo-caption",
		"tabs-demo-panel",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("tabs.html is missing demo class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("tabs.css is missing demo selector .%s", cls)
		}
	}

	for _, cls := range []string{
		"ui-tabs",
		"ui-tabs-list",
		"ui-tabs-item",
		"ui-tab",
		"ui-tab-stacked",
		"ui-tab-icon",
		"ui-tab-label",
		"ui-tab-indicator",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("tabs.html is missing primitive class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("tabs.css is missing primitive selector .%s", cls)
		}
	}

	// The variant is rendered server-side from a closed vocabulary
	// (ui-tab-{{$bar.Variant}}), so the template never hard-codes the class;
	// only the secondary variant needs variant-specific CSS overrides, because
	// primary tabs use the base .ui-tab styles. The Go handler tests prove both
	// variants render with class="ui-tab ui-tab-primary" / ui-tab-secondary.
	if !strings.Contains(tmpl, `ui-tab-{{$bar.Variant}}`) {
		t.Error("tabs.html must build the variant class from the closed $bar.Variant vocabulary")
	}
	if !strings.Contains(css, ".ui-tab-secondary") {
		t.Error("tabs.css must define the secondary variant selector")
	}
	if !strings.Contains(css, `[aria-current="page"]`) {
		t.Error("tabs.css must style selection through aria-current")
	}
}

func TestTabsReducedMotionAndForcedColorsWired(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "tabs.css"), " ")

	if !strings.Contains(css, `@media (prefers-reduced-motion: reduce)`) {
		t.Error("tabs.css must include a reduced-motion media query")
	}
	if !strings.Contains(css, `transition: none;`) {
		t.Error("tabs reduced-motion must disable transitions")
	}
	if !strings.Contains(css, `@media (forced-colors: active)`) {
		t.Error("tabs.css must include a forced-colors media query")
	}
	if !strings.Contains(css, `.ui-tab-indicator { background: CanvasText;`) {
		t.Error("tabs must keep the indicator visible in forced colors")
	}
	if !strings.Contains(css, `.ui-tab[aria-current="page"] { color: Highlight;`) {
		t.Error("tabs must keep the selected tab distinguishable in forced colors")
	}
	if !strings.Contains(css, `.ui-tab:focus-visible { outline-color: Highlight;`) {
		t.Error("tabs focus ring must switch to Highlight in forced colors")
	}
}

func TestEmbeddedCompiledCSSIncludesTabsContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-tabs`,
		`.ui-tabs-list`,
		`.ui-tabs-item`,
		`.ui-tab`,
		`.ui-tab-secondary`,
		`.ui-tab-stacked`,
		`.ui-tab-icon`,
		`.ui-tab-indicator`,
		`.tabs-demo-grid`,
		`.tabs-demo-group`,
		`.tabs-demo-caption`,
		`.tabs-demo-panel`,
		`[aria-current=page]`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled tabs CSS is missing %q", contract)
		}
	}
}

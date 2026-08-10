package web

import (
	"regexp"
	"strings"
	"testing"
)

func TestNavTabPrimitiveCSSMapsMaterialAnatomy(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-tab.css"), " ")

	for _, contract := range []string{
		`.ui-nav-tab {`,
		`flex-direction: column;`,
		`padding: var(--ui-space-2) 0 var(--ui-space-3);`,
		`min-width: 48px;`,
		`font: var(--ui-type-label-sm);`,
		`.ui-nav-tab-icon {`,
		`width: var(--ui-nav-tab-icon-container-width);`,
		`height: var(--ui-nav-tab-icon-container-height);`,
		`.ui-nav-tab-indicator {`,
		`border-radius: var(--ui-radius-full);`,
		`background: var(--ui-nav-tab-indicator-color);`,
		`.ui-nav-tab-glyph {`,
		`width: var(--ui-nav-tab-icon-size);`,
		`.ui-nav-tab-label {`,
		`height: var(--ui-nav-tab-label-height);`,
		`margin-top: var(--ui-nav-tab-label-space);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source navigation-tab CSS is missing contract %q", contract)
		}
	}
}

func TestNavTabActiveInactiveGlyphToggleUsesCSSOnly(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-tab.css"), " ")

	for _, contract := range []string{
		`.ui-nav-tab-glyph--active {`,
		`display: none;`,
		`.ui-nav-tab--active .ui-nav-tab-glyph {`,
		`.ui-nav-tab--active .ui-nav-tab-glyph--active {`,
		`display: inline-block;`,
		`.ui-nav-tab--active .ui-nav-tab-indicator {`,
		`opacity: 1;`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source navigation-tab CSS is missing glyph toggle contract %q", contract)
		}
	}
}

func TestNavTabStatesCoverHoverFocusActiveSelected(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-tab.css"), " ")

	for _, contract := range []string{
		`.ui-nav-tab::before {`,
		`background: var(--ui-nav-tab-state-layer-color);`,
		`.ui-nav-tab:hover::before {`,
		`opacity: var(--ui-nav-tab-hover-opacity);`,
		`.ui-nav-tab:focus-visible::before {`,
		`.ui-nav-tab:active::before {`,
		`opacity: var(--ui-nav-tab-pressed-opacity);`,
		`.ui-nav-tab:focus-visible {`,
		`outline: var(--ui-focus-thickness) solid var(--ui-color-focus-ring);`,
		`.ui-nav-tab:hover .ui-nav-tab-label,`,
		`.ui-nav-tab--active {`,
		`font-weight: var(--ui-nav-tab-active-label-weight);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source navigation-tab CSS is missing state contract %q", contract)
		}
	}
}

func TestNavTabHideInactiveLabelsCollapsesLabels(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-tab.css"), " ")

	for _, contract := range []string{
		`.ui-nav-tab--hide-inactive-label .ui-nav-tab-label {`,
		`height: 0;`,
		`opacity: 0;`,
		`.ui-nav-tab--hide-inactive-label.ui-nav-tab--active .ui-nav-tab-label {`,
		`height: var(--ui-nav-tab-label-height);`,
		`opacity: 1;`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source navigation-tab CSS is missing hide-inactive-label contract %q", contract)
		}
	}
}

func TestNavTabBadgeReuseComposesUiBadge(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-tab.css"), " ")
	tmpl := repositoryFile(t, "web", "templates", "navigation-tab.html")

	for _, contract := range []string{
		`.ui-nav-tab .ui-badge {`,
		`position: absolute;`,
		`inset-inline-start: 50%;`,
		`margin-inline-start: .375rem;`,
		`.ui-nav-tab .ui-badge.ui-badge-large {`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("navigation-tab CSS is missing badge composition contract %q", contract)
		}
	}
	// The template must compose the existing .ui-badge primitive rather than a
	// new nav-tab badge class.
	for _, contract := range []string{
		`class="ui-badge"`,
		`class="ui-badge ui-badge-large"`,
	} {
		if !strings.Contains(tmpl, contract) {
			t.Errorf("navigation-tab template is missing reused .ui-badge %q", contract)
		}
	}
}

// TestNavTabDemoGridClassVocabularyIsClosed is the TDD regression for the closed
// vocabulary rule: demo/preview scaffolding classes in navigation-tab.html MUST
// match the CSS selectors exactly, and the ui- prefix is reserved for the
// component primitives (.ui-nav-tab*), never for the demo grid.
func TestNavTabDemoGridClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "web", "templates", "navigation-tab.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-tab.css"), " ")

	for _, cls := range []string{
		"navigation-tab-demo-grid",
		"navigation-tab-demo-group",
		"navigation-tab-demo-heading",
		"navigation-tab-demo-row",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("navigation-tab.html is missing demo class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("navigation-tab.css is missing demo selector .%s", cls)
		}
	}

	for _, cls := range []string{
		"ui-nav-tab",
		"ui-nav-tab--active",
		"ui-nav-tab--hide-inactive-label",
		"ui-nav-tab-icon",
		"ui-nav-tab-indicator",
		"ui-nav-tab-label",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("navigation-tab.html is missing primitive class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("navigation-tab.css is missing primitive selector .%s", cls)
		}
	}

	// The glyph classes are carried by the trusted SVG constants in
	// navigation_tab.go (rendered inside the aria-hidden icon slot), so the CSS
	// selectors must exist here and the rendered page is asserted in the Go
	// route tests; they cannot appear as template literal text.
	for _, cls := range []string{
		"ui-nav-tab-glyph",
		"ui-nav-tab-glyph--active",
	} {
		if !strings.Contains(css, "."+cls) {
			t.Errorf("navigation-tab.css is missing primitive selector .%s", cls)
		}
	}
}

func TestNavTabReducedMotionAndForcedColorsWired(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-tab.css"), " ")

	if !strings.Contains(css, `@media (prefers-reduced-motion: reduce)`) {
		t.Error("navigation-tab.css must include a reduced-motion media query")
	}
	if !strings.Contains(css, `.ui-nav-tab::before,`) {
		t.Error("navigation-tab reduced-motion must disable the state layer transition")
	}
	if !strings.Contains(css, `@media (forced-colors: active)`) {
		t.Error("navigation-tab.css must include a forced-colors media query")
	}
	if !strings.Contains(css, `.ui-nav-tab--active {`) {
		t.Error("navigation-tab must keep an active foreground in forced colors")
	}
	if !strings.Contains(css, `.ui-nav-tab-indicator { border: 1px solid Highlight;`) {
		t.Error("the active indicator must stay visible in forced colors")
	}
}

func TestEmbeddedCompiledCSSIncludesNavTabContracts(t *testing.T) {
	compiled, err := Assets.ReadFile("static/app.css")
	if err != nil {
		t.Fatalf("read embedded compiled app CSS: %v", err)
	}
	css := string(compiled)
	for _, contract := range []string{
		`.ui-nav-tab{`,
		`.ui-nav-tab--active`,
		`.ui-nav-tab-icon`,
		`.ui-nav-tab-indicator`,
		`.ui-nav-tab-glyph`,
		`.ui-nav-tab-label`,
		`.ui-nav-tab--hide-inactive-label`,
		`.ui-nav-tab .ui-badge`,
		`.navigation-tab-demo-grid`,
		`.navigation-tab-demo-group`,
		`.navigation-tab-demo-heading`,
		`.navigation-tab-demo-row`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled navigation-tab CSS is missing %q", contract)
		}
	}
}

package lib

import (
	"regexp"
	"strings"
	"testing"
)

func TestNavBarPrimitiveCSSMapsMaterialAnatomy(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-bar.css"), " ")

	for _, contract := range []string{
		`.ui-nav-bar {`,
		`height: var(--ui-nav-bar-height);`,
		`background: var(--ui-nav-bar-container);`,
		`box-shadow: var(--ui-nav-bar-container-elevation);`,
		`.ui-nav-bar-list {`,
		`list-style: none;`,
		`.ui-nav-bar-item {`,
		`.ui-nav-bar-destination {`,
		`flex-direction: column;`,
		`padding: var(--ui-space-2) 0 var(--ui-space-3);`,
		`min-width: 48px;`,
		`font: var(--ui-type-label-sm);`,
		`.ui-nav-bar-icon {`,
		`width: var(--ui-nav-bar-icon-container-width);`,
		`height: var(--ui-nav-bar-icon-container-height);`,
		`.ui-nav-bar-indicator {`,
		`border-radius: var(--ui-radius-full);`,
		`background: var(--ui-nav-bar-indicator-color);`,
		`.ui-nav-bar-glyph {`,
		`width: var(--ui-nav-bar-icon-size);`,
		`.ui-nav-bar-label {`,
		`height: var(--ui-nav-bar-label-height);`,
		`margin-top: var(--ui-nav-bar-label-space);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source navigation-bar CSS is missing contract %q", contract)
		}
	}
}

func TestNavBarActiveInactiveGlyphToggleUsesCSSOnly(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-bar.css"), " ")

	for _, contract := range []string{
		`.ui-nav-bar-glyph--active {`,
		`display: none;`,
		`.ui-nav-bar-destination--active .ui-nav-bar-glyph {`,
		`.ui-nav-bar-destination--active .ui-nav-bar-glyph--active {`,
		`display: inline-block;`,
		`.ui-nav-bar-destination--active .ui-nav-bar-indicator {`,
		`opacity: 1;`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source navigation-bar CSS is missing glyph toggle contract %q", contract)
		}
	}
}

func TestNavBarStatesCoverHoverFocusActiveSelected(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-bar.css"), " ")

	for _, contract := range []string{
		`.ui-nav-bar-destination::before {`,
		`background: var(--ui-nav-bar-state-layer-color);`,
		`.ui-nav-bar-destination:hover::before {`,
		`opacity: var(--ui-nav-bar-hover-opacity);`,
		`.ui-nav-bar-destination:focus-visible::before {`,
		`.ui-nav-bar-destination:active::before {`,
		`opacity: var(--ui-nav-bar-pressed-opacity);`,
		`.ui-nav-bar-destination:focus-visible {`,
		`outline: var(--ui-focus-thickness) solid var(--ui-color-focus-ring);`,
		`.ui-nav-bar-destination:hover .ui-nav-bar-label,`,
		`.ui-nav-bar-destination--active {`,
		`font-weight: var(--ui-nav-bar-active-label-weight);`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source navigation-bar CSS is missing state contract %q", contract)
		}
	}
}

func TestNavBarHideInactiveLabelsCollapsesLabels(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-bar.css"), " ")

	for _, contract := range []string{
		`.ui-nav-bar--hide-inactive-labels .ui-nav-bar-label {`,
		`height: 0;`,
		`opacity: 0;`,
		`.ui-nav-bar--hide-inactive-labels .ui-nav-bar-destination--active .ui-nav-bar-label {`,
		`opacity: 1;`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("source navigation-bar CSS is missing hide-inactive-labels contract %q", contract)
		}
	}
}

func TestNavBarBadgeReuseComposesUiBadge(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-bar.css"), " ")
	tmpl := repositoryFile(t, "lib", "templates", "navigation-bar.html")

	for _, contract := range []string{
		`.ui-nav-bar .ui-badge {`,
		`position: absolute;`,
		`inset-inline-start: 50%;`,
		`margin-inline-start: .375rem;`,
		`.ui-nav-bar .ui-badge.ui-badge-large {`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("navigation-bar CSS is missing badge composition contract %q", contract)
		}
	}
	// The template must compose the existing .ui-badge primitive rather than a
	// new nav-bar badge class.
	for _, contract := range []string{
		`class="ui-badge"`,
		`class="ui-badge ui-badge-large"`,
	} {
		if !strings.Contains(tmpl, contract) {
			t.Errorf("navigation-bar template is missing reused .ui-badge %q", contract)
		}
	}
}

// TestNavBarDemoGridClassVocabularyIsClosed is the TDD regression for the closed
// vocabulary rule: demo/preview scaffolding classes in navigation-bar.html MUST
// match the CSS selectors exactly, and the ui- prefix is reserved for the
// component primitives (.ui-nav-bar*), never for the demo grid.
func TestNavBarDemoGridClassVocabularyIsClosed(t *testing.T) {
	tmpl := repositoryFile(t, "lib", "templates", "navigation-bar.html")
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-bar.css"), " ")

	for _, cls := range []string{
		"navigation-bar-demo-grid",
		"navigation-bar-demo-group",
		"navigation-bar-demo-heading",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("navigation-bar.html is missing demo class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("navigation-bar.css is missing demo selector .%s", cls)
		}
	}

	for _, cls := range []string{
		"ui-nav-bar",
		"ui-nav-bar-list",
		"ui-nav-bar-item",
		"ui-nav-bar-destination",
		"ui-nav-bar-destination--active",
		"ui-nav-bar-icon",
		"ui-nav-bar-indicator",
		"ui-nav-bar-label",
		"ui-nav-bar--hide-inactive-labels",
	} {
		if !strings.Contains(tmpl, cls) {
			t.Errorf("navigation-bar.html is missing primitive class %q", cls)
		}
		if !strings.Contains(css, "."+cls) {
			t.Errorf("navigation-bar.css is missing primitive selector .%s", cls)
		}
	}

	// The glyph classes are carried by the trusted SVG constants in
	// navigation_bar.go (rendered inside the aria-hidden icon slot), so the CSS
	// selectors must exist here and the rendered page is asserted in the Go
	// route tests; they cannot appear as template literal text.
	for _, cls := range []string{
		"ui-nav-bar-glyph",
		"ui-nav-bar-glyph--active",
	} {
		if !strings.Contains(css, "."+cls) {
			t.Errorf("navigation-bar.css is missing primitive selector .%s", cls)
		}
	}
}

func TestNavBarReducedMotionAndForcedColorsWired(t *testing.T) {
	css := regexp.MustCompile(`\s+`).ReplaceAllString(sourceComponentCSS(t, "navigation-bar.css"), " ")

	if !strings.Contains(css, `@media (prefers-reduced-motion: reduce)`) {
		t.Error("navigation-bar.css must include a reduced-motion media query")
	}
	if !strings.Contains(css, `.ui-nav-bar-destination::before,`) {
		t.Error("navigation-bar reduced-motion must disable the state layer transition")
	}
	if !strings.Contains(css, `@media (forced-colors: active)`) {
		t.Error("navigation-bar.css must include a forced-colors media query")
	}
	if !strings.Contains(css, `.ui-nav-bar { box-shadow: none; border: 1px solid CanvasText;`) {
		t.Error("navigation-bar must keep a visible boundary in forced colors")
	}
	if !strings.Contains(css, `.ui-nav-bar-indicator { border: 1px solid Highlight;`) {
		t.Error("the active indicator must stay visible in forced colors")
	}
}

func TestEmbeddedCompiledCSSIncludesNavBarContracts(t *testing.T) {
	css := compiledAppCSS(t)
	for _, contract := range []string{
		`.ui-nav-bar{`,
		`.ui-nav-bar-list`,
		`.ui-nav-bar-destination`,
		`.ui-nav-bar-destination--active`,
		`.ui-nav-bar-icon`,
		`.ui-nav-bar-indicator`,
		`.ui-nav-bar-glyph`,
		`.ui-nav-bar-label`,
		`.ui-nav-bar--hide-inactive-labels`,
		`.ui-nav-bar .ui-badge`,
		`.navigation-bar-demo-grid`,
		`.navigation-bar-demo-group`,
		`.navigation-bar-demo-heading`,
		`@media (forced-colors:active)`,
	} {
		if !strings.Contains(css, contract) {
			t.Errorf("embedded compiled navigation-bar CSS is missing %q", contract)
		}
	}
}
